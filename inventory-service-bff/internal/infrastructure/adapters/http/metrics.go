package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// PrometheusMetrics holds Prometheus metric collectors.
type PrometheusMetrics struct {
	latency *prometheus.SummaryVec
}

// NewPrometheusMetrics creates a new PrometheusMetrics instance.
func NewPrometheusMetrics() *PrometheusMetrics {
	latency := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem:  "inventory_bff",
			Name:       "latency_duration_seconds",
			Help:       "latency duration distribution",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.01},
		},
		[]string{"method", "path"},
	)
	return &PrometheusMetrics{latency}
}

// middlewareMetrics returns a middleware that logs requests and records latency.
func middlewareMetrics(log *zap.Logger, p *PrometheusMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			defer func() {
				elapsed := time.Since(start)
				p.latency.WithLabelValues(
					r.Method,
					r.URL.Path,
				).Observe(elapsed.Seconds())
				log.Info("request",
					zap.String("request-id", r.Header.Get(middleware.RequestIDHeader)),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("query", r.URL.RawQuery),
					zap.String("ip", r.RemoteAddr),
					zap.String("user-agent", r.UserAgent()),
					zap.Int("status", ww.Status()),
					zap.Duration("latency", elapsed),
					zap.Int("bytes", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
