.PHONY: build run dev test clean docker-build docker-run migrate

# Variables
APP_NAME=sipodi-api
MAIN_PATH=./cmd/api
DOCKER_IMAGE=sipodi/backend

# Build the application
build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

# Run the application
run: build
	./bin/$(APP_NAME)

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

# Generate swagger docs (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	swag init -g cmd/api/main.go -o docs/swagger

# Database migration (using psql)
migrate:
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f db/db.sql

# Create .env from example
env:
	cp .env.example .env

# Help
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Build and run the application"
	@echo "  make dev           - Run with hot reload"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Lint code"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make migrate       - Run database migrations"
	@echo "  make env           - Create .env from example"
