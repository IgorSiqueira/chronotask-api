package container

import (
	"github.com/gin-gonic/gin"
	deliveryHttp "github.com/igor/chronotask-api/internal/delivery/http"
	"github.com/igor/chronotask-api/internal/delivery/http/middleware"
)

// Delivery contém todas as dependências da camada de entrega
// Inclui: Handlers, Middleware, Router, Engine
type Delivery struct {
	// Handlers
	HealthHandler *deliveryHttp.HealthHandler
	UserHandler   *deliveryHttp.UserHandler
	// HabitHandler    *deliveryHttp.HabitHandler    // Exemplo futuro
	// CharacterHandler *deliveryHttp.CharacterHandler // Exemplo futuro

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// Router e Engine
	Router *deliveryHttp.Router
	Engine *gin.Engine
}

// NewDelivery inicializa toda a camada de entrega
// Para adicionar um novo handler:
// 1. Adicione o campo no struct acima
// 2. Inicialize aqui (1 linha): delivery.NovoHandler = deliveryHttp.NewNovoHandler(app.NovoUseCase)
// 3. Adicione no NewRouter (linha 38)
func NewDelivery(app *Application, infra *Infrastructure) *Delivery {
	// Inicializar handlers
	healthHandler := deliveryHttp.NewHealthHandler()
	userHandler := deliveryHttp.NewUserHandler(
		app.CreateUserUseCase,
		app.LoginUserUseCase,
	)

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

	// Inicializar middleware
	authMiddleware := middleware.NewAuthMiddleware(infra.JWTService)

	// Inicializar router
	router := deliveryHttp.NewRouter(
		healthHandler,
		userHandler,
		authMiddleware,
		// habitHandler,    // Adicionar quando criar
		// characterHandler, // Adicionar quando criar
	)

	// Setup routes
	engine := router.SetupRoutes()

	delivery := &Delivery{
		HealthHandler:  healthHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
		// HabitHandler: habitHandler,
		// CharacterHandler: characterHandler,
		Router: router,
		Engine: engine,
	}

	return delivery
}
