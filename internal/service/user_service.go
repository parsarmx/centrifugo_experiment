package service

import (
	"context"
	"errors"
	"fmt"
	dao "golang_template/api/dto"
	"golang_template/internal/config"
	"golang_template/internal/otp"
	"golang_template/internal/producers"
	"golang_template/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	serviceDefaultTimeoutInSeconds = 10
	serviceDefaultTimeout          = serviceDefaultTimeoutInSeconds * time.Second
	otpDefaultTimeout              = serviceDefaultTimeoutInSeconds * time.Second
	loginDefaultTimeout            = 15 * time.Second
	leaderboardKey                 = "leaderboard:users"
)

type LeaderboardEntry struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Score    float64 `json:"score"`
}

type UserService interface {
	SendOtp(ctx context.Context, phoneNumber string) (userData dao.SendOTPResponseData, err error)
	OTPLogin(ctx context.Context, data dao.OTPLoginRequestData) (*dao.OTPLoginResponseData, error)
}

type userService struct {
	repo         repository.UserRepository
	redis        producers.RedisClient
	logger       *zap.Logger
	otpGenerator otp.OTPGenerator
	authConfig   config.AuthConfig
}

func NewUserService(
	repo repository.UserRepository,
	logger *zap.Logger,
	redis producers.RedisClient,
	authConfig config.AuthConfig,
) UserService {
	otpGenerator, err := otp.NewDefaultOTPGenerator(authConfig.OTPCodeLength, logger)
	if err != nil {
		logger.Error("Failed to create OTP generator", zap.Error(err))
		panic(err)
	}
	return &userService{
		repo:         repo,
		redis:        redis,
		logger:       logger,
		authConfig:   authConfig,
		otpGenerator: otpGenerator,
	}
}

func (s *userService) SendOtp(ctx context.Context, phoneNumber string) (userData dao.SendOTPResponseData, err error) {
	ctx, cancel := context.WithTimeout(ctx, otpDefaultTimeout)
	defer cancel()

	cacheKey := fmt.Sprintf("otp:%s", phoneNumber)

	// Check if OTP already exists
	ttlDuration, err := s.redis.TTL(ctx, cacheKey)
	if err == nil {
		expirationTime := time.Now().Add(ttlDuration)
		if expirationTime.After(time.Now()) {
			return dao.SendOTPResponseData{OtpTtl: ttlDuration.Seconds()}, ErrOTPSent
		}
	}
	if err == producers.ErrFailedToRetrieve {
		s.logger.Error("Error retrieving OTP code TTL", zap.String("phone_number", phoneNumber), zap.Error(err))
		return dao.SendOTPResponseData{}, ErrCacheRetrieval
	}

	// if phoneNumber == s.authConfig.TestUser.PhoneNumber {
	// 	testUserDao, err := s.testUserSendOTP(
	// 		ctx,
	// 		dto.SendOTPRequestData{PhoneNumber: phoneNumber})
	// 	if err != nil {
	// 		return dao.SendOTPResponseData{}, err
	// 	}
	// 	return testUserDao, nil
	// }

	var otpCode string
	otpCode, err = s.otpGenerator.GenerateOTPLowMemoryMode()
	if err != nil {
		s.logger.Debug("Error generating OTP code", zap.String("phone_number", phoneNumber), zap.Error(err))
		return dao.SendOTPResponseData{}, ErrOTPGeneration
	}

	// send OTP
	err = s.redis.Conn().SetArgs(ctx, cacheKey, otpCode, redis.SetArgs{
		Mode: "NX",
		TTL:  time.Duration(s.authConfig.OTPTTL) * time.Second,
	}).Err()

	if errors.Is(err, redis.Nil) {
		// in case of race condition, the OTP code is already sent
		s.logger.Debug("OTP code already sent", zap.String("phone_number", phoneNumber))
		return dao.SendOTPResponseData{}, ErrOTPSent
	}

	if err != nil {
		// u must log with ERROR!!!!!
		s.logger.Error("Error caching OTP code", zap.String("phone_number", phoneNumber), zap.Error(err))
		return dao.SendOTPResponseData{}, ErrCacheSet
	}
	if s.authConfig.UnderDevelopment {
		return dao.SendOTPResponseData{PhoneNumber: phoneNumber, OtpCode: otpCode}, nil
	}
	return dao.SendOTPResponseData{PhoneNumber: phoneNumber}, nil
}

