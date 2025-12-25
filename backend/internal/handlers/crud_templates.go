package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"base-app-service/internal/middleware"
	"base-app-service/pkg/errors"
)

// GetCRUDTemplates returns all available CRUD templates from database
func (h *AdminHandler) GetCRUDTemplates(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	activeOnly := r.URL.Query().Get("active_only") == "true"

	var categoryPtr *string
	if category != "" {
		categoryPtr = &category
	}

	templates, err := h.crudTemplateService.ListTemplates(r.Context(), categoryPtr, activeOnly)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Convert to API format with schema as map
	var result []map[string]interface{}
	for _, template := range templates {
		var schema map[string]interface{}
		if err := json.Unmarshal([]byte(template.Schema), &schema); err != nil {
			continue
		}

		tmpl := map[string]interface{}{
			"id":          template.ID.String(),
			"name":        template.Name,
			"display_name": template.DisplayName,
			"description": template.Description,
			"schema":      schema,
			"icon":        template.Icon,
			"category":    template.Category,
			"created_by":  template.CreatedBy.String(),
			"is_active":   template.IsActive,
			"is_system":   template.IsSystem,
			"created_at":  template.CreatedAt,
			"updated_at":  template.UpdatedAt,
		}
		result = append(result, tmpl)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetCRUDTemplate returns a specific template by name
func (h *AdminHandler) GetCRUDTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["name"]

	template, err := h.crudTemplateService.GetTemplateByName(r.Context(), templateName)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	if template == nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Template not found")
		return
	}

	// Convert schema to map
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(template.Schema), &schema); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid template schema")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":          template.ID.String(),
			"name":        template.Name,
			"display_name": template.DisplayName,
			"description": template.Description,
			"schema":      schema,
			"icon":        template.Icon,
			"category":    template.Category,
			"created_by":  template.CreatedBy.String(),
			"is_active":   template.IsActive,
			"is_system":   template.IsSystem,
			"created_at":  template.CreatedAt,
			"updated_at":  template.UpdatedAt,
		},
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

	// Get template from database
	template, err := h.crudTemplateService.GetTemplateByName(r.Context(), templateName)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if template == nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Template not found")
		return
	}

	// Convert schema JSON to map
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(template.Schema), &schema); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid template schema")
		return
	}

	// Use template values or override
	displayName := template.DisplayName
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}

	description := ""
	if template.Description != nil {
		description = *template.Description
	}
	if req.Description != nil {
		description = *req.Description
	}
	descriptionPtr := &description
	if description == "" {
		descriptionPtr = nil
	}

	// Create entity from template
	entity, err := h.customCRUDService.CreateEntity(
		r.Context(),
		adminID,
		template.Name,
		displayName,
		descriptionPtr,
		schema,
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

// CreateTemplate creates a new CRUD template
func (h *AdminHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserIDFromContext(r.Context())
	
	var req struct {
		Name        string                 `json:"name" validate:"required"`
		DisplayName string                 `json:"display_name" validate:"required"`
		Description *string                `json:"description"`
		Schema      map[string]interface{} `json:"schema" validate:"required"`
		Icon        *string                `json:"icon"`
		Category    *string                `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	template, err := h.crudTemplateService.CreateTemplate(
		r.Context(),
		adminID,
		req.Name,
		req.DisplayName,
		req.Description,
		req.Schema,
		req.Icon,
		req.Category,
	)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	// Convert schema to map for response
	var schema map[string]interface{}
	json.Unmarshal([]byte(template.Schema), &schema)

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":          template.ID.String(),
			"name":        template.Name,
			"display_name": template.DisplayName,
			"description": template.Description,
			"schema":      schema,
			"icon":        template.Icon,
			"category":    template.Category,
			"created_by":  template.CreatedBy.String(),
			"is_active":   template.IsActive,
			"is_system":   template.IsSystem,
			"created_at":  template.CreatedAt,
			"updated_at":  template.UpdatedAt,
		},
	})
}

// UpdateTemplate updates an existing template
func (h *AdminHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid template id")
		return
	}

	var req struct {
		DisplayName *string                `json:"display_name"`
		Description *string                `json:"description"`
		Schema      map[string]interface{} `json:"schema"`
		Icon        *string                `json:"icon"`
		Category    *string                `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	template, err := h.crudTemplateService.UpdateTemplate(
		r.Context(),
		id,
		req.DisplayName,
		req.Description,
		req.Schema,
		req.Icon,
		req.Category,
	)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	// Convert schema to map for response
	var schema map[string]interface{}
	json.Unmarshal([]byte(template.Schema), &schema)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":          template.ID.String(),
			"name":        template.Name,
			"display_name": template.DisplayName,
			"description": template.Description,
			"schema":      schema,
			"icon":        template.Icon,
			"category":    template.Category,
			"created_by":  template.CreatedBy.String(),
			"is_active":   template.IsActive,
			"is_system":   template.IsSystem,
			"created_at":  template.CreatedAt,
			"updated_at":  template.UpdatedAt,
		},
	})
}

// DeleteTemplate deletes a template
func (h *AdminHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid template id")
		return
	}

	if err := h.crudTemplateService.DeleteTemplate(r.Context(), id); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Template deleted successfully",
	})
}

