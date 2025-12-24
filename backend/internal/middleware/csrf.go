package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

const CSRFTokenHeader = "X-CSRF-Token"
const CSRFHeader = "X-CSRF-Required"

func CSRFMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF for GET, HEAD, OPTIONS
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if CSRF is required (set by previous middleware or handler)
			if r.Header.Get(CSRFHeader) == "true" {
				token := r.Header.Get(CSRFTokenHeader)
				cookieToken, _ := r.Cookie("csrf_token")

				if token == "" || cookieToken == nil || token != cookieToken.Value {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"success":false,"error":{"code":"CSRF_TOKEN_INVALID","message":"Invalid CSRF token"}}`))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GenerateCSRFToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return base64.URLEncoding.EncodeToString(token)
}

func SetCSRFTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600, // 1 hour
	})
}

func RequireCSRF() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set(CSRFHeader, "true")
			next.ServeHTTP(w, r)
		})
	}
}

