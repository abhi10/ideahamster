// Package config provides centralized configuration constants and settings.
// All magic numbers and configurable values should be defined here.
package config

import (
	"os"
	"time"
)

// ===========================================
// Server Configuration
// ===========================================

const (
	// DefaultPort is the fallback server port if PORT env var is not set
	DefaultPort = "3001"

	// RequestTimeout is the maximum duration for HTTP handlers
	RequestTimeout = 60 * time.Second

	// CompressionLevel for gzip middleware (1-9, higher = more compression)
	CompressionLevel = 5
)

// ===========================================
// Session & Cookie Configuration
// ===========================================

const (
	// EmailCookieName is the cookie key for storing verified email
	EmailCookieName = "voter_email"

	// CookieMaxAge is how long email verification cookies last (1 year)
	CookieMaxAge = 365 * 24 * 60 * 60
)

// ===========================================
// Rate Limiting Configuration
// ===========================================

const (
	// GlobalRateLimit is requests per minute for general endpoints
	GlobalRateLimit = 100

	// StrictRateLimit is requests per minute for sensitive endpoints (email verification)
	StrictRateLimit = 10

	// VoteRateLimit is votes per minute per IP
	VoteRateLimit = 5

	// SubmitRateLimit is idea submissions per minute per IP
	SubmitRateLimit = 2

	// RateLimitWindow is the time window for rate limiting
	RateLimitWindow = time.Minute
)

// ===========================================
// Security Configuration
// ===========================================

const (
	// IPBlockDuration is how long honeypot-triggered IPs are blocked
	IPBlockDuration = 24 * time.Hour
)

// SuspiciousPaths are endpoints that bots commonly scan for.
// Requests to these paths trigger IP blocking.
var SuspiciousPaths = []string{
	"/.env",
	"/.git",
	"/.git/config",
	"/wp-admin",
	"/wp-login.php",
	"/administrator",
	"/admin.php",
	"/phpmyadmin",
	"/phpMyAdmin",
	"/.well-known/security.txt",
	"/config.php",
	"/setup.php",
	"/install.php",
	"/backup",
	"/db",
	"/database",
}

// ===========================================
// Database Configuration
// ===========================================

const (
	// DBMaxConnections is the maximum connections in the pool
	DBMaxConnections = 10

	// DBMinConnections is the minimum connections to keep open
	DBMinConnections = 2
)

// ===========================================
// Idea Validation Configuration
// ===========================================

const (
	// TitleMinLength is the minimum characters for idea titles
	TitleMinLength = 5

	// TitleMaxLength is the maximum characters for idea titles
	TitleMaxLength = 80

	// DescriptionMinLength is the minimum characters for idea descriptions
	DescriptionMinLength = 20

	// DescriptionMaxLength is the maximum characters for idea descriptions
	DescriptionMaxLength = 500

	// MaxTags is the maximum number of tags per idea
	MaxTags = 5

	// TagMaxLength is the maximum characters per tag
	TagMaxLength = 20

	// VoteThreshold is the minimum votes required for an idea to be "eligible"
	VoteThreshold = 50
)

// ===========================================
// Environment Helpers
// ===========================================

// IsProduction returns true if running in production environment
func IsProduction() bool {
	env := os.Getenv("ENV")
	return env == "production" || env == "prod"
}

// GetEnv returns the environment variable value or a fallback default
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
