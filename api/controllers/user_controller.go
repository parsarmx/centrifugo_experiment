package controllers

import (
	"errors"
	"fmt"
	dao "golang_template/api/dto"
	"golang_template/api/response"
	"golang_template/internal/config"
	"golang_template/internal/pkg"
	"golang_template/internal/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserController interface {
	SendOTP(ctx echo.Context) error
	OTPLogin(ctx echo.Context) error
}

type userController struct {
	service    service.UserService
	logger     *zap.Logger
	authConfig config.AuthConfig
	// fallbackController fiber.Handler
	validator     *pkg.Validator
	bodyValidator *validator.Validate
}

func NewUserController(service service.UserService, logger *zap.Logger, authConfig config.AuthConfig) UserController {
	// fallbackController := NewFallbackController(logger)
	return &userController{service: service,
		logger:        logger,
		authConfig:    authConfig,
		validator:     pkg.NewValidator(),
		bodyValidator: validator.New(validator.WithRequiredStructEnabled()),
		// fallbackController: fallbackController,
	}
}

func (c *userController) SendOTP(ctx echo.Context) error {
	var requestData dao.SendOTPRequestData
	if err := ctx.Bind(&requestData); err != nil {
		return err
	}

	if err := c.validator.ValidateStruct(requestData); err != nil {
		return err
	}

	userData, err := c.service.SendOtp(ctx.Request().Context(), requestData.PhoneNumber)
	if err == nil {
		if c.authConfig.UnderDevelopment {
			return ctx.JSON(http.StatusOK, map[string]interface{}{
				"ok":           true,
				"phone_number": requestData.PhoneNumber,
				"otp_code":     userData.OtpCode,
			})
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"ok":           true,
			"phone_number": requestData.PhoneNumber,
		})
	}

	if errors.Is(err, service.ErrOTPSent) {
		c.logger.Warn("Too many requests", zap.Error(err))
		return ctx.JSON(http.StatusTooManyRequests, map[string]interface{}{
			"error":                   "otp has already been sent",
			"time_to_wait_in_seconds": fmt.Sprintf("%v", userData.OtpTtl),
			"message":                 "otp has already been sent",
			"status":                  http.StatusTooManyRequests,
		})
	}

	return err
}

func (c *userController) OTPLogin(ctx echo.Context) error {

	var requestData dao.OTPLoginRequestData

	if err := ctx.Bind(&requestData); err != nil {
		return err
	}

	if err := c.validator.ValidateStruct(requestData); err != nil {
		return err
	}

	responseData, err := c.service.OTPLogin(ctx.Request().Context(), requestData)
	if err == nil {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"ok":                    true,
			"user_id":               responseData.UserID,
			"session_id":            responseData.SessionID,
			"access_token":          responseData.AccessToken,
			"access_token_exp_time": responseData.AccExpiresAt,
		})
	}

	if errors.Is(err, service.ErrInvalidOTP) {
		return response.SendError(ctx, "OTP is not valid", http.StatusUnauthorized, "OTP is not valid")
	}

	return err
}
