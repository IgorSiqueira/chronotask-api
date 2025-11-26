package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/internal/application/port"
	"github.com/igor/chronotask-api/internal/application/usecase"
	deliveryHttp "github.com/igor/chronotask-api/internal/delivery/http"
	"github.com/igor/chronotask-api/internal/delivery/http/middleware"
	"github.com/igor/chronotask-api/internal/domain/entity"
)

// Mock CharacterRepository for E2E tests
type mockCharacterRepository struct {
	existsByUserIDFunc  func(ctx context.Context, userID string) (bool, error)
	createFunc          func(ctx context.Context, character *entity.Character) error
	findAllByUserIDFunc func(ctx context.Context, userID string) ([]*entity.Character, error)
}

func (m *mockCharacterRepository) Create(ctx context.Context, character *entity.Character) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, character)
	}
	return nil
}

func (m *mockCharacterRepository) FindByID(ctx context.Context, id string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepository) FindByIDAndUserID(ctx context.Context, id string, userID string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepository) FindByUserID(ctx context.Context, userID string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepository) Update(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepository) FindAllByUserID(ctx context.Context, userID string) ([]*entity.Character, error) {
	if m.findAllByUserIDFunc != nil {
		return m.findAllByUserIDFunc(ctx, userID)
	}
	return []*entity.Character{}, nil
}

func (m *mockCharacterRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	if m.existsByUserIDFunc != nil {
		return m.existsByUserIDFunc(ctx, userID)
	}
	return false, nil
}

// Mock JWT Service for testing
type mockJWTService struct{}

func (m *mockJWTService) GenerateAccessToken(userID, email string) (string, error) {
	return "mock_access_token", nil
}

func (m *mockJWTService) GenerateRefreshToken(userID, email string) (string, error) {
	return "mock_refresh_token", nil
}

func (m *mockJWTService) ValidateToken(token string) (*port.TokenClaims, error) {
	if token == "valid_token" {
		return &port.TokenClaims{
			UserID: "test-user-123",
			Email:  "test@example.com",
		}, nil
	}
	return nil, errors.New("invalid token")
}

func (m *mockJWTService) RefreshAccessToken(refreshToken string) (string, error) {
	return "new_access_token", nil
}

// setupTestRouter creates a test router with mocked dependencies
func setupTestRouter(charRepo *mockCharacterRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Create use cases
	createCharacterUseCase := usecase.NewCreateCharacterUseCase(charRepo)
	getUserCharactersUseCase := usecase.NewGetUserCharactersUseCase(charRepo)

	// Create handler
	characterHandler := deliveryHttp.NewCharacterHandler(
		createCharacterUseCase,
		getUserCharactersUseCase,
	)

	// Create auth middleware with mock JWT service
	authMiddleware := middleware.NewAuthMiddleware(&mockJWTService{})

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		authenticated := v1.Group("")
		authenticated.Use(authMiddleware.RequireAuth())
		{
			authenticated.GET("/user/character", characterHandler.GetList)
			authenticated.POST("/character", characterHandler.Create)
		}
	}

	return router
}

func TestCharacterHandler_Create_Success(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, nil // User doesn't have a character
		},
		createFunc: func(ctx context.Context, character *entity.Character) error {
			return nil // Success
		},
	}

	router := setupTestRouter(mockRepo)

	// Prepare request
	requestBody := map[string]interface{}{
		"name": "Warrior King",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/character", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_token")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusCreated {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusCreated)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["name"] != "Warrior King" {
		t.Errorf("response name = %v, want %v", response["name"], "Warrior King")
	}

	if response["level"] != float64(1) {
		t.Errorf("response level = %v, want %v", response["level"], 1)
	}

	if response["currentXp"] != float64(0) {
		t.Errorf("response currentXp = %v, want %v", response["currentXp"], 0)
	}

	if response["userId"] != "test-user-123" {
		t.Errorf("response userId = %v, want %v", response["userId"], "test-user-123")
	}

	if response["id"] == nil || response["id"] == "" {
		t.Error("response id should not be empty")
	}
}

