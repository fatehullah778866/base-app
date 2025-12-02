.PHONY: help build run test clean migrate-up migrate-down docker-build docker-run docker-compose-up docker-compose-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	go build -o bin/server ./cmd/server

run: ## Run the application
	@echo "Running application..."
	go run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	go clean

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@chmod +x scripts/migrate.sh
	./scripts/migrate.sh up

migrate-down: ## Run database migrations down
	@echo "Rolling back migrations..."
	@chmod +x scripts/migrate.sh
	./scripts/migrate.sh down

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@chmod +x scripts/migrate.sh
	./scripts/migrate.sh create $(NAME)

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t base-app-service:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env base-app-service:latest

docker-compose-up: ## Start docker-compose services (PostgreSQL + Redis)
	@echo "Starting docker-compose services..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services are ready!"

docker-compose-down: ## Stop docker-compose services
	@echo "Stopping docker-compose services..."
	docker-compose down

setup: docker-compose-up migrate-up ## Setup local development environment
	@echo "Setup complete! Services are running."
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"

dev: setup run ## Start development environment and run server

