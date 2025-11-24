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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
