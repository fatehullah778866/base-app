package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

var validate = validator.New()

type AuthHandler struct {
	authService *services.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService *services.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

type SignupRequest struct {
	Email           string  `json:"email" validate:"required,email"`
	Password        string  `json:"password" validate:"required,min=8"`
	Name            string  `json:"name" validate:"required"`
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	Phone           *string `json:"phone"`
	MarketingConsent bool   `json:"marketing_consent"`
	TermsAccepted   bool    `json:"terms_accepted" validate:"required"`
	TermsVersion    string  `json:"terms_version" validate:"required"`
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		errors.RespondValidationError(w, err)
		return
	}

	// Extract headers
	signupSource := r.Header.Get("X-Product-Name")
	ipAddress := getIPAddress(r)
	userAgent := r.UserAgent()
	deviceID := r.Header.Get("X-Device-ID")
	deviceName := r.Header.Get("X-Device-Name")

	// Truncate fields to match database constraints
	// signup_source is VARCHAR(100)
	if len(signupSource) > 100 {
		signupSource = signupSource[:100]
	}
	// device_name is VARCHAR(255)
	if len(deviceName) > 255 {
		deviceName = deviceName[:255]
	}
	// device_id is VARCHAR(255)
	if len(deviceID) > 255 {
		deviceID = deviceID[:255]
	}

	var signupSourcePtr *string
	if signupSource != "" {
		signupSourcePtr = &signupSource
	}

	var deviceIDPtr *string
	if deviceID != "" {
		deviceIDPtr = &deviceID
	}

	var deviceNamePtr *string
	if deviceName != "" {
		deviceNamePtr = &deviceName
	}

	// Signup
	serviceReq := services.SignupRequest{
		Email:            req.Email,
		Password:         req.Password,
		Name:             req.Name,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Phone:            req.Phone,
		SignupSource:     signupSourcePtr,
		MarketingConsent: req.MarketingConsent,
		TermsAccepted:    req.TermsAccepted,
		TermsVersion:     req.TermsVersion,
		IPAddress:        &ipAddress,
		UserAgent:        &userAgent,
		DeviceID:         deviceIDPtr,
		DeviceName:       deviceNamePtr,
	}

	user, session, err := h.authService.Signup(r.Context(), serviceReq)
	if err != nil {
		if err.Error() == "email already exists" {
			errors.RespondError(w, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":             user.ID.String(),
				"email":          user.Email,
				"name":           user.Name,
				"email_verified": user.EmailVerified,
				"status":         user.Status,
			},
			"session": map[string]interface{}{
				"token":         session.Token,
				"refresh_token": *session.RefreshToken,
				"expires_at":   session.ExpiresAt.Format(time.RFC3339),
			},
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		errors.RespondValidationError(w, err)
		return
	}

	// Extract headers
	ipAddress := getIPAddress(r)
	userAgent := r.UserAgent()
	deviceID := r.Header.Get("X-Device-ID")
	deviceName := r.Header.Get("X-Device-Name")

	// Truncate fields to match database constraints
	// device_name is VARCHAR(255)
	if len(deviceName) > 255 {
		deviceName = deviceName[:255]
	}
	// device_id is VARCHAR(255)
	if len(deviceID) > 255 {
		deviceID = deviceID[:255]
	}

	var deviceIDPtr *string
	if deviceID != "" {
		deviceIDPtr = &deviceID
	}

	var deviceNamePtr *string
	if deviceName != "" {
		deviceNamePtr = &deviceName
	}

	serviceReq := services.LoginRequest{
		Email:      req.Email,
		Password:   req.Password,
		RememberMe: req.RememberMe,
		IPAddress:  &ipAddress,
		UserAgent:  &userAgent,
		DeviceID:   deviceIDPtr,
		DeviceName: deviceNamePtr,
	}

	user, session, isNewDevice, err := h.authService.Login(r.Context(), serviceReq)
	if err != nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials")
		return
	}

	deviceData := map[string]interface{}{}
	if session.DeviceID != nil {
		deviceData["id"] = *session.DeviceID
		deviceData["is_new_device"] = isNewDevice
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":             user.ID.String(),
				"email":          user.Email,
				"name":           user.Name,
				"email_verified": user.EmailVerified,
				"status":         user.Status,
			},
			"session": map[string]interface{}{
				"id":            session.ID.String(),
				"token":         session.Token,
				"refresh_token": *session.RefreshToken,
				"expires_at":    session.ExpiresAt.Format(time.RFC3339),
			},
			"device": deviceData,
		},
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := validate.Struct(req); err != nil {
		errors.RespondValidationError(w, err)
		return
	}

	session, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"token":      session.Token,
			"expires_at": session.ExpiresAt.Format(time.RFC3339),
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := middleware.GetSessionIDFromContext(r.Context())
	if sessionID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	var req struct {
		RevokeAllSessions bool `json:"revoke_all_sessions"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.authService.Logout(r.Context(), sessionID, req.RevokeAllSessions); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"message":         "Logged out successfully",
			"sessions_revoked": 1,
		},
	})
}

// Helper functions
func getIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	
	// Remove port from IP address (PostgreSQL INET type doesn't accept ports)
	// Handle IPv6 addresses like [::1]:57378
	if idx := strings.LastIndex(ip, "]:"); idx != -1 {
		ip = ip[1:idx] // Remove [ and :port
	} else if idx := strings.LastIndex(ip, ":"); idx != -1 {
		// Handle IPv4 addresses like 127.0.0.1:57378
		ip = ip[:idx]
	}
	
	return ip
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

