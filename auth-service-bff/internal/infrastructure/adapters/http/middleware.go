package httpadapter

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/beabys/wms/auth-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/beabys/wms/pkg/logger"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// corsMiddleware returns a middleware that sets CORS headers.
func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				allowed := false
				for _, o := range allowedOrigins {
					if o == "*" || o == origin {
						allowed = true
						break
					}
				}
				if allowed || len(allowedOrigins) == 0 {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// requestLoggerMiddleware logs each request with method, path, status, and duration.
func requestLoggerMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.Info("http request",
				logger.LogField{Key: "method", Value: r.Method},
				logger.LogField{Key: "path", Value: r.URL.Path},
				logger.LogField{Key: "status", Value: ww.Status()},
				logger.LogField{Key: "duration", Value: time.Since(start).String()},
			)
		})
	}
}

// jwtExtractorMiddleware extracts the Bearer token from the Authorization header
// and stores it in the request context.
func jwtExtractorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		ctx := context.WithValue(r.Context(), httpctx.ContextKeyJWT, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractBearerToken extracts a Bearer token from the Authorization header.
func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return parts[1]
}

// TimeoutMiddleware returns a middleware that sets a context deadline on each request.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
