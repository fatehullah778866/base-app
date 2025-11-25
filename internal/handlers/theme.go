package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/models"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type ThemeHandler struct {
	themeService *services.ThemeService
	logger       *zap.Logger
}

func NewThemeHandler(themeService *services.ThemeService, logger *zap.Logger) *ThemeHandler {
	return &ThemeHandler{
		themeService: themeService,
		logger:       logger,
	}
}

func (h *ThemeHandler) GetTheme(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	productName := r.URL.Query().Get("product")

	var product *string
	if productName != "" {
		product = &productName
	}

	theme, err := h.themeService.GetTheme(r.Context(), userID, product)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	source := "global"
	var productOverride *string
	if product != nil {
		override, err := h.themeService.GetProductOverride(r.Context(), userID, *product)
		if err == nil && override != nil {
			source = "product_override"
			productOverride = product
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"theme":           theme.Theme,
			"contrast":        theme.Contrast,
			"text_direction": theme.TextDirection,
			"brand":           theme.Brand,
			"source":          source,
			"product_override": productOverride,
			"synced_at":       theme.SyncedAt.Format(time.RFC3339),
			"localStorage_keys": map[string]string{
				"theme":          "kompassui-theme",
				"contrast":       "kompassui-contrast",
				"text_direction": "kompassui-text-direction",
			},
		},
	})
}

func (h *ThemeHandler) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	var req struct {
		Theme         *string `json:"theme"`
		Contrast      *string `json:"contrast"`
		TextDirection *string `json:"text_direction"`
		Brand         *string `json:"brand"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	serviceReq := services.ThemeUpdateRequest{
		Theme:         req.Theme,
		Contrast:      req.Contrast,
		TextDirection: req.TextDirection,
		Brand:         req.Brand,
	}

	theme, err := h.themeService.UpdateTheme(r.Context(), userID, serviceReq)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"theme":          theme.Theme,
			"contrast":       theme.Contrast,
			"text_direction": theme.TextDirection,
			"brand":          theme.Brand,
			"synced_at":      theme.SyncedAt.Format(time.RFC3339),
			"message":        "Theme preferences updated successfully",
		},
	})
}

func (h *ThemeHandler) SyncTheme(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	var req struct {
		Theme         string  `json:"theme"`
		Contrast      string  `json:"contrast"`
		TextDirection string  `json:"text_direction"`
		Brand         *string `json:"brand"`
		ClientTimestamp string `json:"client_timestamp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	clientTheme := &models.ThemePreferences{
		UserID:        userID,
		Theme:         req.Theme,
		Contrast:      req.Contrast,
		TextDirection: req.TextDirection,
		Brand:         req.Brand,
	}

	serverTheme, conflicts, err := h.themeService.SyncTheme(r.Context(), userID, clientTheme)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	synced := len(conflicts) == 0

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"synced": synced,
			"server_theme": map[string]interface{}{
				"theme":          serverTheme.Theme,
				"contrast":       serverTheme.Contrast,
				"text_direction": serverTheme.TextDirection,
				"brand":          serverTheme.Brand,
				"synced_at":      serverTheme.SyncedAt.Format(time.RFC3339),
			},
			"conflicts": conflicts,
		},
	})
}

