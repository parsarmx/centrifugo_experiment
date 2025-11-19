package echojwt

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// sendError writes the error response in JSON
func sendError(c echo.Context, body []byte) error {
	return c.JSONBlob(http.StatusUnauthorized, body)
}

// New creates a new JWT middleware handler for Echo.
func New(config MiddlewareConfig) echo.MiddlewareFunc {
	logger := config.Logger
	secret := config.Secret

	bearerPrefix := []byte("Bearer ")
	if len(config.BearerPrefix) > 0 {
		bearerPrefix = config.BearerPrefix
	}
	bearerPrefixLength := len(bearerPrefix)

	if logger == nil {
		panic("JWT Middleware: Logger cannot be nil in MiddlewareConfig")
	}
	if secret == "" {
		logger.Error("JWT Middleware: Secret cannot be empty in MiddlewareConfig")
		panic("JWT Middleware: Secret cannot be empty in MiddlewareConfig")
	}

	if config.AuthHeaderKey == "" {
		config.AuthHeaderKey = echo.HeaderAuthorization
	}
	if config.CookieTokenKey == "" {
		config.CookieTokenKey = "access"
	}

	var (
		authHeaderKey  = config.AuthHeaderKey
		cookieTokenKey = config.CookieTokenKey
	)

	secretBytes := []byte(secret)

	JWTKeyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secretBytes, nil
	}

	// Object pools
	var claimsPool = sync.Pool{
		New: func() interface{} { return new(CustomClaims) },
	}

	tokenBufLength := config.DefaultTokenLength
	if tokenBufLength == 0 {
		tokenBufLength = 256
	}

	var tokenBufPool = sync.Pool{
		New: func() interface{} { return make([]byte, 0, tokenBufLength) },
	}

	// Optimized parser
	var JWTParser = jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithLeeway(5),
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get buffer from pool
			tokenBuf := tokenBufPool.Get().([]byte)
			tokenBytes := tokenBuf[:0]
			defer tokenBufPool.Put(tokenBuf)

			// Check Authorization header
			authHeader := []byte(c.Request().Header.Get(authHeaderKey))
			if len(authHeader) > bearerPrefixLength && bytes.HasPrefix(authHeader, bearerPrefix) {
				tokenBytes = append(tokenBytes, authHeader[bearerPrefixLength:]...)
			} else {
				// Fallback to cookie
				cookie, err := c.Cookie(cookieTokenKey)
				if err == nil && cookie != nil && len(cookie.Value) > 0 {
					tokenBytes = append(tokenBytes, cookie.Value...)
				}
			}

			if len(tokenBytes) == 0 {
				return sendError(c, unauthorizedMissingToken)
			}

			// Claims
			claims := claimsPool.Get().(*CustomClaims)
			*claims = CustomClaims{}
			defer claimsPool.Put(claims)

			// Parse token
			token, err := JWTParser.ParseWithClaims(string(tokenBytes), claims, JWTKeyFunc)
			if err != nil || !token.Valid {
				return sendError(c, unauthorizedInvalidToken)
			}

			if claims.UserID == "" || claims.SessionID == "" {
				return sendError(c, unauthorizedInvalidClaims)
			}

			// Store claims in context
			c.Set(HeaderUserID, claims.UserID)
			c.Set(HeaderSessionID, claims.SessionID)
			c.Set(HeaderRole, claims.Roles)

			return next(c)
		}
	}
}
