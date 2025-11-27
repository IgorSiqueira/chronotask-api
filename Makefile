.PHONY: help run build test test-coverage clean docker-build docker-up docker-down docker-logs docker-clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

run: ## Run the application
	@go run cmd/api/main.go

build: ## Build the application
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go
	@echo "Build complete: bin/api"

test: ## Run tests
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@go test -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

deps: ## Download dependencies
	@go mod download
	@go mod tidy

fmt: ## Format code
	@go fmt ./...

vet: ## Run go vet
	@go vet ./...

lint: fmt vet ## Run formatters and linters

# Docker commands
docker-build: ## Build Docker images
	@docker-compose build --no-cache

docker-up: ## Start Docker containers
	@docker-compose up -d
	@echo "Waiting for services..."
	@sleep 5
	@curl -s http://localhost:8080/health || echo "API starting..."

docker-down: ## Stop Docker containers
	@docker-compose down

docker-restart: docker-down docker-up ## Restart Docker containers

docker-logs: ## Show Docker logs
	@docker-compose logs -f

docker-logs-api: ## Show API logs only
	@docker-compose logs -f api

docker-logs-db: ## Show PostgreSQL logs only
	@docker-compose logs -f postgres

docker-ps: ## Show container status
	@docker-compose ps

docker-clean: ## Stop and remove all containers/volumes
	@docker-compose down -v
	@echo "Docker cleaned!"

docker-shell-api: ## Access API container shell
	@docker-compose exec api /bin/sh

docker-shell-db: ## Access PostgreSQL shell
	@docker-compose exec postgres psql -U postgres -d chronotask

docker-backup: ## Backup PostgreSQL database
	@mkdir -p backups
	@docker-compose exec postgres pg_dump -U postgres chronotask > backups/chronotask_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Backup saved!"

docker-health: ## Check API health
	@curl -s http://localhost:8080/health | jq . || echo "API not responding"
