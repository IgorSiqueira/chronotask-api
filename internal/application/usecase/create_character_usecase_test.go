package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/igor/chronotask-api/internal/application/usecase"
	"github.com/igor/chronotask-api/internal/domain/entity"
)

// Mock CharacterRepository
type mockCharacterRepository struct {
	createFunc         func(ctx context.Context, character *entity.Character) error
	existsByUserIDFunc func(ctx context.Context, userID string) (bool, error)
}

// Mock CharacterAttributeRepository
type mockCharacterAttributeRepository struct {
	createFunc                       func(ctx context.Context, attribute *entity.CharacterAttribute) error
	findByIDFunc                     func(ctx context.Context, id int) (*entity.CharacterAttribute, error)
	findByCharacterIDFunc            func(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error)
	findByCharacterIDAndNameFunc     func(ctx context.Context, characterID string, attributeName string) (*entity.CharacterAttribute, error)
	updateFunc                       func(ctx context.Context, attribute *entity.CharacterAttribute) error
	deleteFunc                       func(ctx context.Context, id int) error
	existsByCharacterIDAndNameFunc   func(ctx context.Context, characterID string, attributeName string) (bool, error)
}

func (m *mockCharacterAttributeRepository) Create(ctx context.Context, attribute *entity.CharacterAttribute) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, attribute)
	}
	return nil
}

