.PHONY: test test-unit test-integration test-coverage build run migrate clean lint fmt

# Test commands
test:
	go test ./... -v

test-unit:
	go test ./tests/unit/... -v

test-integration:
	go test ./tests/integration/... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Build
build:
	go build -o bin/server ./cmd/server

# Run
run:
	go run ./cmd/server/main.go

# Database
migrate-up:
	go run ./scripts/migrate.go up

migrate-down:
	go run ./scripts/migrate.go down

# Code quality
lint:
	golangci-lint run

fmt:
	go fmt ./...

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run all checks
check: fmt lint test

