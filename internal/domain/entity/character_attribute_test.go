package entity_test

import (
	"testing"
	"time"

	"github.com/igor/chronotask-api/internal/domain/entity"
)

func TestNewCharacterAttribute_ValidAttribute(t *testing.T) {
	attribute, err := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	if err != nil {
		t.Fatalf("NewCharacterAttribute() error = %v, want nil", err)
	}

	if attribute.ID() != "attr-123" {
		t.Errorf("ID() = %v, want %v", attribute.ID(), "attr-123")
	}

	if attribute.AttributeName() != "Strength" {
		t.Errorf("AttributeName() = %v, want %v", attribute.AttributeName(), "Strength")
	}

	if attribute.Value() != 10 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 10)
	}

	if attribute.CharacterID() != "char-456" {
		t.Errorf("CharacterID() = %v, want %v", attribute.CharacterID(), "char-456")
	}

	if attribute.CreatedAt().IsZero() {
		t.Error("CreatedAt() should not be zero")
	}
}

func TestNewCharacterAttribute_ZeroValue(t *testing.T) {
	attribute, err := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		0,
		"char-456",
	)

	if err != nil {
		t.Fatalf("NewCharacterAttribute() error = %v, want nil", err)
	}

	if attribute.Value() != 0 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 0)
	}
}

func TestNewCharacterAttribute_InvalidAttributeName(t *testing.T) {
	tests := []struct {
		name          string
		attributeName string
	}{
		{"empty", ""},
		{"too short", "A"},
		{"only spaces", "   "},
		{"too long", "ThisIsAVeryLongAttributeNameThatExceedsFiftyCharactersLimit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.NewCharacterAttribute(
				"attr-123",
				tt.attributeName,
				10,
				"char-456",
			)
			if err == nil {
				t.Error("NewCharacterAttribute() error = nil, want error for invalid attribute name")
			}
		})
	}
}

func TestNewCharacterAttribute_InvalidID(t *testing.T) {
	_, err := entity.NewCharacterAttribute(
		"",
		"Strength",
		10,
		"char-456",
	)

	if err == nil {
		t.Error("NewCharacterAttribute() error = nil, want error for empty ID")
	}
}

func TestNewCharacterAttribute_InvalidCharacterID(t *testing.T) {
	_, err := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"",
	)

	if err == nil {
		t.Error("NewCharacterAttribute() error = nil, want error for empty CharacterID")
	}
}

func TestNewCharacterAttribute_NegativeValue(t *testing.T) {
	_, err := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		-5,
		"char-456",
	)

	if err == nil {
		t.Error("NewCharacterAttribute() error = nil, want error for negative value")
	}
}

func TestCharacterAttribute_UpdateValue(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.UpdateValue(25)
	if err != nil {
		t.Fatalf("UpdateValue() error = %v, want nil", err)
	}

	if attribute.Value() != 25 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 25)
	}
}

func TestCharacterAttribute_UpdateValue_ToZero(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.UpdateValue(0)
	if err != nil {
		t.Fatalf("UpdateValue() error = %v, want nil", err)
	}

	if attribute.Value() != 0 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 0)
	}
}

func TestCharacterAttribute_UpdateValue_Negative(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.UpdateValue(-5)
	if err == nil {
		t.Error("UpdateValue() error = nil, want error for negative value")
	}

	// Value should remain unchanged
	if attribute.Value() != 10 {
		t.Errorf("Value() = %v, want %v (should remain unchanged)", attribute.Value(), 10)
	}
}

func TestCharacterAttribute_IncrementValue(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.IncrementValue(5)
	if err != nil {
		t.Fatalf("IncrementValue() error = %v, want nil", err)
	}

	if attribute.Value() != 15 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 15)
	}
}

func TestCharacterAttribute_IncrementValue_ByZero(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.IncrementValue(0)
	if err != nil {
		t.Fatalf("IncrementValue() error = %v, want nil", err)
	}

	if attribute.Value() != 10 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 10)
	}
}

func TestCharacterAttribute_IncrementValue_Negative(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.IncrementValue(-5)
	if err == nil {
		t.Error("IncrementValue() error = nil, want error for negative amount")
	}

	// Value should remain unchanged
	if attribute.Value() != 10 {
		t.Errorf("Value() = %v, want %v (should remain unchanged)", attribute.Value(), 10)
	}
}

func TestCharacterAttribute_DecrementValue(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.DecrementValue(3)
	if err != nil {
		t.Fatalf("DecrementValue() error = %v, want nil", err)
	}

	if attribute.Value() != 7 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 7)
	}
}

func TestCharacterAttribute_DecrementValue_ToZero(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.DecrementValue(10)
	if err != nil {
		t.Fatalf("DecrementValue() error = %v, want nil", err)
	}

	if attribute.Value() != 0 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 0)
	}
}

func TestCharacterAttribute_DecrementValue_BelowZero(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.DecrementValue(15)
	if err == nil {
		t.Error("DecrementValue() error = nil, want error when trying to go below zero")
	}

	// Value should remain unchanged
	if attribute.Value() != 10 {
		t.Errorf("Value() = %v, want %v (should remain unchanged)", attribute.Value(), 10)
	}
}

func TestCharacterAttribute_DecrementValue_ByZero(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.DecrementValue(0)
	if err != nil {
		t.Fatalf("DecrementValue() error = %v, want nil", err)
	}

	if attribute.Value() != 10 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 10)
	}
}

func TestCharacterAttribute_DecrementValue_Negative(t *testing.T) {
	attribute, _ := entity.NewCharacterAttribute(
		"attr-123",
		"Strength",
		10,
		"char-456",
	)

	err := attribute.DecrementValue(-5)
	if err == nil {
		t.Error("DecrementValue() error = nil, want error for negative amount")
	}

	// Value should remain unchanged
	if attribute.Value() != 10 {
		t.Errorf("Value() = %v, want %v (should remain unchanged)", attribute.Value(), 10)
	}
}

func TestReconstituteCharacterAttribute(t *testing.T) {
	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	attribute := entity.ReconstituteCharacterAttribute(
		"attr-123",
		"Strength",
		50,
		"char-456",
		createdAt,
	)

	if attribute == nil {
		t.Fatal("ReconstituteCharacterAttribute() returned nil")
	}

	if attribute.ID() != "attr-123" {
		t.Errorf("ID() = %v, want %v", attribute.ID(), "attr-123")
	}

	if attribute.AttributeName() != "Strength" {
		t.Errorf("AttributeName() = %v, want %v", attribute.AttributeName(), "Strength")
	}

	if attribute.Value() != 50 {
		t.Errorf("Value() = %v, want %v", attribute.Value(), 50)
	}

	if attribute.CharacterID() != "char-456" {
		t.Errorf("CharacterID() = %v, want %v", attribute.CharacterID(), "char-456")
	}

	if attribute.CreatedAt() != createdAt {
		t.Errorf("CreatedAt() = %v, want %v", attribute.CreatedAt(), createdAt)
	}
}
