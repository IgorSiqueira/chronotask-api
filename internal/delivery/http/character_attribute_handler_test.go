package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/internal/application/usecase"
	deliveryHttp "github.com/igor/chronotask-api/internal/delivery/http"
	"github.com/igor/chronotask-api/internal/delivery/http/middleware"
	"github.com/igor/chronotask-api/internal/domain/entity"
)

// Mock CharacterAttributeRepository for E2E tests
type mockCharacterAttributeRepository struct {
	findByCharacterIDFunc func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error)
}

func (m *mockCharacterAttributeRepository) Create(ctx context.Context, attribute *entity.CharacterAttribute) error {
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) FindByID(ctx context.Context, id string) (*entity.CharacterAttribute, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) FindByCharacterID(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
	if m.findByCharacterIDFunc != nil {
		return m.findByCharacterIDFunc(ctx, characterID)
	}
	return []*entity.CharacterAttribute{}, nil
}

func (m *mockCharacterAttributeRepository) FindByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (*entity.CharacterAttribute, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) Update(ctx context.Context, attribute *entity.CharacterAttribute) error {
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) ExistsByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (bool, error) {
	return false, errors.New("not implemented")
}

// Mock CharacterRepository for attribute tests
type mockCharacterRepositoryForAttributeTests struct {
	findByIDFunc          func(ctx context.Context, id string) (*entity.Character, error)
	findByIDAndUserIDFunc func(ctx context.Context, id string, userID string) (*entity.Character, error)
}

func (m *mockCharacterRepositoryForAttributeTests) Create(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributeTests) FindByID(ctx context.Context, id string) (*entity.Character, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockCharacterRepositoryForAttributeTests) FindByIDAndUserID(ctx context.Context, id string, userID string) (*entity.Character, error) {
	if m.findByIDAndUserIDFunc != nil {
		return m.findByIDAndUserIDFunc(ctx, id, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockCharacterRepositoryForAttributeTests) FindByUserID(ctx context.Context, userID string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributeTests) Update(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributeTests) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributeTests) FindAllByUserID(ctx context.Context, userID string) ([]*entity.Character, error) {
	return []*entity.Character{}, errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributeTests) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	return false, errors.New("not implemented")
}

// setupTestRouterForAttributes creates a test router with attribute endpoints
func setupTestRouterForAttributes(charRepo *mockCharacterRepositoryForAttributeTests, attrRepo *mockCharacterAttributeRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Create use case
	getAttributesUseCase := usecase.NewGetCharacterAttributesUseCase(charRepo, attrRepo)

	// Create handler
	attributeHandler := deliveryHttp.NewCharacterAttributeHandler(getAttributesUseCase)

	// Create auth middleware with mock JWT service
	authMiddleware := middleware.NewAuthMiddleware(&mockJWTService{})

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		authenticated := v1.Group("")
		authenticated.Use(authMiddleware.RequireAuth())
		{
			authenticated.GET("/character/:characterId/attribute", attributeHandler.GetByCharacterID)
		}
	}

	return router
}

