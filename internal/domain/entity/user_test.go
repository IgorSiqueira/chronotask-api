package entity_test

import (
	"testing"
	"time"

	"github.com/igor/chronotask-api/internal/domain/entity"
	"github.com/igor/chronotask-api/internal/domain/valueobject"
)

func TestNewUser_ValidUser(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	user, err := entity.NewUser(
		"user-123",
		"John Doe",
		email,
		"hashed_password_here",
		birthDate,
		true,
	)

	if err != nil {
		t.Fatalf("NewUser() error = %v, want nil", err)
	}

	if user.ID() != "user-123" {
		t.Errorf("ID() = %v, want %v", user.ID(), "user-123")
	}

	if user.FullName() != "John Doe" {
		t.Errorf("FullName() = %v, want %v", user.FullName(), "John Doe")
	}

	if !user.Email().Equals(email) {
		t.Errorf("Email() = %v, want %v", user.Email(), email)
	}

	if user.Password() != "hashed_password_here" {
		t.Errorf("Password() = %v, want %v", user.Password(), "hashed_password_here")
	}

	if user.BirthDate() != birthDate {
		t.Errorf("BirthDate() = %v, want %v", user.BirthDate(), birthDate)
	}

	if !user.AcceptTerms() {
		t.Error("AcceptTerms() = false, want true")
	}

	if user.CreatedAt().IsZero() {
		t.Error("CreatedAt() should not be zero")
	}

	if user.UpdatedAt().IsZero() {
		t.Error("UpdatedAt() should not be zero")
	}
}

func TestNewUser_InvalidFullName(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		fullName string
	}{
		{"empty", ""},
		{"too short", "A"},
		{"only spaces", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entity.NewUser(
				"user-123",
				tt.fullName,
				email,
				"hashed_password",
				birthDate,
				true,
			)
			if err == nil {
				t.Error("NewUser() error = nil, want error for invalid full name")
			}
		})
	}
}

func TestNewUser_FutureBirthDate(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	futureBirthDate := time.Now().AddDate(1, 0, 0)

	_, err := entity.NewUser(
		"user-123",
		"John Doe",
		email,
		"hashed_password",
		futureBirthDate,
		true,
	)

	if err == nil {
		t.Error("NewUser() error = nil, want error for future birth date")
	}
}

func TestNewUser_TooYoung(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	recentBirthDate := time.Now().AddDate(-10, 0, 0) // 10 years old

	_, err := entity.NewUser(
		"user-123",
		"John Doe",
		email,
		"hashed_password",
		recentBirthDate,
		true,
	)

	if err == nil {
		t.Error("NewUser() error = nil, want error for user under 13")
	}
}

func TestNewUser_TermsNotAccepted(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := entity.NewUser(
		"user-123",
		"John Doe",
		email,
		"hashed_password",
		birthDate,
		false, // Terms not accepted
	)

	if err == nil {
		t.Error("NewUser() error = nil, want error for terms not accepted")
	}
}

func TestUser_UpdateFullName(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	user, _ := entity.NewUser(
		"user-123",
		"John Doe",
		email,
		"hashed_password",
		birthDate,
		true,
	)

	originalUpdatedAt := user.UpdatedAt()
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	err := user.UpdateFullName("Jane Smith")
	if err != nil {
		t.Fatalf("UpdateFullName() error = %v, want nil", err)
	}

	if user.FullName() != "Jane Smith" {
		t.Errorf("FullName() = %v, want %v", user.FullName(), "Jane Smith")
	}

	if !user.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt() should be updated after UpdateFullName()")
	}
}

func TestUser_Age(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")

	// Create user born 25 years ago
	birthDate := time.Now().AddDate(-25, 0, 0)

	user, _ := entity.NewUser(
		"user-123",
		"John Doe",
		email,
		"hashed_password",
		birthDate,
		true,
	)

	age := user.Age()
	if age != 25 {
		t.Errorf("Age() = %v, want %v", age, 25)
	}
}

func TestReconstituteUser(t *testing.T) {
	email, _ := valueobject.NewEmail("user@example.com")
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)

	user := entity.ReconstituteUser(
		"user-123",
		"John Doe",
		email,
		"hashed_password",
		birthDate,
		true,
		createdAt,
		updatedAt,
	)

	if user == nil {
		t.Fatal("ReconstituteUser() returned nil")
	}

	if user.ID() != "user-123" {
		t.Errorf("ID() = %v, want %v", user.ID(), "user-123")
	}

	if user.CreatedAt() != createdAt {
		t.Errorf("CreatedAt() = %v, want %v", user.CreatedAt(), createdAt)
	}

	if user.UpdatedAt() != updatedAt {
		t.Errorf("UpdatedAt() = %v, want %v", user.UpdatedAt(), updatedAt)
	}
}
