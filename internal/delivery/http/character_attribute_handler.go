package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/internal/application/usecase"
	"github.com/igor/chronotask-api/internal/delivery/dto"
	"github.com/igor/chronotask-api/internal/delivery/http/middleware"
)

// CharacterAttributeHandler handles character attribute-related HTTP requests
type CharacterAttributeHandler struct {
	getCharacterAttributesUseCase *usecase.GetCharacterAttributesUseCase
}

// NewCharacterAttributeHandler creates a new CharacterAttributeHandler
func NewCharacterAttributeHandler(
	getCharacterAttributesUseCase *usecase.GetCharacterAttributesUseCase,
) *CharacterAttributeHandler {
	return &CharacterAttributeHandler{
		getCharacterAttributesUseCase: getCharacterAttributesUseCase,
	}
}

// GetByCharacterID handles GET /character/:characterId/attributes - gets all attributes for a character
// This is a protected route that requires authentication
func (h *CharacterAttributeHandler) GetByCharacterID(c *gin.Context) {
	// Get character ID from URL parameter
	characterID := c.Param("characterId")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "character id is required",
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

	// Execute use case (it validates character ownership)
	output, err := h.getCharacterAttributesUseCase.Execute(c.Request.Context(), usecase.GetCharacterAttributesInput{
		CharacterID: characterID,
		UserID:      userID,
	})

	if err != nil {
		// Check if character not found or doesn't belong to user (both return same error from repo)
		if strings.Contains(err.Error(), "character not found") || strings.Contains(err.Error(), "does not belong") {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "character not found or you are not authorized to access it",
			})
			return
		}

		// Generic error
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "failed_to_fetch_attributes",
			Message: err.Error(),
		})
		return
	}

	// Convert use case output to DTOs
	attributeDTOs := make([]dto.CharacterAttributeResponse, len(output.Attributes))
	for i, attr := range output.Attributes {
		attributeDTOs[i] = dto.CharacterAttributeResponse{
			ID:            attr.ID,
			AttributeName: attr.AttributeName,
			Value:         attr.Value,
			CharacterID:   attr.CharacterID,
			CreatedAt:     attr.CreatedAt,
		}
	}

	// Return response
	c.JSON(http.StatusOK, dto.GetCharacterAttributesResponse{
		CharacterID: output.CharacterID,
		Attributes:  attributeDTOs,
	})
}
