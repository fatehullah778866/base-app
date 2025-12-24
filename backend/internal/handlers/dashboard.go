package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

var validateDashboard = validator.New()

type DashboardHandler struct {
	dashboardService *services.DashboardService
	logger           *zap.Logger
}

func NewDashboardHandler(dashboardService *services.DashboardService, logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		logger:           logger,
	}
}

type CreateDashboardItemRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Priority    int     `json:"priority"`
	Metadata    *string `json:"metadata"`
}

// CreateItem creates a new dashboard item
func (h *DashboardHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var req CreateDashboardItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := validateDashboard.Struct(req); err != nil {
		errors.RespondValidationError(w, err)
		return
	}

	item, err := h.dashboardService.CreateItem(r.Context(), userID, req.Title, req.Description, req.Category, req.Metadata)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    item,
	})
}

// GetItem retrieves a dashboard item by ID
func (h *DashboardHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "id parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id format")
		return
	}

	item, err := h.dashboardService.GetItem(r.Context(), id)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Item not found")
		return
	}

	// Verify ownership
	if item.UserID != userID {
		errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", "You don't have access to this item")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    item,
	})
}

// ListItems retrieves all dashboard items for a user
func (h *DashboardHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	status := r.URL.Query().Get("status") // Optional: active, archived, deleted

	items, err := h.dashboardService.ListItems(r.Context(), userID, status)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    items,
		"count":   len(items),
	})
}

// UpdateItem updates a dashboard item
func (h *DashboardHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "id parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id format")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	item, err := h.dashboardService.UpdateItem(r.Context(), id, userID, updates)
	if err != nil {
		if err.Error() == "unauthorized: item does not belong to user" {
			errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    item,
		"message": "Item updated successfully",
	})
}

// DeleteItem permanently deletes a dashboard item
func (h *DashboardHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "id parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id format")
		return
	}

	if err := h.dashboardService.DeleteItem(r.Context(), id, userID); err != nil {
		if err.Error() == "unauthorized: item does not belong to user" {
			errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Item deleted successfully",
	})
}

// SoftDeleteItem soft deletes a dashboard item
func (h *DashboardHandler) SoftDeleteItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "id parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id format")
		return
	}

	if err := h.dashboardService.SoftDeleteItem(r.Context(), id, userID); err != nil {
		if err.Error() == "unauthorized: item does not belong to user" {
			errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
			return
		}
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Item archived successfully",
	})
}

