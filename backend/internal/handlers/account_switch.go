package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type AccountSwitchHandler struct {
	accountSwitchService *services.AccountSwitchService
	logger               *zap.Logger
}

func NewAccountSwitchHandler(accountSwitchService *services.AccountSwitchService, logger *zap.Logger) *AccountSwitchHandler {
	return &AccountSwitchHandler{
		accountSwitchService: accountSwitchService,
		logger:               logger,
	}
}

func (h *AccountSwitchHandler) SwitchAccount(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var req struct {
		SwitchedToUserID *uuid.UUID `json:"switched_to_user_id"`
		SwitchedToRole   *string    `json:"switched_to_role"`
		Reason           *string    `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	ipAddress := getIPAddress(r)
	userAgent := r.UserAgent()

	switchRecord, err := h.accountSwitchService.SwitchAccount(r.Context(), userID, req.SwitchedToUserID, req.SwitchedToRole, req.Reason, &ipAddress, &userAgent)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    switchRecord,
		"message": "Account switched successfully",
	})
}

func (h *AccountSwitchHandler) GetSwitchHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	history, err := h.accountSwitchService.GetSwitchHistory(r.Context(), userID, 20)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    history,
	})
}

