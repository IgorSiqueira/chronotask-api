## 📄 ChronoTask API - AI Context (Go Edition)

### 🎯 Project Overview

**ChronoTask API** is a Habit Tracking and Character Evolution backend service. It transforms real-life activities into in-game progression (XP, Stats) to enable simulated PVP/PVE battles.

### 🛠️ Project Information

- **Project Type**: Backend API (High-Complexity Domain)
- **Primary Language**: **Go (Golang)**
- **Architecture**: **Clean Architecture (Ports and Adapters) & Domain-Driven Design (DDD)**
- **Goal**: Implement User, Habit, and Character management, focusing on the core evolution logic.
- **Frameworks/Libs**: Minimalist libraries only. **Gin** for HTTP routing. Avoid monolithic frameworks.
- **Database**: **PostgreSQL** with **pgx/v5** driver (connection pooling via pgxpool)

### 🏛️ Architectural Standards (Critical)

1.  **SOLID Principles**: Rigorous application, with **DIP (Dependency Inversion Principle)** being key for decoupling the Application Layer.
2.  **Clean Architecture (CA)**: Strict separation into `Domain`, `Application`, `Infrastructure`, and `Delivery`. Dependencies must point inwards.
3.  **Domain-Driven Design (DDD)**: Essential for modeling complex business logic. Focus on the following core concepts:

### 🧩 Domain Layer (DDD Implementation)

#### Entities
Entities are domain objects with unique identity and lifecycle. They encapsulate business rules and maintain invariants.

**Structure** (`internal/domain/entity/`):
```go
type User struct {
    id          string              // Unique identifier
    fullName    string
    email       valueobject.Email   // Value Object
    password    string              // Hashed password
    birthDate   time.Time
    acceptTerms bool
    createdAt   time.Time
    updatedAt   time.Time
}

// Constructor with validation (NewUser)
func NewUser(...) (*User, error) {
    // Validate all invariants
    // Return entity or error
}

// Getters for read-only access
func (u *User) ID() string { return u.id }

// Business methods that maintain invariants
func (u *User) UpdateFullName(name string) error {
    // Validate
    // Update
    // Track updatedAt
}

// Reconstitute for loading from persistence
func ReconstituteUser(...) *User {
    // No validation - data already validated
}
```

**Key Principles**:
- Private fields with getter methods (encapsulation)
- Constructor validates all invariants (`NewUser`)
- Business methods enforce rules
- `Reconstitute*` function for repository loading (bypasses validation)
- No framework dependencies (JSON tags, ORM annotations, etc.)

#### Value Objects
Value Objects have no identity - they're defined by their values. They're immutable and contain validation logic.

**Structure** (`internal/domain/valueobject/`):
```go
type Email struct {
    value string  // Private
}

func NewEmail(email string) (Email, error) {
    // Validate format, length, etc.
    // Return value object or error
}

func (e Email) Value() string { return e.value }
func (e Email) Equals(other Email) bool { return e.value == other.value }
```

**Key Principles**:
- Immutable (no setters)
- Validation in constructor
- Comparison by value (Equals method)
- Self-contained validation rules

#### Repository Interfaces (Ports)
Repositories are defined in the domain layer as interfaces (ports).

**Structure** (`internal/domain/repository/`):
```go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)
}
```

**Key Principles**:
- Interface defined in domain (not infrastructure)
- Accepts and returns domain entities/value objects
- Implementation is in `infrastructure/persistence/`
- Use cases depend on the interface, not implementation

#### Domain Services
For business logic that doesn't belong to a single entity.

**When to use**:
- Logic involving multiple entities
- Complex domain calculations
- Business rules that span aggregates

### ✍️ Coding Standards - Go Specific

- **Error Handling**: Use the standard Go pattern of returning `(Result, error)`.
- **Interfaces**: Define small, focused Go interfaces (`interface` keyword) as **Ports** in the *Application Layer*.
- **Domain Purity**: The Domain Layer **must be completely agnostic** of databases, HTTP, and external frameworks.

### 🌐 HTTP Routing System (Gin Framework)

#### Route Organization
- **Router Setup**: Centralized in `internal/delivery/http/router.go`
- **Handler Pattern**: Each handler is a struct with methods corresponding to HTTP endpoints
- **Dependency Injection**: Handlers are injected into the Router via constructor

#### File Structure
```
internal/delivery/http/
├── router.go           # Central router setup and route configuration
├── health_handler.go   # Health check handler
└── *_handler.go        # Future handlers (user_handler, habit_handler, etc.)
```

