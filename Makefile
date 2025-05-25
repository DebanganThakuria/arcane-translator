# Novel Translation App Makefile

.PHONY: help install start-frontend start dev stop clean build-backend

# Default target
help:
	@echo "Available commands:"
	@echo "  make install       - Install all dependencies"
	@echo "  make start         - Start both backend and frontend"
	@echo "  make dev           - Start both in development mode"
	@echo "  make start-frontend- Start only the frontend"
	@echo "  make start-backend - Build the Go backend"
	@echo "  make stop          - Stop all running processes"
	@echo "  make clean         - Clean all dependencies and build files"

# Install dependencies
install:
	@echo "Installing frontend dependencies..."
	npm install
	@echo "Installing backend dependencies..."
	cd backend_golang && go mod tidy
	@echo "All dependencies installed!"

# Start both backend and frontend
start: start-backend start-frontend

# Development mode (same as start for now)
dev: start

# Build backend
build-backend:
	@echo "Building Go backend..."
	cd backend_golang && go build -o ../bin/server main.go

# Start backend only
start-backend: build-backend
	@echo "Starting Go backend..."
	./bin/server &
	@echo "Backend started on http://localhost:8088"

# Start frontend only  
start-frontend:
	@echo "Starting React frontend..."
	npm run dev &
	@echo "Frontend started on http://localhost:5173"

# Stop all processes
stop:
	@echo "Stopping all processes..."
	@if [ -f server.pid ]; then \
		echo "Stopping Go backend (PID: $$(cat server.pid))"; \
		kill $$(cat server.pid) 2>/dev/null || true; \
		rm -f server.pid; \
	else \
		echo "No PID file found, attempting fallback method"; \
		pkill -f "server" || true; \
	fi
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
