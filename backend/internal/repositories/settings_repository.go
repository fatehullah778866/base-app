package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type SettingsRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.ComprehensiveSettings, error)
	Create(ctx context.Context, settings *models.ComprehensiveSettings) error
	Update(ctx context.Context, settings *models.ComprehensiveSettings) error
	GetUsername(ctx context.Context, username string) (*models.ComprehensiveSettings, error)
}

