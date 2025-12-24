package middleware

import (
	"net/http"
	"strconv"
)

func RequestSizeLimitMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set max request body size
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			
			// Check Content-Length header
			if r.ContentLength > maxSize {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				w.Write([]byte(`{"success":false,"error":{"code":"REQUEST_TOO_LARGE","message":"Request body too large. Maximum size: ` + strconv.FormatInt(maxSize, 10) + ` bytes"}}`))
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

