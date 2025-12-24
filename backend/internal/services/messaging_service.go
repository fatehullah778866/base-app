package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type MessagingService struct {
	messageRepo repositories.MessageRepository
	userRepo    repositories.UserRepository
	logger      *zap.Logger
}

func NewMessagingService(
	messageRepo repositories.MessageRepository,
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) *MessagingService {
	return &MessagingService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

// SendMessage sends a message from one user to another
func (s *MessagingService) SendMessage(ctx context.Context, senderID, recipientID uuid.UUID, subject *string, content string, metadata *string) (*models.Message, error) {
	if content == "" {
		return nil, errors.New("message content is required")
	}

	// Verify recipient exists
	recipient, err := s.userRepo.GetByID(ctx, recipientID)
	if err != nil || recipient == nil {
		return nil, errors.New("recipient not found")
	}

	// Get or create conversation
	_, err = s.messageRepo.GetOrCreateConversation(ctx, senderID, recipientID)
	if err != nil {
		return nil, err
	}

	// Create message
	message := &models.Message{
		ID:          uuid.New(),
		SenderID:    senderID,
		RecipientID: recipientID,
		Subject:     subject,
		Content:     content,
		IsRead:      false,
		IsArchived:  false,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Update conversation (this would ideally be done in a transaction)
	// For now, we'll just log it
	s.logger.Info("Message sent", zap.String("message_id", message.ID.String()))

	return message, nil
}

// GetConversations retrieves all conversations for a user
func (s *MessagingService) GetConversations(ctx context.Context, userID uuid.UUID) ([]*models.Conversation, error) {
	return s.messageRepo.GetConversationsByUserID(ctx, userID)
}

// GetMessages retrieves messages for a conversation
func (s *MessagingService) GetMessages(ctx context.Context, conversationID uuid.UUID, limit int) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.messageRepo.GetByConversation(ctx, conversationID, limit)
}

// MarkAsRead marks a message as read
func (s *MessagingService) MarkAsRead(ctx context.Context, messageID uuid.UUID) error {
	return s.messageRepo.MarkAsRead(ctx, messageID)
}

// ArchiveMessage archives a message
func (s *MessagingService) ArchiveMessage(ctx context.Context, messageID uuid.UUID) error {
	return s.messageRepo.MarkAsArchived(ctx, messageID)
}

// GetUnreadCount returns the count of unread messages
func (s *MessagingService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.messageRepo.GetUnreadCount(ctx, userID)
}

