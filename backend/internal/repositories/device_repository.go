package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error)
	GetByDeviceID(ctx context.Context, userID uuid.UUID, deviceID *string) (*models.Device, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Device, error)
	Update(ctx context.Context, device *models.Device) error
	Delete(ctx context.Context, id uuid.UUID) error
}
