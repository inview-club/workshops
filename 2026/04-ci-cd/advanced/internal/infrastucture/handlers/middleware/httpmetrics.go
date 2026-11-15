package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"products/internal/infrastucture/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		method := r.Method
		status := rw.StatusCodeStr()
		metrics.HttpRequestTotal.WithLabelValues(method, path, status).Inc()
		metrics.HttpRequestsDurationHistogram.WithLabelValues(method, path).Observe(duration)
	})
}
