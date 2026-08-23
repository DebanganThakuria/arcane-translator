package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/config"
	"backend/handler"

	_ "backend/provider/llm"
	_ "backend/provider/webscraper"
	_ "backend/repo"
	_ "backend/service"
)

func main() {
	cfg := config.Get()

	// Fail here, with a message that says what to set, rather than panicking
	// inside a package initialiser where nothing can report it usefully.
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v\n\nSee .env.example for the settings this server expects.", err)
	}

	server := setupServer(cfg)

	go func() {
		log.Printf("Arcane Translator listening on port %d", cfg.Port)
		log.Printf("Translating with the %s provider, database at %s", cfg.Provider, cfg.DBPath)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}

// setupServer configures and returns an HTTP server with the middleware stack.
func setupServer(cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	handler.RegisterHandlers(mux)

	if handler.RegisterStatic(mux, cfg.WebDir) {
		log.Printf("Serving the frontend from %s", cfg.WebDir)
	}

	// Applied in reverse order, so logging runs first and CORS runs last.
	var stack http.Handler = mux
	stack = handler.CorsMiddleware(stack)
	stack = handler.SecurityMiddleware(stack)
	stack = handler.ErrorHandlingMiddleware(stack)
	stack = handler.LoggingMiddleware(stack)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: stack,
		// Translating a chapter is a slow LLM call, so the write timeout is
		// generous while the header timeout stays short to shed slow-loris
		// connections early.
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      180 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
