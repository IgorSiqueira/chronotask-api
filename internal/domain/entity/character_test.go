package entity_test

import (
	"testing"
	"time"

	"github.com/igor/chronotask-api/internal/domain/entity"
)

func TestNewCharacter_ValidCharacter(t *testing.T) {
	character, err := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	if err != nil {
		t.Fatalf("NewCharacter() error = %v, want nil", err)
	}

	if character.ID() != "char-123" {
		t.Errorf("ID() = %v, want %v", character.ID(), "char-123")
	}

	if character.Name() != "Warrior King" {
		t.Errorf("Name() = %v, want %v", character.Name(), "Warrior King")
	}

	if character.UserID() != "user-456" {
		t.Errorf("UserID() = %v, want %v", character.UserID(), "user-456")
	}

	// New characters should start at level 1 with 0 XP
	if character.Level() != 1 {
		t.Errorf("Level() = %v, want %v", character.Level(), 1)
	}

	if character.CurrentXp() != 0 {
		t.Errorf("CurrentXp() = %v, want %v", character.CurrentXp(), 0)
	}

	if character.TotalXp() != 0 {
		t.Errorf("TotalXp() = %v, want %v", character.TotalXp(), 0)
	}

	if character.CreatedAt().IsZero() {
		t.Error("CreatedAt() should not be zero")
	}
}

func TestNewCharacter_InvalidName(t *testing.T) {
	tests := []struct {
		name          string
		characterName string
	}{
		{"empty", ""},
		{"too short", "A"},
		{"only spaces", "   "},
		{"too long", "ThisIsAVeryLongNameThatExceedsFiftyCharactersLimit123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.NewCharacter(
				"char-123",
				tt.characterName,
				"user-456",
			)
			if err == nil {
				t.Error("NewCharacter() error = nil, want error for invalid name")
			}
		})
	}
}

func TestNewCharacter_InvalidID(t *testing.T) {
	_, err := entity.NewCharacter(
		"",
		"Warrior King",
		"user-456",
	)

	if err == nil {
		t.Error("NewCharacter() error = nil, want error for empty ID")
	}
}

func TestNewCharacter_InvalidUserID(t *testing.T) {
	_, err := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"",
	)

	if err == nil {
		t.Error("NewCharacter() error = nil, want error for empty UserID")
	}
}

func TestCharacter_UpdateName(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	err := character.UpdateName("Mighty Wizard")
	if err != nil {
		t.Fatalf("UpdateName() error = %v, want nil", err)
	}

	if character.Name() != "Mighty Wizard" {
		t.Errorf("Name() = %v, want %v", character.Name(), "Mighty Wizard")
	}
}

func TestCharacter_UpdateName_Invalid(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	tests := []struct {
		name string
		newName string
	}{
		{"empty", ""},
		{"too short", "A"},
		{"only spaces", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := character.UpdateName(tt.newName)
			if err == nil {
				t.Error("UpdateName() error = nil, want error for invalid name")
			}
		})
	}
}

func TestCharacter_AddXp_NoLevelUp(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	// Add 50 XP (not enough to level up from level 1)
	levelsGained, err := character.AddXp(50)

	if err != nil {
		t.Fatalf("AddXp() error = %v, want nil", err)
	}

	if levelsGained != 0 {
		t.Errorf("levelsGained = %v, want %v", levelsGained, 0)
	}

	if character.CurrentXp() != 50 {
		t.Errorf("CurrentXp() = %v, want %v", character.CurrentXp(), 50)
	}

	if character.TotalXp() != 50 {
		t.Errorf("TotalXp() = %v, want %v", character.TotalXp(), 50)
	}

	if character.Level() != 1 {
		t.Errorf("Level() = %v, want %v", character.Level(), 1)
	}
}

func TestCharacter_AddXp_SingleLevelUp(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	// Level 1 requires 100 XP for level 2
	levelsGained, err := character.AddXp(150)

	if err != nil {
		t.Fatalf("AddXp() error = %v, want nil", err)
	}

	if levelsGained != 1 {
		t.Errorf("levelsGained = %v, want %v", levelsGained, 1)
	}

	if character.Level() != 2 {
		t.Errorf("Level() = %v, want %v", character.Level(), 2)
	}

	if character.TotalXp() != 150 {
		t.Errorf("TotalXp() = %v, want %v", character.TotalXp(), 150)
	}

	// CurrentXp should be the remainder after leveling up
	// 150 - 100 (cost of level 1->2) = 50
	if character.CurrentXp() != 50 {
		t.Errorf("CurrentXp() = %v, want %v", character.CurrentXp(), 50)
	}
}

