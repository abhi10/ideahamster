package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// RequestLogger logs HTTP requests with timing information
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status
		wrapped := wrapResponseWriter(w)

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log request
		duration := time.Since(start)
		log.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			wrapped.status,
			duration,
			getClientIP(r),
		)
	})
}

// Recoverer recovers from panics and logs the stack trace
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log panic with stack trace
				log.Printf(
					"PANIC: %v\n%s",
					err,
					debug.Stack(),
				)

				// Return 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
