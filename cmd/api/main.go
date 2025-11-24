package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/igor/chronotask-api/config"
	"github.com/igor/chronotask-api/internal/application/usecase"
	deliveryHttp "github.com/igor/chronotask-api/internal/delivery/http"
	"github.com/igor/chronotask-api/internal/infrastructure/persistence"
	"github.com/igor/chronotask-api/internal/infrastructure/service"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Dependency Injection - Manual wiring following Clean Architecture
	// 1. Initialize infrastructure layer (database, external services)
	db, err := persistence.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database health
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Health(ctx); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	// Initialize infrastructure services
	hasherService := service.NewBcryptHasher(bcrypt.DefaultCost)

	// Initialize repositories
	userRepo := persistence.NewPostgresUserRepository(db)

	// 2. Initialize application layer (use cases)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo, hasherService)

	// 3. Initialize delivery layer (HTTP handlers)
	healthHandler := deliveryHttp.NewHealthHandler()
	userHandler := deliveryHttp.NewUserHandler(createUserUseCase)

	// Initialize router with handlers
	router := deliveryHttp.NewRouter(healthHandler, userHandler)

	// Setup routes
	engine := router.SetupRoutes()

	// Start server
	log.Printf("Starting ChronoTask API server on port %s", cfg.Server.Port)

	// Graceful shutdown
	go func() {
		if err := engine.Run(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