func TestCharacterHandler_Create_MissingAuthorization(t *testing.T) {
	mockRepo := &mockCharacterRepository{}
	router := setupTestRouter(mockRepo)

	requestBody := map[string]interface{}{
		"name": "Warrior King",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/character", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestCharacterHandler_Create_InvalidToken(t *testing.T) {
	mockRepo := &mockCharacterRepository{}
	router := setupTestRouter(mockRepo)

	requestBody := map[string]interface{}{
		"name": "Warrior King",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/character", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestCharacterHandler_Create_InvalidRequestBody(t *testing.T) {
	mockRepo := &mockCharacterRepository{}
	router := setupTestRouter(mockRepo)

	tests := []struct {
		name        string
		requestBody map[string]interface{}
	}{
		{"missing name", map[string]interface{}{}},
		{"empty name", map[string]interface{}{"name": ""}},
		{"name too short", map[string]interface{}{"name": "A"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)

			req, _ := http.NewRequest("POST", "/api/v1/character", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Status code = %v, want %v", w.Code, http.StatusBadRequest)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if response["error"] != "invalid_request" {
				t.Errorf("error = %v, want %v", response["error"], "invalid_request")
			}
		})
	}
}

func TestCharacterHandler_Create_UserAlreadyHasCharacter(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return true, nil // User already has a character
		},
	}

	router := setupTestRouter(mockRepo)

	requestBody := map[string]interface{}{
		"name": "Warrior King",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/character", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusConflict)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "character_already_exists" {
		t.Errorf("error = %v, want %v", response["error"], "character_already_exists")
	}
}

func TestCharacterHandler_Create_RepositoryError(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		createFunc: func(ctx context.Context, character *entity.Character) error {
			return errors.New("database error")
		},
	}

	router := setupTestRouter(mockRepo)

	requestBody := map[string]interface{}{
		"name": "Warrior King",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/v1/character", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusUnprocessableEntity)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "character_creation_failed" {
		t.Errorf("error = %v, want %v", response["error"], "character_creation_failed")
	}
}

func TestCharacterHandler_Create_InvalidJSON(t *testing.T) {
	mockRepo := &mockCharacterRepository{}
	router := setupTestRouter(mockRepo)

	req, _ := http.NewRequest("POST", "/api/v1/character", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestCharacterHandler_GetList_Success_MultipleCharacters(t *testing.T) {
	char1 := entity.ReconstituteCharacter(
		"char-1",
		"Warrior",
		5,
		50,
		500,
		"test-user-123",
		time.Now(),
	)
	char2 := entity.ReconstituteCharacter(
		"char-2",
		"Mage",
		3,
		30,
		300,
		"test-user-123",
		time.Now(),
	)

	mockRepo := &mockCharacterRepository{
		findAllByUserIDFunc: func(ctx context.Context, userID string) ([]*entity.Character, error) {
			if userID == "test-user-123" {
				return []*entity.Character{char1, char2}, nil
			}
			return []*entity.Character{}, nil
		},
	}

	router := setupTestRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/api/v1/user/character", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	characters, ok := response["characters"].([]interface{})
	if !ok {
		t.Fatal("response should contain 'characters' array")
	}

	if len(characters) != 2 {
		t.Errorf("len(characters) = %v, want %v", len(characters), 2)
	}

	// Verify first character
	char1Map := characters[0].(map[string]interface{})
	if char1Map["id"] != "char-1" {
		t.Errorf("character[0].id = %v, want %v", char1Map["id"], "char-1")
	}
	if char1Map["name"] != "Warrior" {
		t.Errorf("character[0].name = %v, want %v", char1Map["name"], "Warrior")
	}
	if char1Map["level"] != float64(5) {
		t.Errorf("character[0].level = %v, want %v", char1Map["level"], 5)
	}

	// Verify second character
	char2Map := characters[1].(map[string]interface{})
	if char2Map["id"] != "char-2" {
		t.Errorf("character[1].id = %v, want %v", char2Map["id"], "char-2")
	}
	if char2Map["name"] != "Mage" {
		t.Errorf("character[1].name = %v, want %v", char2Map["name"], "Mage")
	}
}

func TestCharacterHandler_GetList_Success_EmptyList(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		findAllByUserIDFunc: func(ctx context.Context, userID string) ([]*entity.Character, error) {
			return []*entity.Character{}, nil
		},
	}

	router := setupTestRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/api/v1/user/character", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	characters, ok := response["characters"].([]interface{})
	if !ok {
		t.Fatal("response should contain 'characters' array")
	}

	if len(characters) != 0 {
		t.Errorf("len(characters) = %v, want %v", len(characters), 0)
	}
}

func TestCharacterHandler_GetList_MissingAuthorization(t *testing.T) {
	mockRepo := &mockCharacterRepository{}
	router := setupTestRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/api/v1/user/character", nil)
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestCharacterHandler_GetList_InvalidToken(t *testing.T) {
	mockRepo := &mockCharacterRepository{}
	router := setupTestRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/api/v1/user/character", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestCharacterHandler_GetList_RepositoryError(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		findAllByUserIDFunc: func(ctx context.Context, userID string) ([]*entity.Character, error) {
			return nil, errors.New("database connection failed")
		},
	}

	router := setupTestRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/api/v1/user/character", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "failed_to_fetch_characters" {
		t.Errorf("error = %v, want %v", response["error"], "failed_to_fetch_characters")
	}
}
