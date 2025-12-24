package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, token *models.PasswordResetToken) error
	GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context, before time.Time) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

