.PHONY: help dev build test clean install templ css run

# Default target
help:
	@echo "🐹 Idea Hamster - Available Commands"
	@echo ""
	@echo "  make install    - Install all dependencies (Go, Node, Templ)"
	@echo "  make dev        - Run development server with hot reload"
	@echo "  make templ      - Generate Templ templates"
	@echo "  make css        - Build Tailwind CSS (minified)"
	@echo "  make css-watch  - Watch and rebuild CSS on changes"
	@echo "  make build      - Build production binary"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make run        - Run the application"
	@echo ""

# Install all dependencies
install:
	@echo "📦 Installing Go dependencies..."
	go mod download
	go install github.com/a-h/templ/cmd/templ@latest
	@echo "📦 Installing Node dependencies..."
	npm install
	@echo "✅ All dependencies installed!"

# Generate Templ templates
templ:
	@echo "🔨 Generating Templ templates..."
	templ generate

# Build Tailwind CSS (production)
css:
	@echo "🎨 Building Tailwind CSS..."
	npm run build:css

# Watch and rebuild CSS
css-watch:
	@echo "👀 Watching CSS changes..."
	npm run dev:css

# Development server (with hot reload)
dev: templ css
	@echo "🚀 Starting development server..."
	@echo "💡 Run 'make css-watch' in another terminal for CSS hot reload"
	go run cmd/server/main.go

# Build production binary
build: templ css
	@echo "🏗️  Building production binary..."
	go build -o bin/idea-hamster -ldflags="-s -w" cmd/server/main.go
	@echo "✅ Binary built: bin/idea-hamster"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf web/static/css/output.css
	rm -f **/*_templ.go
	@echo "✅ Cleaned!"

# Run the application
run: templ css
	@echo "🐹 Starting Idea Hamster..."
	go run cmd/server/main.go

# Database migrations (for future use)
migrate-up:
	@echo "⬆️  Running migrations..."
	@echo "Please run migrations/001_initial_schema.sql in your Supabase SQL editor"

# Format code
fmt:
	@echo "💅 Formatting code..."
	go fmt ./...
	templ fmt .

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Watch for changes and rebuild (requires air)
watch:
	@echo "👀 Watching for changes..."
	air
