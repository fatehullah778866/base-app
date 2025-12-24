package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type DashboardRepository interface {
	Create(ctx context.Context, item *models.DashboardItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.DashboardItem, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, status string) ([]*models.DashboardItem, error)
	Update(ctx context.Context, item *models.DashboardItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

