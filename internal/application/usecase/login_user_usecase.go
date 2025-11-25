package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/igor/chronotask-api/internal/application/port"
	"github.com/igor/chronotask-api/internal/domain/repository"
	"github.com/igor/chronotask-api/internal/domain/valueobject"
)

var (
	// ErrInvalidCredentials is returned when email or password is incorrect
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// LoginUserInput represents the input for the LoginUser use case
type LoginUserInput struct {
	Email    string
	Password string
}

// LoginUserOutput represents the output of the LoginUser use case
type LoginUserOutput struct {
	UserID       string
	Email        string
	FullName     string
	AccessToken  string
	RefreshToken string
}

// LoginUserUseCase handles user authentication
type LoginUserUseCase struct {
	userRepo      repository.UserRepository
	hasherService port.HasherService
	jwtService    port.JWTService
}

// NewLoginUserUseCase creates a new LoginUserUseCase instance
func NewLoginUserUseCase(
	userRepo repository.UserRepository,
	hasherService port.HasherService,
	jwtService port.JWTService,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:      userRepo,
		hasherService: hasherService,
		jwtService:    jwtService,
	}
}

// Execute performs the user login operation
func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error) {
	// 1. Validate and create email value object
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return nil, ErrInvalidCredentials
	}

	// 3. Verify password
	if err := uc.hasherService.Compare(user.Password(), input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 4. Generate access token
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID(), user.Email().String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// 5. Generate refresh token
	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID(), user.Email().String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 6. Return login output
	return &LoginUserOutput{
		UserID:       user.ID(),
		Email:        user.Email().String(),
		FullName:     user.FullName(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
