package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/repositories"
	"base-app-service/pkg/auth"
	"base-app-service/pkg/errors"
)

type UserHandler struct {
	userRepo    repositories.UserRepository
	requestRepo repositories.AccessRequestRepository
	themeRepo   repositories.ThemeRepository
	sessionRepo repositories.SessionRepository
	logger      *zap.Logger
}

func NewUserHandler(
	userRepo repositories.UserRepository,
	requestRepo repositories.AccessRequestRepository,
	themeRepo repositories.ThemeRepository,
	sessionRepo repositories.SessionRepository,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userRepo:    userRepo,
		requestRepo: requestRepo,
		themeRepo:   themeRepo,
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":             user.ID.String(),
			"email":          user.Email,
			"email_verified": user.EmailVerified,
			"name":           user.Name,
			"first_name":     user.FirstName,
			"last_name":      user.LastName,
			"photo_url":      user.PhotoURL,
			"phone":          user.Phone,
			"phone_verified": user.PhoneVerified,
			"status":         user.Status,
			"role":           user.Role,
			"created_at":     user.CreatedAt,
			"updated_at":     user.UpdatedAt,
		},
	})
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	var req struct {
		Name      *string `json:"name"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Phone     *string `json:"phone"`
		PhotoURL  *string `json:"photo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.PhotoURL != nil {
		user.PhotoURL = req.PhotoURL
	}

	// Note: Password updates are not allowed here - use change password endpoint

	if err := h.userRepo.Update(r.Context(), user); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":         user.ID.String(),
			"email":      user.Email,
			"name":       user.Name,
			"updated_at": user.UpdatedAt,
		},
	})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Both current and new password are required")
		return
	}

	if !auth.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Current password is incorrect")
		return
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	changedAt := time.Now()
	if err := h.userRepo.UpdatePassword(r.Context(), userID, newHash, changedAt); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})
}

// ExportData lets a user download their data snapshot.
func (h *UserHandler) ExportData(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	requests, err := h.requestRepo.ListByUser(r.Context(), userID.String())
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	theme, _ := h.themeRepo.GetGlobalTheme(r.Context(), userID)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"profile":         user,
			"theme":           theme,
			"access_requests": requests,
		},
	})
}

// RequestDeletion performs a soft delete, revokes sessions, and schedules purge after retention.
func (h *UserHandler) RequestDeletion(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}
	if user.Status == "deleted" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Account deletion already requested")
		return
	}

	if err := h.userRepo.MarkDeleted(r.Context(), userID); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Revoke all active sessions for safety.
	if err := h.sessionRepo.RevokeAllForUser(r.Context(), userID); err != nil {
		h.logger.Warn("failed to revoke sessions after deletion request", zap.Error(err))
	}

	purgeAt := time.Now().Add(5 * 24 * time.Hour)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Account scheduled for deletion in 5 days. You will be signed out.",
		"data": map[string]interface{}{
			"purge_at": purgeAt.Format(time.RFC3339),
			"status":   "deleted",
		},
	})
}
