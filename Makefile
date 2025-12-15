.PHONY: build test lint clean install fmt vet

# Binary name
BINARY_NAME := toepub
BINARY_DIR := ./cmd/toepub

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOCLEAN := $(GOCMD) clean
GOVET := $(GOCMD) vet
GOFMT := gofmt
GOINSTALL := $(GOCMD) install

# Build flags
LDFLAGS := -ldflags "-s -w"

# Default target
all: build

# Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(BINARY_DIR)

# Build with debug symbols
build-debug:
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_DIR)

# Install binary to GOPATH/bin
install:
	$(GOINSTALL) $(BINARY_DIR)

# Run all tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run unit tests only
test-unit:
	$(GOTEST) -v ./internal/...

# Run integration tests only
test-integration:
	$(GOTEST) -v ./tests/integration/...

# Run contract tests only
test-contract:
	$(GOTEST) -v ./tests/contract/...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	$(GOFMT) -w .
	goimports -w .

# Run go vet
vet:
	$(GOVET) ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Cross-compile for all platforms
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(BINARY_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 $(BINARY_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(BINARY_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 $(BINARY_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(BINARY_DIR)

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  build-debug    - Build with debug symbols"
	@echo "  install        - Install binary to GOPATH/bin"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-contract  - Run contract tests only"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code with gofmt and goimports"
	@echo "  vet            - Run go vet"
	@echo "  clean          - Remove build artifacts"
	@echo "  build-all      - Cross-compile for all platforms"
