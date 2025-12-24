package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type SearchHandler struct {
	searchService *services.SearchService
	logger        *zap.Logger
}

func NewSearchHandler(searchService *services.SearchService, logger *zap.Logger) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		logger:        logger,
	}
}

// Search handles search requests with advanced filtering
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	// Support both query parameters and JSON body
	var req services.SearchRequest

	// Try JSON body first
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			// Successfully parsed JSON body
		} else {
			// Fall back to query parameters
			h.parseQueryParams(r, &req)
		}
	} else {
		// Parse query parameters
		h.parseQueryParams(r, &req)
	}

	// Validate and set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.Type == "" {
		req.Type = "all"
	}

	result, err := h.searchService.Search(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Search error", zap.Error(err))
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// parseQueryParams parses search parameters from query string
func (h *SearchHandler) parseQueryParams(r *http.Request, req *services.SearchRequest) {
	req.Query = r.URL.Query().Get("q")
	req.Type = r.URL.Query().Get("type")
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			req.Limit = l
		}
	}
	
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			req.Offset = o
		}
	}

	// Location filters
	if location := r.URL.Query().Get("location"); location != "" {
		req.Location = &location
	}
	if country := r.URL.Query().Get("country"); country != "" {
		req.Country = &country
	}
	if city := r.URL.Query().Get("city"); city != "" {
		req.City = &city
	}

	// Date filters
	if dateFromStr := r.URL.Query().Get("date_from"); dateFromStr != "" {
		if t, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			req.DateFrom = &t
		}
	}
	if dateToStr := r.URL.Query().Get("date_to"); dateToStr != "" {
		if t, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			req.DateTo = &t
		}
	}

	// Category and status filters
	if category := r.URL.Query().Get("category"); category != "" {
		req.Category = &category
	}
	if status := r.URL.Query().Get("status"); status != "" {
		req.Status = &status
	}

	// Entity ID for custom CRUD search
	if entityIDStr := r.URL.Query().Get("entity_id"); entityIDStr != "" {
		if id, err := uuid.Parse(entityIDStr); err == nil {
			req.EntityID = &id
		}
	}
}

// GetSearchHistory retrieves user's search history
func (h *SearchHandler) GetSearchHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := h.searchService.GetSearchHistory(r.Context(), userID, limit)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    history,
	})
}

// ClearSearchHistory clears user's search history
func (h *SearchHandler) ClearSearchHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	if err := h.searchService.ClearSearchHistory(r.Context(), userID); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Search history cleared",
	})
}
