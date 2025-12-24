package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	GetByConversation(ctx context.Context, conversationID uuid.UUID, limit int) ([]*models.Message, error)
	GetConversationsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Conversation, error)
	GetOrCreateConversation(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAsArchived(ctx context.Context, id uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}

