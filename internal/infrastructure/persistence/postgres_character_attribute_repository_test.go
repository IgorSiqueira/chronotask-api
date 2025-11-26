package persistence_test

import (
	"context"
	"testing"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/repository"
	"github.com/igor/chronotask-api/internal/infrastructure/persistence"
)

// createTestCharacter creates a character for testing attribute operations
func createTestCharacter(t *testing.T, userRepo repository.UserRepository, charRepo repository.CharacterRepository) *entity.Character {
	t.Helper()

	// Create user first
	user := createTestUser(t, userRepo)

	// Create character
	character, err := entity.NewCharacter(
		"test-char-id",
		"Test Warrior",
		user.ID(),
	)
	if err != nil {
		t.Fatalf("Failed to create test character entity: %v", err)
	}

	err = charRepo.Create(context.Background(), character)
	if err != nil {
		t.Fatalf("Failed to save test character: %v", err)
	}

	return character
}

func TestPostgresCharacterAttributeRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	// Create test character
	character := createTestCharacter(t, userRepo, charRepo)

	// Create attribute
	attribute, err := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		character.ID(),
	)
	if err != nil {
		t.Fatalf("Failed to create attribute entity: %v", err)
	}

	err = attrRepo.Create(context.Background(), attribute)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	// Verify creation
	found, err := attrRepo.FindByID(context.Background(), attribute.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v, want nil", err)
	}

	if found.ID() != attribute.ID() {
		t.Errorf("found.ID() = %v, want %v", found.ID(), attribute.ID())
	}

	if found.AttributeName() != attribute.AttributeName() {
		t.Errorf("found.AttributeName() = %v, want %v", found.AttributeName(), attribute.AttributeName())
	}

	if found.Value() != attribute.Value() {
		t.Errorf("found.Value() = %v, want %v", found.Value(), attribute.Value())
	}
}

func TestPostgresCharacterAttributeRepository_FindByID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	_, err := attrRepo.FindByID(context.Background(), "non-existent-id")
	if err == nil {
		t.Error("FindByID() error = nil, want error for non-existent attribute")
	}
}

func TestPostgresCharacterAttributeRepository_FindByCharacterID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	// Create multiple attributes
	attr1, _ := entity.NewCharacterAttribute("attr-1", "Strength", 10, character.ID())
	attr2, _ := entity.NewCharacterAttribute("attr-2", "Agility", 15, character.ID())
	attr3, _ := entity.NewCharacterAttribute("attr-3", "Intelligence", 20, character.ID())

	attrRepo.Create(context.Background(), attr1)
	attrRepo.Create(context.Background(), attr2)
	attrRepo.Create(context.Background(), attr3)

	// Find all attributes for character
	attributes, err := attrRepo.FindByCharacterID(context.Background(), character.ID())
	if err != nil {
		t.Fatalf("FindByCharacterID() error = %v, want nil", err)
	}

	if len(attributes) != 3 {
		t.Errorf("len(attributes) = %v, want %v", len(attributes), 3)
	}

	// Verify attributes are sorted by name (Agility, Intelligence, Strength)
	expectedNames := []string{"Agility", "Intelligence", "Strength"}
	for i, attr := range attributes {
		if attr.AttributeName() != expectedNames[i] {
			t.Errorf("attributes[%d].AttributeName() = %v, want %v", i, attr.AttributeName(), expectedNames[i])
		}
	}
}

func TestPostgresCharacterAttributeRepository_FindByCharacterID_Empty(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	// Find attributes for character with no attributes
	attributes, err := attrRepo.FindByCharacterID(context.Background(), character.ID())
	if err != nil {
		t.Fatalf("FindByCharacterID() error = %v, want nil", err)
	}

	if len(attributes) != 0 {
		t.Errorf("len(attributes) = %v, want 0", len(attributes))
	}
}

func TestPostgresCharacterAttributeRepository_FindByCharacterIDAndName(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	attribute, _ := entity.NewCharacterAttribute("attr-123", "Strength", 10, character.ID())
	attrRepo.Create(context.Background(), attribute)

	// Find specific attribute
	found, err := attrRepo.FindByCharacterIDAndName(context.Background(), character.ID(), "Strength")
	if err != nil {
		t.Fatalf("FindByCharacterIDAndName() error = %v, want nil", err)
	}

	if found.AttributeName() != "Strength" {
		t.Errorf("found.AttributeName() = %v, want %v", found.AttributeName(), "Strength")
	}

	if found.Value() != 10 {
		t.Errorf("found.Value() = %v, want %v", found.Value(), 10)
	}
}

