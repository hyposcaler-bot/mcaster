.PHONY: build clean test install lint fmt vet run-send run-receive help

BINARY_NAME=mcaster
MAIN_PATH=./cmd/mcaster
BUILD_DIR=bin
VERSION ?= dev
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(MAIN_PATH)

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

## lint: Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## mod-tidy: Clean up go.mod
mod-tidy:
	@echo "Tidying go.mod..."
	@go mod tidy

## run-send: Run sender with development settings
run-send:
	@echo "Starting sender..."
	@go run $(MAIN_PATH) send

## run-receive: Run receiver with development settings
run-receive:
	@echo "Starting receiver..."
	@go run $(MAIN_PATH) receive

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

## release: Build release binaries for multiple platforms
release:
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)/release
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## [a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = "## "}; {printf "  %-15s %s\n", $$2}' | \
		sed 's/: / - /'

.DEFAULT_GOAL := build
