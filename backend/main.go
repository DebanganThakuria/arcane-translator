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

	"backend/provider/webscraper"

	"backend/handler"
	_ "backend/provider/gemini"
	_ "backend/provider/webscraper"
	_ "backend/repo"
	_ "backend/service"
)

func main() {
	// Create a new HTTP server
	port := 8088
	server := setupServer()

	// Start the server in a goroutine
	go func() {
		log.Printf("🚀 Starting Arcane Translator server on port %d...", port)
		log.Printf("📚 Server ready to handle novel translation requests")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for the interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close the web scraper service connections
	webscraper.GetScraperService().ClosOpenCon()

	// Attempt a graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}

	log.Println("✅ Server gracefully stopped")
}

// setupServer configures and returns an HTTP server with enhanced middleware
func setupServer() *http.Server {
	// Create a new router
	mux := http.NewServeMux()

	// Register API handlers
	handler.RegisterHandlers(mux)

	// Create middleware stack
	var middlewareStack http.Handler = mux

	// Apply middleware in reverse order (last applied = first executed)
	middlewareStack = handler.CorsMiddleware(middlewareStack)
	middlewareStack = handler.SecurityMiddleware(middlewareStack)
	middlewareStack = handler.ErrorHandlingMiddleware(middlewareStack)
	middlewareStack = handler.LoggingMiddleware(middlewareStack)

	// Create and configure the server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8088),
		Handler:      middlewareStack,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		// Enhanced server configuration
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	return server
}
