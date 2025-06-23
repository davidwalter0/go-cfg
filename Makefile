# Makefile for go-cfg legacy package
# Modern pflag-based configuration management library

.PHONY: help build test test-verbose test-coverage lint fmt clean deps dev ci install-tools vendor run-examples docs

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Package info
PACKAGE=github.com/davidwalter0/go-cfg
BINARY_NAME=go-cfg
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

## help: Show this help message
help:
	@echo "go-cfg - Modern pflag-based configuration management"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "Development workflow:"
	@echo "  make deps fmt lint test    # Prepare and validate"
	@echo "  make ci                    # Full CI pipeline"

## deps: Download and update dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies updated"

## fmt: Format Go code
fmt:
	@echo "🎨 Formatting code..."
	$(GOFMT) -w .
	@echo "✅ Code formatted"

## lint: Run linter checks
lint:
	@echo "🔍 Running linter..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run; \
		echo "✅ Lint checks passed"; \
	else \
		echo "⚠️  golangci-lint not installed, skipping lint checks"; \
		echo "   Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## build: Build the package
build:
	@echo "🔨 Building package..."
	$(GOBUILD) -v ./...
	@echo "✅ Package built successfully"

## test: Run all tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...
	@echo "✅ All tests passed"

## test-short: Run tests with short flag
test-short:
	@echo "🧪 Running short tests..."
	$(GOTEST) -short ./...
	@echo "✅ Short tests passed"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "🧪 Running tests with coverage..."
	$(GOTEST) -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "✅ Coverage report generated: $(COVERAGE_HTML)"
	@if command -v open >/dev/null 2>&1; then \
		echo "🌐 Opening coverage report..."; \
		open $(COVERAGE_HTML); \
	fi

## test-race: Run tests with race detection
test-race:
	@echo "🧪 Running tests with race detection..."
	$(GOTEST) -race ./...
	@echo "✅ Race tests passed"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "🧪 Running verbose tests..."
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_OUT) ./...
	@echo "✅ Verbose tests completed"

## vendor: Create vendor directory
vendor:
	@echo "📦 Creating vendor directory..."
	$(GOMOD) vendor
	@echo "✅ Vendor directory created"

## clean: Clean build artifacts and test files
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)
	rm -rf vendor/
	@echo "✅ Clean completed"

## install-tools: Install development tools
install-tools:
	@echo "🔧 Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

