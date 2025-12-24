package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get or generate request ID
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			
			// Add to response header
			w.Header().Set("X-Request-ID", requestID)
			
			// Add to context
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

