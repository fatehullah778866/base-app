package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type AccountSwitchRepository interface {
	Create(ctx context.Context, switchRecord *models.AccountSwitch) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.AccountSwitch, error)
}