func TestCharacterAttributeHandler_GetByCharacterID_Success(t *testing.T) {
	mockChar := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		5,
		50,
		500,
		"test-user-123", // Must match JWT mock userID
		time.Now(),
	)

	mockAttributes := []*entity.CharacterAttribute{
		entity.ReconstituteCharacterAttribute("attr-1", "Strength", 10, "char-123", time.Now()),
		entity.ReconstituteCharacterAttribute("attr-2", "Agility", 15, "char-123", time.Now()),
		entity.ReconstituteCharacterAttribute("attr-3", "Intelligence", 20, "char-123", time.Now()),
	}

	mockCharRepo := &mockCharacterRepositoryForAttributeTests{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			if id == "char-123" && userID == "test-user-123" {
				return mockChar, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{
		findByCharacterIDFunc: func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
			if characterID == "char-123" {
				return mockAttributes, nil
			}
			return []*entity.CharacterAttribute{}, nil
		},
	}

	router := setupTestRouterForAttributes(mockCharRepo, mockAttrRepo)

	req, _ := http.NewRequest("GET", "/api/v1/character/char-123/attribute", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["characterId"] != "char-123" {
		t.Errorf("response characterId = %v, want %v", response["characterId"], "char-123")
	}

	attributes := response["attributes"].([]interface{})
	if len(attributes) != 3 {
		t.Errorf("len(attributes) = %v, want %v", len(attributes), 3)
	}

	// Verify first attribute
	attr1 := attributes[0].(map[string]interface{})
	if attr1["attributeName"] != "Strength" {
		t.Errorf("attributes[0].attributeName = %v, want %v", attr1["attributeName"], "Strength")
	}
	if attr1["value"] != float64(10) {
		t.Errorf("attributes[0].value = %v, want %v", attr1["value"], 10)
	}
}

func TestCharacterAttributeHandler_GetByCharacterID_CharacterNotFound(t *testing.T) {
	mockCharRepo := &mockCharacterRepositoryForAttributeTests{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{}

	router := setupTestRouterForAttributes(mockCharRepo, mockAttrRepo)

	req, _ := http.NewRequest("GET", "/api/v1/character/non-existent/attribute", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusForbidden)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "forbidden" {
		t.Errorf("error = %v, want %v", response["error"], "forbidden")
	}
}

func TestCharacterAttributeHandler_GetByCharacterID_EmptyAttributes(t *testing.T) {
	mockChar := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		1,
		0,
		0,
		"test-user-123", // Must match JWT mock userID
		time.Now(),
	)

	mockCharRepo := &mockCharacterRepositoryForAttributeTests{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			if id == "char-123" && userID == "test-user-123" {
				return mockChar, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{
		findByCharacterIDFunc: func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
			return []*entity.CharacterAttribute{}, nil // Empty list
		},
	}

	router := setupTestRouterForAttributes(mockCharRepo, mockAttrRepo)

	req, _ := http.NewRequest("GET", "/api/v1/character/char-123/attribute", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	attributes := response["attributes"].([]interface{})
	if len(attributes) != 0 {
		t.Errorf("len(attributes) = %v, want %v", len(attributes), 0)
	}
}

func TestCharacterAttributeHandler_GetByCharacterID_MissingAuthorization(t *testing.T) {
	mockCharRepo := &mockCharacterRepositoryForAttributeTests{}
	mockAttrRepo := &mockCharacterAttributeRepository{}

	router := setupTestRouterForAttributes(mockCharRepo, mockAttrRepo)

	req, _ := http.NewRequest("GET", "/api/v1/character/char-123/attribute", nil)
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestCharacterAttributeHandler_GetByCharacterID_InvalidToken(t *testing.T) {
	mockCharRepo := &mockCharacterRepositoryForAttributeTests{}
	mockAttrRepo := &mockCharacterAttributeRepository{}

	router := setupTestRouterForAttributes(mockCharRepo, mockAttrRepo)

	req, _ := http.NewRequest("GET", "/api/v1/character/char-123/attribute", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestCharacterAttributeHandler_GetByCharacterID_RepositoryError(t *testing.T) {
	mockChar := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		1,
		0,
		0,
		"test-user-123", // Same user as JWT token
		time.Now(),
	)

	mockCharRepo := &mockCharacterRepositoryForAttributeTests{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			if id == "char-123" && userID == "test-user-123" {
				return mockChar, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{
		findByCharacterIDFunc: func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
			return nil, errors.New("database connection failed")
		},
	}

	router := setupTestRouterForAttributes(mockCharRepo, mockAttrRepo)

	req, _ := http.NewRequest("GET", "/api/v1/character/char-123/attribute", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "failed_to_fetch_attributes" {
		t.Errorf("error = %v, want %v", response["error"], "failed_to_fetch_attributes")
	}
}

func TestCharacterAttributeHandler_GetByCharacterID_Forbidden(t *testing.T) {
	// Character belongs to user-999, but token is from test-user-123
	mockChar := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		1,
		0,
		0,
		"user-999", // Different user than token
		time.Now(),
	)

	mockCharRepo := &mockCharacterRepositoryForAttributeTests{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			// Character belongs to user-999, token is from test-user-123
			// So FindByIDAndUserID will fail to match
			if id == "char-123" && userID == "user-999" {
				return mockChar, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{}

	router := setupTestRouterForAttributes(mockCharRepo, mockAttrRepo)

	req, _ := http.NewRequest("GET", "/api/v1/character/char-123/attribute", nil)
	req.Header.Set("Authorization", "Bearer valid_token") // Token belongs to test-user-123

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusForbidden)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "forbidden" {
		t.Errorf("error = %v, want %v", response["error"], "forbidden")
	}
}
