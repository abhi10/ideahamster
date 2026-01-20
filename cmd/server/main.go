package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abhi10/idea-hamster/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(60 * time.Second))

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Routes
	r.Get("/", handlers.HandleHome)
	r.Get("/leaderboard", handlers.HandleLeaderboard)

	// API Routes
	r.Post("/api/vote/{ideaID}", handlers.HandleVote)
	r.Post("/api/verify-email", handlers.HandleVerifyEmail)
	r.Get("/idea/{ideaID}/expand", handlers.HandleExpandIdea)

	// Start server
	log.Printf("🐹 Idea Hamster starting on http://localhost:%s", port)
	log.Printf("Environment: %s", getEnv("ENV", "development"))

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
