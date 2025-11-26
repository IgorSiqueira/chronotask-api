package dto

// CharacterAttributeResponse represents a character attribute in the response
type CharacterAttributeResponse struct {
	ID            string `json:"id"`
	AttributeName string `json:"attributeName"`
	Value         int    `json:"value"`
	CharacterID   string `json:"characterId"`
	CreatedAt     string `json:"createdAt"`
}

// GetCharacterAttributesResponse represents the response when fetching all attributes
type GetCharacterAttributesResponse struct {
	CharacterID string                       `json:"characterId"`
	Attributes  []CharacterAttributeResponse `json:"attributes"`
}
