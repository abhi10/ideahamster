package database

import (
	"context"
	"fmt"
	"os"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global database connection pool
var DB *pgxpool.Pool

// Connect establishes a connection to the database
func Connect() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Set connection pool settings from config
	poolConfig.MaxConns = int32(config.DBMaxConnections)
	poolConfig.MinConns = int32(config.DBMinConnections)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = pool
	return nil
}

// Close closes the database connection pool
func Close() {
	if DB != nil {
		DB.Close()
	}
}
