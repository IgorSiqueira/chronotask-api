package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/valueobject"
	"github.com/jackc/pgx/v5"
)

// PostgresUserRepository implements the UserRepository interface
type PostgresUserRepository struct {
	db *PostgresDB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *PostgresDB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Create persists a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, full_name, email, password, birth_date, accept_terms, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID(),
		user.FullName(),
		user.Email().Value(),
		user.Password(),
		user.BirthDate(),
		user.AcceptTerms(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID retrieves a user by their ID
func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, full_name, email, password, birth_date, accept_terms, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var (
		userID      string
		fullName    string
		emailStr    string
		password    string
		birthDate   time.Time
		acceptTerms bool
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&userID,
		&fullName,
		&emailStr,
		&password,
		&birthDate,
		&acceptTerms,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	user := entity.ReconstituteUser(
		userID,
		fullName,
		email,
		password,
		birthDate,
		acceptTerms,
		createdAt,
		updatedAt,
	)

	return user, nil
}

// FindByEmail retrieves a user by their email
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	query := `
		SELECT id, full_name, email, password, birth_date, accept_terms, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var (
		userID      string
		fullName    string
		emailStr    string
		password    string
		birthDate   time.Time
		acceptTerms bool
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := r.db.Pool.QueryRow(ctx, query, email.Value()).Scan(
		&userID,
		&fullName,
		&emailStr,
		&password,
		&birthDate,
		&acceptTerms,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	emailVO, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	user := entity.ReconstituteUser(
		userID,
		fullName,
		emailVO,
		password,
		birthDate,
		acceptTerms,
		createdAt,
		updatedAt,
	)

	return user, nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET full_name = $2, email = $3, password = $4, birth_date = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		user.ID(),
		user.FullName(),
		user.Email().Value(),
		user.Password(),
		user.BirthDate(),
		user.UpdatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete removes a user (hard delete)
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ExistsByEmail checks if a user with the given email already exists
func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, email.Value()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists, nil
}