#### Handler Convention
```go
type ExampleHandler struct {
    useCase ExampleUseCase  // Injected from Application Layer
}

func NewExampleHandler(useCase ExampleUseCase) *ExampleHandler {
    return &ExampleHandler{useCase: useCase}
}

func (h *ExampleHandler) HandleAction(c *gin.Context) {
    // 1. Parse/validate request (DTOs)
    // 2. Call use case
    // 3. Return response
}
```

#### Router Pattern
- Router holds references to all handlers
- `SetupRoutes()` method returns configured `*gin.Engine`
- Group routes by version (`/api/v1`)
- Keep handlers thin - delegate to use cases

#### Example Route Setup
```go
router := gin.Default()
router.GET("/health", h.healthHandler.Check)

v1 := router.Group("/api/v1")
{
    v1.POST("/users", h.userHandler.Create)
    v1.GET("/habits", h.habitHandler.List)
}
```

### 💉 Dependency Injection (Factory Pattern)

#### Container Organization
O projeto usa **Factory Pattern** para gerenciar dependências sem bibliotecas externas. Todos os containers estão em `cmd/api/container/`:

**Estrutura:**
```
cmd/api/container/
├── infrastructure.go  # Camada de Infraestrutura (DB, Repos, Services)
├── application.go     # Camada de Aplicação (Use Cases)
├── delivery.go        # Camada de Entrega (Handlers, Router)
└── container.go       # Container Principal (orquestra tudo)
```

#### Como Funciona

**1. Infrastructure Container** (`infrastructure.go`)
```go
type Infrastructure struct {
    DB                 *PostgresDB
    HasherService      port.HasherService
    UserRepository     repository.UserRepository
    // Adicionar novos repositórios aqui
}

func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
    db, _ := persistence.NewPostgresDB(&cfg.Database)
    return &Infrastructure{
        DB:              db,
        HasherService:   service.NewBcryptHasher(bcrypt.DefaultCost),
        UserRepository:  persistence.NewPostgresUserRepository(db),
        // Adicionar inicializações aqui (1 linha por repo)
    }, nil
}
```

**2. Application Container** (`application.go`)
```go
type Application struct {
    CreateUserUseCase *usecase.CreateUserUseCase
    // Adicionar novos use cases aqui
}

func NewApplication(infra *Infrastructure) *Application {
    return &Application{
        CreateUserUseCase: usecase.NewCreateUserUseCase(
            infra.UserRepository,
            infra.HasherService,
        ),
        // Adicionar inicializações aqui (1-3 linhas por use case)
    }
}
```

**3. Delivery Container** (`delivery.go`)
```go
type Delivery struct {
    HealthHandler *deliveryHttp.HealthHandler
    UserHandler   *deliveryHttp.UserHandler
    Router        *deliveryHttp.Router
    Engine        *gin.Engine
}

func NewDelivery(app *Application) *Delivery {
    healthHandler := deliveryHttp.NewHealthHandler()
    userHandler := deliveryHttp.NewUserHandler(app.CreateUserUseCase)
    router := deliveryHttp.NewRouter(healthHandler, userHandler)

    return &Delivery{
        HealthHandler: healthHandler,
        UserHandler:   userHandler,
        Router:        router,
        Engine:        router.SetupRoutes(),
    }
}
```

**4. Main Container** (`container.go`)
```go
type Container struct {
    Config         *config.Config
    Infrastructure *Infrastructure
    Application    *Application
    Delivery       *Delivery
}

func New(cfg *config.Config) (*Container, error) {
    infra, _ := NewInfrastructure(cfg)
    app := NewApplication(infra)
    delivery := NewDelivery(app)

    return &Container{
        Config:         cfg,
        Infrastructure: infra,
        Application:    app,
        Delivery:       delivery,
    }, nil
}
```

#### Uso no main.go

```go
func main() {
    cfg, _ := config.Load()

    // UMA LINHA para inicializar TUDO!
    container, _ := container.New(cfg)
    defer container.Close()

    engine := container.GetEngine()
    engine.Run(":8080")
}
```

#### Adicionando Nova Entidade (Exemplo: Habit)

**Passo 1:** Criar domain, repository interface, use case (como usual)

**Passo 2:** Adicionar no `infrastructure.go`:
```go
type Infrastructure struct {
    // ...
    HabitRepository repository.HabitRepository  // ← Adicionar
}

func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
    // ...
    habitRepo := persistence.NewPostgresHabitRepository(db) // ← Adicionar
    return &Infrastructure{
        // ...
        HabitRepository: habitRepo, // ← Adicionar
    }, nil
}
```

**Passo 3:** Adicionar no `application.go`:
```go
type Application struct {
    // ...
    CreateHabitUseCase *usecase.CreateHabitUseCase // ← Adicionar
}

func NewApplication(infra *Infrastructure) *Application {
    return &Application{
        // ...
        CreateHabitUseCase: usecase.NewCreateHabitUseCase( // ← Adicionar
            infra.HabitRepository,
        ),
    }
}
```

