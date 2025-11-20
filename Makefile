.PHONY: build run test clean docker-build docker-run docker-stop deps lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=api
BINARY_UNIX=$(BINARY_NAME)_unix

api-reb: 
	@echo "Rebuilding the Docker image..."
	docker-compose up -d --build api
	@echo "Rebuild complete"

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/app/main.go

# Build for linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v cmd/app/main.go

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/app/main.go
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
deps-update:
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run security check
security:
	gosec ./...

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME):latest .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

# Development setup
dev-setup: deps
	@echo "Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest


# Run all checks
check: fmt lint test security


reb-api:
	docker-compose up -d --build apis

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  build-linux   - Build for Linux"
	@echo "  run           - Build and run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build files"
	@echo "  deps          - Download dependencies"
	@echo "  deps-update   - Update dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  security      - Run security check"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker Compose"
	@echo "  docker-logs   - View application logs"
	@echo "  dev-setup     - Setup development environment"
	@echo "  check         - Run all checks (format, lint, test, security)"
	@echo "  help          - Show this help"
