package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type NotificationService struct {
	notificationRepo repositories.NotificationRepository
	logger            *zap.Logger
}

func NewNotificationService(notificationRepo repositories.NotificationRepository, logger *zap.Logger) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(ctx context.Context, userID uuid.UUID, notificationType, title, message string, link *string, metadata *string) (*models.Notification, error) {
	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Link:      link,
		IsRead:    false,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// GetNotifications retrieves notifications for a user
func (s *NotificationService) GetNotifications(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit int) ([]*models.Notification, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	return s.notificationRepo.GetByUserID(ctx, userID, unreadOnly, limit)
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(ctx, id)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.Delete(ctx, id)
}

// GetUnreadCount returns the count of unread notifications
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.notificationRepo.GetUnreadCount(ctx, userID)
}

