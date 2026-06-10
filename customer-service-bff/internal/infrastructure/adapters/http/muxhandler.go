package httpadapter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/beabys/wms/customer-service-bff/internal/api/v1"
	"github.com/beabys/wms/customer-service-bff/internal/infrastructure/adapters/http/context"
	"github.com/beabys/wms/pkg/logger"
)

// NewMuxHandler creates a chi router with all middleware and registers the OpenAPI routes.
func NewMuxHandler(server v1.ServerInterface, log logger.Logger, allowedOrigins []string) (http.Handler, error) {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(requestLoggerMiddleware(log))
	r.Use(corsMiddleware(allowedOrigins))
	r.Use(TimeoutMiddleware(30 * time.Second))

	// JWT extractor: puts Bearer token into context
	r.Use(jwtExtractorMiddleware)

	// Health check endpoint (not in OpenAPI spec)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		successResponseJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Register all OpenAPI routes
	v1.HandlerWithOptions(server, v1.ChiServerOptions{
		BaseRouter: r,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			errorResponseJSON(w, http.StatusBadRequest, fmt.Errorf("%s", err.Error()))
		},
	})

	return r, nil
}

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
