# Firefox Bookmarks Converter - Makefile
# This Makefile contains all useful commands for building, testing, and managing the project

.PHONY: help build build-all test test-cover test-integration benchmark clean install-deps example run-example deps-python

# Default target
help: ## Show this help message
	@echo "Firefox Bookmarks Converter - Available Commands"
	@echo "================================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make build        # Build for current platform"
	@echo "  make build-all    # Build for all platforms"
	@echo "  make test         # Run all tests"
	@echo "  make benchmark    # Compare Go vs Python performance"

# Build commands
build: ## Build binary for current platform
	@echo "Building Firefox Bookmarks Converter for current platform..."
	go build -ldflags "-s -w" -o firefox-bookmarks
	@echo "✓ Build complete: firefox-bookmarks"

build-all: ## Build binaries for all platforms
	@echo "Building Firefox Bookmarks Converter for multiple platforms..."
	@mkdir -p builds
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_windows_amd64.exe
	@echo "Building for Windows (arm64)..."
	GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_windows_arm64.exe
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_macos_intel
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_macos_arm64
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_linux_amd64
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o builds/firefox-bookmarks_linux_arm64
	@echo ""
	@echo "✓ Build complete! Binaries created:"
	@ls -l builds/

# Dependencies
install-deps: ## Install Go dependencies
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Go dependencies installed"

deps-python: ## Install Python dependencies for comparison scripts
	@echo "Installing Python dependencies..."
	@if command -v python3 >/dev/null 2>&1; then \
		pip3 install -r python/requirements.txt; \
		echo "✓ Python dependencies installed"; \
	else \
		echo "❌ Python3 not found. Please install Python 3.x"; \
		exit 1; \
	fi

# Testing commands
test: ## Run all tests
	@echo "Running all tests..."
	go test -v ./...
	@echo "✓ All tests passed"

test-cover: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -cover ./...
	@echo "Running detailed coverage for bookmarks package..."
	go test ./bookmarks -cover

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	go test ./bookmarks -run Integration -v
	go test ./tests -v

# Performance and benchmarking
benchmark: build deps-python ## Compare Go vs Python performance
	@echo "Firefox Bookmarks Converter - Performance Comparison"
	@echo "==================================================="
	@echo ""
	@if [ ! -f "example/test_bookmarks.json" ]; then \
		echo "❌ Error: Test file example/test_bookmarks.json not found"; \
		exit 1; \
	fi
	@echo "Testing with file: example/test_bookmarks.json"
	@echo ""
	@echo "Python version:"
	@if [ -f "python/ff_bookmarks.py" ]; then \
		time python3 python/ff_bookmarks.py example/test_bookmarks.json -o test_python_perf.html; \
		PYTHON_SIZE=$$(du -h test_python_perf.html | cut -f1); \
		echo "Output size: $$PYTHON_SIZE"; \
	else \
		echo "Python version (python/ff_bookmarks.py) not found"; \
	fi
	@echo ""
	@echo "Go version:"
	@if [ -f "firefox-bookmarks" ]; then \
		time ./firefox-bookmarks -o test_go_perf.html example/test_bookmarks.json; \
		GO_SIZE=$$(du -h test_go_perf.html | cut -f1); \
		echo "Output size: $$GO_SIZE"; \
	else \
		echo "Go version not found - run 'make build' first"; \
	fi
	@echo ""
	@if [ -f "test_python_perf.html" ] && [ -f "test_go_perf.html" ]; then \
		echo "Comparing outputs:"; \
		if diff test_python_perf.html test_go_perf.html > /dev/null; then \
			echo "✅ Outputs are identical"; \
		else \
			echo "❌ Outputs differ"; \
		fi; \
		rm -f test_python_perf.html test_go_perf.html; \
	fi
	@echo ""
	@echo "Binary sizes:"
	@if [ -f "firefox-bookmarks" ]; then \
		echo "Go binary: $$(du -h firefox-bookmarks | cut -f1)"; \
	fi
	@echo "Python script: $$(du -h python/ff_bookmarks.py | cut -f1)"

# Example and demo commands
example: build ## Run the example usage demonstration
	@echo "Running example usage demonstration..."
	@cd example && ./example_usage.sh

run-example: build ## Run converter with the test file
	@echo "Running converter with test file..."
	@if [ -f "example/test_bookmarks.json" ]; then \
		./firefox-bookmarks -o example/output.html example/test_bookmarks.json; \
		echo "✓ Conversion complete: example/output.html"; \
	else \
		echo "❌ Test file example/test_bookmarks.json not found"; \
	fi

# Utility commands
clean: ## Clean build artifacts and temporary files
	@echo "Cleaning build artifacts..."
	rm -f firefox-bookmarks
	rm -rf builds/
	rm -f test_*.html
	rm -f example/output.html
	@echo "✓ Clean complete"

format: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "✓ Code formatted"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "✓ Linting complete"; \
	else \
		echo "❌ golangci-lint not installed. Install with: brew install golangci-lint"; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "✓ Vet complete"

# Release commands
check: test vet ## Run all checks (tests + vet)
	@echo "✓ All checks passed"

release-build: clean build-all ## Clean build for release
	@echo "✓ Release build complete"

# Development helpers
dev: install-deps build test ## Full development setup
	@echo "✓ Development environment ready"

quick: build run-example ## Quick build and test
	@echo "✓ Quick test complete"

all: clean install-deps build-all test benchmark ## Do everything
	@echo "✓ Full build and test cycle complete"

# File size information
info: ## Show project information
	@echo "Firefox Bookmarks Converter - Project Info"
	@echo "=========================================="
	@echo "Go version: $$(go version)"
	@echo "Module: $$(head -1 go.mod)"
	@echo ""
	@if [ -f "firefox-bookmarks" ]; then \
		echo "Current binary size: $$(du -h firefox-bookmarks | cut -f1)"; \
	else \
		echo "Binary not built yet (run 'make build')"; \
	fi
	@echo ""
	@echo "Available build targets:"
	@ls -la builds/ 2>/dev/null || echo "No builds yet (run 'make build-all')"
