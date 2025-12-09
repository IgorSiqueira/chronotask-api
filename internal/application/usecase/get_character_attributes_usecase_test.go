package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/igor/chronotask-api/internal/application/usecase"
	"github.com/igor/chronotask-api/internal/domain/entity"
)

// Mock CharacterAttributeRepository for GetCharacterAttributes tests
type mockCharacterAttributeRepositoryGet struct {
	findByCharacterIDFunc func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error)
}

func (m *mockCharacterAttributeRepositoryGet) Create(ctx context.Context, attribute *entity.CharacterAttribute) error {
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepositoryGet) FindByID(ctx context.Context, id int) (*entity.CharacterAttribute, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterAttributeRepositoryGet) FindByCharacterID(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
	if m.findByCharacterIDFunc != nil {
		return m.findByCharacterIDFunc(ctx, characterID)
	}
	return []*entity.CharacterAttribute{}, nil
}

func (m *mockCharacterAttributeRepositoryGet) FindByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (*entity.CharacterAttribute, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterAttributeRepositoryGet) Update(ctx context.Context, attribute *entity.CharacterAttribute) error {
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepositoryGet) Delete(ctx context.Context, id int) error {
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepositoryGet) ExistsByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (bool, error) {
	return false, errors.New("not implemented")
}

// Mock CharacterRepository for this test
type mockCharacterRepositoryForAttributes struct {
	findByIDFunc           func(ctx context.Context, id string) (*entity.Character, error)
	findByIDAndUserIDFunc  func(ctx context.Context, id string, userID string) (*entity.Character, error)
}

func (m *mockCharacterRepositoryForAttributes) Create(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributes) FindByID(ctx context.Context, id string) (*entity.Character, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockCharacterRepositoryForAttributes) FindByIDAndUserID(ctx context.Context, id string, userID string) (*entity.Character, error) {
	if m.findByIDAndUserIDFunc != nil {
		return m.findByIDAndUserIDFunc(ctx, id, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockCharacterRepositoryForAttributes) FindByUserID(ctx context.Context, userID string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributes) FindAllByUserID(ctx context.Context, userID string) ([]*entity.Character, error) {
	return []*entity.Character{}, errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributes) Update(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributes) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForAttributes) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	return false, errors.New("not implemented")
}

func TestGetCharacterAttributesUseCase_Execute_Success(t *testing.T) {
	mockCharacter := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		5,
		50,
		500,
		"user-123",
		time.Now(),
	)

	mockAttributes := []*entity.CharacterAttribute{
		entity.ReconstituteCharacterAttribute(1, "Força", 10, "char-123", time.Now()),
		entity.ReconstituteCharacterAttribute(2, "Destreza", 15, "char-123", time.Now()),
		entity.ReconstituteCharacterAttribute(3, "Inteligência", 20, "char-123", time.Now()),
	}

	mockCharRepo := &mockCharacterRepositoryForAttributes{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			if id == "char-123" && userID == "user-123" {
				return mockCharacter, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepositoryGet{
		findByCharacterIDFunc: func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
			if characterID == "char-123" {
				return mockAttributes, nil
			}
			return []*entity.CharacterAttribute{}, nil
		},
	}

	useCase := usecase.NewGetCharacterAttributesUseCase(mockCharRepo, mockAttrRepo)

	input := usecase.GetCharacterAttributesInput{
		CharacterID: "char-123",
		UserID:      "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() output = nil, want non-nil")
	}

	if output.CharacterID != "char-123" {
		t.Errorf("output.CharacterID = %v, want %v", output.CharacterID, "char-123")
	}

	if len(output.Attributes) != 3 {
		t.Errorf("len(output.Attributes) = %v, want %v", len(output.Attributes), 3)
	}

	// Verify first attribute
	if output.Attributes[0].AttributeName != "Força" {
		t.Errorf("output.Attributes[0].AttributeName = %v, want %v", output.Attributes[0].AttributeName, "Força")
	}

	if output.Attributes[0].Value != 10 {
		t.Errorf("output.Attributes[0].Value = %v, want %v", output.Attributes[0].Value, 10)
	}
}

