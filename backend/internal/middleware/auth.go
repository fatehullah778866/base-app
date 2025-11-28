package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/pkg/auth"
	"base-app-service/pkg/errors"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	SessionIDKey contextKey = "session_id"
	UserRoleKey  contextKey = "user_role"
)

func AuthMiddleware(jwtSecret string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
				return
			}

			token := parts[1]
			claims, err := auth.ValidateToken(token, jwtSecret)
			if err != nil {
				errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
				return
			}

			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
				return
			}

			sessionID, err := uuid.Parse(claims.SessionID)
			if err != nil {
				errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session ID")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, SessionIDKey, sessionID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}

func GetSessionIDFromContext(ctx context.Context) uuid.UUID {
	sessionID, ok := ctx.Value(SessionIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return sessionID
}

func GetUserRoleFromContext(ctx context.Context) string {
	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return ""
	}
	return role
}

// RequireRole ensures the authenticated user has a matching role.
func RequireRole(role string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRoleFromContext(r.Context())
			if userRole == "" {
				errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", "Missing role")
				return
			}
			if strings.ToLower(userRole) != strings.ToLower(role) {
				errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

