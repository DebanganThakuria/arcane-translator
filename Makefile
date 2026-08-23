# Arcane Translator
#
# backend/  Go API and SQLite store
# web/      SvelteKit frontend (TypeScript)

.PHONY: help install dev backend frontend check fmt build build-backend build-frontend run stop clean

help:
	@echo "Arcane Translator"
	@echo
	@echo "  make install   Install backend and frontend dependencies"
	@echo "  make dev       Run the API and the frontend dev server together"
	@echo "  make backend   Run only the Go API on :8088"
	@echo "  make frontend  Run only the frontend dev server on :8080"
	@echo "  make check     Vet the backend and typecheck the frontend"
	@echo "  make fmt       Format the backend with gofumpt and goimports"
	@echo "  make build     Build the binary and the frontend bundle"
	@echo "  make run       Build, then serve everything from the binary on :8088"
	@echo "  make stop      Stop running dev processes"
	@echo "  make clean     Remove dependencies and build output"

install:
	cd web && npm install
	cd backend && go mod download
	@echo "Dependencies installed. Copy .env.example to .env and set a provider."

# Backend first so the frontend has an API to talk to on first paint.
dev:
	@$(MAKE) backend
	@sleep 1
	cd web && npm run dev

backend:
	@cd backend && go run . &
	@echo "API on http://localhost:8088"

frontend:
	cd web && npm run dev

check:
	cd backend && go vet ./...
	cd web && npm run check

fmt:
	cd backend && go run mvdan.cc/gofumpt@latest -w . \
		&& go run golang.org/x/tools/cmd/goimports@latest -w -local backend .

build: build-backend build-frontend

build-backend:
	cd backend && go build -ldflags "-s -w" -o ../bin/arcane-translator .

build-frontend:
	cd web && npm run build

# One process serving the API and the built frontend, as a packaged install does.
run: build
	ARCANE_WEB_DIR=$(CURDIR)/web/build ./bin/arcane-translator

stop:
	@pkill -f "vite dev" || true
	@pkill -f "go run \." || true
	@pkill -f "arcane-translator" || true
	@echo "Stopped."

clean:
	rm -rf bin web/build web/.svelte-kit web/node_modules
	@echo "Clean. Your library in data/ was left alone."