**Passo 4:** Adicionar no `delivery.go`:
```go
type Delivery struct {
    // ...
    HabitHandler *deliveryHttp.HabitHandler // ← Adicionar
}

func NewDelivery(app *Application) *Delivery {
    // ...
    habitHandler := deliveryHttp.NewHabitHandler(app.CreateHabitUseCase) // ← Adicionar
    router := deliveryHttp.NewRouter(
        healthHandler,
        userHandler,
        habitHandler, // ← Adicionar no router
    )
    // ...
}
```

**Total: ~10 linhas** para integrar uma nova entidade completa!

#### Vantagens desta Abordagem

✅ **Zero dependências externas** - Sem libs de DI
✅ **Performance máxima** - Sem reflection
✅ **Type-safe** - Erros em compile-time
✅ **Escalável** - Funciona perfeitamente com 100+ entidades
✅ **Testável** - Fácil mockar containers inteiros
✅ **Explícito** - Código claro e auditável
✅ **Clean Architecture compliant** - Separação por camadas

### 🔐 JWT Authentication System

#### Overview
The project implements JWT (JSON Web Token) authentication following Clean Architecture principles with:
- Access tokens (short-lived, 15m default)
- Refresh tokens (long-lived, 7 days default)
- Authentication middleware for protected routes
- Stateless authentication

#### Architecture Components

**1. JWT Service Interface (Port)** - `internal/application/port/jwt_service.go`
```go
type JWTService interface {
    GenerateAccessToken(userID, email string) (string, error)
    GenerateRefreshToken(userID, email string) (string, error)
    ValidateToken(token string) (*TokenClaims, error)
    RefreshAccessToken(refreshToken string) (string, error)
}
```

**2. JWT Service Implementation** - `internal/infrastructure/service/jwt_service.go`
- Uses `github.com/golang-jwt/jwt/v5`
- HMAC-SHA256 signing algorithm
- Configurable token durations
- Custom claims with userID and email

**3. Login Use Case** - `internal/application/usecase/login_user_usecase.go`
```go
type LoginUserUseCase struct {
    userRepo      repository.UserRepository
    hasherService port.HasherService
    jwtService    port.JWTService
}

// Execute validates credentials and returns JWT tokens
func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (*LoginUserOutput, error)
```

**4. Authentication Middleware** - `internal/delivery/http/middleware/auth_middleware.go`
```go
type AuthMiddleware struct {
    jwtService port.JWTService
}

// RequireAuth validates JWT token from Authorization header
// Stores user_id and user_email in Gin context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc
```

#### Configuration

Add to `.env`:
```env
# JWT Configuration
JWT_SECRET=your-super-secret-key-min-32-chars
JWT_ACCESS_TOKEN_DURATION=15m      # 15 minutes
JWT_REFRESH_TOKEN_DURATION=168h    # 7 days (use hours, not days)
```

**Important**: Go's `time.ParseDuration` doesn't support "d" for days. Use "h" for hours (e.g., 168h = 7 days).

#### API Endpoints

**1. Login (Public Route)**
```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK)**:
```json
{
  "userId": "uuid",
  "email": "user@example.com",
  "fullName": "User Name",
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc..."
}
```

**Response (401 Unauthorized)**:
```json
{
  "error": "authentication_failed",
  "message": "invalid email or password"
}
```

**2. Protected Routes**
Add `Authorization: Bearer <access_token>` header:
```http
GET /api/v1/user/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Middleware Responses**:
- **401 Unauthorized** - Missing/invalid/expired token
- **200 OK** - Valid token, user_id and user_email stored in context

#### Implementation Guide

**Step 1**: Add JWT configuration to containers

`cmd/api/container/infrastructure.go`:
```go
type Infrastructure struct {
    JWTService port.JWTService  // Add this
    // ...
}

func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
    jwtService, err := service.NewJWTService(
        cfg.JWT.Secret,
        cfg.JWT.AccessTokenDuration,
        cfg.JWT.RefreshTokenDuration,
    )
    // ...
}
```

**Step 2**: Add LoginUserUseCase

`cmd/api/container/application.go`:
```go
type Application struct {
    LoginUserUseCase *usecase.LoginUserUseCase  // Add this
    // ...
}

func NewApplication(infra *Infrastructure) *Application {
    return &Application{
        LoginUserUseCase: usecase.NewLoginUserUseCase(
            infra.UserRepository,
            infra.HasherService,
            infra.JWTService,
        ),
    }
}
```

**Step 3**: Add AuthMiddleware and update UserHandler

