package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/internal/application/port"
	"github.com/igor/chronotask-api/internal/delivery/dto"
)

const (
	// AuthorizationHeader is the name of the authorization header
	AuthorizationHeader = "Authorization"
	// BearerPrefix is the prefix for bearer tokens
	BearerPrefix = "Bearer "
	// UserIDKey is the context key for the user ID
	UserIDKey = "user_id"
	// UserEmailKey is the context key for the user email
	UserEmailKey = "user_email"
)

// AuthMiddleware is a middleware that validates JWT tokens
type AuthMiddleware struct {
	jwtService port.JWTService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService port.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// RequireAuth validates the JWT token from the Authorization header
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract token from header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "missing_authorization",
				Message: "authorization header is required",
			})
			c.Abort()
			return
		}

		// 2. Check if token has Bearer prefix
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "invalid_authorization",
				Message: "authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// 3. Extract token (remove "Bearer " prefix)
		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "missing_token",
				Message: "bearer token is missing",
			})
			c.Abort()
			return
		}

		// 4. Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "invalid_token",
				Message: "token is invalid or expired",
			})
			c.Abort()
			return
		}

		// 5. Store user info in context for use in handlers
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		// 6. Continue to next handler
		c.Next()
	}
}

// GetUserID extracts the user ID from the Gin context
// This is a helper function for use in handlers that require authentication
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserEmail extracts the user email from the Gin context
// This is a helper function for use in handlers that require authentication
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	emailStr, ok := email.(string)
	return emailStr, ok
}
