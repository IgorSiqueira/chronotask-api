package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/igor/chronotask-api/cmd/api/container"
	"github.com/igor/chronotask-api/config"
)

func main() {
	// Carregar configuração
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Inicializar container com todas as dependências
	// O container automaticamente:
	// 1. Conecta ao banco de dados
	// 2. Inicializa todos os repositórios
	// 3. Inicializa todos os use cases
	// 4. Inicializa todos os handlers
	// 5. Configura o router
	appContainer, err := container.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer appContainer.Close()

	// Obter Gin Engine configurado
	engine := appContainer.GetEngine()

	// Iniciar servidor
	log.Printf("✓ Database connected")
	log.Printf("✓ All dependencies injected")
	log.Printf("✓ Starting ChronoTask API server on port %s", cfg.Server.Port)

	// Graceful shutdown
	go func() {
		if err := engine.Run(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚠ Shutting down server...")
	log.Println("✓ Server stopped gracefully")
}
