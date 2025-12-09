package repository

import (
	"context"

	"github.com/igor/chronotask-api/internal/domain/entity"
)

// CharacterAttributeRepository defines the interface for character attribute persistence (Port)
// This is defined in the domain layer, but implemented in the infrastructure layer
type CharacterAttributeRepository interface {
	// Create persists a new character attribute
	Create(ctx context.Context, attribute *entity.CharacterAttribute) error

	// FindByID retrieves a character attribute by its ID
	FindByID(ctx context.Context, id int) (*entity.CharacterAttribute, error)

	// FindByCharacterID retrieves all attributes for a character
	FindByCharacterID(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error)

	// FindByCharacterIDAndName retrieves a specific attribute by character ID and attribute name
	FindByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (*entity.CharacterAttribute, error)

	// Update updates an existing character attribute
	Update(ctx context.Context, attribute *entity.CharacterAttribute) error

	// Delete removes a character attribute
	Delete(ctx context.Context, id int) error

	// ExistsByCharacterIDAndName checks if an attribute with the given name already exists for a character
	ExistsByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (bool, error)
}