func TestPostgresCharacterAttributeRepository_FindByCharacterIDAndName_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	_, err := attrRepo.FindByCharacterIDAndName(context.Background(), character.ID(), "NonExistent")
	if err == nil {
		t.Error("FindByCharacterIDAndName() error = nil, want error for non-existent attribute")
	}
}

func TestPostgresCharacterAttributeRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	attribute, _ := entity.NewCharacterAttribute("attr-123", "Strength", 10, character.ID())
	attrRepo.Create(context.Background(), attribute)

	// Update attribute value
	attribute.UpdateValue(25)

	err := attrRepo.Update(context.Background(), attribute)
	if err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}

	// Verify update
	found, _ := attrRepo.FindByID(context.Background(), attribute.ID())

	if found.Value() != 25 {
		t.Errorf("found.Value() = %v, want %v", found.Value(), 25)
	}
}

func TestPostgresCharacterAttributeRepository_Update_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	// Try to update non-existent attribute
	attribute, _ := entity.NewCharacterAttribute("non-existent", "Strength", 10, character.ID())

	err := attrRepo.Update(context.Background(), attribute)
	if err == nil {
		t.Error("Update() error = nil, want error for non-existent attribute")
	}
}

func TestPostgresCharacterAttributeRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	attribute, _ := entity.NewCharacterAttribute("attr-123", "Strength", 10, character.ID())
	attrRepo.Create(context.Background(), attribute)

	err := attrRepo.Delete(context.Background(), attribute.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	// Verify deletion
	_, err = attrRepo.FindByID(context.Background(), attribute.ID())
	if err == nil {
		t.Error("FindByID() after Delete() should return error")
	}
}

func TestPostgresCharacterAttributeRepository_Delete_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	err := attrRepo.Delete(context.Background(), "non-existent-id")
	if err == nil {
		t.Error("Delete() error = nil, want error for non-existent attribute")
	}
}

func TestPostgresCharacterAttributeRepository_ExistsByCharacterIDAndName(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	// Should not exist initially
	exists, err := attrRepo.ExistsByCharacterIDAndName(context.Background(), character.ID(), "Strength")
	if err != nil {
		t.Fatalf("ExistsByCharacterIDAndName() error = %v, want nil", err)
	}
	if exists {
		t.Error("ExistsByCharacterIDAndName() = true, want false before creation")
	}

	// Create attribute
	attribute, _ := entity.NewCharacterAttribute("attr-123", "Strength", 10, character.ID())
	attrRepo.Create(context.Background(), attribute)

	// Should exist now
	exists, err = attrRepo.ExistsByCharacterIDAndName(context.Background(), character.ID(), "Strength")
	if err != nil {
		t.Fatalf("ExistsByCharacterIDAndName() error = %v, want nil", err)
	}
	if !exists {
		t.Error("ExistsByCharacterIDAndName() = false, want true after creation")
	}
}

func TestPostgresCharacterAttributeRepository_ForeignKeyConstraint(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	// Try to create attribute with non-existent character ID
	attribute, _ := entity.NewCharacterAttribute("attr-123", "Strength", 10, "non-existent-character")

	err := attrRepo.Create(context.Background(), attribute)
	if err == nil {
		t.Error("Create() with non-existent character should fail due to foreign key constraint")
	}
}

func TestPostgresCharacterAttributeRepository_UniqueConstraint(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	// Create first attribute
	attribute1, _ := entity.NewCharacterAttribute("attr-123", "Strength", 10, character.ID())
	err := attrRepo.Create(context.Background(), attribute1)
	if err != nil {
		t.Fatalf("First Create() error = %v, want nil", err)
	}

	// Try to create second attribute with same name for same character
	attribute2, _ := entity.NewCharacterAttribute("attr-456", "Strength", 15, character.ID())
	err = attrRepo.Create(context.Background(), attribute2)
	if err == nil {
		t.Error("Second Create() should fail due to unique constraint on (character_id, attribute_name)")
	}
}

func TestPostgresCharacterAttributeRepository_CascadeDelete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)
	attrRepo := persistence.NewPostgresCharacterAttributeRepository(db)

	character := createTestCharacter(t, userRepo, charRepo)

	// Create attribute
	attribute, _ := entity.NewCharacterAttribute("attr-123", "Strength", 10, character.ID())
	attrRepo.Create(context.Background(), attribute)

	// Delete character
	err := charRepo.Delete(context.Background(), character.ID())
	if err != nil {
		t.Fatalf("Delete character error = %v, want nil", err)
	}

	// Attribute should be automatically deleted due to CASCADE
	_, err = attrRepo.FindByID(context.Background(), attribute.ID())
	if err == nil {
		t.Error("FindByID() after character deletion should return error (cascade delete)")
	}
}
