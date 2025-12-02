package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"base-app-service/pkg/errors"
)

func ErrorRecovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered", zap.Any("error", err))
					errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
