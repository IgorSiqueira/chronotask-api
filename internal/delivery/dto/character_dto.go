package dto

// CreateCharacterRequest represents the request to create a new character
type CreateCharacterRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

// CreateCharacterResponse represents the response after creating a character
type CreateCharacterResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	CurrentXp int    `json:"currentXp"`
	TotalXp   int    `json:"totalXp"`
	UserID    string `json:"userId"`
	CreatedAt string `json:"createdAt"`
}

// CharacterItemResponse represents a character in a list
type CharacterItemResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	CurrentXp int    `json:"currentXp"`
	TotalXp   int    `json:"totalXp"`
	CreatedAt string `json:"createdAt"`
}

// GetUserCharactersResponse represents the response when fetching user's characters
type GetUserCharactersResponse struct {
	Characters []CharacterItemResponse `json:"characters"`
}
