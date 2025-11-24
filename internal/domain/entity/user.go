package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/igor/chronotask-api/internal/domain/valueobject"
)

// User represents a user in the system (Domain Entity)
type User struct {
	id           string
	fullName     string
	email        valueobject.Email
	password     string // Hashed password
	birthDate    time.Time
	acceptTerms  bool
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUser creates a new User entity with validation
func NewUser(
	id string,
	fullName string,
	email valueobject.Email,
	hashedPassword string,
	birthDate time.Time,
	acceptTerms bool,
) (*User, error) {
	// Validate full name
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return nil, fmt.Errorf("full name cannot be empty")
	}
	if len(fullName) < 2 {
		return nil, fmt.Errorf("full name must be at least 2 characters")
	}
	if len(fullName) > 255 {
		return nil, fmt.Errorf("full name cannot exceed 255 characters")
	}

	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("user id cannot be empty")
	}

	// Validate hashed password
	if hashedPassword == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Validate birth date (must be in the past)
	if birthDate.After(time.Now()) {
		return nil, fmt.Errorf("birth date cannot be in the future")
	}

	// Validate minimum age (e.g., 13 years old)
	minAge := time.Now().AddDate(-13, 0, 0)
	if birthDate.After(minAge) {
		return nil, fmt.Errorf("user must be at least 13 years old")
	}

	// Validate terms acceptance
	if !acceptTerms {
		return nil, fmt.Errorf("user must accept terms and conditions")
	}

	now := time.Now()

	return &User{
		id:           id,
		fullName:     fullName,
		email:        email,
		password:     hashedPassword,
		birthDate:    birthDate,
		acceptTerms:  acceptTerms,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Getters (Read-only access to ensure encapsulation)

func (u *User) ID() string {
	return u.id
}

func (u *User) FullName() string {
	return u.fullName
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) BirthDate() time.Time {
	return u.birthDate
}

func (u *User) AcceptTerms() bool {
	return u.acceptTerms
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// Business Methods

// UpdateFullName updates the user's full name
func (u *User) UpdateFullName(fullName string) error {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	if len(fullName) < 2 {
		return fmt.Errorf("full name must be at least 2 characters")
	}
	if len(fullName) > 255 {
		return fmt.Errorf("full name cannot exceed 255 characters")
	}

	u.fullName = fullName
	u.updatedAt = time.Now()
	return nil
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email valueobject.Email) {
	u.email = email
	u.updatedAt = time.Now()
}

// UpdatePassword updates the user's hashed password
func (u *User) UpdatePassword(hashedPassword string) error {
	if hashedPassword == "" {
		return fmt.Errorf("password cannot be empty")
	}

	u.password = hashedPassword
	u.updatedAt = time.Now()
	return nil
}

// Age calculates the user's age in years
func (u *User) Age() int {
	now := time.Now()
	age := now.Year() - u.birthDate.Year()

	// Adjust if birthday hasn't occurred this year
	if now.Month() < u.birthDate.Month() ||
		(now.Month() == u.birthDate.Month() && now.Day() < u.birthDate.Day()) {
		age--
	}

	return age
}

// Reconstitute creates a User from existing data (for repository loading)
func ReconstituteUser(
	id string,
	fullName string,
	email valueobject.Email,
	hashedPassword string,
	birthDate time.Time,
	acceptTerms bool,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:           id,
		fullName:     fullName,
		email:        email,
		password:     hashedPassword,
		birthDate:    birthDate,
		acceptTerms:  acceptTerms,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
