package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/igor/chronotask-api/internal/application/usecase"
	"github.com/igor/chronotask-api/internal/domain/entity"
)

// Mock CharacterRepository for list tests
type mockCharacterRepositoryForList struct {
	findAllByUserIDFunc func(ctx context.Context, userID string) ([]*entity.Character, error)
}

func (m *mockCharacterRepositoryForList) Create(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForList) FindByID(ctx context.Context, id string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepositoryForList) FindByIDAndUserID(ctx context.Context, id string, userID string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepositoryForList) FindByUserID(ctx context.Context, userID string) (*entity.Character, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCharacterRepositoryForList) FindAllByUserID(ctx context.Context, userID string) ([]*entity.Character, error) {
	if m.findAllByUserIDFunc != nil {
		return m.findAllByUserIDFunc(ctx, userID)
	}
	return []*entity.Character{}, nil
}

func (m *mockCharacterRepositoryForList) Update(ctx context.Context, character *entity.Character) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForList) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockCharacterRepositoryForList) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	return false, errors.New("not implemented")
}

func TestGetUserCharactersUseCase_Execute_SingleCharacter(t *testing.T) {
	mockCharacter := entity.ReconstituteCharacter(
		"char-123",
		"Warrior King",
		5,
		50,
		500,
		"user-123",
		time.Now(),
	)

	mockRepo := &mockCharacterRepositoryForList{
		findAllByUserIDFunc: func(ctx context.Context, userID string) ([]*entity.Character, error) {
			if userID == "user-123" {
				return []*entity.Character{mockCharacter}, nil
			}
			return []*entity.Character{}, nil
		},
	}

	useCase := usecase.NewGetUserCharactersUseCase(mockRepo)

	input := usecase.GetUserCharactersInput{
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() output = nil, want non-nil")
	}

	if len(output.Characters) != 1 {
		t.Errorf("len(output.Characters) = %v, want %v", len(output.Characters), 1)
	}

	if output.Characters[0].ID != "char-123" {
		t.Errorf("output.Characters[0].ID = %v, want %v", output.Characters[0].ID, "char-123")
	}

	if output.Characters[0].Name != "Warrior King" {
		t.Errorf("output.Characters[0].Name = %v, want %v", output.Characters[0].Name, "Warrior King")
	}

	if output.Characters[0].Level != 5 {
		t.Errorf("output.Characters[0].Level = %v, want %v", output.Characters[0].Level, 5)
	}
}

func TestGetUserCharactersUseCase_Execute_EmptyList(t *testing.T) {
	mockRepo := &mockCharacterRepositoryForList{
		findAllByUserIDFunc: func(ctx context.Context, userID string) ([]*entity.Character, error) {
			return []*entity.Character{}, nil // Empty list
		},
	}

	useCase := usecase.NewGetUserCharactersUseCase(mockRepo)

	input := usecase.GetUserCharactersInput{
		UserID: "user-123",
	}

	output, err := useCase.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if output == nil {
		t.Fatal("Execute() output = nil, want non-nil")
	}

	if len(output.Characters) != 0 {
		t.Errorf("len(output.Characters) = %v, want %v", len(output.Characters), 0)
	}
}

func TestGetUserCharactersUseCase_Execute_RepositoryError(t *testing.T) {
	mockRepo := &mockCharacterRepositoryForList{
		findAllByUserIDFunc: func(ctx context.Context, userID string) ([]*entity.Character, error) {
			return nil, errors.New("database connection failed")
		},
	}

	useCase := usecase.NewGetUserCharactersUseCase(mockRepo)

	input := usecase.GetUserCharactersInput{
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