func TestGetCharacterAttributesUseCase_Execute_CharacterNotFound(t *testing.T) {
	mockCharRepo := &mockCharacterRepositoryForAttributes{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepositoryGet{}

	useCase := usecase.NewGetCharacterAttributesUseCase(mockCharRepo, mockAttrRepo)

	input := usecase.GetCharacterAttributesInput{
		CharacterID: "non-existent",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error for character not found")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}
}

func TestGetCharacterAttributesUseCase_Execute_EmptyAttributes(t *testing.T) {
	mockCharacter := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		1,
		0,
		0,
		"user-123",
		time.Now(),
	)

	mockCharRepo := &mockCharacterRepositoryForAttributes{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			if id == "char-123" && userID == "user-123" {
				return mockCharacter, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepositoryGet{
		findByCharacterIDFunc: func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
			return []*entity.CharacterAttribute{}, nil // Empty list
		},
	}

	useCase := usecase.NewGetCharacterAttributesUseCase(mockCharRepo, mockAttrRepo)

	input := usecase.GetCharacterAttributesInput{
		CharacterID: "char-123",
		UserID:      "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() output = nil, want non-nil")
	}

	if len(output.Attributes) != 0 {
		t.Errorf("len(output.Attributes) = %v, want %v", len(output.Attributes), 0)
	}
}

func TestGetCharacterAttributesUseCase_Execute_RepositoryError(t *testing.T) {
	mockCharacter := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		1,
		0,
		0,
		"user-123",
		time.Now(),
	)

	mockCharRepo := &mockCharacterRepositoryForAttributes{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			if id == "char-123" && userID == "user-123" {
				return mockCharacter, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepositoryGet{
		findByCharacterIDFunc: func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
			return nil, errors.New("database connection failed")
		},
	}

	useCase := usecase.NewGetCharacterAttributesUseCase(mockCharRepo, mockAttrRepo)

	input := usecase.GetCharacterAttributesInput{
		CharacterID: "char-123",
		UserID:      "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}
}

func TestGetCharacterAttributesUseCase_Execute_ContextCancellation(t *testing.T) {
	mockCharRepo := &mockCharacterRepositoryForAttributes{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			return nil, ctx.Err() // Return context error
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepositoryGet{}

	useCase := usecase.NewGetCharacterAttributesUseCase(mockCharRepo, mockAttrRepo)

	input := usecase.GetCharacterAttributesInput{
		CharacterID: "char-123",
		UserID:      "user-123",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	output, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error for cancelled context")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}
}

func TestGetCharacterAttributesUseCase_Execute_UnauthorizedAccess(t *testing.T) {
	// Character belongs to user-123
	mockCharacter := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		5,
		50,
		500,
		"user-123", // Owner
		time.Now(),
	)

	mockCharRepo := &mockCharacterRepositoryForAttributes{
		findByIDAndUserIDFunc: func(ctx context.Context, id string, userID string) (*entity.Character, error) {
			// Only return character if both ID matches AND user matches
			if id == "char-123" && userID == "user-123" {
				return mockCharacter, nil
			}
			return nil, errors.New("character not found or does not belong to user")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepositoryGet{}

	useCase := usecase.NewGetCharacterAttributesUseCase(mockCharRepo, mockAttrRepo)

	// User-456 tries to access user-123's character
	input := usecase.GetCharacterAttributesInput{
		CharacterID: "char-123",
		UserID:      "user-456", // Different user trying to access
	}

	output, err := useCase.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error for unauthorized access")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}

	// Verify error message indicates authorization issue
	if !strings.Contains(err.Error(), "not authorized") && !strings.Contains(err.Error(), "does not belong") {
		t.Errorf("error message = %v, want authorization error", err.Error())
	}
}
