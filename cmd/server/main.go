package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/abhi10/ideahamster/internal/handlers"
	"github.com/abhi10/ideahamster/internal/middleware"
	"github.com/abhi10/ideahamster/internal/repository"
	"github.com/abhi10/ideahamster/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// ===========================================
	// DEPENDENCY INJECTION
	// Initialize all dependencies without global state
	// ===========================================

	// Create repositories (swap implementations for testing)
	ideaRepo := repository.NewMockIdeaRepository()
	voteRepo := repository.NewMockVoteRepository()

	// Create services with injected dependencies
	ideaService := service.NewIdeaService(ideaRepo, voteRepo)

	// Create handler dependencies
	deps := &handlers.Dependencies{
		IdeaService: ideaService,
	}

	// ===========================================
	// ROUTER SETUP
	// ===========================================

	r := chi.NewRouter()

	// Security middleware (order matters!)
	r.Use(middleware.Honeypot)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.RateLimiter())
	r.Use(middleware.CSRFProtection())

	// Standard middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(chimiddleware.Compress(config.CompressionLevel))
	r.Use(chimiddleware.Timeout(config.RequestTimeout))

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Health check
	r.Get("/health", handlers.HandleHealth)

	// Public routes
	r.Get("/", handlers.HandleHome)

	// Leaderboard disabled for Phase 1
	r.Get("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	// API routes with dependencies
	r.Route("/api", func(r chi.Router) {
		r.With(middleware.VoteRateLimiter()).Post("/vote/{ideaID}", deps.HandleVote)
		r.With(middleware.StrictRateLimiter()).Post("/verify-email", handlers.HandleVerifyEmail)
	})

	// ===========================================
	// SERVER WITH GRACEFUL SHUTDOWN
	// ===========================================

	port := config.GetEnv("PORT", config.DefaultPort)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Channel for server errors
	serverErr := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Printf("Idea Hamster starting on http://localhost:%s", port)
		log.Printf("Environment: %s", config.GetEnv("ENV", "development"))
		if config.IsProduction() {
			log.Println("Production mode: Secure cookies enabled")
		}
		serverErr <- server.ListenAndServe()
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErr:
		if err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("Received signal %v, initiating graceful shutdown...", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v, forcing close", err)
			if err := server.Close(); err != nil {
				log.Fatalf("Server close error: %v", err)
			}
		}

		log.Println("Server stopped gracefully")
	}
}
