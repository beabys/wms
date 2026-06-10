package httpadapter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	v1 "github.com/beabys/wms/auth-service-bff/internal/api/v1"
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
