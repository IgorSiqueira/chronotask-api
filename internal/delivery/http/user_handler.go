package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/internal/application/usecase"
	"github.com/igor/chronotask-api/internal/delivery/dto"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	createUserUseCase *usecase.CreateUserUseCase
	loginUserUseCase  *usecase.LoginUserUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	createUserUseCase *usecase.CreateUserUseCase,
	loginUserUseCase *usecase.LoginUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
		loginUserUseCase:  loginUserUseCase,
	}
}

// Create handles POST /user - creates a new user
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_birth_date",
			Message: "birth date must be in format YYYY-MM-DD",
		})
		return
	}

	// Execute use case
	output, err := h.createUserUseCase.Execute(c.Request.Context(), usecase.CreateUserInput{
		FullName:    req.FullName,
		Email:       req.Email,
		Password:    req.Password,
		BirthDate:   birthDate,
		AcceptTerms: req.AcceptTerms,
	})

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error:   "user_creation_failed",
			Message: err.Error(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, dto.CreateUserResponse{
		ID:        output.ID,
		FullName:  output.FullName,
		Email:     output.Email,
		BirthDate: output.BirthDate.Format("2006-01-02"),
		CreatedAt: output.CreatedAt,
	})
}

// Login handles POST /login - authenticates a user
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Execute use case
	output, err := h.loginUserUseCase.Execute(c.Request.Context(), usecase.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// Check if it's an invalid credentials error
		if err == usecase.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "authentication_failed",
				Message: "invalid email or password",
			})
			return
		}

		// Internal server error for other errors
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "login_failed",
			Message: "failed to authenticate user",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.LoginResponse{
		UserID:       output.UserID,
		Email:        output.Email,
		FullName:     output.FullName,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}

// GetProfile handles GET /user/profile - gets the authenticated user's profile
// This is a protected route that requires authentication
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user info from middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
		return
	}

	email, _ := c.Get("user_email")

	// Return user profile (simplified for testing)
	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"email":  email,
		"message": "This is a protected route - you are authenticated!",
	})
}