func TestCharacter_AddXp_MultipleLevelUps(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	// Add enough XP to gain multiple levels
	// Level 1->2: 100 XP
	// Level 2->3: ~282 XP
	// Total needed for 2 levels: ~382 XP
	levelsGained, err := character.AddXp(500)

	if err != nil {
		t.Fatalf("AddXp() error = %v, want nil", err)
	}

	if levelsGained < 2 {
		t.Errorf("levelsGained = %v, want at least 2", levelsGained)
	}

	if character.Level() < 3 {
		t.Errorf("Level() = %v, want at least 3", character.Level())
	}

	if character.TotalXp() != 500 {
		t.Errorf("TotalXp() = %v, want %v", character.TotalXp(), 500)
	}
}

func TestCharacter_AddXp_NegativeValue(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	_, err := character.AddXp(-50)

	if err == nil {
		t.Error("AddXp() error = nil, want error for negative XP")
	}
}

func TestCharacter_AddXp_Zero(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	levelsGained, err := character.AddXp(0)

	if err != nil {
		t.Fatalf("AddXp() error = %v, want nil", err)
	}

	if levelsGained != 0 {
		t.Errorf("levelsGained = %v, want %v", levelsGained, 0)
	}
}

func TestCharacter_XpForNextLevel(t *testing.T) {
	tests := []struct {
		level          int
		expectedXpMin  int
		expectedXpMax  int
	}{
		{1, 99, 101},      // Level 1->2: ~100 XP
		{2, 281, 283},     // Level 2->3: ~282 XP
		{3, 518, 520},     // Level 3->4: ~519 XP
		{10, 3161, 3163},  // Level 10->11: ~3162 XP
	}

	for _, tt := range tests {
		t.Run("level_"+string(rune(tt.level+'0')), func(t *testing.T) {
			character, _ := entity.NewCharacter("char-123", "Warrior", "user-456")

			// Manually set level for testing (using Reconstitute)
			character = entity.ReconstituteCharacter(
				"char-123",
				"Warrior",
				tt.level,
				0,
				0,
				"user-456",
				time.Now(),
			)

			xpNeeded := character.XpForNextLevel()

			if xpNeeded < tt.expectedXpMin || xpNeeded > tt.expectedXpMax {
				t.Errorf("XpForNextLevel() = %v, want between %v and %v",
					xpNeeded, tt.expectedXpMin, tt.expectedXpMax)
			}
		})
	}
}

func TestCharacter_XpProgress(t *testing.T) {
	character, _ := entity.NewCharacter(
		"char-123",
		"Warrior King",
		"user-456",
	)

	// Level 1 needs 100 XP
	character.AddXp(50) // 50/100 = 50%

	progress := character.XpProgress()

	if progress < 49.0 || progress > 51.0 {
		t.Errorf("XpProgress() = %v, want ~50.0", progress)
	}
}

func TestReconstituteCharacter(t *testing.T) {
	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	character := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		10,
		250,
		5000,
		"user-456",
		createdAt,
	)

	if character == nil {
		t.Fatal("ReconstituteCharacter() returned nil")
	}

	if character.ID() != "char-123" {
		t.Errorf("ID() = %v, want %v", character.ID(), "char-123")
	}

	if character.Name() != "Warrior King" {
		t.Errorf("Name() = %v, want %v", character.Name(), "Warrior King")
	}

	if character.Level() != 10 {
		t.Errorf("Level() = %v, want %v", character.Level(), 10)
	}

	if character.CurrentXp() != 250 {
		t.Errorf("CurrentXp() = %v, want %v", character.CurrentXp(), 250)
	}

	if character.TotalXp() != 5000 {
		t.Errorf("TotalXp() = %v, want %v", character.TotalXp(), 5000)
	}

	if character.UserID() != "user-456" {
		t.Errorf("UserID() = %v, want %v", character.UserID(), "user-456")
	}

	if character.CreatedAt() != createdAt {
		t.Errorf("CreatedAt() = %v, want %v", character.CreatedAt(), createdAt)
	}
}
