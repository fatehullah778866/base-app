package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, search string) ([]*models.User, error)
	SetStatus(ctx context.Context, id uuid.UUID, status string) error
	MarkDeleted(ctx context.Context, id uuid.UUID) error
	PurgeDeletedBefore(ctx context.Context, cutoff time.Time) error
}
