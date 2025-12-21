.PHONY: help build install clean test lint fmt

# Variables
BINARY_DIR := bin
CMD_DIR := cmd
CMDS := $(shell ls $(CMD_DIR))

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all binaries
	@echo "Building binaries..."
	@mkdir -p $(BINARY_DIR)
	@for cmd in $(CMDS); do \
		echo "Building $$cmd..."; \
		go build -o $(BINARY_DIR)/$$cmd ./$(CMD_DIR)/$$cmd; \
	done
	@echo "✓ Build complete. Binaries in $(BINARY_DIR)/"

build-%: ## Build a specific binary (e.g., make build-ddd-gen)
	@echo "Building $*..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$* ./$(CMD_DIR)/$*
	@echo "✓ Built $(BINARY_DIR)/$*"

install: ## Install all binaries to GOPATH/bin
	@echo "Installing binaries..."
	@for cmd in $(CMDS); do \
		echo "Installing $$cmd..."; \
		go install ./$(CMD_DIR)/$$cmd; \
	done
	@echo "✓ Installation complete"

install-%: ## Install a specific binary (e.g., make install-ddd-gen)
	@echo "Installing $*..."
	@go install ./$(CMD_DIR)/$*
	@echo "✓ Installed $*"

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@go clean
	@echo "✓ Clean complete"

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests complete"

test-coverage: test ## Run tests and show coverage
	@go tool cover -html=coverage.out

test-pkg: ## Run tests for pkg directory only
	@echo "Running pkg tests..."
	@go test -v -race ./pkg/...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .
	@echo "✓ Formatting complete"

tidy: ## Tidy go modules
	@echo "Tidying go.mod..."
	@go mod tidy
	@echo "✓ Tidy complete"

verify: fmt lint test ## Format, lint, and test

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

upgrade-deps: ## Upgrade all dependencies
	@echo "Upgrading dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✓ Dependencies upgraded"

list-cmds: ## List all available commands
	@echo "Available commands:"
	@for cmd in $(CMDS); do \
		echo "  - $$cmd"; \
	done

version: ## Show Go version
	@go version

# Development helpers
dev-ddd-gen: ## Build and run ddd-gen with example args
	@go run ./cmd/ddd-gen --domain=example --output=./tmp

watch: ## Watch for changes and rebuild (requires entr)
	@echo "Watching for changes..."
	@find . -name '*.go' | entr -r make build