func (m *mockCharacterAttributeRepository) FindByID(ctx context.Context, id int) (*entity.CharacterAttribute, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) FindByCharacterID(ctx context.Context, characterID string) ([]*entity.CharacterAttribute, error) {
	if m.findByCharacterIDFunc != nil {
		return m.findByCharacterIDFunc(ctx, characterID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) FindByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (*entity.CharacterAttribute, error) {
	if m.findByCharacterIDAndNameFunc != nil {
		return m.findByCharacterIDAndNameFunc(ctx, characterID, attributeName)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) Update(ctx context.Context, attribute *entity.CharacterAttribute) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, attribute)
	}
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func (m *mockCharacterAttributeRepository) ExistsByCharacterIDAndName(ctx context.Context, characterID string, attributeName string) (bool, error) {
	if m.existsByCharacterIDAndNameFunc != nil {
		return m.existsByCharacterIDAndNameFunc(ctx, characterID, attributeName)
	}
	return false, nil
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

func (m *mockCharacterRepository) FindAllByUserID(ctx context.Context, userID string) ([]*entity.Character, error) {
	return []*entity.Character{}, errors.New("not implemented")
}

func (m *mockCharacterRepository) Update(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	if m.existsByUserIDFunc != nil {
		return m.existsByUserIDFunc(ctx, userID)
	}
	return false, nil
}

func TestCreateCharacterUseCase_Execute_Success(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, nil // User doesn't have a character yet
		},
		createFunc: func(ctx context.Context, character *entity.Character) error {
			return nil // Success
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{
		createFunc: func(ctx context.Context, attribute *entity.CharacterAttribute) error {
			return nil // Success
		},
	}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	input := usecase.CreateCharacterInput{
		Name:   "Warrior King",
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() output = nil, want non-nil")
	}

	if output.Name != "Warrior King" {
		t.Errorf("output.Name = %v, want %v", output.Name, "Warrior King")
	}

	if output.UserID != "user-123" {
		t.Errorf("output.UserID = %v, want %v", output.UserID, "user-123")
	}

	if output.Level != 1 {
		t.Errorf("output.Level = %v, want %v", output.Level, 1)
	}

	if output.CurrentXp != 0 {
		t.Errorf("output.CurrentXp = %v, want %v", output.CurrentXp, 0)
	}

	if output.TotalXp != 0 {
		t.Errorf("output.TotalXp = %v, want %v", output.TotalXp, 0)
	}

	if output.ID == "" {
		t.Error("output.ID should not be empty")
	}

	if output.CreatedAt == "" {
		t.Error("output.CreatedAt should not be empty")
	}
}

func TestCreateCharacterUseCase_Execute_UserAlreadyHasCharacter(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return true, nil // User already has a character
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	input := usecase.CreateCharacterInput{
		Name:   "Warrior King",
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error for user already has character")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}

	expectedError := "user already has a character"
	if err.Error() != expectedError {
		t.Errorf("error message = %v, want %v", err.Error(), expectedError)
	}
}

func TestCreateCharacterUseCase_Execute_ExistsByUserIDError(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, errors.New("database error")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	input := usecase.CreateCharacterInput{
		Name:   "Warrior King",
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}
}

func TestCreateCharacterUseCase_Execute_InvalidCharacterName(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	tests := []struct {
		name          string
		characterName string
	}{
		{"empty", ""},
		{"too short", "A"},
		{"only spaces", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := usecase.CreateCharacterInput{
				Name:   tt.characterName,
				UserID: "user-123",
			}

			output, err := useCase.Execute(context.Background(), input)

			if err == nil {
				t.Error("Execute() error = nil, want error for invalid character name")
			}

			if output != nil {
				t.Errorf("Execute() output = %v, want nil", output)
			}
		})
	}
}

func TestCreateCharacterUseCase_Execute_CreateRepositoryError(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		createFunc: func(ctx context.Context, character *entity.Character) error {
			return errors.New("database connection failed")
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	input := usecase.CreateCharacterInput{
		Name:   "Warrior King",
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}
}

func TestCreateCharacterUseCase_Execute_ContextCancellation(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, ctx.Err() // Return context error
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	input := usecase.CreateCharacterInput{
		Name:   "Warrior King",
		UserID: "user-123",
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

func TestCreateCharacterUseCase_Execute_CreatesBaseAttributes(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		createFunc: func(ctx context.Context, character *entity.Character) error {
			return nil
		},
	}

	// Track how many attributes are created
	attributesCreated := 0
	expectedAttributes := map[string]int{
		"Força":         5,
		"Constituição":  5,
		"Vontade":       5,
		"Sabedoria":     5,
		"Inteligência":  5,
		"Carisma":       5,
		"Destreza":      5,
	}

	mockAttrRepo := &mockCharacterAttributeRepository{
		createFunc: func(ctx context.Context, attribute *entity.CharacterAttribute) error {
			attributesCreated++

			// Validate that the attribute is one of the expected base attributes
			expectedValue, exists := expectedAttributes[attribute.AttributeName()]
			if !exists {
				t.Errorf("Unexpected attribute created: %s", attribute.AttributeName())
			}

			if attribute.Value() != expectedValue {
				t.Errorf("Attribute %s has value %d, want %d",
					attribute.AttributeName(), attribute.Value(), expectedValue)
			}

			return nil
		},
	}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	input := usecase.CreateCharacterInput{
		Name:   "Warrior King",
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() output = nil, want non-nil")
	}

	// Verify that all 7 base attributes were created
	expectedCount := 7
	if attributesCreated != expectedCount {
		t.Errorf("Created %d attributes, want %d", attributesCreated, expectedCount)
	}
}

func TestCreateCharacterUseCase_Execute_AttributeCreationError(t *testing.T) {
	mockRepo := &mockCharacterRepository{
		existsByUserIDFunc: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		createFunc: func(ctx context.Context, character *entity.Character) error {
			return nil
		},
	}

	mockAttrRepo := &mockCharacterAttributeRepository{
		createFunc: func(ctx context.Context, attribute *entity.CharacterAttribute) error {
			return errors.New("failed to create attribute")
		},
	}

	useCase := usecase.NewCreateCharacterUseCase(mockRepo, mockAttrRepo)

	input := usecase.CreateCharacterInput{
		Name:   "Warrior King",
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("Execute() error = nil, want error for attribute creation failure")
	}

	if output != nil {
		t.Errorf("Execute() output = %v, want nil", output)
	}
}
