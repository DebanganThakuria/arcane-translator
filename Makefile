# Novel Translation App Makefile

.PHONY: help install start-frontend dev stop clean build-backend

# Default target
help:
	@echo "Available commands:"
	@echo "  make install       - Install all dependencies"
	@echo "  make dev           - Start both in development mode"
	@echo "  make start-frontend- Start only the frontend"
	@echo "  make build-backend - Build the Go backend"
	@echo "  make stop          - Stop all running processes"
	@echo "  make clean         - Clean all dependencies and build files"

# Install dependencies
install:
	@echo "Installing frontend dependencies..."
	npm install
	@echo "Installing backend dependencies..."
	cd backend && go mod tidy
	@echo "All dependencies installed!"

# Development mode (same as start for now)
dev: start

# Build backend
build-backend:
	@echo "Building Go backend..."
	cd backend && go build -o ../bin/server main.go

# Start frontend only  
start-frontend: stop
	@echo "Starting React frontend..."
	npm run dev &
	@echo "Frontend started on http://localhost:5173"

# Stop all processes
stop:
	@echo "Stopping all processes..."
	pkill -f "npm run dev" || true
	pkill -f "vite" || true
	@echo "All processes stopped!"

# Clean everything
clean:
	@echo "Cleaning frontend..."
	rm -rf node_modules
	rm -rf dist
	@echo "Cleaning backend..."
	rm -rf bin
	@echo "Cleanup complete!"
