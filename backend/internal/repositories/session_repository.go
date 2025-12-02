package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	GetByToken(ctx context.Context, token string) (*models.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