func (s *userService) OTPLogin(ctx context.Context, data dao.OTPLoginRequestData) (*dao.OTPLoginResponseData, error) {
	ctx, cancel := context.WithTimeout(ctx, loginDefaultTimeout)
	defer cancel()
	otpCode := data.OTPCode
	phoneNumber := data.PhoneNumber

	err := s.validateOtpWithCacheOnly(ctx, phoneNumber, otpCode)
	if err != nil {
		return nil, err
	}

	userData, err := s.repo.UserDataByPhoneNumber(ctx, phoneNumber)

	issuedAt := time.Now()
	accExpiresAt := issuedAt.Add(time.Duration(s.authConfig.AccTokenExpTime) * time.Minute)
	sessionID := uuid.New()
	var jti uuid.UUID
	accessTokenData := dao.AccessToken{
		UserID:    userData.ID,
		SessionID: sessionID,
		Role:      "", // TODO retrieve and pass the associated role
		JTI:       jti,
		IssuedAt:  issuedAt,
		ExpiresAt: accExpiresAt,
	}
	accessToken, err := s.CreateAccessToken(accessTokenData)
	if err != nil {
		s.logger.Debug("Error creating access token",
			zap.String("phoneNumber", phoneNumber),
			zap.String("otpCode", otpCode),
			zap.Error(err))
		return nil, err
	}

	// Return response data
	return &dao.OTPLoginResponseData{
			RefreshTokenID: jti,
			UserID:         userData.ID,
			SessionID:      sessionID,
			AccessToken:    accessToken,
			AccExpiresAt:   accExpiresAt,
		},
		nil

}

func (s *userService) validateOtpWithCacheOnly(ctx context.Context, phoneNumber string, otpCode string) error {
	cacheKey := fmt.Sprintf("otp:%s", phoneNumber)
	cacheOtp, err := s.redis.Conn().Get(ctx, cacheKey).Result()

	if errors.Is(err, redis.Nil) {
		// otp not found in cache, expired or does not exist
		s.logger.Debug("OTP code not found in cache", zap.String("phone_number", phoneNumber))
		return ErrInvalidOTP
	}
	if err != nil {
		// redis error
		s.logger.Error("Error getting OTP code from cache", zap.String("phone_number", phoneNumber), zap.Error(err))
		return err
	}

	// otp found then check if it is valid
	if cacheOtp != otpCode {
		// otp is not valid
		s.logger.Debug("OTP code is not valid", zap.String("phone_number", phoneNumber))
		return ErrInvalidOTP
	}

	// delete otp from cache
	isDeletedCode, err := s.redis.Conn().Del(ctx, cacheKey).Result()
	if err != nil {
		s.logger.Error("Error deleting OTP code from cache", zap.String("phone_number", phoneNumber), zap.Error(err))
		return err
	}
	if isDeletedCode != 1 {
		// otp not found in cache, during validation another request has deleted it, so we can not authenticate the user
		s.logger.Debug("Error deleting OTP code from cache", zap.String("phone_number", phoneNumber), zap.Error(err))
		return ErrInvalidOTP
	}
	return nil
}

func (s *userService) CreateAccessToken(accessTokenData dao.AccessToken) (string, error) {
	// Ensure the expiration time is in the future
	if accessTokenData.ExpiresAt.Before(time.Now()) {
		s.logger.Error("Expiration time must be in the future", zap.String("accessTokenData", accessTokenData.ExpiresAt.String()))
		return "", fmt.Errorf("expiration time must be in the future")
	}

	claims := dao.Claims{
		UserID:    accessTokenData.UserID,
		Role:      accessTokenData.Role,
		SessionID: accessTokenData.SessionID,
		JTI:       accessTokenData.JTI,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(accessTokenData.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(accessTokenData.ExpiresAt),
		},
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	accessToken, err := token.SignedString([]byte(s.authConfig.JWTSecret))
	if err != nil {
		s.logger.Error("Failed to create access token", zap.String("accessTokenData", accessTokenData.ExpiresAt.String()), zap.Error(err))
		return "", fmt.Errorf("failed to sign access token: %s", err.Error())
	}

	return accessToken, nil
}
