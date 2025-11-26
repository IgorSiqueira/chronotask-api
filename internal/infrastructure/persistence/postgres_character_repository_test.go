package persistence_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/igor/chronotask-api/config"
	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/repository"
	"github.com/igor/chronotask-api/internal/domain/valueobject"
	"github.com/igor/chronotask-api/internal/infrastructure/persistence"
)

// setupTestDB creates a test database connection
// Set TEST_DATABASE_URL environment variable to run integration tests
// Example: TEST_DATABASE_URL=postgres://user:pass@localhost:5432/chronotask_test
func setupTestDB(t *testing.T) (*persistence.PostgresDB, func()) {
	t.Helper()

	// Skip if no test database is configured
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("Skipping integration test: TEST_DATABASE_URL not set")
	}

	cfg := &config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME") + "_test",
		SSLMode:  "disable",
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	if cfg.User == "" {
		cfg.User = "postgres"
	}
	if cfg.Password == "" {
		cfg.Password = "postgres"
	}
	if cfg.DBName == "_test" {
		cfg.DBName = "chronotask_test"
	}

	db, err := persistence.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		ctx := context.Background()
		db.Pool.Exec(ctx, "DELETE FROM characters")
		db.Pool.Exec(ctx, "DELETE FROM users")
		db.Close()
	}

	return db, cleanup
}

// createTestUser creates a user for testing character operations
func createTestUser(t *testing.T, userRepo repository.UserRepository) *entity.User {
	t.Helper()

	email, err := valueobject.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create test email: %v", err)
	}

	user, err := entity.NewUser(
		"test-user-id",
		"Test User",
		email,
		"hashed_password",
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		true,
	)
	if err != nil {
		t.Fatalf("Failed to create test user entity: %v", err)
	}

	err = userRepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}

	return user
}

func TestPostgresCharacterRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	// Create test user first
	user := createTestUser(t, userRepo)

	// Create character
	character, err := entity.NewCharacter(
		"char-123",
		"Warrior King",
		user.ID(),
	)
	if err != nil {
		t.Fatalf("Failed to create character entity: %v", err)
	}

	err = charRepo.Create(context.Background(), character)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	// Verify creation
	found, err := charRepo.FindByID(context.Background(), character.ID())
	if err != nil {
		t.Fatalf("FindByID() error = %v, want nil", err)
	}

	if found.ID() != character.ID() {
		t.Errorf("found.ID() = %v, want %v", found.ID(), character.ID())
	}

	if found.Name() != character.Name() {
		t.Errorf("found.Name() = %v, want %v", found.Name(), character.Name())
	}
}

func TestPostgresCharacterRepository_FindByID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	charRepo := persistence.NewPostgresCharacterRepository(db)

	_, err := charRepo.FindByID(context.Background(), "non-existent-id")
	if err == nil {
		t.Error("FindByID() error = nil, want error for non-existent character")
	}
}

func TestPostgresCharacterRepository_FindByUserID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	user := createTestUser(t, userRepo)

	character, _ := entity.NewCharacter("char-123", "Warrior King", user.ID())
	charRepo.Create(context.Background(), character)

	found, err := charRepo.FindByUserID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("FindByUserID() error = %v, want nil", err)
	}

	if found.UserID() != user.ID() {
		t.Errorf("found.UserID() = %v, want %v", found.UserID(), user.ID())
	}
}

func TestPostgresCharacterRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	user := createTestUser(t, userRepo)

	character, _ := entity.NewCharacter("char-123", "Warrior King", user.ID())
	charRepo.Create(context.Background(), character)

	// Add XP and level up
	character.AddXp(150)
	character.UpdateName("Mighty Warrior")

	err := charRepo.Update(context.Background(), character)
	if err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}

	// Verify update
	found, _ := charRepo.FindByID(context.Background(), character.ID())

	if found.Name() != "Mighty Warrior" {
		t.Errorf("found.Name() = %v, want %v", found.Name(), "Mighty Warrior")
	}

	if found.Level() != character.Level() {
		t.Errorf("found.Level() = %v, want %v", found.Level(), character.Level())
	}

	if found.TotalXp() != 150 {
		t.Errorf("found.TotalXp() = %v, want %v", found.TotalXp(), 150)
	}
}

func TestPostgresCharacterRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	user := createTestUser(t, userRepo)

	character, _ := entity.NewCharacter("char-123", "Warrior King", user.ID())
	charRepo.Create(context.Background(), character)

	err := charRepo.Delete(context.Background(), character.ID())
	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	// Verify deletion
	_, err = charRepo.FindByID(context.Background(), character.ID())
	if err == nil {
		t.Error("FindByID() after Delete() should return error")
	}
}

func TestPostgresCharacterRepository_ExistsByUserID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	user := createTestUser(t, userRepo)

	// Should not exist initially
	exists, err := charRepo.ExistsByUserID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("ExistsByUserID() error = %v, want nil", err)
	}
	if exists {
		t.Error("ExistsByUserID() = true, want false before creation")
	}

	// Create character
	character, _ := entity.NewCharacter("char-123", "Warrior King", user.ID())
	charRepo.Create(context.Background(), character)

	// Should exist now
	exists, err = charRepo.ExistsByUserID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("ExistsByUserID() error = %v, want nil", err)
	}
	if !exists {
		t.Error("ExistsByUserID() = false, want true after creation")
	}
}

func TestPostgresCharacterRepository_ForeignKeyConstraint(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	charRepo := persistence.NewPostgresCharacterRepository(db)

	// Try to create character with non-existent user ID
	character, _ := entity.NewCharacter("char-123", "Warrior King", "non-existent-user")

	err := charRepo.Create(context.Background(), character)
	if err == nil {
		t.Error("Create() with non-existent user should fail due to foreign key constraint")
	}
}

func TestPostgresCharacterRepository_UniqueUserConstraint(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	user := createTestUser(t, userRepo)

	// Create first character
	character1, _ := entity.NewCharacter("char-123", "Warrior King", user.ID())
	err := charRepo.Create(context.Background(), character1)
	if err != nil {
		t.Fatalf("First Create() error = %v, want nil", err)
	}

	// Try to create second character for same user
	character2, _ := entity.NewCharacter("char-456", "Another Warrior", user.ID())
	err = charRepo.Create(context.Background(), character2)
	if err == nil {
		t.Error("Second Create() should fail due to unique user constraint")
	}
}

func TestPostgresCharacterRepository_FindByIDAndUserID_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	user := createTestUser(t, userRepo)

	character, _ := entity.NewCharacter("char-123", "Warrior King", user.ID())
	charRepo.Create(context.Background(), character)

	// Find character with correct user ID
	found, err := charRepo.FindByIDAndUserID(context.Background(), character.ID(), user.ID())
	if err != nil {
		t.Fatalf("FindByIDAndUserID() error = %v, want nil", err)
	}

	if found.ID() != character.ID() {
		t.Errorf("found.ID() = %v, want %v", found.ID(), character.ID())
	}

	if found.UserID() != user.ID() {
		t.Errorf("found.UserID() = %v, want %v", found.UserID(), user.ID())
	}
}

func TestPostgresCharacterRepository_FindByIDAndUserID_WrongUser(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := persistence.NewPostgresUserRepository(db)
	charRepo := persistence.NewPostgresCharacterRepository(db)

	user := createTestUser(t, userRepo)

	character, _ := entity.NewCharacter("char-123", "Warrior King", user.ID())
	charRepo.Create(context.Background(), character)

	// Try to find character with wrong user ID
	_, err := charRepo.FindByIDAndUserID(context.Background(), character.ID(), "wrong-user-id")
	if err == nil {
		t.Error("FindByIDAndUserID() with wrong user should return error")
	}
}

func TestPostgresCharacterRepository_FindByIDAndUserID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	charRepo := persistence.NewPostgresCharacterRepository(db)

	// Try to find non-existent character
	_, err := charRepo.FindByIDAndUserID(context.Background(), "non-existent", "any-user")
	if err == nil {
		t.Error("FindByIDAndUserID() for non-existent character should return error")
	}
}
