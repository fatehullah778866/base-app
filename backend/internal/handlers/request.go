package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type RequestHandler struct {
	requestService *services.RequestService
	logger         *zap.Logger
}

func NewRequestHandler(requestService *services.RequestService, logger *zap.Logger) *RequestHandler {
	return &RequestHandler{
		requestService: requestService,
		logger:         logger,
	}
}

type createRequestPayload struct {
	Title   *string `json:"title"`
	Details *string `json:"details"`
}

func (h *RequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	req, err := h.requestService.Create(r.Context(), userID, services.CreateRequestInput{
		Title:   body.Title,
		Details: body.Details,
	})
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    req,
	})
}

func (h *RequestHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	reqs, err := h.requestService.ListForUser(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    reqs,
	})
}
