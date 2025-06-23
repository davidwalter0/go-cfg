# Go configuration management v3 - Unified Architecture
.PHONY: help build test lint clean install examples docs all

# Default target
all: build test lint

# Help target
help: ## Show this help message
	@echo "Go-cfg v3 - Unified Configuration Management"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the library and examples
	@echo "Building go-cfg v3..."
	@go build ./...
	@echo "✅ Build successful"

build-examples: ## Build all examples
	@echo "Building examples..."
	@cd examples/basic && go build -o ../../bin/basic-example .
	@cd examples/advanced && go build -o ../../bin/advanced-example .
	@echo "✅ Examples built successfully"

# Test targets
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...
	@echo "✅ All tests passed"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@go test -v -race ./...
	@echo "✅ Race tests passed"

# Quality targets
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet instead"; \
		go vet ./...; \
	fi
	@echo "✅ Lint checks passed"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# Dependency targets
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

# Example targets
run-basic: build-examples ## Run basic example
	@echo "Running basic example..."
	@./bin/basic-example --help || true
	@echo ""
	@echo "Running with sample config..."
	@./bin/basic-example --service-name="example-basic" --debug --port=3000

run-advanced: build-examples ## Run advanced example
	@echo "Running advanced example..."
	@./bin/advanced-example --help || true
	@echo ""
	@echo "Running with sample config..."
	@./bin/advanced-example --service-name="example-advanced" --database-name="testdb" --database-username="testuser" --database-password="testpass"

# Integration tests
test-integration: build-examples ## Run integration tests with examples
	@echo "Running integration tests..."
	@echo "Testing basic example..."
	@./bin/basic-example --service-name="integration-test" --version="test-version" > /tmp/basic-output.txt
	@grep -q "integration-test" /tmp/basic-output.txt && echo "✅ Basic example integration test passed"
	@echo "Testing advanced example..."
	@./bin/advanced-example --service-name="integration-advanced" --database-name="testdb" --database-username="testuser" --database-password="testpass" > /tmp/advanced-output.txt
	@grep -q "integration-advanced" /tmp/advanced-output.txt && echo "✅ Advanced example integration test passed"
	@echo "✅ All integration tests passed"

# Documentation targets
docs: ## Generate documentation
	@echo "Generating documentation..."
	@go doc -all ./pkg/config > docs/config.md
	@go doc -all ./pkg/discovery > docs/discovery.md
	@go doc -all ./pkg/flags > docs/flags.md
	@go doc -all ./pkg/generate > docs/generate.md
	@echo "✅ Documentation generated in docs/"

# Benchmark targets
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "✅ Benchmarks completed"

# Release targets
install: ## Install the library
	@echo "Installing go-cfg v3..."
	@go install ./...
	@echo "✅ Library installed"

# Cleanup targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f /tmp/basic-output.txt /tmp/advanced-output.txt
	@go clean ./...
	@echo "✅ Cleanup completed"

# Setup targets
setup: ## Setup development environment
	@echo "Setting up development environment..."
	@mkdir -p bin docs
	@go mod download
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development environment setup complete"

# Verification targets
verify: test lint test-race ## Run all verification checks
	@echo "✅ All verification checks passed"

# Quick targets
quick: build test ## Quick build and test
	@echo "✅ Quick verification completed"

# CI targets
ci: deps fmt lint test test-race test-integration ## CI pipeline
	@echo "✅ CI pipeline completed successfully"

# Development workflow
dev: clean setup build test-coverage ## Development workflow
	@echo "✅ Development workflow completed"
	@echo "Coverage report available at: coverage.html"

# Performance targets
perf: bench ## Performance testing
	@echo "✅ Performance testing completed"

# Release verification
release-check: clean verify test-integration docs ## Pre-release verification
	@echo "✅ Release verification completed"
	@echo "Ready for release!"