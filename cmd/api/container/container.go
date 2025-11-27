package container

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/config"
)

// Container é o container principal que orquestra todas as camadas
// Segue Clean Architecture: Infrastructure <- Application <- Delivery
type Container struct {
	Config         *config.Config
	Infrastructure *Infrastructure
	Application    *Application
	Delivery       *Delivery
}

// New cria e inicializa o container completo com todas as dependências
// Esta função:
// 1. Inicializa a camada de Infraestrutura (DB, Repos, Services)
// 2. Inicializa a camada de Aplicação (Use Cases)
// 3. Inicializa a camada de Entrega (Handlers, Router)
// 4. Verifica a saúde do sistema (DB health check)
func New(cfg *config.Config) (*Container, error) {
	// 1. Inicializar camada de Infraestrutura
	infra, err := NewInfrastructure(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	// Verificar saúde do banco de dados
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := infra.DB.Health(ctx); err != nil {
		infra.Close()
		return nil, fmt.Errorf("database health check failed: %w", err)
	}

	// 2. Inicializar camada de Aplicação (depende da Infraestrutura)
	app := NewApplication(infra)

	// 3. Inicializar camada de Entrega (depende da Aplicação e Infraestrutura)
	delivery := NewDelivery(app, infra, cfg)

	container := &Container{
		Config:         cfg,
		Infrastructure: infra,
		Application:    app,
		Delivery:       delivery,
	}

	return container, nil
}

// GetEngine retorna o Gin Engine configurado e pronto para uso
func (c *Container) GetEngine() *gin.Engine {
	return c.Delivery.Engine
}

// Close libera todos os recursos e encerra conexões
// Deve ser chamado quando a aplicação for encerrada (defer container.Close())
func (c *Container) Close() error {
	if c.Infrastructure != nil {
		return c.Infrastructure.Close()
	}
	return nil
}
