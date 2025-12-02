package repositories

import (
	"context"

	"base-app-service/internal/models"
)

type AccessRequestRepository interface {
	Create(ctx context.Context, request *models.AccessRequest) error
	List(ctx context.Context, status *string) ([]*models.AccessRequest, error)
	ListByUser(ctx context.Context, userID string) ([]*models.AccessRequest, error)
	UpdateStatus(ctx context.Context, id string, status string, feedback *string) (*models.AccessRequest, error)
	GetByID(ctx context.Context, id string) (*models.AccessRequest, error)
}
