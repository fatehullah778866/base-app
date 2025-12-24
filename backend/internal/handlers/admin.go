package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type AdminHandler struct {
	adminService         *services.AdminService
	adminSettingsService *services.AdminSettingsService
	customCRUDService    *services.CustomCRUDService
	logger               *zap.Logger
}

func NewAdminHandler(adminService *services.AdminService, adminSettingsService *services.AdminSettingsService, customCRUDService *services.CustomCRUDService, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		adminService:         adminService,
		adminSettingsService: adminSettingsService,
		customCRUDService:    customCRUDService,
		logger:               logger,
	}
}

type adminLoginRequest struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
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

// Admin Settings handlers
func (h *AdminHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserIDFromContext(r.Context())
	settings, err := h.adminSettingsService.GetSettings(r.Context(), adminID)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    settings,
	})
}

func (h *AdminHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserIDFromContext(r.Context())
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	if err := h.adminSettingsService.UpdateSettings(r.Context(), adminID, updates); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Settings updated successfully",
	})
}

// Custom CRUD Entity handlers
func (h *AdminHandler) CreateCRUDEntity(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserIDFromContext(r.Context())
	var req struct {
		EntityName  string                 `json:"entity_name" validate:"required"`
		DisplayName string                 `json:"display_name" validate:"required"`
		Description *string                `json:"description"`
		Schema      map[string]interface{} `json:"schema" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	entity, err := h.customCRUDService.CreateEntity(r.Context(), adminID, req.EntityName, req.DisplayName, req.Description, req.Schema)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    entity,
	})
}

func (h *AdminHandler) ListCRUDEntities(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"
	entities, err := h.customCRUDService.ListEntities(r.Context(), nil, activeOnly)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    entities,
	})
}

func (h *AdminHandler) GetCRUDEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id")
		return
	}
	entity, err := h.customCRUDService.GetEntity(r.Context(), id)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    entity,
	})
}

func (h *AdminHandler) UpdateCRUDEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id")
		return
	}
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	entity, err := h.customCRUDService.UpdateEntity(r.Context(), id, updates)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    entity,
	})
}

func (h *AdminHandler) DeleteCRUDEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id")
		return
	}
	if err := h.customCRUDService.DeleteEntity(r.Context(), id); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Entity deleted successfully",
	})
}

// Custom CRUD Data handlers
func (h *AdminHandler) CreateCRUDData(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	entityID, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid entity id")
		return
	}
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	crudData, err := h.customCRUDService.CreateData(r.Context(), entityID, adminID, data)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    crudData,
	})
}

func (h *AdminHandler) ListCRUDData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entityID, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid entity id")
		return
	}
	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}
	data, err := h.customCRUDService.ListData(r.Context(), entityID, limit, offset)
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (h *AdminHandler) GetCRUDData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id")
		return
	}
	data, err := h.customCRUDService.GetData(r.Context(), id)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (h *AdminHandler) UpdateCRUDData(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserIDFromContext(r.Context())
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id")
		return
	}
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	crudData, err := h.customCRUDService.UpdateData(r.Context(), id, adminID, data)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    crudData,
	})
}

func (h *AdminHandler) DeleteCRUDData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid id")
		return
	}
	if err := h.customCRUDService.DeleteData(r.Context(), id); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Data deleted successfully",
	})
}

// VerifyAdminCode - Public endpoint to verify admin creation code (no auth required)
func (h *AdminHandler) VerifyAdminCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VerificationCode string `json:"verification_code" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	
	// Get expected code
	defaultCode := "Kompasstech2025@"
	expectedCode := defaultCode
	
	systemCode, err := h.adminSettingsService.GetSystemVerificationCode(r.Context())
	if err == nil {
		expectedCode = systemCode
	}
	
	if req.VerificationCode != expectedCode {
		errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", "Invalid verification code")
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Verification code is valid",
	})
}

// CreateAdminPublic - Public endpoint for creating admin with verification code (no auth required)
func (h *AdminHandler) CreateAdminPublic(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,min=8"`
		Name             string `json:"name" validate:"required"`
		VerificationCode string `json:"verification_code" validate:"required"` // Required for public admin creation
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	
	// Verify code
	defaultCode := "Kompasstech2025@"
	expectedCode := defaultCode
	
	systemCode, err := h.adminSettingsService.GetSystemVerificationCode(r.Context())
	if err == nil {
		expectedCode = systemCode
	}
	
	if req.VerificationCode != expectedCode {
		errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", "Invalid verification code")
		return
	}
	
	// Create admin user
	user, err := h.adminService.CreateUser(r.Context(), uuid.Nil, services.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Role:     "admin",
		Status:   "active",
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			errors.RespondError(w, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":    user.ID.String(),
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// Enhanced User CRUD handlers
func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,min=8"`
		Name             string `json:"name" validate:"required"`
		Role             string `json:"role"`
		Status           string `json:"status"`
		VerificationCode string `json:"verification_code"` // For admin creation
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	
	// If creating admin, verify code (only for public admin creation)
	if req.Role == "admin" {
		adminID := middleware.GetUserIDFromContext(r.Context())
		if adminID == uuid.Nil {
			// Public admin creation - check verification code
			// Try to get verification code from first admin's settings, or use default
			defaultCode := "Kompasstech2025@"
			expectedCode := defaultCode
			
			// Try to get from any admin's settings (use first admin found)
			// For simplicity, we'll use the default code
			// In production, you might want to query for the first admin and get their settings
			systemCode, err := h.adminSettingsService.GetSystemVerificationCode(r.Context())
			if err == nil {
				expectedCode = systemCode
			}
			
			if req.VerificationCode != expectedCode {
				errors.RespondError(w, http.StatusForbidden, "FORBIDDEN", "Invalid verification code")
				return
			}
		}
		// If adminID != uuid.Nil, an admin is logged in creating another admin - no verification needed
	}
	
	if req.Role == "" {
		req.Role = "user"
	}
	if req.Status == "" {
		req.Status = "active"
	}
	
	// Create user via admin service
	actorID := middleware.GetUserIDFromContext(r.Context())
	user, err := h.adminService.CreateUser(r.Context(), actorID, services.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Role:     req.Role,
		Status:   req.Status,
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			errors.RespondError(w, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":    user.ID.String(),
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid user id")
		return
	}
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User update endpoint - implement via admin service",
		"user_id": userID.String(),
	})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid user id")
		return
	}
	adminID := middleware.GetUserIDFromContext(r.Context())
	if err := h.adminService.SetUserStatus(r.Context(), adminID, userID, "deleted"); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

func (h *AdminHandler) GetUserSessions(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Get user sessions endpoint - implement via session repository",
	})
}

func (h *AdminHandler) RevokeUserSessions(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Revoke user sessions endpoint - implement via session repository",
	})
}
