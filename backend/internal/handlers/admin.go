package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type AdminHandler struct {
	adminService *services.AdminService
	logger       *zap.Logger
}

func NewAdminHandler(adminService *services.AdminService, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		logger:       logger,
	}
}

type adminLoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceID   *string `json:"device_id"`
	DeviceName *string `json:"device_name"`
}

type statusUpdateRequest struct {
	Status string `json:"status"`
}

type createAdminRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type updateRequestStatusRequest struct {
	Status   string  `json:"status"`
	Feedback *string `json:"feedback"`
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req adminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	ip := getIPAddress(r)
	agent := r.UserAgent()

	user, session, err := h.adminService.Login(r.Context(), services.AdminLoginRequest{
		Email:      req.Email,
		Password:   req.Password,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		IPAddress:  &ip,
		UserAgent:  &agent,
	})
	if err != nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"admin": map[string]interface{}{
				"id":    user.ID.String(),
				"email": user.Email,
				"name":  user.Name,
				"role":  user.Role,
			},
			"session": map[string]interface{}{
				"token":         session.Token,
				"refresh_token": session.RefreshToken,
				"expires_at":    session.ExpiresAt.Format(time.RFC3339),
			},
		},
	})
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	users, err := h.adminService.ListUsers(r.Context(), search)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    users,
	})
}

func (h *AdminHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid user id")
		return
	}

	user, err := h.adminService.GetUser(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid user id")
		return
	}

	var req statusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	adminID := middleware.GetUserIDFromContext(r.Context())
	if err := h.adminService.SetUserStatus(r.Context(), adminID, userID, req.Status); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "status updated",
	})
}

func (h *AdminHandler) AddAdmin(w http.ResponseWriter, r *http.Request) {
	var req createAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	adminID := middleware.GetUserIDFromContext(r.Context())
	admin, err := h.adminService.AddAdmin(r.Context(), adminID, services.CreateAdminRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":    admin.ID.String(),
			"email": admin.Email,
			"name":  admin.Name,
		},
	})
}

func (h *AdminHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	limit := 200
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			limit = parsed
		}
	}
	logs, err := h.adminService.ListLogs(r.Context(), limit)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    logs,
	})
}

func (h *AdminHandler) ListRequests(w http.ResponseWriter, r *http.Request) {
	var status *string
	if raw := r.URL.Query().Get("status"); raw != "" {
		status = &raw
	}
	reqs, err := h.adminService.ListRequests(r.Context(), status)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    reqs,
	})
}

func (h *AdminHandler) UpdateRequestStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	var req updateRequestStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	adminID := middleware.GetUserIDFromContext(r.Context())
	updated, err := h.adminService.UpdateRequestStatus(r.Context(), adminID, requestID, req.Status, req.Feedback)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    updated,
	})
}

func (h *AdminHandler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	admins, err := h.adminService.ListAdmins(r.Context(), search)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    admins,
	})
}
