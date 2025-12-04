.PHONY: help build run test clean docker-build docker-run migrate-up migrate-down seed

# Variables
APP_NAME=hospital-emr-backend
VERSION=1.0.0
BUILD_DIR=./bin
GO=go
DOCKER_IMAGE=hospital-emr/backend

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	$(GO) run cmd/api/main.go

dev: ## Run in development mode with hot reload
	@echo "Starting development mode..."
	air

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	$(GO) test -v -short ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GO) test -v -run Integration ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.txt coverage.out coverage.html
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...
	goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

security-scan: ## Run security scan
	@echo "Running security scan..."
	gosec ./...

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	$(GO) run cmd/migrate/main.go up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	$(GO) run cmd/migrate/main.go down

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	@echo "Creating migration: $(name)"
	$(GO) run cmd/migrate/main.go create $(name)

seed: ## Seed database with initial data
	@echo "Seeding database..."
	$(GO) run cmd/seed/main.go

reset-db: ## Reset database (drop all tables, migrate up, and seed)
	@echo "Resetting database..."
	$(MAKE) migrate-down
	$(MAKE) migrate-up
	$(MAKE) seed
	@echo "Database reset complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

docker-compose-up: ## Start all services with docker-compose
	@echo "Starting services with docker-compose..."
	docker-compose up -d

docker-compose-down: ## Stop all services
	@echo "Stopping services..."
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	swag init -g cmd/api/main.go -o api/docs

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest

.DEFAULT_GOAL := help
