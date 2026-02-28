package middleware
package middleware

import (
    "net/http"
    "time"

    "github.com/rs/zerolog/log"
)

// responseWriter is a wrapper around http.ResponseWriter that captures the status code and response size.
type responseWriter struct {
    http.ResponseWriter
    status      int
    size        int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    size, err := rw.ResponseWriter.Write(b)
    rw.size += size
    return size, err
}

// LoggerMiddleware logs each HTTP request with method, path, status, duration, and size.
func LoggerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap the response writer to capture status and size
        wrapped := &responseWriter{
            ResponseWriter: w,
            status:         http.StatusOK, // default status if not set
        }

        // Process the request
        next.ServeHTTP(wrapped, r)

        // Log after the request is processed
        log.Info().
            Str("method", r.Method).
            Str("path", r.URL.Path).
            Str("remote_addr", r.RemoteAddr).
            Int("status", wrapped.status).
            Dur("duration", time.Since(start)).
            Int("size", wrapped.size).
            Str("user_agent", r.UserAgent()).
            Msg("HTTP request")
    })
}