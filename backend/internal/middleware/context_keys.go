package middleware

// Shared context key type for all middleware
type contextKey string

// Context keys used across middleware
const (
	UserIDKey    contextKey = "user_id"
	SessionIDKey contextKey = "session_id"
	UserRoleKey  contextKey = "user_role"
	RequestIDKey contextKey = "request_id"
)

