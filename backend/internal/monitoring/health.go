package monitoring

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"base-app-service/internal/database"
)

type HealthChecker struct {
	db     *database.DB
	logger *zap.Logger
}

func NewHealthChecker(db *database.DB, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		db:     db,
		logger: logger,
	}
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func (hc *HealthChecker) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks:    make(map[string]string),
	}

	// Check database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := hc.db.PingContext(ctx); err != nil {
		status.Status = "unhealthy"
		status.Checks["database"] = "unhealthy: " + err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status.Checks["database"] = "healthy"
	}

	status.Checks["api"] = "healthy"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (hc *HealthChecker) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Checks:    make(map[string]string),
	}

	// Check database readiness
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := hc.db.PingContext(ctx); err != nil {
		status.Status = "not_ready"
		status.Checks["database"] = "not_ready: " + err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status.Checks["database"] = "ready"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (hc *HealthChecker) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "alive",
		"timestamp": time.Now(),
	})
}

