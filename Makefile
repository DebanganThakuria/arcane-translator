
# Novel Translation App Makefile

.PHONY: help install start-backend start-frontend start dev stop clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make install       - Install all dependencies"
	@echo "  make start         - Start both backend and frontend"
	@echo "  make dev           - Start both in development mode"
	@echo "  make start-backend - Start only the backend"
	@echo "  make start-frontend- Start only the frontend"
	@echo "  make stop          - Stop all running processes"
	@echo "  make clean         - Clean all dependencies and build files"

# Install dependencies
install:
	@echo "Installing frontend dependencies..."
	npm install
	@echo "Installing backend dependencies..."
	cd backend && pip install -r requirements.txt
	@echo "All dependencies installed!"

# Start both backend and frontend
start: start-backend start-frontend

# Development mode (same as start for now)
dev: start

# Start backend only
start-backend:
	@echo "Starting Python backend..."
	cd backend && python main.py &
	@echo "Backend started on http://localhost:8000"

# Start frontend only  
start-frontend:
	@echo "Starting React frontend..."
	npm run dev &
	@echo "Frontend started on http://localhost:8080"

# Stop all processes
stop:
	@echo "Stopping all processes..."
	pkill -f "python main.py" || true
	pkill -f "npm run dev" || true
	pkill -f "vite" || true
	@echo "All processes stopped!"

# Clean everything
clean:
	@echo "Cleaning frontend..."
	rm -rf node_modules
	rm -rf dist
	@echo "Cleaning backend..."
	cd backend && find . -name "__pycache__" -type d -exec rm -rf {} + || true
	cd backend && find . -name "*.pyc" -delete || true
	cd backend && rm -f novels.db || true
	@echo "Cleanup complete!"
