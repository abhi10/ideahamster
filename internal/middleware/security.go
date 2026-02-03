package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/abhi10/ideahamster/internal/config"
	"github.com/go-chi/httprate"
	"github.com/justinas/nosurf"
)

// blockedIPs stores temporarily blocked IPs (in production, use Redis)
var (
	blockedIPs   = make(map[string]time.Time)
	blockedIPsMu sync.RWMutex
)

// RateLimiter returns a rate limiting middleware for general requests
func RateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(config.GlobalRateLimit, config.RateLimitWindow)
}

// StrictRateLimiter returns a stricter rate limiter for sensitive endpoints
func StrictRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(config.StrictRateLimit, config.RateLimitWindow)
}

// VoteRateLimiter specifically for vote endpoints
func VoteRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(config.VoteRateLimit, config.RateLimitWindow)
}

// SubmitRateLimiter specifically for idea submission
func SubmitRateLimiter() func(http.Handler) http.Handler {
	return httprate.LimitByIP(config.SubmitRateLimit, config.RateLimitWindow)
}

// CSRFProtection returns CSRF middleware configured for HTMX
func CSRFProtection() func(http.Handler) http.Handler {
	csrfHandler := nosurf.New
	return func(next http.Handler) http.Handler {
		handler := csrfHandler(next)

		// Configure for HTMX - look for token in header
		handler.SetBaseCookie(http.Cookie{
			Name:     "csrf_token",
			Path:     "/",
			HttpOnly: true,
			Secure:   config.IsProduction(),
			SameSite: http.SameSiteStrictMode,
		})

		// Custom failure handler
		handler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CSRF failure: %s %s from %s", r.Method, r.URL.Path, getClientIP(r))
			http.Error(w, "CSRF token mismatch", http.StatusForbidden)
		}))

		return handler
	}
}

// GetCSRFToken extracts the CSRF token for templates
func GetCSRFToken(r *http.Request) string {
	return nosurf.Token(r)
}

// Honeypot blocks requests to suspicious paths and blacklists the IP
func Honeypot(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		// Check if IP is blocked
		if isIPBlocked(clientIP) {
			log.Printf("Blocked IP attempted access: %s -> %s", clientIP, r.URL.Path)
			w.WriteHeader(http.StatusTeapot) // 418 - confuse the bot
			return
		}

		// Check for suspicious paths
		path := strings.ToLower(r.URL.Path)
		for _, suspicious := range config.SuspiciousPaths {
			if strings.HasPrefix(path, suspicious) || strings.Contains(path, suspicious) {
				log.Printf("Honeypot triggered: %s -> %s", clientIP, r.URL.Path)
				blockIP(clientIP)
				w.WriteHeader(http.StatusTeapot)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// isIPBlocked checks if an IP is blocked and cleans up expired blocks
func isIPBlocked(ip string) bool {
	blockedIPsMu.RLock()
	blockTime, isBlocked := blockedIPs[ip]
	blockedIPsMu.RUnlock()

	if !isBlocked {
		return false
	}

	// Check if block has expired
	if time.Since(blockTime) >= config.IPBlockDuration {
		// Block expired, remove it (using write lock with defer)
		blockedIPsMu.Lock()
		defer blockedIPsMu.Unlock()
		delete(blockedIPs, ip)
		return false
	}

	return true
}

// blockIP adds an IP to the blocked list
func blockIP(ip string) {
	blockedIPsMu.Lock()
	defer blockedIPsMu.Unlock()
	blockedIPs[ip] = time.Now()
}

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// XSS protection (legacy but still useful)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://unpkg.com; "+
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"img-src 'self' data: https://web.archive.org; "+
				"connect-src 'self'")

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the real client IP, accounting for proxies
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (from reverse proxies like Cloudflare)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// GetBlockedIPCount returns the number of currently blocked IPs (for monitoring)
func GetBlockedIPCount() int {
	blockedIPsMu.RLock()
	defer blockedIPsMu.RUnlock()
	return len(blockedIPs)
}
