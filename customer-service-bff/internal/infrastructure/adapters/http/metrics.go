package httpadapter

import (
	"net/http"
	"time"

	"github.com/beabys/wms/pkg/logger"
)

type PrometheusMetrics struct {
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{}
}

func middlewareMetrics(log logger.Logger, p *PrometheusMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info("request",
				logger.LogField{Key: "method", Value: r.Method},
				logger.LogField{Key: "path", Value: r.URL.Path},
				logger.LogField{Key: "latency", Value: time.Since(start)},
			)
		})
	}
}
