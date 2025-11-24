package valueobject_test

import (
	"testing"

	"github.com/igor/chronotask-api/internal/domain/valueobject"
)

func TestNewEmail_ValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "user@example.com", "user@example.com"},
		{"uppercase converted", "USER@EXAMPLE.COM", "user@example.com"},
		{"with plus", "user+tag@example.com", "user+tag@example.com"},
		{"with dash", "user-name@example.com", "user-name@example.com"},
		{"with dots", "user.name@example.com", "user.name@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := valueobject.NewEmail(tt.input)
			if err != nil {
				t.Errorf("NewEmail() error = %v, want nil", err)
				return
			}
			if email.Value() != tt.want {
				t.Errorf("NewEmail() = %v, want %v", email.Value(), tt.want)
			}
		})
	}
}

func TestNewEmail_InvalidEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"no at", "userexample.com"},
		{"no domain", "user@"},
		{"no local", "@example.com"},
		{"multiple at", "user@@example.com"},
		{"no tld", "user@example"},
		{"spaces", "user @example.com"},
		{"invalid chars", "user<>@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := valueobject.NewEmail(tt.input)
			if err == nil {
				t.Errorf("NewEmail() error = nil, want error")
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := valueobject.NewEmail("user@example.com")
	email2, _ := valueobject.NewEmail("user@example.com")
	email3, _ := valueobject.NewEmail("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Equal emails should be equal")
	}

	if email1.Equals(email3) {
		t.Error("Different emails should not be equal")
	}
}
