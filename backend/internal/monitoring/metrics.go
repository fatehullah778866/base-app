package monitoring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Metrics struct {
	requestsTotal    int64
	requestsDuration map[string]time.Duration
	errorsTotal      int64
	logger           *zap.Logger
}

func NewMetrics(logger *zap.Logger) *Metrics {
	return &Metrics{
		requestsDuration: make(map[string]time.Duration),
		logger:           logger,
	}
}

func (m *Metrics) RecordRequest(method, path string, duration time.Duration, statusCode int) {
	m.requestsTotal++
	
	key := fmt.Sprintf("%s:%s", method, path)
	m.requestsDuration[key] = duration
	
	if statusCode >= 400 {
		m.errorsTotal++
	}
}

func (m *Metrics) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"requests_total": m.requestsTotal,
		"errors_total":   m.errorsTotal,
		"uptime_seconds": time.Since(startTime).Seconds(),
	}
}

var startTime = time.Now()

func MetricsMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(wrapped, r)
			
			duration := time.Since(start)
			metrics.RecordRequest(r.Method, r.URL.Path, duration, wrapped.statusCode)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (m *Metrics) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	metrics := m.GetMetrics()
	
	// Simple JSON response (can be enhanced with Prometheus format)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"metrics": metrics,
	})
}

