package entity

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Character represents a user's game character (Domain Entity)
type Character struct {
	id        string
	name      string
	level     int
	currentXp int
	totalXp   int
	userID    string
	createdAt time.Time
}

// NewCharacter creates a new Character entity with validation
func NewCharacter(
	id string,
	name string,
	userID string,
) (*Character, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("character id cannot be empty")
	}

	// Validate name
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("character name cannot be empty")
	}
	if len(name) < 2 {
		return nil, fmt.Errorf("character name must be at least 2 characters")
	}
	if len(name) > 50 {
		return nil, fmt.Errorf("character name cannot exceed 50 characters")
	}

	// Validate user ID
	if userID == "" {
		return nil, fmt.Errorf("user id cannot be empty")
	}

	return &Character{
		id:        id,
		name:      name,
		level:     1,      // Characters start at level 1
		currentXp: 0,      // Start with 0 XP
		totalXp:   0,      // Start with 0 total XP
		userID:    userID,
		createdAt: time.Now(),
	}, nil
}

// Getters (Read-only access to ensure encapsulation)

func (c *Character) ID() string {
	return c.id
}

func (c *Character) Name() string {
	return c.name
}

func (c *Character) Level() int {
	return c.level
}

func (c *Character) CurrentXp() int {
	return c.currentXp
}

func (c *Character) TotalXp() int {
	return c.totalXp
}

func (c *Character) UserID() string {
	return c.userID
}

func (c *Character) CreatedAt() time.Time {
	return c.createdAt
}

// Business Methods

// UpdateName updates the character's name
func (c *Character) UpdateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("character name cannot be empty")
	}
	if len(name) < 2 {
		return fmt.Errorf("character name must be at least 2 characters")
	}
	if len(name) > 50 {
		return fmt.Errorf("character name cannot exceed 50 characters")
	}

	c.name = name
	return nil
}

// AddXp adds experience points to the character and handles level-ups
// Returns the number of levels gained (0 if no level up)
func (c *Character) AddXp(xp int) (int, error) {
	if xp < 0 {
		return 0, fmt.Errorf("xp cannot be negative")
	}

	if xp == 0 {
		return 0, nil
	}

	c.currentXp += xp
	c.totalXp += xp

	// Check for level-ups
	levelsGained := 0
	for c.currentXp >= c.XpForNextLevel() {
		c.currentXp -= c.XpForNextLevel()
		c.level++
		levelsGained++
	}

	return levelsGained, nil
}

// XpForNextLevel calculates the XP required to reach the next level
// Uses a progressive formula: 100 * level^1.5
func (c *Character) XpForNextLevel() int {
	return int(math.Round(100 * math.Pow(float64(c.level), 1.5)))
}

// XpProgress returns the percentage of XP progress towards the next level (0-100)
func (c *Character) XpProgress() float64 {
	xpNeeded := c.XpForNextLevel()
	if xpNeeded == 0 {
		return 0
	}
	return (float64(c.currentXp) / float64(xpNeeded)) * 100
}

// ReconstituteCharacter creates a Character from existing data (for repository loading)
func ReconstituteCharacter(
	id string,
	name string,
	level int,
	currentXp int,
	totalXp int,
	userID string,
	createdAt time.Time,
) *Character {
	return &Character{
		id:        id,
		name:      name,
		level:     level,
		currentXp: currentXp,
		totalXp:   totalXp,
		userID:    userID,
		createdAt: createdAt,
	}
}
