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