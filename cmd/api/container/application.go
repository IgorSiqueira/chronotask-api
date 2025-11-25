package container

import (
	"github.com/igor/chronotask-api/internal/application/usecase"
)

// Application contém todas as dependências da camada de aplicação
// Inclui: Use Cases
type Application struct {
	// User Use Cases
	CreateUserUseCase *usecase.CreateUserUseCase
	LoginUserUseCase  *usecase.LoginUserUseCase
	// UpdateUserUseCase *usecase.UpdateUserUseCase   // Exemplo futuro
	// DeleteUserUseCase *usecase.DeleteUserUseCase   // Exemplo futuro
	// GetUserUseCase    *usecase.GetUserUseCase      // Exemplo futuro

	// Habit Use Cases (futuro)
	// CreateHabitUseCase *usecase.CreateHabitUseCase
	// UpdateHabitUseCase *usecase.UpdateHabitUseCase
	// ListHabitsUseCase  *usecase.ListHabitsUseCase

	// Character Use Cases (futuro)
	// CreateCharacterUseCase *usecase.CreateCharacterUseCase
	// LevelUpCharacterUseCase *usecase.LevelUpCharacterUseCase
}

// NewApplication inicializa toda a camada de aplicação
// Para adicionar um novo use case:
// 1. Adicione o campo no struct acima
// 2. Inicialize aqui (1-3 linhas):
//    app.NovoUseCase = usecase.NewNovoUseCase(infra.Repo1, infra.Service1)
func NewApplication(infra *Infrastructure) *Application {
	app := &Application{
		// User Use Cases
		CreateUserUseCase: usecase.NewCreateUserUseCase(
			infra.UserRepository,
			infra.HasherService,
		),
		LoginUserUseCase: usecase.NewLoginUserUseCase(
			infra.UserRepository,
			infra.HasherService,
			infra.JWTService,
		),

		// Futuro: adicionar novos use cases aqui
		// UpdateUserUseCase: usecase.NewUpdateUserUseCase(infra.UserRepository),
		// DeleteUserUseCase: usecase.NewDeleteUserUseCase(infra.UserRepository),
		// GetUserUseCase: usecase.NewGetUserUseCase(infra.UserRepository),

		// Habit Use Cases
		// CreateHabitUseCase: usecase.NewCreateHabitUseCase(
		//     infra.HabitRepository,
		//     infra.UserRepository,
		// ),

		// Character Use Cases
		// CreateCharacterUseCase: usecase.NewCreateCharacterUseCase(
		//     infra.CharacterRepository,
		//     infra.UserRepository,
		// ),
	}

	return app
}
