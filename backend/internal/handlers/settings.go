package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

var validateSettings = validator.New()

type SettingsHandler struct {
	settingsService *services.SettingsService
	sessionRepo     repositories.SessionRepository
	logger          *zap.Logger
}

func NewSettingsHandler(settingsService *services.SettingsService, sessionRepo repositories.SessionRepository, logger *zap.Logger) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		sessionRepo:     sessionRepo,
		logger:          logger,
	}
}

// GetSettings retrieves all user settings
func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	settings, err := h.settingsService.GetSettings(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    settings,
	})
}

// GetActiveSessions retrieves all active sessions for the user
func (h *SettingsHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	sessions, err := h.sessionRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Remove sensitive data
	safeSessions := make([]map[string]interface{}, len(sessions))
	for i, s := range sessions {
		safeSessions[i] = map[string]interface{}{
			"id":              s.ID.String(),
			"device_name":     s.DeviceName,
			"device_type":     s.DeviceType,
			"os":              s.OS,
			"browser":         s.Browser,
			"ip_address":      s.IPAddress,
			"location_country": s.LocationCountry,
			"location_city":   s.LocationCity,
			"created_at":      s.CreatedAt,
			"last_used_at":    s.LastUsedAt,
			"is_current":      false, // Will be set by frontend based on current token
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    safeSessions,
	})
}

// LogoutAllDevices logs out from all devices
func (h *SettingsHandler) LogoutAllDevices(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	if err := h.sessionRepo.RevokeAllForUser(r.Context(), userID); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged out from all devices successfully",
	})
}

// UpdateProfileSettings updates profile settings
func (h *SettingsHandler) UpdateProfileSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.settingsService.UpdateProfileSettings(r.Context(), userID, updates); err != nil {
		if err.Error() == "username already taken" {
			errors.RespondError(w, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile settings updated successfully",
	})
}

// UpdateSecuritySettings updates security settings
func (h *SettingsHandler) UpdateSecuritySettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.settingsService.UpdateSecuritySettings(r.Context(), userID, updates); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Security settings updated successfully",
	})
}

// UpdatePrivacySettings updates privacy settings
func (h *SettingsHandler) UpdatePrivacySettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.settingsService.UpdatePrivacySettings(r.Context(), userID, updates); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Privacy settings updated successfully",
	})
}

// UpdateNotificationSettings updates notification settings
func (h *SettingsHandler) UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.settingsService.UpdateNotificationSettings(r.Context(), userID, updates); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Notification settings updated successfully",
	})
}

// UpdateAccountPreferences updates account preferences
func (h *SettingsHandler) UpdateAccountPreferences(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.settingsService.UpdateAccountPreferences(r.Context(), userID, updates); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Account preferences updated successfully",
	})
}

// AddConnectedAccount adds a connected third-party account
func (h *SettingsHandler) AddConnectedAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var account models.ConnectedAccount
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.settingsService.AddConnectedAccount(r.Context(), userID, account); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Connected account added successfully",
	})
}

// RemoveConnectedAccount removes a connected account
func (h *SettingsHandler) RemoveConnectedAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var req struct {
		Provider string `json:"provider" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := validateSettings.Struct(req); err != nil {
		errors.RespondValidationError(w, err)
		return
	}

	if err := h.settingsService.RemoveConnectedAccount(r.Context(), userID, req.Provider); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Connected account removed successfully",
	})
}

// RequestAccountDeletion schedules account deletion
func (h *SettingsHandler) RequestAccountDeletion(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var req struct {
		DaysUntilDeletion int `json:"days_until_deletion" validate:"required,min=1,max=30"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := validateSettings.Struct(req); err != nil {
		errors.RespondValidationError(w, err)
		return
	}

	if err := h.settingsService.RequestAccountDeletion(r.Context(), userID, req.DaysUntilDeletion); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Account deletion scheduled successfully",
	})
}

// DeactivateAccount temporarily deactivates account
func (h *SettingsHandler) DeactivateAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	if err := h.settingsService.DeactivateAccount(r.Context(), userID); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Account deactivated successfully",
	})
}

// ReactivateAccount reactivates a deactivated account
func (h *SettingsHandler) ReactivateAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	if err := h.settingsService.ReactivateAccount(r.Context(), userID); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Account reactivated successfully",
	})
}
