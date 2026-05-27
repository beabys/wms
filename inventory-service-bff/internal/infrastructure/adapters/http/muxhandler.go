package http

import (
	"errors"
	"fmt"
	"net/http"

	v1 "github.com/beabys/wms/inventory-service-bff/internal/api/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewMuxHandler creates a new chi router with all middleware and routes.
func NewMuxHandler(server *HttpServer) (http.Handler, error) {
	swagger, err := v1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("error loading swagger spec: %w", err)
	}

	swagger.Servers = nil

	r := chi.NewRouter()

	prometheusMetrics := NewPrometheusMetrics()

	r.NotFound(notFound)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middlewareMetrics(server.Logger, prometheusMetrics))

	r.Group(func(r chi.Router) {
		r.Mount("/metrics", promhttp.Handler())
	})

	r.Group(func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders: []string{
				"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
				"Access-Control-Allow-Headers", "X-Requested-With",
				"Access-Control-Request-Method", "Access-Control-Request-Headers",
			},
			MaxAge: 300,
		}))

		v1.HandlerWithOptions(server, v1.ChiServerOptions{
			BaseRouter: r,
			ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				errorResponseJSON(w, http.StatusInternalServerError, err)
			},
		})
	})

	return r, nil
}

func notFound(w http.ResponseWriter, r *http.Request) {
	errorResponseJSON(w, http.StatusNotFound, errors.New("not found"))
}

// DefaultError is a default error handler.
func DefaultError(w http.ResponseWriter, r *http.Request, err error) {
	errorResponseJSON(w, http.StatusInternalServerError, err)
}

// JsonContentType sets Content-Type to application/json.
func JsonContentType(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
