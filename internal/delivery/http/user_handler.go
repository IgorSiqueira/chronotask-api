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
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(createUserUseCase *usecase.CreateUserUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
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