`cmd/api/container/delivery.go`:
```go
func NewDelivery(app *Application, infra *Infrastructure) *Delivery {
    // Initialize middleware
    authMiddleware := middleware.NewAuthMiddleware(infra.JWTService)

    // Update user handler with LoginUseCase
    userHandler := deliveryHttp.NewUserHandler(
        app.CreateUserUseCase,
        app.LoginUserUseCase,  // Add this
    )

    // Pass middleware to router
    router := deliveryHttp.NewRouter(
        healthHandler,
        userHandler,
        authMiddleware,  // Add this
    )
    // ...
}
```

**Step 4**: Configure routes with middleware

`internal/delivery/http/router.go`:
```go
func (r *Router) SetupRoutes() *gin.Engine {
    router := gin.Default()

    v1 := router.Group("/api/v1")
    {
        // Public routes
        v1.POST("/user", r.userHandler.Create)
        v1.POST("/login", r.userHandler.Login)

        // Protected routes
        authenticated := v1.Group("")
        authenticated.Use(r.authMiddleware.RequireAuth())
        {
            authenticated.GET("/user/profile", r.userHandler.GetProfile)
            authenticated.POST("/habit", r.habitHandler.Create)
        }
    }

    return router
}
```

#### Security Best Practices

1. **Secret Key**: Use strong random strings (min 32 chars)
2. **Token Duration**: Short-lived access tokens (15m), longer refresh tokens (7 days)
3. **HTTPS Only**: Always use HTTPS in production
4. **Token Storage**: Store tokens securely on client (HttpOnly cookies recommended)
5. **Error Messages**: Don't reveal user existence in login errors
6. **Password Validation**: Already handled by HasherService with bcrypt

#### Testing Authentication

**Test 1: Login with valid credentials**
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

**Test 2: Access protected route without token**
```bash
curl -X GET http://localhost:8080/api/v1/user/profile
# Expected: 401 Unauthorized
```

**Test 3: Access protected route with token**
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <access_token>"
# Expected: 200 OK with user data
```

#### Helper Functions

Extract user info from protected routes:
```go
import "github.com/igor/chronotask-api/internal/delivery/http/middleware"

func (h *Handler) ProtectedEndpoint(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    email, _ := middleware.GetUserEmail(c)

    // Use userID and email
}
```

### 🗄️ Database Infrastructure (PostgreSQL)

#### Database Setup
- **Driver**: `github.com/jackc/pgx/v5` (modern, performant PostgreSQL driver)
- **Connection Pooling**: `pgxpool` for efficient connection management
- **Location**: `internal/infrastructure/persistence/postgres.go`

#### Configuration
Database configuration is loaded from environment variables via `.env` file:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=chronotask
DB_SSL_MODE=disable
```

#### Connection Pool Settings
- **MaxConns**: 25 (maximum connections)
- **MinConns**: 5 (minimum idle connections)
- **MaxConnLifetime**: 1 hour
- **MaxConnIdleTime**: 30 minutes

#### Repository Pattern (Clean Architecture)
Repositories must follow the Ports and Adapters pattern:

1. **Define Port (Interface)** in `internal/domain/repository/`:
```go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
}
```

2. **Implement Adapter** in `internal/infrastructure/persistence/`:
```go
type PostgresUserRepository struct {
    db *PostgresDB
}

func NewPostgresUserRepository(db *PostgresDB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}
```

3. **Use in Application Layer**: Use cases depend on the interface (Port), not the implementation
```go
type CreateUserUseCase struct {
    userRepo repository.UserRepository  // Interface, not concrete type
}
```

#### Database Initialization Flow
1. Load configuration (`config/config.go`)
2. Create PostgresDB connection pool (`infrastructure/persistence/postgres.go`)
3. Verify database health with ping
4. Create repository implementations
5. Inject repositories into use cases
6. Graceful shutdown on SIGINT/SIGTERM

### 🔒 Security Considerations

- **Password Hashing**: Enforced via a required `HasherService` interface (DIP).
- **Input Validation**: Strict validation in both the Delivery Layer (DTOs) and the Domain Layer (Value Objects).

### 🧪 Testing

- **Unit Tests**: High coverage required for **Domain** and **Application Layer (Use Cases)** logic.
- **Mocking**: Used extensively in the Application layer tests to isolate business logic from infrastructure.

### 🤖 AI Assistant Guidelines

When working on this project, prioritize the following:

1.  **DDD First**: Model the **Domain** (Entities, Value Objects, Services) before implementing Use Cases.
2.  **DI Required**: Always assemble components using Dependency Injection.
3.  **Separation of Concerns**: Ensure no database/HTTP code leaks into the Domain or Application layers.