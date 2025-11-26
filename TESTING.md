# 🧪 Testing Guide - ChronoTask API

Este guia documenta a estratégia de testes do projeto e como executá-los.

## 📊 Cobertura de Testes

O projeto possui **3 níveis de testes** para garantir qualidade e confiabilidade:

### 1. **Testes Unitários** (Unit Tests)
- **Camada Domain**: Entidades e Value Objects
- **Camada Application**: Use Cases (com mocks)
- **Objetivo**: Validar lógica de negócio isoladamente
- **Velocidade**: Muito rápidos (< 1s)

### 2. **Testes de Integração** (Integration Tests)
- **Camada Infrastructure**: Repositórios com banco de dados real
- **Objetivo**: Validar interação com PostgreSQL
- **Velocidade**: Moderados (2-5s)

### 3. **Testes End-to-End** (E2E Tests)
- **Camada Delivery**: Handlers HTTP completos
- **Objetivo**: Validar fluxo completo da API
- **Velocidade**: Rápidos (< 1s com mocks)

---

## 🚀 Como Executar os Testes

### Executar TODOS os testes unitários e E2E

```bash
go test ./... -short
```

### Executar testes de uma camada específica

```bash
# Domain Layer (Entities)
go test ./internal/domain/entity -v

# Application Layer (Use Cases)
go test ./internal/application/usecase -v

# Delivery Layer (HTTP Handlers)
go test ./internal/delivery/http -v
```

### Executar testes de uma entidade específica

```bash
# Testes da entidade Character
go test ./internal/domain/entity -v -run Character

# Testes do CreateCharacterUseCase
go test ./internal/application/usecase -v -run CreateCharacter

# Testes do CharacterHandler (E2E)
go test ./internal/delivery/http -v -run Character
```

---

## 🗄️ Testes de Integração (com PostgreSQL)

Os testes de integração requerem um banco de dados PostgreSQL de testes.

### Setup do Banco de Testes

1. **Criar banco de dados de testes:**

```bash
createdb chronotask_test
```

2. **Executar migrations:**

```bash
psql -d chronotask_test -f internal/infrastructure/persistence/migrations/001_create_users_table.sql
psql -d chronotask_test -f internal/infrastructure/persistence/migrations/002_create_characters_table.sql
```

3. **Configurar variáveis de ambiente:**

```bash
export TEST_DATABASE_URL=postgres://postgres:password@localhost:5432/chronotask_test
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=chronotask
```

4. **Executar testes de integração:**

```bash
go test ./internal/infrastructure/persistence -v
```

**Nota**: Se `TEST_DATABASE_URL` não estiver configurado, os testes de integração serão **automaticamente pulados** (skip).

---

## 📈 Verificar Cobertura de Testes

### Cobertura geral

```bash
go test ./... -cover
```

### Cobertura detalhada com HTML

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Cobertura de uma camada específica

```bash
# Domain Layer
go test ./internal/domain/entity -cover

# Application Layer
go test ./internal/application/usecase -cover

# Delivery Layer
go test ./internal/delivery/http -cover
```

---

## ✅ Checklist de Testes para Novas Features

Ao adicionar uma nova entidade ou feature, garanta que possui:

- [ ] **Testes unitários da entidade** (`*_test.go` no domain/entity)
  - Testes de criação válida
  - Testes de validações (campos inválidos)
  - Testes de métodos de negócio
  - Teste de `Reconstitute*` function

- [ ] **Testes unitários do Use Case** (`*_test.go` no application/usecase)
  - Cenário de sucesso
  - Cenários de erro (validações, regras de negócio)
  - Erros de dependências (repository, services)
  - Contexto cancelado

- [ ] **Testes de integração do Repository** (`*_test.go` no infrastructure/persistence)
  - Create, FindByID, Update, Delete
  - Constraints do banco (Foreign Keys, Unique)
  - Queries específicas (FindByUserID, etc)

- [ ] **Testes E2E do Handler** (`*_test.go` no delivery/http)
  - Requisição válida (201 Created)
  - Autenticação (401 Unauthorized)
  - Validações de input (400 Bad Request)
  - Regras de negócio (409 Conflict, etc)
  - Erros internos (422 Unprocessable Entity)

---

## 🎯 Exemplo: Testes da Entidade Character

