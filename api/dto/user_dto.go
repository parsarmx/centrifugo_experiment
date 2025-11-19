package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SendOTPRequestData struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
}

type SendOTPResponseData struct {
	PhoneNumber string  `json:"phone_number,omitempty"`
	OtpTtl      float64 `json:"otp_ttl,omitempty"`
	OtpCode     string  `json:"otp_code,omitempty"`
}

type OTPLoginRequestData struct {
	OTPCode     string `json:"otp_code" validate:"required,numeric"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
}

type OTPLoginResponseData struct {
	RefreshTokenID uuid.UUID `json:"refresh_token_id"`
	UserID         uuid.UUID `json:"user_id"`
	SessionID      uuid.UUID `json:"session_id"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	AccExpiresAt   time.Time `json:"access_token_expires_at"`
	RefExpiresAt   time.Time `json:"refresh_token_expires_at"`
}

type RefreshToken struct {
	ID                       uuid.UUID
	UserID                   uuid.UUID
	HashedRefreshToken       string
	SessionID                uuid.UUID
	IssuedAt                 time.Time
	ExpiresAt                time.Time
	LastAccessExpirationTime time.Time
	IsRevoked                bool
}

type AccessToken struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Role      string
	JTI       uuid.UUID
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"roles"`
	SessionID uuid.UUID `json:"session_id"`
	JTI       uuid.UUID `json:"jti"`
	jwt.RegisteredClaims
}

type UpdateScoreRequest struct {
	Username string  `json:"username" validate:"required"`
	Score    float64 `json:"score" validate:"required"`
}
