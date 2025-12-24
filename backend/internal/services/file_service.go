package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FileService struct {
	uploadDir string
	maxSize   int64
	logger    *zap.Logger
}

type FileInfo struct {
	ID          uuid.UUID `json:"id"`
	OriginalName string   `json:"original_name"`
	StoredName   string   `json:"stored_name"`
	Path         string   `json:"path"`
	Size         int64    `json:"size"`
	ContentType  string   `json:"content_type"`
	UploadedBy   uuid.UUID `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type FileUploadConfig struct {
	UploadDir string
	MaxSize   int64 // in bytes
}

func NewFileService(config FileUploadConfig, logger *zap.Logger) *FileService {
	// Create upload directory if it doesn't exist
	if config.UploadDir == "" {
		config.UploadDir = "uploads"
	}
	os.MkdirAll(config.UploadDir, 0755)

	return &FileService{
		uploadDir: config.UploadDir,
		maxSize:   config.MaxSize,
		logger:    logger,
	}
}

var (
	allowedImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	allowedDocumentTypes = map[string]bool{
		"application/pdf": true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"text/plain": true,
	}
)

func (fs *FileService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID uuid.UUID, allowedTypes map[string]bool) (*FileInfo, error) {
	// Check file size
	if header.Size > fs.maxSize {
		return nil, errors.New("file size exceeds maximum allowed size")
	}

	// Check content type
	contentType := header.Header.Get("Content-Type")
	if allowedTypes != nil && !allowedTypes[contentType] {
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomString := hex.EncodeToString(randomBytes)
	storedName := fmt.Sprintf("%s_%d%s", randomString, time.Now().Unix(), ext)

	// Create file path
	filePath := filepath.Join(fs.uploadDir, storedName)

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	fileInfo := &FileInfo{
		ID:           uuid.New(),
		OriginalName: header.Filename,
		StoredName:   storedName,
		Path:         filePath,
		Size:         header.Size,
		ContentType:  contentType,
		UploadedBy:   userID,
		CreatedAt:    time.Now(),
	}

	fs.logger.Info("File uploaded", zap.String("file_id", fileInfo.ID.String()))
	return fileInfo, nil
}

func (fs *FileService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID uuid.UUID) (*FileInfo, error) {
	return fs.UploadFile(ctx, file, header, userID, allowedImageTypes)
}

func (fs *FileService) UploadDocument(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID uuid.UUID) (*FileInfo, error) {
	return fs.UploadFile(ctx, file, header, userID, allowedDocumentTypes)
}

func (fs *FileService) DeleteFile(ctx context.Context, filePath string) error {
	if !strings.HasPrefix(filePath, fs.uploadDir) {
		return errors.New("invalid file path")
	}

	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	fs.logger.Info("File deleted", zap.String("path", filePath))
	return nil
}

func (fs *FileService) GetFile(ctx context.Context, filePath string) (*os.File, error) {
	if !strings.HasPrefix(filePath, fs.uploadDir) {
		return nil, errors.New("invalid file path")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

