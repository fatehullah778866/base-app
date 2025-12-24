package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type SearchRepository interface {
	SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error
	SearchDashboardItems(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.DashboardItem, error)
	SearchMessages(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Message, error)
	SearchUsers(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.User, error)
}