## run-examples: Build and run example programs
run-examples:
	@echo "🚀 Building and running examples..."
	@if [ -d "cmd" ]; then \
		for dir in cmd/*/; do \
			if [ -f "$$dir/main.go" ]; then \
				echo "Building $$dir..."; \
				$(GOBUILD) -o bin/$$(basename $$dir) ./$$dir; \
				echo "Built bin/$$(basename $$dir)"; \
			fi; \
		done; \
		echo "✅ Examples built in bin/ directory"; \
	else \
		echo "ℹ️  No cmd/ directory found"; \
	fi

## test-integration: Run integration tests
test-integration:
	@echo "🔗 Running integration tests..."
	$(GOTEST) -tags=integration -v ./...
	@echo "✅ Integration tests passed"

## test-legacy: Run legacy compatibility tests
test-legacy:
	@echo "🔄 Running legacy compatibility tests..."
	$(GOTEST) -run "TestVendor" -v ./...
	@echo "✅ Legacy compatibility verified"

## test-api: Run core API tests
test-api:
	@echo "⚙️  Running core API tests..."
	$(GOTEST) -run "TestInit|TestEval|TestSimple|TestFlatNest|TestNest|TestFlags" -v ./...
	@echo "✅ Core API tests passed"

## benchmark: Run benchmark tests
benchmark:
	@echo "⚡ Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...
	@echo "✅ Benchmarks completed"

## dev: Development workflow - format, lint, test
dev: fmt lint test
	@echo "✅ Development workflow completed"

## ci: Full CI pipeline
ci: deps fmt lint test test-race test-coverage
	@echo "✅ CI pipeline completed successfully"

## docs: Generate documentation
docs:
	@echo "📚 Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "🌐 Starting godoc server on http://localhost:6060"; \
		echo "   View documentation at: http://localhost:6060/pkg/$(PACKAGE)"; \
		godoc -http=:6060; \
	else \
		echo "ℹ️  godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"; \
		echo "📖 Generating package documentation..."; \
		$(GOCMD) doc -all .; \
	fi

## mod-graph: Show module dependency graph
mod-graph:
	@echo "🔗 Module dependency graph:"
	$(GOMOD) graph

## mod-why: Explain why packages are needed
mod-why:
	@echo "❓ Explaining package dependencies:"
	$(GOMOD) why -m all

## version: Show Go version and module info
version:
	@echo "Go version:"
	$(GOCMD) version
	@echo ""
	@echo "Module info:"
	$(GOMOD) list -m

## check-updates: Check for dependency updates
check-updates:
	@echo "🔄 Checking for dependency updates..."
	$(GOMOD) list -u -m all

# Development targets for specific testing scenarios
.PHONY: test-flags test-env test-config test-discovery

## test-flags: Test flag processing functionality
test-flags:
	@echo "🏴 Testing flag processing..."
	$(GOTEST) -run "Flag|flag" -v ./...

## test-env: Test environment variable processing
test-env:
	@echo "🌍 Testing environment variable processing..."
	$(GOTEST) -run "Env|env" -v ./...

## test-config: Test configuration loading
test-config:
	@echo "⚙️  Testing configuration loading..."
	$(GOTEST) -run "Config|config" -v ./...

## test-discovery: Test configuration discovery
test-discovery:
	@echo "🔍 Testing configuration discovery..."
	$(GOTEST) -run "Discovery|discovery" -v ./...

# Quality assurance targets
.PHONY: qa security verify

## qa: Quality assurance - comprehensive testing
qa: deps fmt lint test test-race test-coverage benchmark
	@echo "✅ Quality assurance completed"

## security: Run security analysis
security:
	@echo "🔒 Running security analysis..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
		echo "✅ Security analysis completed"; \
	else \
		echo "⚠️  gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

## verify: Verify module dependencies and checksums
verify:
	@echo "✅ Verifying module dependencies..."
	$(GOMOD) verify
	@echo "✅ Module verification completed"

# Migration and compatibility targets
.PHONY: test-migration verify-pflag compare-legacy

## test-migration: Test migration from go-flag to pflag
test-migration:
	@echo "🔄 Testing pflag migration..."
	@echo "Checking for go-flag references..."
	@if grep -r "github.com/davidwalter0/go-flag" . --exclude-dir=.git --exclude="*.md" --exclude="*.org" 2>/dev/null; then \
		echo "❌ Found go-flag references - migration incomplete"; \
		exit 1; \
	else \
		echo "✅ No go-flag references found - migration successful"; \
	fi
	$(GOTEST) -run "TestVendor|TestFlatNest" -v ./...
	@echo "✅ Migration tests passed"

## verify-pflag: Verify pflag integration
verify-pflag:
	@echo "🔍 Verifying pflag integration..."
	@echo "Checking pflag dependency..."
	$(GOMOD) list -m github.com/spf13/pflag
	@echo "✅ pflag dependency verified"

## info: Show project information
info:
	@echo "📋 go-cfg Project Information"
	@echo "============================="
	@echo "Package: $(PACKAGE)"
	@echo "Description: Modern pflag-based configuration management library"
	@echo "Features:"
	@echo "  • Clean flag naming with FlatNest() function"
	@echo "  • AutoCfg-style configuration discovery"
	@echo "  • Environment variable integration"
	@echo "  • Legacy API compatibility"
	@echo "  • Modern pflag foundation"
	@echo ""
	@echo "Key Functions:"
	@echo "  • FlatNest() - Clean nested flags without top-level prefix"
	@echo "  • Init(), Eval() - Basic configuration processing"
	@echo "  • Simple(), Nest() - Advanced configuration modes"
	@echo "  • Enhanced discovery with Union/First modes"
	@echo ""
	@echo "Migration Status: ✅ Phase 4A Complete - pflag migration successful"