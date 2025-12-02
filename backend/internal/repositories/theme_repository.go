package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type ThemeRepository interface {
	GetGlobalTheme(ctx context.Context, userID uuid.UUID) (*models.ThemePreferences, error)
	UpdateGlobalTheme(ctx context.Context, theme *models.ThemePreferences) error
	GetProductOverride(ctx context.Context, userID uuid.UUID, productName string) (*models.ProductThemeOverride, error)
	UpsertProductOverride(ctx context.Context, override *models.ProductThemeOverride) error
	DeleteProductOverride(ctx context.Context, userID uuid.UUID, productName string) error
}
