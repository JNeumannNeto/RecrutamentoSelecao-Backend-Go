.PHONY: build run test clean docker-up docker-down migrate-up migrate-down

# Build all services
build:
	@echo "Building all services..."
	@cd services/auth-service && go build -o ../../bin/auth-service ./cmd/main.go
	@cd services/job-service && go build -o ../../bin/job-service ./cmd/main.go
	@cd services/candidate-service && go build -o ../../bin/candidate-service ./cmd/main.go

# Run all services locally
run-all:
	@echo "Starting all services..."
	@make run-auth &
	@make run-job &
	@make run-candidate &
	@wait

run-auth:
	@cd services/auth-service && go run ./cmd/main.go

run-job:
	@cd services/job-service && go run ./cmd/main.go

run-candidate:
	@cd services/candidate-service && go run ./cmd/main.go

# Test all services
test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Docker commands
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

docker-build:
	@echo "Building Docker images..."
	@docker-compose build

docker-logs:
	@docker-compose logs -f

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@docker-compose exec postgres psql -U postgres -d recruitment_db -f /docker-entrypoint-initdb.d/001_init.sql

migrate-down:
	@echo "Rolling back migrations..."
	@docker-compose exec postgres psql -U postgres -d recruitment_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@docker-compose down -v
	@docker system prune -f

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Create directories
setup:
	@echo "Setting up project structure..."
	@mkdir -p bin
	@mkdir -p services/auth-service/cmd
	@mkdir -p services/auth-service/internal/domain
	@mkdir -p services/auth-service/internal/application
	@mkdir -p services/auth-service/internal/infrastructure
	@mkdir -p services/auth-service/internal/interfaces
	@mkdir -p services/job-service/cmd
	@mkdir -p services/job-service/internal/domain
	@mkdir -p services/job-service/internal/application
	@mkdir -p services/job-service/internal/infrastructure
	@mkdir -p services/job-service/internal/interfaces
	@mkdir -p services/candidate-service/cmd
	@mkdir -p services/candidate-service/internal/domain
	@mkdir -p services/candidate-service/internal/application
	@mkdir -p services/candidate-service/internal/infrastructure
	@mkdir -p services/candidate-service/internal/interfaces
	@mkdir -p shared/database
	@mkdir -p shared/middleware
	@mkdir -p shared/utils
	@mkdir -p migrations

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build all services"
	@echo "  run-all       - Run all services locally"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  docker-up     - Start Docker services"
	@echo "  docker-down   - Stop Docker services"
	@echo "  docker-build  - Build Docker images"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  clean         - Clean up build artifacts and Docker"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  setup         - Create project directories"
	@echo "  help          - Show this help message"
