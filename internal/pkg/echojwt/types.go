package echojwt

import (
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	// ErrMissingOrMalformed is returned when the JWT is missing from the Authorization header or is malformed.
	rrrMissingOrMalformed = errors.New("missing or malformed JWT token")
	// ErrInvalidToken is returned when the JWT is invalid (e.g., signature mismatch, invalid claims).
	ErrInvalidToken = errors.New("invalid token")
)

var (
	// Renamed for clarity as it covers header and cookie
	MISSING_TOKEN_ERROR = echo.Map{
		"error":   "Unauthorized",
		"message": "Missing or malformed JWT token",
	}
	INVALID_TOKEN_ERROR = echo.Map{
		"error":   "Unauthorized",
		"message": "Invalid token",
	}
	INVALID_CLAIMS_ERROR = echo.Map{
		"error":   "Unauthorized",
		"message": "Invalid claims",
	}
)

var (
	unauthorizedMissingToken, _  = json.Marshal(MISSING_TOKEN_ERROR)
	unauthorizedInvalidToken, _  = json.Marshal(INVALID_TOKEN_ERROR)
	unauthorizedInvalidClaims, _ = json.Marshal(INVALID_CLAIMS_ERROR)
)

// CustomClaims defines the structure of the JWT claims we expect.
type CustomClaims struct {
	UserID    string `json:"user_id"`
	Roles     string `json:"roles"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

// MiddlewareConfig defines the configuration for the JWT middleware.
type MiddlewareConfig struct {
	// Secret is the secret key to validate the JWT signature. This is required.
	Secret             string
	Logger             *zap.Logger
	BearerPrefix       []byte
	AuthHeaderKey      string
	CookieTokenKey     string
	DefaultTokenLength int
}

const (
	// HeaderUserID is the header name for User ID.
	HeaderUserID = "X-User-ID"
	// HeaderSessionID is the header name for Session ID.
	HeaderSessionID = "X-Session-ID"
	// HeaderRole is the header name for Role.
	HeaderRole = "X-Role"
)
