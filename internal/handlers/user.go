package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/repositories"
	"base-app-service/pkg/errors"
)

type UserHandler struct {
	userRepo repositories.UserRepository
	logger   *zap.Logger
}

func NewUserHandler(userRepo repositories.UserRepository, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":             user.ID.String(),
			"email":          user.Email,
			"email_verified": user.EmailVerified,
			"name":           user.Name,
			"first_name":     user.FirstName,
			"last_name":      user.LastName,
			"photo_url":      user.PhotoURL,
			"phone":          user.Phone,
			"phone_verified": user.PhoneVerified,
			"status":         user.Status,
			"created_at":     user.CreatedAt,
			"updated_at":     user.UpdatedAt,
		},
	})
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID.String() == "00000000-0000-0000-0000-000000000000" {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	var req struct {
		Name      *string `json:"name"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Phone     *string `json:"phone"`
		PhotoURL  *string `json:"photo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.PhotoURL != nil {
		user.PhotoURL = req.PhotoURL
	}

	// Note: Password updates are not allowed here - use change password endpoint

	if err := h.userRepo.Update(r.Context(), user); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":        user.ID.String(),
			"email":     user.Email,
			"name":      user.Name,
			"updated_at": user.UpdatedAt,
		},
	})
}

