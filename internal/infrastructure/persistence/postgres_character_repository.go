package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/jackc/pgx/v5"
)

// PostgresCharacterRepository implements the CharacterRepository interface
type PostgresCharacterRepository struct {
	db *PostgresDB
}

// NewPostgresCharacterRepository creates a new PostgresCharacterRepository
func NewPostgresCharacterRepository(db *PostgresDB) *PostgresCharacterRepository {
	return &PostgresCharacterRepository{
		db: db,
	}
}

// Create persists a new character
func (r *PostgresCharacterRepository) Create(ctx context.Context, character *entity.Character) error {
	query := `
		INSERT INTO characters (id, name, level, current_xp, total_xp, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		character.ID(),
		character.Name(),
		character.Level(),
		character.CurrentXp(),
		character.TotalXp(),
		character.UserID(),
		character.CreatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to create character: %w", err)
	}

	return nil
}

// FindByID retrieves a character by their ID
func (r *PostgresCharacterRepository) FindByID(ctx context.Context, id string) (*entity.Character, error) {
	query := `
		SELECT id, name, level, current_xp, total_xp, user_id, created_at
		FROM characters
		WHERE id = $1
	`

	var (
		characterID string
		name        string
		level       int
		currentXp   int
		totalXp     int
		userID      string
		createdAt   time.Time
	)

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&characterID,
		&name,
		&level,
		&currentXp,
		&totalXp,
		&userID,
		&createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	character := entity.ReconstituteCharacter(
		characterID,
		name,
		level,
		currentXp,
		totalXp,
		userID,
		createdAt,
	)

	return character, nil
}

// FindByIDAndUserID retrieves a character by ID and validates ownership
// Returns error if character doesn't exist OR doesn't belong to the user
func (r *PostgresCharacterRepository) FindByIDAndUserID(ctx context.Context, id string, userID string) (*entity.Character, error) {
	query := `
		SELECT id, name, level, current_xp, total_xp, user_id, created_at
		FROM characters
		WHERE id = $1 AND user_id = $2
	`

	var (
		characterID string
		name        string
		level       int
		currentXp   int
		totalXp     int
		userIDVal   string
		createdAt   time.Time
	)

	err := r.db.Pool.QueryRow(ctx, query, id, userID).Scan(
		&characterID,
		&name,
		&level,
		&currentXp,
		&totalXp,
		&userIDVal,
		&createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("character not found or does not belong to user")
		}
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	character := entity.ReconstituteCharacter(
		characterID,
		name,
		level,
		currentXp,
		totalXp,
		userIDVal,
		createdAt,
	)

	return character, nil
}

// FindByUserID retrieves a character by their user ID
func (r *PostgresCharacterRepository) FindByUserID(ctx context.Context, userID string) (*entity.Character, error) {
	query := `
		SELECT id, name, level, current_xp, total_xp, user_id, created_at
		FROM characters
		WHERE user_id = $1
	`

	var (
		characterID string
		name        string
		level       int
		currentXp   int
		totalXp     int
		userIDVal   string
		createdAt   time.Time
	)

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&characterID,
		&name,
		&level,
		&currentXp,
		&totalXp,
		&userIDVal,
		&createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("character not found for user")
		}
		return nil, fmt.Errorf("failed to find character by user: %w", err)
	}

	character := entity.ReconstituteCharacter(
		characterID,
		name,
		level,
		currentXp,
		totalXp,
		userIDVal,
		createdAt,
	)

	return character, nil
}

// FindAllByUserID retrieves all characters for a user
func (r *PostgresCharacterRepository) FindAllByUserID(ctx context.Context, userID string) ([]*entity.Character, error) {
	query := `
		SELECT id, name, level, current_xp, total_xp, user_id, created_at
		FROM characters
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %w", err)
	}
	defer rows.Close()

	var characters []*entity.Character

	for rows.Next() {
		var (
			characterID string
			name        string
			level       int
			currentXp   int
			totalXp     int
			userIDVal   string
			createdAt   time.Time
		)

		err := rows.Scan(
			&characterID,
			&name,
			&level,
			&currentXp,
			&totalXp,
			&userIDVal,
			&createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}

		character := entity.ReconstituteCharacter(
			characterID,
			name,
			level,
			currentXp,
			totalXp,
			userIDVal,
			createdAt,
		)

		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating characters: %w", err)
	}

	return characters, nil
}

// Update updates an existing character
func (r *PostgresCharacterRepository) Update(ctx context.Context, character *entity.Character) error {
	query := `
		UPDATE characters
		SET name = $2, level = $3, current_xp = $4, total_xp = $5
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		character.ID(),
		character.Name(),
		character.Level(),
		character.CurrentXp(),
		character.TotalXp(),
	)

	if err != nil {
		return fmt.Errorf("failed to update character: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("character not found")
	}

	return nil
}

// Delete removes a character
func (r *PostgresCharacterRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM characters WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("character not found")
	}

	return nil
}

// ExistsByUserID checks if a user already has a character
func (r *PostgresCharacterRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM characters WHERE user_id = $1)`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if character exists for user: %w", err)
	}

	return exists, nil
}
