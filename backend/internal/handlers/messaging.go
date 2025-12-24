package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type MessagingHandler struct {
	messagingService *services.MessagingService
	logger           *zap.Logger
}

func NewMessagingHandler(messagingService *services.MessagingService, logger *zap.Logger) *MessagingHandler {
	return &MessagingHandler{
		messagingService: messagingService,
		logger:           logger,
	}
}

func (h *MessagingHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var req struct {
		RecipientID uuid.UUID `json:"recipient_id" validate:"required"`
		Subject     *string   `json:"subject"`
		Content     string    `json:"content" validate:"required"`
		Metadata    *string   `json:"metadata"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	message, err := h.messagingService.SendMessage(r.Context(), userID, req.RecipientID, req.Subject, req.Content, req.Metadata)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    message,
	})
}

func (h *MessagingHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	conversations, err := h.messagingService.GetConversations(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    conversations,
	})
}

func (h *MessagingHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	conversationIDStr := r.URL.Query().Get("conversation_id")
	if conversationIDStr == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "conversation_id is required")
		return
	}

	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid conversation_id")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	messages, err := h.messagingService.GetMessages(r.Context(), conversationID, limit)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    messages,
	})
}

func (h *MessagingHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var req struct {
		MessageID uuid.UUID `json:"message_id" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.messagingService.MarkAsRead(r.Context(), req.MessageID); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Message marked as read",
	})
}

func (h *MessagingHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	count, err := h.messagingService.GetUnreadCount(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   count,
	})
}

