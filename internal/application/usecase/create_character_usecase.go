package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/repository"
)

// Base attributes that every new character starts with
var baseAttributes = []struct {
	name  string
	value int
}{
	{"Força", 5},
	{"Constituição", 5},
	{"Vontade", 5},
	{"Sabedoria", 5},
	{"Inteligência", 5},
	{"Carisma", 5},
	{"Destreza", 5},
}

// CreateCharacterInput represents the input for creating a character
type CreateCharacterInput struct {
	Name   string
	UserID string
}

// CreateCharacterOutput represents the output after creating a character
type CreateCharacterOutput struct {
	ID        string
	Name      string
	Level     int
	CurrentXp int
	TotalXp   int
	UserID    string
	CreatedAt string
}

// CreateCharacterUseCase handles the creation of new characters
type CreateCharacterUseCase struct {
	characterRepo          repository.CharacterRepository
	characterAttributeRepo repository.CharacterAttributeRepository
}

// NewCreateCharacterUseCase creates a new CreateCharacterUseCase
func NewCreateCharacterUseCase(
	characterRepo repository.CharacterRepository,
	characterAttributeRepo repository.CharacterAttributeRepository,
) *CreateCharacterUseCase {
	return &CreateCharacterUseCase{
		characterRepo:          characterRepo,
		characterAttributeRepo: characterAttributeRepo,
	}
}

// Execute creates a new character for a user
func (uc *CreateCharacterUseCase) Execute(ctx context.Context, input CreateCharacterInput) (*CreateCharacterOutput, error) {
	// Check if user already has a character (business rule: 1 character per user)
	exists, err := uc.characterRepo.ExistsByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if character exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user already has a character")
	}

	// Generate unique ID
	characterID := uuid.New().String()

	// Create character entity (with domain validation)
	character, err := entity.NewCharacter(
		characterID,
		input.Name,
		input.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create character: %w", err)
	}

	// Persist character
	if err := uc.characterRepo.Create(ctx, character); err != nil {
		return nil, fmt.Errorf("failed to save character: %w", err)
	}

	// Create base attributes for the new character
	for _, baseAttr := range baseAttributes {
		attribute, err := entity.NewCharacterAttribute(
			baseAttr.name,
			baseAttr.value,
			character.ID(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create base attribute '%s': %w", baseAttr.name, err)
		}

		if err := uc.characterAttributeRepo.Create(ctx, attribute); err != nil {
			return nil, fmt.Errorf("failed to save base attribute '%s': %w", baseAttr.name, err)
		}
	}

	// Return output
	return &CreateCharacterOutput{
		ID:        character.ID(),
		Name:      character.Name(),
		Level:     character.Level(),
		CurrentXp: character.CurrentXp(),
		TotalXp:   character.TotalXp(),
		UserID:    character.UserID(),
		CreatedAt: character.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
