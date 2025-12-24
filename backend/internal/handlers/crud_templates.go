package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

// GetCRUDTemplates returns all available CRUD templates
func (h *AdminHandler) GetCRUDTemplates(w http.ResponseWriter, r *http.Request) {
	templates := services.GetCRUDTemplates()
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    templates,
	})
}

// GetCRUDTemplate returns a specific template by name
func (h *AdminHandler) GetCRUDTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["name"]

	template, err := services.GetTemplateByName(templateName)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	if template == nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Template not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    template,
	})
}

// CreateEntityFromTemplate creates a CRUD entity from a template
func (h *AdminHandler) CreateEntityFromTemplate(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	templateName := vars["name"]
	
	var req struct {
		DisplayName *string `json:"display_name"` // Optional override
		Description *string `json:"description"`  // Optional override
	}

	// Decode optional overrides
	json.NewDecoder(r.Body).Decode(&req)

	// Get template
	template, err := services.GetTemplateByName(templateName)
	if err != nil || template == nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Template not found")
		return
	}

	// Use template values or override
	displayName := template.DisplayName
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}

	description := template.Description
	if req.Description != nil {
		description = *req.Description
	}

	// Create entity from template
	entity, err := h.customCRUDService.CreateEntity(
		r.Context(),
		adminID,
		template.Name,
		displayName,
		&description,
		template.Schema,
	)

	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    entity,
		"message": "Entity created from template successfully",
	})
}

