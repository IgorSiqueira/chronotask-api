package usecase

import (
	"context"
	"fmt"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/repository"
)

// GetCharacterAttributesInput represents the input for getting character attributes
type GetCharacterAttributesInput struct {
	CharacterID string
	UserID      string // User ID from authentication token
}

// CharacterAttributeOutput represents a single character attribute in the output
type CharacterAttributeOutput struct {
	ID            int
	AttributeName string
	Value         int
	CharacterID   string
	CreatedAt     string
}

// GetCharacterAttributesOutput represents the output after getting character attributes
type GetCharacterAttributesOutput struct {
	CharacterID string
	Attributes  []CharacterAttributeOutput
}

// GetCharacterAttributesUseCase handles fetching all attributes for a character
type GetCharacterAttributesUseCase struct {
	characterRepo          repository.CharacterRepository
	characterAttributeRepo repository.CharacterAttributeRepository
}

// NewGetCharacterAttributesUseCase creates a new GetCharacterAttributesUseCase
func NewGetCharacterAttributesUseCase(
	characterRepo repository.CharacterRepository,
	characterAttributeRepo repository.CharacterAttributeRepository,
) *GetCharacterAttributesUseCase {
	return &GetCharacterAttributesUseCase{
		characterRepo:          characterRepo,
		characterAttributeRepo: characterAttributeRepo,
	}
}

// Execute retrieves all attributes for a character
func (uc *GetCharacterAttributesUseCase) Execute(ctx context.Context, input GetCharacterAttributesInput) (*GetCharacterAttributesOutput, error) {
	// Validate character exists AND belongs to the authenticated user (in one query)
	_, err := uc.characterRepo.FindByIDAndUserID(ctx, input.CharacterID, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("character not found or does not belong to user: %w", err)
	}

	// Fetch all attributes for the character
	attributes, err := uc.characterAttributeRepo.FindByCharacterID(ctx, input.CharacterID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch character attributes: %w", err)
	}

	// Convert entities to output
	attributeOutputs := make([]CharacterAttributeOutput, len(attributes))
	for i, attr := range attributes {
		attributeOutputs[i] = mapEntityToOutput(attr)
	}

	return &GetCharacterAttributesOutput{
		CharacterID: input.CharacterID,
		Attributes:  attributeOutputs,
	}, nil
}

// mapEntityToOutput converts a CharacterAttribute entity to output format
func mapEntityToOutput(attr *entity.CharacterAttribute) CharacterAttributeOutput {
	return CharacterAttributeOutput{
		ID:            attr.ID(),
		AttributeName: attr.AttributeName(),
		Value:         attr.Value(),
		CharacterID:   attr.CharacterID(),
		CreatedAt:     attr.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
