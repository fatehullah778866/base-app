package handlers

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/middleware"
	"base-app-service/internal/services"
	"base-app-service/pkg/errors"
)

type FileUploadHandler struct {
	fileService *services.FileService
	logger      *zap.Logger
}

func NewFileUploadHandler(fileService *services.FileService, logger *zap.Logger) *FileUploadHandler {
	return &FileUploadHandler{
		fileService: fileService,
		logger:      logger,
	}
}

func (h *FileUploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "No file provided")
		return
	}
	defer file.Close()

	fileInfo, err := h.fileService.UploadImage(r.Context(), file, header, userID)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    fileInfo,
	})
}

func (h *FileUploadHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	err := r.ParseMultipartForm(50 << 20) // 50MB for documents
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "No file provided")
		return
	}
	defer file.Close()

	fileInfo, err := h.fileService.UploadDocument(r.Context(), file, header, userID)
	if err != nil {
		errors.RespondError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    fileInfo,
	})
}

func (h *FileUploadHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "File path is required")
		return
	}

	file, err := h.fileService.GetFile(r.Context(), filePath)
	if err != nil {
		errors.RespondError(w, http.StatusNotFound, "NOT_FOUND", "File not found")
		return
	}
	defer file.Close()

	// Get file info
	stat, err := file.Stat()
	if err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Set headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", string(rune(stat.Size())))

	// Stream file
	io.Copy(w, file)
}

func (h *FileUploadHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		errors.RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid session")
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		errors.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "File path is required")
		return
	}

	if err := h.fileService.DeleteFile(r.Context(), filePath); err != nil {
		errors.RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "File deleted successfully",
	})
}

