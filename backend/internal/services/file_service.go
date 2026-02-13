package services

import (
	"bytes"
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

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/config"
)

type FileService struct {
	uploadDir string
	maxSize   int64
	logger    *zap.Logger
	s3Client  *s3.Client
	s3Bucket  string
	s3Region  string
	useS3     bool
}

type FileInfo struct {
	ID           uuid.UUID `json:"id"`
	OriginalName string    `json:"original_name"`
	StoredName   string    `json:"stored_name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadedBy   uuid.UUID `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type FileUploadConfig struct {
	UploadDir string
	MaxSize   int64 // in bytes
	S3        config.S3Config
}

func NewFileService(cfg FileUploadConfig, logger *zap.Logger) *FileService {
	if cfg.UploadDir == "" {
		cfg.UploadDir = "uploads"
	}
	os.MkdirAll(cfg.UploadDir, 0755)

	fs := &FileService{
		uploadDir: cfg.UploadDir,
		maxSize:   cfg.MaxSize,
		logger:    logger,
	}

	if cfg.S3.Enabled() {
		client, err := initS3Client(context.Background(), cfg.S3)
		if err != nil {
			logger.Warn("Failed to initialize S3 client; falling back to local uploads", zap.Error(err))
		} else {
			fs.s3Client = client
			fs.s3Bucket = cfg.S3.Bucket
			fs.s3Region = cfg.S3.Region
			fs.useS3 = true
		}
	}

	return fs
}

func initS3Client(ctx context.Context, cfg config.S3Config) (*s3.Client, error) {
	if !cfg.Enabled() {
		return nil, errors.New("incomplete S3 configuration")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg), nil
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
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"text/plain": true,
	}
)

func (fs *FileService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID uuid.UUID, allowedTypes map[string]bool) (*FileInfo, error) {
	if header.Size > fs.maxSize {
		return nil, errors.New("file size exceeds maximum allowed size")
	}

	contentType := header.Header.Get("Content-Type")
	if allowedTypes != nil && contentType != "" {
		if !allowedTypes[contentType] {
			return nil, fmt.Errorf("file type %s is not allowed", contentType)
		}
	}

	ext := filepath.Ext(header.Filename)
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomString := hex.EncodeToString(randomBytes)
	storedName := fmt.Sprintf("%s_%d%s", randomString, time.Now().Unix(), ext)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if int64(len(data)) > fs.maxSize {
		return nil, errors.New("file size exceeds maximum allowed size")
	}

	var filePath string

	if fs.useS3 {
		url, err := fs.uploadToS3(ctx, storedName, data, contentType)
		if err != nil {
			return nil, err
		}
		filePath = url
	} else {
		filePath = filepath.Join(fs.uploadDir, storedName)
		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		defer dst.Close()

		if _, err := dst.Write(data); err != nil {
			os.Remove(filePath)
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	fileInfo := &FileInfo{
		ID:           uuid.New(),
		OriginalName: header.Filename,
		StoredName:   storedName,
		Path:         filePath,
		Size:         int64(len(data)),
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

func (fs *FileService) uploadToS3(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	if fs.s3Client == nil {
		return "", errors.New("S3 client not configured")
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := fs.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(fs.s3Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", fs.s3Bucket, fs.s3Region, key), nil
}

func (fs *FileService) DeleteFile(ctx context.Context, filePath string) error {
	if fs.useS3 && fs.isS3URL(filePath) {
		return fs.deleteFromS3(ctx, filePath)
	}

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

func (fs *FileService) deleteFromS3(ctx context.Context, url string) error {
	if fs.s3Client == nil {
		return errors.New("S3 client not configured")
	}

	key := strings.TrimPrefix(url, fs.s3URLPrefix())
	if key == "" {
		return errors.New("invalid S3 URL")
	}

	_, err := fs.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(fs.s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	fs.logger.Info("File deleted", zap.String("path", url))
	return nil
}

func (fs *FileService) GetFile(ctx context.Context, filePath string) (*os.File, error) {
	if fs.useS3 && fs.isS3URL(filePath) {
		return fs.downloadFromS3(ctx, filePath)
	}

	if !strings.HasPrefix(filePath, fs.uploadDir) {
		return nil, errors.New("invalid file path")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fs *FileService) downloadFromS3(ctx context.Context, url string) (*os.File, error) {
	if fs.s3Client == nil {
		return nil, errors.New("S3 client not configured")
	}

	key := strings.TrimPrefix(url, fs.s3URLPrefix())
	if key == "" {
		return nil, errors.New("invalid S3 URL")
	}

	resp, err := fs.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(fs.s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp("", "download-*.tmp")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}

	if _, err := tempFile.Seek(0, io.SeekStart); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}

	return tempFile, nil
}

func (fs *FileService) isS3URL(path string) bool {
	return fs.useS3 && strings.HasPrefix(path, fs.s3URLPrefix())
}

func (fs *FileService) s3URLPrefix() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", fs.s3Bucket, fs.s3Region)
}
