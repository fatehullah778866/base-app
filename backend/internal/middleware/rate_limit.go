package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
	GetRemaining(ctx context.Context, key string, limit int) (int, error)
}

type inMemoryRateLimiter struct {
	store map[string]*rateLimitEntry
}

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

func NewInMemoryRateLimiter() *inMemoryRateLimiter {
	rl := &inMemoryRateLimiter{
		store: make(map[string]*rateLimitEntry),
	}
	// Cleanup goroutine
	go rl.cleanup()
	return rl
}

func (rl *inMemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		for key, entry := range rl.store {
			if now.After(entry.resetTime) {
				delete(rl.store, key)
			}
		}
	}
}

func (rl *inMemoryRateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now()
	entry, exists := rl.store[key]

	if !exists || now.After(entry.resetTime) {
		rl.store[key] = &rateLimitEntry{
			count:     1,
			resetTime: now.Add(window),
		}
		return true, nil
	}

	if entry.count >= limit {
		return false, nil
	}

	entry.count++
	return true, nil
}

func (rl *inMemoryRateLimiter) GetRemaining(ctx context.Context, key string, limit int) (int, error) {
	entry, exists := rl.store[key]
	if !exists {
		return limit, nil
	}
	remaining := limit - entry.count
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}

func RateLimitMiddleware(limiter RateLimiter, limit int, window time.Duration, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get identifier (IP address or user ID)
			identifier := GetIPAddress(r)
			
			// If user is authenticated, use user ID instead
			userID := GetUserIDFromContext(r.Context())
			if userID.String() != "" {
				identifier = fmt.Sprintf("user:%s", userID.String())
			}

			allowed, err := limiter.Allow(r.Context(), identifier, limit, window)
			if err != nil {
				logger.Error("Rate limiter error", zap.Error(err))
				// On error, allow the request (fail open)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				remaining, _ := limiter.GetRemaining(r.Context(), identifier, limit)
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
				w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
				w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
				w.Header().Set("Retry-After", strconv.FormatInt(int64(window.Seconds()), 10))
				http.Error(w, `{"success":false,"error":{"code":"RATE_LIMIT_EXCEEDED","message":"Too many requests. Please try again later."}}`, http.StatusTooManyRequests)
				return
			}

			remaining, _ := limiter.GetRemaining(r.Context(), identifier, limit)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

			next.ServeHTTP(w, r)
		})
	}
}

// Per-endpoint rate limiting
func RateLimitByEndpoint(limiter RateLimiter, limits map[string]RateLimitConfig, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			method := r.Method
			key := fmt.Sprintf("%s:%s", method, path)

			// Find matching rate limit config
			var config RateLimitConfig
			found := false
			for pattern, cfg := range limits {
				if matchPattern(path, pattern) {
					config = cfg
					found = true
					break
				}
			}

			if !found {
				// Use default
				config = RateLimitConfig{Limit: 100, Window: 1 * time.Minute}
			}

			identifier := GetIPAddress(r)
			userID := GetUserIDFromContext(r.Context())
			if userID.String() != "" {
				identifier = fmt.Sprintf("user:%s", userID.String())
			}

			key = fmt.Sprintf("%s:%s", key, identifier)

			allowed, err := limiter.Allow(r.Context(), key, config.Limit, config.Window)
			if err != nil {
				logger.Error("Rate limiter error", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				remaining, _ := limiter.GetRemaining(r.Context(), key, config.Limit)
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.Limit))
				w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
				w.Header().Set("Retry-After", strconv.FormatInt(int64(config.Window.Seconds()), 10))
				http.Error(w, `{"success":false,"error":{"code":"RATE_LIMIT_EXCEEDED","message":"Too many requests. Please try again later."}}`, http.StatusTooManyRequests)
				return
			}

			remaining, _ := limiter.GetRemaining(r.Context(), key, config.Limit)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.Limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			next.ServeHTTP(w, r)
		})
	}
}

type RateLimitConfig struct {
	Limit  int
	Window time.Duration
}

func matchPattern(path, pattern string) bool {
	// Simple pattern matching - can be enhanced
	if pattern == path {
		return true
	}
	// Support wildcards
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(path) >= len(prefix) && path[:len(prefix)] == prefix
	}
	return false
}

