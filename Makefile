# Match-Me Project Makefile

# Variables
CLIENT_DIR = frontend
SERVER_DIR = server
DIST_DIR = $(CLIENT_DIR)/dist

# Default target
.DEFAULT_GOAL := help

# Build client for production
build-client:
	@echo "Building client for production..."
	cd $(CLIENT_DIR) && npm ci && npm run build
	@echo "Client built successfully! Output in $(DIST_DIR)/"

# Install client dependencies
install-client:
	@echo "Installing client dependencies..."
	cd $(CLIENT_DIR) && npm install

# Development server for client
dev-client:
	@echo "Starting client development server..."
	cd $(CLIENT_DIR) && npm i && npm run dev

# Build and run server
build-server:
	@echo "Building server..."
	cd $(SERVER_DIR) && go build -o bin/server .

# Run server
run-server:
	@echo "Running server..."
	cd $(SERVER_DIR) && go run .



# Full production build
build-all: build-client build-server
	@echo "Full build completed!"

# Development setup
dev-setup: install-client
	@echo "Development environment setup complete!"

# Help target
help:
	@echo "Available targets:"
	@echo "  build-client    - Build client for production (output: client/dist/)"
	@echo "  install-client  - Install client dependencies"  
	@echo "  dev-client      - Start client development server"
	@echo "  build-server    - Build server binary"
	@echo "  run-server      - Run server"
	@echo "  build-all       - Build both client and server"
	@echo "  dev-setup       - Setup development environment"
	@echo "  help            - Show this help message"

.PHONY: build-client install-client dev-client build-server run-server build-all dev-setup help