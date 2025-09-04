# Makefile for faceit-cli

# Variables
BINARY_NAME=faceit-cli
VERSION?=dev
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"

# Default target
.PHONY: all
all: test build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 main.go
	@echo "Build complete! Binaries are in the dist/ directory"

# Run unit tests (fast, no external dependencies)
.PHONY: test
test:
	@echo "Running unit tests..."
	go test -v -race -short ./...

# Run integration tests (requires FACEIT_API_KEY)
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@if [ -z "$$FACEIT_API_KEY" ]; then \
		echo "Error: FACEIT_API_KEY environment variable is required for integration tests"; \
		echo "Set it with: export FACEIT_API_KEY=your_api_key"; \
		exit 1; \
	fi
	go test -v -race ./internal/repository/

# Run all tests (unit + integration)
.PHONY: test-all
test-all: test test-integration

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -short ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	@if [ -z "$$FACEIT_API_KEY" ]; then \
		echo "Error: FACEIT_API_KEY environment variable is required for benchmarks"; \
		echo "Set it with: export FACEIT_API_KEY=your_api_key"; \
		exit 1; \
	fi
	go test -v -bench=. ./internal/repository/

# Run linting
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -rf dist/
	rm -f coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Show version
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Git commit: $(GIT_COMMIT)"

# Docker Compose targets
.PHONY: kafka-up
kafka-up:
	@echo "Starting Kafka infrastructure..."
	docker-compose up -d

.PHONY: kafka-down
kafka-down:
	@echo "Stopping Kafka infrastructure..."
	docker-compose down

.PHONY: kafka-logs
kafka-logs:
	@echo "Showing Kafka logs..."
	docker-compose logs -f

.PHONY: kafka-ui
kafka-ui:
	@echo "Kafka UI available at: http://localhost:8080"

.PHONY: kafka-topics
kafka-topics:
	@echo "Listing Kafka topics..."
	docker exec -it faceit-cli-kafka kafka-topics --bootstrap-server localhost:9092 --list

.PHONY: kafka-create-topic
kafka-create-topic:
	@echo "Creating faceit-cli-logs topic..."
	docker exec -it faceit-cli-kafka kafka-topics --bootstrap-server localhost:9092 --create --if-not-exists --topic faceit-cli-logs --partitions 3 --replication-factor 1

# Run with Kafka logging enabled
.PHONY: run-kafka
run-kafka:
	@echo "Running with Kafka logging enabled..."
	KAFKA_ENABLED=true LOG_LEVEL=debug go run main.go

# Run in production mode (no stdout logs, only Kafka)
.PHONY: run-production
run-production:
	@echo "Running in production mode..."
	PRODUCTION_MODE=true LOG_TO_STDOUT=false KAFKA_ENABLED=true LOG_LEVEL=info go run main.go

# Run in production mode with Kafka
.PHONY: run-prod-kafka
run-prod-kafka:
	@echo "Running in production mode with Kafka..."
	PRODUCTION_MODE=true LOG_TO_STDOUT=false KAFKA_ENABLED=true LOG_LEVEL=info make kafka-up && go run main.go

# Run with caching enabled
.PHONY: run-cache
run-cache:
	@echo "Running with caching enabled..."
	CACHE_ENABLED=true CACHE_TTL=30 go run main.go

# Run with all optimizations
.PHONY: run-optimized
run-optimized:
	@echo "Running with all optimizations..."
	CACHE_ENABLED=true CACHE_TTL=30 MATCHES_PER_PAGE=10 MAX_MATCHES_TO_LOAD=50 go run main.go

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  test          - Run unit tests (fast, no external dependencies)"
	@echo "  test-integration - Run integration tests (requires FACEIT_API_KEY)"
	@echo "  test-all      - Run all tests (unit + integration)"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  benchmark     - Run benchmarks (requires FACEIT_API_KEY)"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  install-tools - Install development tools"
	@echo "  run           - Build and run the application"
	@echo "  run-kafka     - Run with Kafka logging enabled"
	@echo "  run-production - Run in production mode (no stdout logs)"
	@echo "  run-prod-kafka - Run in production mode with Kafka"
	@echo "  run-cache     - Run with caching enabled"
	@echo "  run-optimized - Run with all optimizations (cache + pagination)"
	@echo "  kafka-up      - Start Kafka infrastructure (Kafka KRaft, Kafka UI)"
	@echo "  kafka-down    - Stop Kafka infrastructure"
	@echo "  kafka-logs    - Show Kafka infrastructure logs"
	@echo "  kafka-ui      - Show Kafka UI URL (http://localhost:8080)"
	@echo "  kafka-topics  - List all Kafka topics"
	@echo "  kafka-create-topic - Create faceit-cli-logs topic manually"
	@echo "  version       - Show version information"
	@echo "  help          - Show this help message"
