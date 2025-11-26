package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/jackc/pgx/v5"
)

// PostgresCharacterAttributeRepository implements the CharacterAttributeRepository interface
type PostgresCharacterAttributeRepository struct {
	db *PostgresDB
}

// NewPostgresCharacterAttributeRepository creates a new PostgresCharacterAttributeRepository
func NewPostgresCharacterAttributeRepository(db *PostgresDB) *PostgresCharacterAttributeRepository {
	return &PostgresCharacterAttributeRepository{
		db: db,
	}
}

// Create persists a new character attribute
func (r *PostgresCharacterAttributeRepository) Create(ctx context.Context, attribute *entity.CharacterAttribute) error {
	query := `
		INSERT INTO character_attributes (id, attribute_name, value, character_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		attribute.ID(),
		attribute.AttributeName(),
		attribute.Value(),
		attribute.CharacterID(),
		attribute.CreatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to create character attribute: %w", err)
	}

	return nil
}

// FindByID retrieves a character attribute by its ID
func (r *PostgresCharacterAttributeRepository) FindByID(ctx context.Context, id string) (*entity.CharacterAttribute, error) {
	query := `
		SELECT id, attribute_name, value, character_id, created_at
		FROM character_attributes
		WHERE id = $1
	`

	var (
		attributeID   string
		attributeName string
		value         int
		characterID   string
		createdAt     time.Time
	)

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&attributeID,
		&attributeName,
		&value,
		&characterID,
		&createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("character attribute not found")
		}
		return nil, fmt.Errorf("failed to find character attribute: %w", err)
	}

	attribute := entity.ReconstituteCharacterAttribute(
		attributeID,
		attributeName,
		value,
		characterID,
		createdAt,
	)

	return attribute, nil
}

// FindByCharacterID retrieves all attributes for a character
func (r *PostgresCharacterAttributeRepository) FindByCharacterID(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
	query := `
		SELECT id, attribute_name, value, character_id, created_at
		FROM character_attributes
		WHERE character_id = $1
		ORDER BY attribute_name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to find character attributes: %w", err)
	}
	defer rows.Close()

	var attributes []*entity.CharacterAttribute

	for rows.Next() {
		var (
			attributeID   string
			attributeName string
			value         int
			charID        string
			createdAt     time.Time
		)

		err := rows.Scan(
			&attributeID,
			&attributeName,
			&value,
			&charID,
			&createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan character attribute: %w", err)
		}

		attribute := entity.ReconstituteCharacterAttribute(
			attributeID,
			attributeName,
			value,
			charID,
			createdAt,
		)

		attributes = append(attributes, attribute)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating character attributes: %w", err)
	}

	return attributes, nil
}

// FindByCharacterIDAndName retrieves a specific attribute by character ID and attribute name
func (r *PostgresCharacterAttributeRepository) FindByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (*entity.CharacterAttribute, error) {
	query := `
		SELECT id, attribute_name, value, character_id, created_at
		FROM character_attributes
		WHERE character_id = $1 AND attribute_name = $2
	`

	var (
		attributeID   string
		attrName      string
		value         int
		charID        string
		createdAt     time.Time
	)

	err := r.db.Pool.QueryRow(ctx, query, characterID, attributeName).Scan(
		&attributeID,
		&attrName,
		&value,
		&charID,
		&createdAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("character attribute not found")
		}
		return nil, fmt.Errorf("failed to find character attribute: %w", err)
	}

	attribute := entity.ReconstituteCharacterAttribute(
		attributeID,
		attrName,
		value,
		charID,
		createdAt,
	)

	return attribute, nil
}

// Update updates an existing character attribute
func (r *PostgresCharacterAttributeRepository) Update(ctx context.Context, attribute *entity.CharacterAttribute) error {
	query := `
		UPDATE character_attributes
		SET attribute_name = $2, value = $3
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		attribute.ID(),
		attribute.AttributeName(),
		attribute.Value(),
	)

	if err != nil {
		return fmt.Errorf("failed to update character attribute: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("character attribute not found")
	}

	return nil
}

// Delete removes a character attribute
func (r *PostgresCharacterAttributeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM character_attributes WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete character attribute: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("character attribute not found")
	}

	return nil
}

// ExistsByCharacterIDAndName checks if an attribute with the given name already exists for a character
func (r *PostgresCharacterAttributeRepository) ExistsByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM character_attributes WHERE character_id = $1 AND attribute_name = $2)`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, characterID, attributeName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if character attribute exists: %w", err)
	}

	return exists, nil
}
