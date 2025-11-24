package repository

import (
	"context"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/valueobject"
)

// UserRepository defines the interface for user persistence (Port)
// This is defined in the domain layer, but implemented in the infrastructure layer
type UserRepository interface {
	// Create persists a new user
	Create(ctx context.Context, user *entity.User) error

	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id string) (*entity.User, error)

	// FindByEmail retrieves a user by their email
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user (soft delete recommended)
	Delete(ctx context.Context, id string) error

	// ExistsByEmail checks if a user with the given email already exists
	ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)
}
