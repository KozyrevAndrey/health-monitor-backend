.PHONY: help build test lint run clean fmt vet tidy docker-build docker-run dev

# Variables
BINARY_NAME=health-monitor
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Directories
BIN_DIR=bin
DATA_DIR=data

help: ## Display this help message
	@echo "Health Monitor - Build Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/server
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)-linux-amd64"

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/server
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/server
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/server
	@echo "All builds complete"

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete"

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
		echo "Lint complete"; \
	else \
		echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Format complete"

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

tidy: ## Tidy go modules
	@echo "Tidying modules..."
	$(GOMOD) tidy
	@echo "Tidy complete"

run: build ## Build and run the application
	@echo "Starting $(BINARY_NAME)..."
	@mkdir -p $(DATA_DIR)
	./$(BIN_DIR)/$(BINARY_NAME) --config configs/example.yaml

dev: ## Run in development mode (without building)
	@echo "Starting in development mode..."
	@mkdir -p $(DATA_DIR)
	$(GOCMD) run ./cmd/server --config configs/example.yaml

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

clean-all: clean ## Clean everything including data
	@rm -rf $(DATA_DIR)
	@echo "Deep clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .
	@echo "Docker build complete"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 -v $$(pwd)/configs:/configs -v $$(pwd)/data:/data $(BINARY_NAME):latest

docker-compose-up: ## Start with docker-compose
	@echo "Starting with docker-compose..."
	docker-compose up -d
	@echo "Services started"

docker-compose-down: ## Stop docker-compose services
	@echo "Stopping docker-compose services..."
	docker-compose down
	@echo "Services stopped"

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

install-deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	fi
	@echo "Dependencies installed"

validate-config: ## Validate configuration file
	@echo "Validating configuration..."
	@$(GOCMD) run ./cmd/server --config configs/example.yaml --version >/dev/null 2>&1 && echo " Configuration valid" || echo " Configuration invalid"

openapi-gen: ## Generate code from OpenAPI specification
	@./scripts/generate-openapi.sh

version: ## Show version information
	@./$(BIN_DIR)/$(BINARY_NAME) --version 2>/dev/null || echo "Build the binary first with 'make build'"

.DEFAULT_GOAL := help