### Testes Unitários (Domain)

✅ **15 testes** cobrindo:
- Criação válida
- Validações (nome, ID, userID)
- Sistema de XP e level-up
- Cálculo de XP necessário
- Atualização de nome
- Reconstitution

**Executar:**
```bash
go test ./internal/domain/entity -v -run Character
```

### Testes Unitários (Use Case)

✅ **7 testes** cobrindo:
- Criação bem-sucedida
- Usuário já tem personagem
- Validações de nome
- Erros de repository
- Contexto cancelado

**Executar:**
```bash
go test ./internal/application/usecase -v -run CreateCharacter
```

### Testes de Integração (Repository)

✅ **8 testes** cobrindo:
- CRUD completo
- Foreign Key constraint
- Unique user constraint
- ExistsByUserID

**Executar:**
```bash
export TEST_DATABASE_URL=postgres://...
go test ./internal/infrastructure/persistence -v -run Character
```

### Testes E2E (HTTP Handler)

✅ **7 testes** cobrindo:
- POST /api/v1/character com sucesso
- Autenticação (token válido/inválido)
- Validações de input
- Usuário já tem personagem (409)
- Erros de repositório

**Executar:**
```bash
go test ./internal/delivery/http -v -run Character
```

---

## 🛠️ Ferramentas de Teste

### Bibliotecas Utilizadas

- **testing** - Framework nativo do Go
- **httptest** - Testes HTTP sem servidor real
- **gin** (TestMode) - Framework HTTP em modo de teste

### Mocks

Mocks são criados manualmente seguindo as interfaces do domain:
- `mockCharacterRepository` - Mock do repository
- `mockJWTService` - Mock do serviço JWT
- `mockHasherService` - Mock do hasher (bcrypt)

---

## 📝 Convenções de Nomenclatura

### Arquivos de Teste
```
entity_name_test.go        # Testes da entidade
usecase_name_test.go       # Testes do use case
handler_name_test.go       # Testes E2E do handler
repository_name_test.go    # Testes de integração
```

### Funções de Teste
```go
func TestEntityName_MethodName_Scenario(t *testing.T)
func TestUseCaseName_Execute_Scenario(t *testing.T)
func TestHandler_Action_Scenario(t *testing.T)
```

**Exemplos:**
```go
TestNewCharacter_ValidCharacter
TestCreateCharacterUseCase_Execute_Success
TestCharacterHandler_Create_MissingAuthorization
```

---

## 🎓 Boas Práticas

1. **AAA Pattern**: Arrange, Act, Assert
2. **Table-Driven Tests**: Use slices para múltiplos casos
3. **Test Helpers**: Funções auxiliares com `t.Helper()`
4. **Mocks Simples**: Preferir mocks manuais a frameworks complexos
5. **Isolamento**: Cada teste deve ser independente
6. **Limpeza**: Use `defer cleanup()` em testes de integração
7. **Nomenclatura Clara**: Nome do teste deve descrever o cenário

---

## 📊 Métricas de Qualidade

### Objetivos de Cobertura

- **Domain Layer**: > 90%
- **Application Layer**: > 85%
- **Infrastructure Layer**: > 70%
- **Delivery Layer**: > 80%

### Execução em CI/CD

```yaml
# .github/workflows/test.yml (exemplo)
- name: Run tests
  run: go test ./... -short -cover

- name: Run integration tests
  run: |
    docker-compose up -d postgres
    go test ./internal/infrastructure/persistence -v
```

---

## 🐛 Debugging de Testes

### Ver logs detalhados

```bash
go test ./... -v
```

### Executar um teste específico

```bash
go test ./internal/domain/entity -v -run TestNewCharacter_ValidCharacter
```

### Executar com race detector

```bash
go test ./... -race
```

### Executar apenas testes rápidos

```bash
go test ./... -short
```

---

## ✨ Resultados Atuais

```
✅ Domain Layer (Entity):        15/15 passed
✅ Application Layer (UseCase):   7/7 passed
✅ Delivery Layer (Handler E2E):  7/7 passed
✅ Infrastructure Layer (Repo):   8/8 passed (with TEST_DATABASE_URL)

Total: 37 testes, 100% de sucesso
```

**Todas as camadas da entidade Character estão completamente testadas!** 🎉
