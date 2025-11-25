package dto

import "time"

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	FullName    string `json:"fullName" binding:"required,min=2,max=255"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	BirthDate   string `json:"birthDate" binding:"required"` // Format: YYYY-MM-DD
	AcceptTerms bool   `json:"acceptTerms" binding:"required"`
}

// CreateUserResponse represents the response after creating a user
type CreateUserResponse struct {
	ID        string    `json:"id"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birthDate"`
	CreatedAt time.Time `json:"createdAt"`
}

// LoginRequest represents the request to authenticate a user
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful authentication
type LoginResponse struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	FullName     string `json:"fullName"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
