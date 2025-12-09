package entity

import (
	"fmt"
	"strings"
	"time"
)

// CharacterAttribute represents a character's attribute (Domain Entity)
type CharacterAttribute struct {
	id            int
	attributeName string
	value         int
	characterID   string
	createdAt     time.Time
}

// NewCharacterAttribute creates a new CharacterAttribute entity with validation
func NewCharacterAttribute(
	attributeName string,
	value int,
	characterID string,
) (*CharacterAttribute, error) {
	// Validate attribute name
	attributeName = strings.TrimSpace(attributeName)
	if attributeName == "" {
		return nil, fmt.Errorf("attribute name cannot be empty")
	}
	if len(attributeName) < 2 {
		return nil, fmt.Errorf("attribute name must be at least 2 characters")
	}
	if len(attributeName) > 50 {
		return nil, fmt.Errorf("attribute name cannot exceed 50 characters")
	}

	// Validate value (attributes must be non-negative)
	if value < 0 {
		return nil, fmt.Errorf("attribute value cannot be negative")
	}

	// Validate character ID
	if characterID == "" {
		return nil, fmt.Errorf("character id cannot be empty")
	}

	return &CharacterAttribute{
		id:            0, // Will be set by database sequence
		attributeName: attributeName,
		value:         value,
		characterID:   characterID,
		createdAt:     time.Now(),
	}, nil
}

// Getters (Read-only access to ensure encapsulation)

func (ca *CharacterAttribute) ID() int {
	return ca.id
}

func (ca *CharacterAttribute) AttributeName() string {
	return ca.attributeName
}

func (ca *CharacterAttribute) Value() int {
	return ca.value
}

func (ca *CharacterAttribute) CharacterID() string {
	return ca.characterID
}

func (ca *CharacterAttribute) CreatedAt() time.Time {
	return ca.createdAt
}

// Business Methods

// UpdateValue updates the attribute's value
func (ca *CharacterAttribute) UpdateValue(newValue int) error {
	if newValue < 0 {
		return fmt.Errorf("attribute value cannot be negative")
	}

	ca.value = newValue
	return nil
}

// IncrementValue increments the attribute's value by the given amount
func (ca *CharacterAttribute) IncrementValue(amount int) error {
	if amount < 0 {
		return fmt.Errorf("increment amount cannot be negative")
	}

	ca.value += amount
	return nil
}

// DecrementValue decrements the attribute's value by the given amount
// Will not allow the value to go below 0
func (ca *CharacterAttribute) DecrementValue(amount int) error {
	if amount < 0 {
		return fmt.Errorf("decrement amount cannot be negative")
	}

	newValue := ca.value - amount
	if newValue < 0 {
		return fmt.Errorf("cannot decrement below 0 (current: %d, decrement: %d)", ca.value, amount)
	}

	ca.value = newValue
	return nil
}

// ReconstituteCharacterAttribute creates a CharacterAttribute from existing data (for repository loading)
func ReconstituteCharacterAttribute(
	id int,
	attributeName string,
	value int,
	characterID string,
	createdAt time.Time,
) *CharacterAttribute {
	return &CharacterAttribute{
		id:            id,
		attributeName: attributeName,
		value:         value,
		characterID:   characterID,
		createdAt:     createdAt,
	}
}
