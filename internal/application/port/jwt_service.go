package port

import "time"

// TokenClaims represents the JWT token claims
type TokenClaims struct {
	UserID   string
	Email    string
	IssuedAt time.Time
	ExpiresAt time.Time
}

// JWTService defines the interface for JWT token operations
// This is a Port in Hexagonal Architecture - the domain defines what it needs
// and the infrastructure will provide the implementation
type JWTService interface {
	// GenerateAccessToken creates a new access token for the given user
	GenerateAccessToken(userID, email string) (string, error)

	// GenerateRefreshToken creates a new refresh token for the given user
	GenerateRefreshToken(userID, email string) (string, error)

	// ValidateToken validates a token and returns the claims if valid
	ValidateToken(token string) (*TokenClaims, error)

	// RefreshAccessToken generates a new access token from a valid refresh token
	RefreshAccessToken(refreshToken string) (string, error)
}
