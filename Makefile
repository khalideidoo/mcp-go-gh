.PHONY: all generate build test clean install help lint lint-fix

# Default target
all: generate build

# Generate Go code from YAML definitions
generate:
	@echo "Generating Go code from YAML definitions..."
	@go run tools/gen/main.go tools/gen/codegen.go tools/gen/parser.go tools/gen/templates.go tools/gen/types.go

# Build the MCP server binary
build:
	@echo "Building mcp-go-gh..."
	@mkdir -p bin
	@go build -o bin/mcp-go-gh cmd/mcp-go-gh/main.go
	@echo "Built bin/mcp-go-gh"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f internal/commands/generated/*_gen.go
	@echo "Cleaned build artifacts"

# Install the binary to GOPATH/bin
install: build
	@echo "Installing to $(shell go env GOPATH)/bin..."
	@cp bin/mcp-go-gh $(shell go env GOPATH)/bin/
	@echo "Installed mcp-go-gh"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running golangci-lint v2..."
	@golangci-lint run

# Run linter with auto-fix
lint-fix:
	@echo "Running golangci-lint v2 with auto-fix..."
	@golangci-lint run --fix

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -o bin/mcp-go-gh-linux-amd64 cmd/mcp-go-gh/main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/mcp-go-gh-darwin-amd64 cmd/mcp-go-gh/main.go
	@GOOS=darwin GOARCH=arm64 go build -o bin/mcp-go-gh-darwin-arm64 cmd/mcp-go-gh/main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/mcp-go-gh-windows-amd64.exe cmd/mcp-go-gh/main.go
	@echo "Built binaries for multiple platforms in bin/"

# Show help
help:
	@echo "Available targets:"
	@echo "  all         - Generate code and build (default)"
	@echo "  generate    - Generate Go code from YAML definitions"
	@echo "  build       - Build the MCP server binary"
	@echo "  test        - Run tests"
	@echo "  clean       - Remove build artifacts"
	@echo "  install     - Install binary to GOPATH/bin"
	@echo "  deps        - Install and tidy dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run golangci-lint v2"
	@echo "  lint-fix    - Run golangci-lint v2 with auto-fix"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  help        - Show this help message"
