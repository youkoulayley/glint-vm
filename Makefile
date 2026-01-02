# glint-vm Makefile

# Build variables
BINARY_NAME := glint-vm
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go variables
GO := go
GOFLAGS := -trimpath
LDFLAGS :=

# Directories
BUILD_DIR := build
DIST_DIR := dist

# Targets
.PHONY: all build install test clean fmt lint help cross-compile

all: build ## Build the binary

build: ## Build for current platform
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(DIST_DIR)
	@echo 'package version' > internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo '// This file is auto-generated during build. Do not edit manually.' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'const (' >> internal/version/build.go
	@echo '	versionValue = "$(VERSION)"' >> internal/version/build.go
	@echo '	commitValue  = "$(COMMIT)"' >> internal/version/build.go
	@echo '	dateValue    = "$(BUILD_DATE)"' >> internal/version/build.go
	@echo ')' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'func getVersion() string {' >> internal/version/build.go
	@echo '	return versionValue' >> internal/version/build.go
	@echo '}' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'func getCommit() string {' >> internal/version/build.go
	@echo '	return commitValue' >> internal/version/build.go
	@echo '}' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'func getDate() string {' >> internal/version/build.go
	@echo '	return dateValue' >> internal/version/build.go
	@echo '}' >> internal/version/build.go
	$(GO) build $(GOFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) ./cmd/glint-vm
	@echo "✓ Built $(DIST_DIR)/$(BINARY_NAME)"

install: build ## Install to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo install -m 755 $(DIST_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installed $(BINARY_NAME)"

uninstall: ## Uninstall from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled $(BINARY_NAME)"

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "✓ Code formatted"

lint: build ## Run golangci-lint via glint-vm
	@echo "Running linter with glint-vm..."
	@eval "$$(./$(DIST_DIR)/$(BINARY_NAME) use v2.7.2)" && golangci-lint run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BINARY_NAME) $(BUILD_DIR) $(DIST_DIR) coverage.out coverage.html
	@echo "✓ Cleaned"

cross-compile: ## Build for all platforms
	@echo "Cross-compiling for all platforms..."
	@mkdir -p $(DIST_DIR)
	@echo 'package version' > internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo '// This file is auto-generated during build. Do not edit manually.' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'const (' >> internal/version/build.go
	@echo '	versionValue = "$(VERSION)"' >> internal/version/build.go
	@echo '	commitValue  = "$(COMMIT)"' >> internal/version/build.go
	@echo '	dateValue    = "$(BUILD_DATE)"' >> internal/version/build.go
	@echo ')' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'func getVersion() string {' >> internal/version/build.go
	@echo '	return versionValue' >> internal/version/build.go
	@echo '}' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'func getCommit() string {' >> internal/version/build.go
	@echo '	return commitValue' >> internal/version/build.go
	@echo '}' >> internal/version/build.go
	@echo '' >> internal/version/build.go
	@echo 'func getDate() string {' >> internal/version/build.go
	@echo '	return dateValue' >> internal/version/build.go
	@echo '}' >> internal/version/build.go

	@echo "Building for Linux amd64..."
	@GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/glint-vm

	@echo "Building for Linux arm64..."
	@GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/glint-vm

	@echo "Building for macOS amd64..."
	@GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/glint-vm

	@echo "Building for macOS arm64..."
	@GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/glint-vm

	@echo "Building for Windows amd64..."
	@GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/glint-vm

	@echo "✓ Cross-compilation complete. Binaries in $(DIST_DIR)/"

release: clean test cross-compile ## Create release builds
	@echo "Creating release archives..."
	@cd $(DIST_DIR) && \
	for binary in $(BINARY_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			tar czf "$$binary.tar.gz" "$$binary" && \
			echo "✓ Created $$binary.tar.gz"; \
		fi; \
	done
	@echo "✓ Release builds complete"

run: build ## Build and run
	./$(DIST_DIR)/$(BINARY_NAME)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download
	@echo "✓ Dependencies downloaded"

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	$(GO) mod tidy
	@echo "✓ go.mod tidied"

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	$(GO) mod verify
	@echo "✓ Dependencies verified"

help: ## Show this help
	@echo "glint-vm Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
