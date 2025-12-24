package middleware

import (
	"net/http"
	"strings"
)

func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Remove port from IP address
	if idx := strings.LastIndex(ip, "]:"); idx != -1 {
		ip = ip[1:idx] // Remove [ and :port
	} else if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return ip
}

