package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/igor/chronotask-api/internal/application/port"
	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/repository"
	"github.com/igor/chronotask-api/internal/domain/valueobject"
)

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	FullName    string
	Email       string
	Password    string
	BirthDate   time.Time
	AcceptTerms bool
}

// CreateUserOutput represents the output after creating a user
type CreateUserOutput struct {
	ID        string
	FullName  string
	Email     string
	BirthDate time.Time
	CreatedAt time.Time
}

// CreateUserUseCase handles the creation of new users
type CreateUserUseCase struct {
	userRepo      repository.UserRepository
	hasherService port.HasherService
}

// NewCreateUserUseCase creates a new CreateUserUseCase
func NewCreateUserUseCase(
	userRepo repository.UserRepository,
	hasherService port.HasherService,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:      userRepo,
		hasherService: hasherService,
	}
}

// Execute creates a new user
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// Validate email format
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Check if user already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email.Value())
	}

	// Hash password
	hashedPassword, err := uc.hasherService.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate unique ID
	userID := uuid.New().String()

	// Create user entity (with domain validation)
	user, err := entity.NewUser(
		userID,
		input.FullName,
		email,
		hashedPassword,
		input.BirthDate,
		input.AcceptTerms,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Persist user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Return output
	return &CreateUserOutput{
		ID:        user.ID(),
		FullName:  user.FullName(),
		Email:     user.Email().Value(),
		BirthDate: user.BirthDate(),
		CreatedAt: user.CreatedAt(),
	}, nil
}
