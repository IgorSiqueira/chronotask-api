package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/internal/application/usecase"
	"github.com/igor/chronotask-api/internal/delivery/dto"
	"github.com/igor/chronotask-api/internal/delivery/http/middleware"
)

// CharacterHandler handles character-related HTTP requests
type CharacterHandler struct {
	createCharacterUseCase    *usecase.CreateCharacterUseCase
	getUserCharactersUseCase  *usecase.GetUserCharactersUseCase
}

// NewCharacterHandler creates a new CharacterHandler
func NewCharacterHandler(
	createCharacterUseCase *usecase.CreateCharacterUseCase,
	getUserCharactersUseCase *usecase.GetUserCharactersUseCase,
) *CharacterHandler {
	return &CharacterHandler{
		createCharacterUseCase:   createCharacterUseCase,
		getUserCharactersUseCase: getUserCharactersUseCase,
	}
}

// Create handles POST /character - creates a new character for the authenticated user
// This is a protected route that requires authentication
func (h *CharacterHandler) Create(c *gin.Context) {
	var req dto.CreateCharacterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Get authenticated user ID from middleware
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
		return
	}

	// Execute use case
	output, err := h.createCharacterUseCase.Execute(c.Request.Context(), usecase.CreateCharacterInput{
		Name:   req.Name,
		UserID: userID,
	})

	if err != nil {
		// Check specific error messages
		errMsg := err.Error()
		if errMsg == "user already has a character" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "character_already_exists",
				Message: "you already have a character",
			})
			return
		}

		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error:   "character_creation_failed",
			Message: err.Error(),
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, dto.CreateCharacterResponse{
		ID:        output.ID,
		Name:      output.Name,
		Level:     output.Level,
		CurrentXp: output.CurrentXp,
		TotalXp:   output.TotalXp,
		UserID:    output.UserID,
		CreatedAt: output.CreatedAt,
	})
}

// GetList handles GET /user/characters - gets all characters for the authenticated user
// This is a protected route that requires authentication
func (h *CharacterHandler) GetList(c *gin.Context) {
	// Get authenticated user ID from middleware
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "user not authenticated",
		})
		return
	}

	// Execute use case
	output, err := h.getUserCharactersUseCase.Execute(c.Request.Context(), usecase.GetUserCharactersInput{
		UserID: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "failed_to_fetch_characters",
			Message: err.Error(),
		})
		return
	}

	// Convert use case output to DTOs
	characterDTOs := make([]dto.CharacterItemResponse, len(output.Characters))
	for i, char := range output.Characters {
		characterDTOs[i] = dto.CharacterItemResponse{
			ID:        char.ID,
			Name:      char.Name,
			Level:     char.Level,
			CurrentXp: char.CurrentXp,
			TotalXp:   char.TotalXp,
			CreatedAt: char.CreatedAt,
		}
	}

	// Return response
	c.JSON(http.StatusOK, dto.GetUserCharactersResponse{
		Characters: characterDTOs,
	})
}
