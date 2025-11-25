package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/igor/chronotask-api/internal/application/port"
)

// JWTServiceImpl implements the JWTService interface using golang-jwt
type JWTServiceImpl struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// customClaims defines the JWT claims structure
type customClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secret string, accessDuration, refreshDuration string) (*JWTServiceImpl, error) {
	accessDur, err := time.ParseDuration(accessDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid access token duration: %w", err)
	}

	refreshDur, err := time.ParseDuration(refreshDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token duration: %w", err)
	}

	return &JWTServiceImpl{
		secretKey:            []byte(secret),
		accessTokenDuration:  accessDur,
		refreshTokenDuration: refreshDur,
	}, nil
}

// GenerateAccessToken creates a new access token for the given user
func (s *JWTServiceImpl) GenerateAccessToken(userID, email string) (string, error) {
	return s.generateToken(userID, email, s.accessTokenDuration)
}

// GenerateRefreshToken creates a new refresh token for the given user
func (s *JWTServiceImpl) GenerateRefreshToken(userID, email string) (string, error) {
	return s.generateToken(userID, email, s.refreshTokenDuration)
}

// generateToken creates a token with the specified duration
func (s *JWTServiceImpl) generateToken(userID, email string, duration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	claims := customClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a token and returns the claims if valid
func (s *JWTServiceImpl) ValidateToken(tokenString string) (*port.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &port.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (s *JWTServiceImpl) RefreshAccessToken(refreshToken string) (string, error) {
	// Validate the refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new access token with the same user info
	return s.GenerateAccessToken(claims.UserID, claims.Email)
}
