# DocLoom Makefile

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Build variables
BINARY_NAME := docloom
AGENT_BINARY := docloom-agent-csharp
BUILD_DIR := build
MAIN_PACKAGE := ./cmd/docloom
AGENT_PACKAGE := ./cmd/docloom-agent-csharp

# Go build flags
LDFLAGS := -ldflags "-X github.com/karolswdev/docloom/internal/version.Version=$(VERSION) \
	-X github.com/karolswdev/docloom/internal/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/karolswdev/docloom/internal/version.BuildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "Building $(AGENT_BINARY)..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_PACKAGE)
	@echo "Agent binary built: $(BUILD_DIR)/$(AGENT_BINARY)"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -race -count=1 ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo "Coverage report saved to coverage.out"
	@echo "Run 'go tool cover -html=coverage.out' to view HTML report"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .
	goimports -w .

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

# Run all quality checks
.PHONY: ci
ci: fmt vet lint test

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

# Install the binary to $GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PACKAGE)

# Run the application
.PHONY: run
run:
	go run $(LDFLAGS) $(MAIN_PACKAGE) $(ARGS)

# Display version info
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint"
	@echo "  ci            - Run all quality checks (fmt, vet, lint, test)"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install the binary to GOPATH/bin"
	@echo "  run           - Run the application"
	@echo "  version       - Display version information"
	@echo "  help          - Display this help message"