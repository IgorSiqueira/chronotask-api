package container

import (
	"github.com/gin-gonic/gin"
	deliveryHttp "github.com/igor/chronotask-api/internal/delivery/http"
)

// Delivery contém todas as dependências da camada de entrega
// Inclui: Handlers, Router, Engine
type Delivery struct {
	// Handlers
	HealthHandler *deliveryHttp.HealthHandler
	UserHandler   *deliveryHttp.UserHandler
	// HabitHandler    *deliveryHttp.HabitHandler    // Exemplo futuro
	// CharacterHandler *deliveryHttp.CharacterHandler // Exemplo futuro

	// Router e Engine
	Router *deliveryHttp.Router
	Engine *gin.Engine
}

// NewDelivery inicializa toda a camada de entrega
// Para adicionar um novo handler:
// 1. Adicione o campo no struct acima
// 2. Inicialize aqui (1 linha): delivery.NovoHandler = deliveryHttp.NewNovoHandler(app.NovoUseCase)
// 3. Adicione no NewRouter (linha 38)
func NewDelivery(app *Application) *Delivery {
	// Inicializar handlers
	healthHandler := deliveryHttp.NewHealthHandler()
	userHandler := deliveryHttp.NewUserHandler(app.CreateUserUseCase)

	// Futuro: adicionar novos handlers aqui
	// habitHandler := deliveryHttp.NewHabitHandler(
	//     app.CreateHabitUseCase,
	//     app.UpdateHabitUseCase,
	//     app.ListHabitsUseCase,
	// )
	// characterHandler := deliveryHttp.NewCharacterHandler(
	//     app.CreateCharacterUseCase,
	//     app.LevelUpCharacterUseCase,
	// )

	// Inicializar router
	router := deliveryHttp.NewRouter(
		healthHandler,
		userHandler,
		// habitHandler,    // Adicionar quando criar
		// characterHandler, // Adicionar quando criar
	)

	// Setup routes
	engine := router.SetupRoutes()

	delivery := &Delivery{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
		// HabitHandler: habitHandler,
		// CharacterHandler: characterHandler,
		Router: router,
		Engine: engine,
	}

	return delivery
}
