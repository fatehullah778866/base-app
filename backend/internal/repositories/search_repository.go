package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type SearchRepository interface {
	SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error
	GetSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.SearchHistory, error)
	ClearSearchHistory(ctx context.Context, userID uuid.UUID) error
	SearchDashboardItems(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.DashboardItem, error)
	SearchMessages(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Message, error)
	SearchUsers(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.User, error)
	SearchUsersByLocation(ctx context.Context, userID uuid.UUID, country, city *string, limit int) ([]*models.User, error)
	SearchNotifications(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Notification, error)
}

