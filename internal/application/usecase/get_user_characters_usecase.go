package usecase

import (
	"context"
	"fmt"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/repository"
)

// GetUserCharactersInput represents the input for getting user's characters
type GetUserCharactersInput struct {
	UserID string // User ID from authentication token
}

// CharacterOutput represents a single character in the output
type CharacterOutput struct {
	ID        string
	Name      string
	Level     int
	CurrentXp int
	TotalXp   int
	UserID    string
	CreatedAt string
}

// GetUserCharactersOutput represents the output after getting user's characters
type GetUserCharactersOutput struct {
	Characters []CharacterOutput
}

// GetUserCharactersUseCase handles fetching all characters for a user
type GetUserCharactersUseCase struct {
	characterRepo repository.CharacterRepository
}

// NewGetUserCharactersUseCase creates a new GetUserCharactersUseCase
func NewGetUserCharactersUseCase(
	characterRepo repository.CharacterRepository,
) *GetUserCharactersUseCase {
	return &GetUserCharactersUseCase{
		characterRepo: characterRepo,
	}
}

// Execute retrieves all characters for a user
func (uc *GetUserCharactersUseCase) Execute(ctx context.Context, input GetUserCharactersInput) (*GetUserCharactersOutput, error) {
	// Fetch all characters for the user
	characters, err := uc.characterRepo.FindAllByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user characters: %w", err)
	}

	// Convert entities to output
	characterOutputs := make([]CharacterOutput, len(characters))
	for i, char := range characters {
		characterOutputs[i] = mapCharacterEntityToOutput(char)
	}

	return &GetUserCharactersOutput{
		Characters: characterOutputs,
	}, nil
}

// mapCharacterEntityToOutput converts a Character entity to output format
func mapCharacterEntityToOutput(char *entity.Character) CharacterOutput {
	return CharacterOutput{
		ID:        char.ID(),
		Name:      char.Name(),
		Level:     char.Level(),
		CurrentXp: char.CurrentXp(),
		TotalXp:   char.TotalXp(),
		UserID:    char.UserID(),
		CreatedAt: char.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
