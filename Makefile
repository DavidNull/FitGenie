# FitGenie Makefile
.PHONY: all build test clean docker-run docker-down lint help

# Variables
BINARY_NAME=server
BUILD_DIR=./build
CMD_PATH=./cmd/server/main.go
DOCKER_COMPOSE=docker-compose

# Default target
all: build

## build: Build the application binary
build:
	@echo "Building FitGenie server..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## test: Run all tests with coverage
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## test-short: Run tests without coverage (faster)
test-short:
	@echo "Running tests (short)..."
	@go test -v -race ./...

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@golangci-lint run --timeout=5m ./...

## lint-fix: Run golangci-lint with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	@golangci-lint run --fix --timeout=5m ./...

## docker-run: Start all services with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	@$(DOCKER_COMPOSE) up --build -d
	@echo "Services starting..."
	@echo "API: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"
	@echo "LocalStack S3: http://localhost:4566"

## docker-down: Stop all Docker Compose services
docker-down:
	@echo "Stopping Docker Compose services..."
	@$(DOCKER_COMPOSE) down

## docker-logs: View logs from all services
docker-logs:
	@$(DOCKER_COMPOSE) logs -f

## docker-clean: Remove all containers, volumes, and images
docker-clean:
	@echo "Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) down -v --rmi all

## run: Build and run locally (requires local PostgreSQL)
run: build
	@echo "Starting server..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run with hot reload (requires air)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air -c .air.toml

## clean: Remove build artifacts and coverage files
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

## deps: Download and verify dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

## tidy: Tidy and vendor dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@go mod vendor

## migrate: Run database migrations (requires running PostgreSQL)
migrate:
	@echo "Running migrations..."
	@go run $(CMD_PATH) migrate

## generate: Generate any auto-generated code
generate:
	@echo "Generating code..."
	@go generate ./...

## ci: Run all CI checks (lint + test)
ci: lint test

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/## //g'
