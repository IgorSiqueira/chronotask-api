# ChronoTask API

> Sistema de rastreamento de hábitos com evolução de personagem que transforma atividades da vida real em progressão dentro do jogo (XP, Stats) para permitir batalhas PVP/PVE simuladas.

## 📋 Índice

- [Sobre o Projeto](#sobre-o-projeto)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Configuração](#configuração)
- [Como Executar](#como-executar)
- [API Endpoints](#api-endpoints)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Testes](#testes)
- [Princípios Arquiteturais](#princípios-arquiteturais)
- [Desenvolvimento](#desenvolvimento)

## 🎯 Sobre o Projeto

**ChronoTask API** é uma API backend desenvolvida em Go que implementa um sistema gamificado de rastreamento de hábitos. O projeto transforma atividades cotidianas em progressão de personagem, permitindo que usuários evoluam seus avatares através de hábitos saudáveis e batalhem contra outros jogadores ou contra o sistema.

### Funcionalidades Principais

- ✅ Gerenciamento de usuários com autenticação segura
- 🎮 Sistema de evolução de personagem baseado em hábitos
- 📊 Rastreamento de atividades e progresso
- ⚔️ Sistema de batalhas PVP/PVE (em desenvolvimento)
- 📈 Sistema de XP e estatísticas

## 🚀 Tecnologias

Este projeto foi construído utilizando as seguintes tecnologias:

### Core

- **[Go](https://golang.org/)** (1.21+) - Linguagem de programação
- **[Gin](https://gin-gonic.com/)** - Framework HTTP web minimalista
- **[PostgreSQL](https://www.postgresql.org/)** - Banco de dados relacional
- **[pgx/v5](https://github.com/jackc/pgx)** - Driver PostgreSQL de alto desempenho

### Bibliotecas e Ferramentas

- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Hash seguro de senhas
- **[godotenv](https://github.com/joho/godotenv)** - Gerenciamento de variáveis de ambiente
- **[uuid](https://github.com/google/uuid)** - Geração de identificadores únicos

## 🏗️ Arquitetura

O projeto segue rigorosamente os princípios de **Clean Architecture** e **Domain-Driven Design (DDD)** com separação estrita de responsabilidades.

### Camadas da Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                    Delivery Layer (HTTP)                     │
│                   Handlers, DTOs, Router                     │
└────────────────────────┬────────────────────────────────────┘
                         │ Depende de ↓
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                          │
│              Use Cases, Ports (Interfaces)                   │
└────────────────────────┬────────────────────────────────────┘
                         │ Depende de ↓
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                             │
│         Entities, Value Objects, Repository Interfaces       │
│              (Sem dependências externas)                     │
└────────────────────────┬────────────────────────────────────┘
                         │ Implementa ↑
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                        │
│         Database, External Services, Implementations         │
└─────────────────────────────────────────────────────────────┘
```

### Estrutura de Diretórios

```
chronotask-api/
├── cmd/
│   └── api/                    # Ponto de entrada da aplicação
│       └── main.go            # Bootstrap e injeção de dependências
│
├── internal/                   # Código privado da aplicação
│   ├── domain/                # Camada de Domínio (Core Business)
│   │   ├── entity/           # Entidades (User, Habit, Character)
│   │   ├── valueobject/      # Objetos de Valor (Email, etc)
│   │   ├── service/          # Serviços de Domínio
│   │   └── repository/       # Interfaces de Repositório (Ports)
│   │
│   ├── application/           # Camada de Aplicação (Casos de Uso)
│   │   ├── usecase/          # Implementação dos Use Cases
│   │   └── port/             # Interfaces de serviços externos
│   │
│   ├── infrastructure/        # Camada de Infraestrutura
│   │   ├── persistence/      # Implementações de repositórios
│   │   │   ├── migrations/  # Migrações SQL
│   │   │   ├── postgres.go  # Conexão com PostgreSQL
│   │   │   └── *_repository.go
│   │   └── service/          # Implementações de serviços
│   │
│   └── delivery/              # Camada de Entrega (Interfaces Externas)
│       ├── http/             # Handlers HTTP
│       │   ├── router.go    # Configuração de rotas
│       │   └── *_handler.go # Handlers por recurso
│       └── dto/              # Data Transfer Objects
│
├── config/                    # Configurações da aplicação
│   └── config.go             # Carregamento de configs
│
├── scripts/                   # Scripts auxiliares
│   └── migrate.sh            # Script de migração
│
├── .env.example              # Exemplo de variáveis de ambiente
├── .gitignore
├── go.mod                    # Dependências Go
├── go.sum
├── Makefile                  # Comandos úteis
├── CLAUDE.md                 # Contexto para IA
└── README.md                 # Este arquivo
```

## 📦 Pré-requisitos

Antes de começar, você precisará ter instalado em sua máquina:

- **Go** 1.21 ou superior - [Download](https://golang.org/dl/)
- **PostgreSQL** 13+ - [Download](https://www.postgresql.org/download/)
- **Make** (opcional, mas recomendado)
- **Git**

## 💻 Instalação

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/chronotask-api.git
cd chronotask-api
```

### 2. Instale as dependências

```bash
go mod download
go mod tidy
```

### 3. Configure o banco de dados

```bash
# Acesse o PostgreSQL
psql -U postgres

# Crie o banco de dados
CREATE DATABASE chronotask;

# Saia do psql
\q
```

## ⚙️ Configuração

### Variáveis de Ambiente

Copie o arquivo de exemplo e configure suas variáveis:

```bash
cp .env.example .env
```

Edite o arquivo `.env` com suas configurações:

```env
# Configuração do Servidor
SERVER_PORT=8080

# Configuração do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua_senha_aqui
DB_NAME=chronotask
DB_SSL_MODE=disable
```

### Execute as Migrações

```bash
# Torne o script executável (apenas uma vez)
chmod +x scripts/migrate.sh

# Execute as migrações
./scripts/migrate.sh
```

## 🎮 Como Executar

### Usando Make (Recomendado)

```bash
# Executar a aplicação
make run

# Compilar
make build

# Executar testes
make test

# Ver todos os comandos disponíveis
make help
```

### Diretamente com Go

```bash
# Executar em modo desenvolvimento
go run cmd/api/main.go

# Compilar para produção
go build -o bin/api cmd/api/main.go

# Executar o binário compilado
./bin/api
```

### Verificar se está funcionando

```bash
# Health check
curl http://localhost:8080/health
```

Resposta esperada:
```json
{"status":"ok"}
```

## 🔌 API Endpoints

### Health Check

```http
GET /health
```

Verifica se a API está funcionando.

**Resposta de Sucesso (200)**
```json
{
  "status": "ok"
}
```

---

### Criar Usuário

```http
POST /api/v1/user
Content-Type: application/json
```

Cria um novo usuário no sistema.

**Corpo da Requisição**
```json
{
  "fullName": "João da Silva",
  "email": "joao.silva@email.com",
  "password": "SenhaSegura123",
  "birthDate": "1990-05-15",
  "acceptTerms": true
}
```

**Validações**
- `fullName`: Obrigatório, mínimo 2 caracteres, máximo 255
- `email`: Obrigatório, formato de email válido
- `password`: Obrigatório, mínimo 8 caracteres
- `birthDate`: Obrigatório, formato YYYY-MM-DD, idade mínima 13 anos
- `acceptTerms`: Obrigatório, deve ser `true`

**Resposta de Sucesso (201)**
```json
{
  "id": "3db4e681-fc6a-42d0-b9c8-20b44bd55291",
  "fullName": "João da Silva",
  "email": "joao.silva@email.com",
  "birthDate": "1990-05-15",
  "createdAt": "2025-11-24T20:32:14.744128-03:00"
}
```

**Resposta de Erro - Email Duplicado (422)**
```json
{
  "error": "user_creation_failed",
  "message": "user with email joao.silva@email.com already exists"
}
```

**Resposta de Erro - Validação (400)**
```json
{
  "error": "invalid_request",
  "message": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

### Exemplo de Uso com cURL

```bash
curl -X POST http://localhost:8080/api/v1/user \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "João da Silva",
    "email": "joao.silva@email.com",
    "password": "SenhaSegura123",
    "birthDate": "1990-05-15",
    "acceptTerms": true
  }'
```

## 🧪 Testes

### Executar todos os testes

```bash
# Todos os testes
go test ./...

# Com saída verbosa
go test -v ./...

# Apenas testes da camada de domínio
go test ./internal/domain/...
```

### Cobertura de Testes

```bash
# Gerar relatório de cobertura
make test-coverage

# Ou manualmente
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Estrutura de Testes

- **Testes Unitários**: Cobertura de entidades, value objects e use cases
- **Testes de Integração**: Testes de repositórios (em desenvolvimento)
- **Testes E2E**: Testes de endpoints (em desenvolvimento)

## 📐 Princípios Arquiteturais

### SOLID

- **S**ingle Responsibility Principle: Cada struct tem uma única responsabilidade
- **O**pen/Closed Principle: Aberto para extensão, fechado para modificação
- **L**iskov Substitution Principle: Interfaces podem ser substituídas por implementações
- **I**nterface Segregation Principle: Interfaces pequenas e focadas
- **D**ependency Inversion Principle: Dependências apontam para abstrações

### Clean Architecture

1. **Domain Layer** (Centro)
   - Entidades com regras de negócio
   - Value Objects imutáveis
   - Interfaces de repositório
   - **Sem dependências externas**

2. **Application Layer**
   - Use Cases com lógica de aplicação
   - Interfaces de serviços (Ports)
   - Orquestração do domínio

3. **Infrastructure Layer**
   - Implementações de repositórios (Adapters)
   - Serviços externos (bcrypt, etc)
   - Conexões com banco de dados

4. **Delivery Layer** (Externa)
   - HTTP Handlers
   - DTOs de request/response
   - Roteamento

### Domain-Driven Design (DDD)

- **Entities**: Objetos com identidade única (User, Habit)
- **Value Objects**: Objetos definidos por valores (Email)
- **Repositories**: Abstração de persistência
- **Use Cases**: Casos de uso da aplicação
- **Aggregates**: Clusters de objetos de domínio (em desenvolvimento)

### Padrões de Design

- **Repository Pattern**: Abstração de acesso a dados
- **Dependency Injection**: Inversão de controle manual
- **DTO Pattern**: Separação entre domínio e apresentação
- **Factory Pattern**: Criação de entidades com validação

## 🛠️ Desenvolvimento

### Adicionando uma Nova Feature

1. **Modelar o Domínio** (`internal/domain/`)
   ```go
   // entity/habit.go
   type Habit struct {
       id string
       // ...
   }
   ```

2. **Criar Interface do Repositório** (`internal/domain/repository/`)
   ```go
   type HabitRepository interface {
       Create(ctx context.Context, habit *entity.Habit) error
   }
   ```

3. **Implementar Use Case** (`internal/application/usecase/`)
   ```go
   type CreateHabitUseCase struct {
       habitRepo repository.HabitRepository
   }
   ```

4. **Implementar Repositório** (`internal/infrastructure/persistence/`)
   ```go
   type PostgresHabitRepository struct {
       db *PostgresDB
   }
   ```

5. **Criar Handler HTTP** (`internal/delivery/http/`)
   ```go
   func (h *HabitHandler) Create(c *gin.Context) {
       // ...
   }
   ```

6. **Registrar Rota** (`internal/delivery/http/router.go`)
   ```go
   v1.POST("/habit", h.habitHandler.Create)
   ```

7. **Injetar Dependências** (`cmd/api/main.go`)
   ```go
   habitRepo := persistence.NewPostgresHabitRepository(db)
   createHabitUC := usecase.NewCreateHabitUseCase(habitRepo)
   habitHandler := deliveryHttp.NewHabitHandler(createHabitUC)
   ```

### Comandos Úteis

```bash
# Formatar código
make fmt
# ou
go fmt ./...

# Análise estática
make vet
# ou
go vet ./...

# Linting
make lint

# Limpar builds
make clean

# Visualizar dependências
go mod graph
```

### Boas Práticas

- ✅ Escreva testes antes de implementar (TDD)
- ✅ Mantenha o domínio puro (sem dependências externas)
- ✅ Use injeção de dependências
- ✅ Valide no domínio E na camada de entrega
- ✅ Sempre retorne erros descritivos
- ✅ Use context.Context para propagação de cancelamento
- ✅ Documente funções públicas
- ✅ Commits descritivos em português

## 📄 Licença

Este projeto está sob a licença [A DEFINIR].

---

Desenvolvido com ❤️ usando Clean Architecture e DDD
