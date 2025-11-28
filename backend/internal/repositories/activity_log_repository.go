package repositories

import (
	"context"

	"base-app-service/internal/models"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *models.ActivityLog) error
	List(ctx context.Context, limit int) ([]*models.ActivityLog, error)
}
