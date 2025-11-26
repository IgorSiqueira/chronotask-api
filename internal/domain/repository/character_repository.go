package repository

import (
	"context"

	"github.com/igor/chronotask-api/internal/domain/entity"
)

// CharacterRepository defines the interface for character persistence (Port)
// This is defined in the domain layer, but implemented in the infrastructure layer
type CharacterRepository interface {
	// Create persists a new character
	Create(ctx context.Context, character *entity.Character) error

	// FindByID retrieves a character by their ID
	FindByID(ctx context.Context, id string) (*entity.Character, error)

	// FindByIDAndUserID retrieves a character by ID and validates ownership
	// Returns error if character doesn't exist OR doesn't belong to the user
	FindByIDAndUserID(ctx context.Context, id string, userID string) (*entity.Character, error)

	// FindByUserID retrieves a character by their user ID
	FindByUserID(ctx context.Context, userID string) (*entity.Character, error)

	// FindAllByUserID retrieves all characters for a user (returns list for future support of multiple characters)
	FindAllByUserID(ctx context.Context, userID string) ([]*entity.Character, error)

	// Update updates an existing character
	Update(ctx context.Context, character *entity.Character) error

	// Delete removes a character
	Delete(ctx context.Context, id string) error

	// ExistsByUserID checks if a user already has a character
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
}
