package container

import (
	"fmt"

	"github.com/igor/chronotask-api/config"
	"github.com/igor/chronotask-api/internal/application/port"
	"github.com/igor/chronotask-api/internal/domain/repository"
	"github.com/igor/chronotask-api/internal/infrastructure/persistence"
	"github.com/igor/chronotask-api/internal/infrastructure/service"
	"golang.org/x/crypto/bcrypt"
)

// Infrastructure contém todas as dependências da camada de infraestrutura
// Inclui: database, repositories, e serviços externos
type Infrastructure struct {
	// Database
	DB *persistence.PostgresDB

	// External Services
	HasherService port.HasherService
	JWTService    port.JWTService

	// Repositories
	// Adicione novos repositórios aqui conforme necessário
	UserRepository repository.UserRepository
	// HabitRepository repository.HabitRepository     // Exemplo futuro
	// CharacterRepository repository.CharacterRepository // Exemplo futuro
}

// NewInfrastructure inicializa toda a camada de infraestrutura
// Para adicionar um novo repositório:
// 1. Adicione o campo no struct acima
// 2. Inicialize aqui (1 linha): infra.NovoRepo = persistence.NewNovoRepo(db)
func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
	// Conectar ao banco de dados
	db, err := persistence.NewPostgresDB(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Inicializar serviços externos
	hasherService := service.NewBcryptHasher(bcrypt.DefaultCost)

	jwtService, err := service.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWT service: %w", err)
	}

	// Inicializar repositórios
	userRepo := persistence.NewPostgresUserRepository(db)

	// Futuro: adicionar novos repositórios aqui
	// habitRepo := persistence.NewPostgresHabitRepository(db)
	// characterRepo := persistence.NewPostgresCharacterRepository(db)

	infra := &Infrastructure{
		DB:             db,
		HasherService:  hasherService,
		JWTService:     jwtService,
		UserRepository: userRepo,
		// HabitRepository: habitRepo,
		// CharacterRepository: characterRepo,
	}

	return infra, nil
}

// Close encerra conexões e libera recursos
func (i *Infrastructure) Close() error {
	if i.DB != nil {
		i.DB.Close()
	}
	return nil
}
