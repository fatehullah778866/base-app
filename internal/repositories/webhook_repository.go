package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type WebhookRepository interface {
	// Event methods
	CreateEvent(ctx context.Context, event *models.WebhookEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*models.WebhookEvent, error)
	UpdateEvent(ctx context.Context, event *models.WebhookEvent) error
	GetEventByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error)

	// Subscription methods
	GetActiveSubscriptions(ctx context.Context, eventType string) ([]*models.WebhookSubscription, error)
	GetSubscriptionByURL(ctx context.Context, url string) (*models.WebhookSubscription, error)
	CreateSubscription(ctx context.Context, sub *models.WebhookSubscription) error
	UpdateSubscription(ctx context.Context, sub *models.WebhookSubscription) error
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
}

