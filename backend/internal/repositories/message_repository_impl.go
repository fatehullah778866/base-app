package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type messageRepository struct {
	db *database.DB
}

func NewMessageRepository(db *database.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (id, sender_id, recipient_id, subject, content, is_read, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		message.ID.String(),
		message.SenderID.String(),
		message.RecipientID.String(),
		message.Subject,
		message.Content,
		message.IsRead,
		message.Metadata,
		message.CreatedAt,
		message.UpdatedAt,
	)
	return err
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var m models.Message
	var senderIDStr, recipientIDStr, idStr string
	var subject, metadata sql.NullString
	var readAt, archivedAt sql.NullTime

	query := `SELECT id, sender_id, recipient_id, subject, content, is_read, read_at, is_archived, archived_at, metadata, created_at, updated_at
		FROM messages WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &senderIDStr, &recipientIDStr, &subject, &m.Content,
		&m.IsRead, &readAt, &m.IsArchived, &archivedAt, &metadata,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	m.ID, _ = uuid.Parse(idStr)
	m.SenderID, _ = uuid.Parse(senderIDStr)
	m.RecipientID, _ = uuid.Parse(recipientIDStr)
	if subject.Valid {
		m.Subject = &subject.String
	}
	if readAt.Valid {
		m.ReadAt = &readAt.Time
	}
	if archivedAt.Valid {
		m.ArchivedAt = &archivedAt.Time
	}
	if metadata.Valid {
		m.Metadata = &metadata.String
	}

	return &m, nil
}

func (r *messageRepository) GetByConversation(ctx context.Context, conversationID uuid.UUID, limit int) ([]*models.Message, error) {
	// Get conversation participants
	var conv models.Conversation
	var p1IDStr, p2IDStr, convIDStr string
	var lastMsgID sql.NullString
	var lastMsgAt sql.NullTime

	query := `SELECT id, participant1_id, participant2_id, last_message_id, last_message_at,
		participant1_unread_count, participant2_unread_count, created_at, updated_at
		FROM conversations WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, conversationID.String()).Scan(
		&convIDStr, &p1IDStr, &p2IDStr, &lastMsgID, &lastMsgAt,
		&conv.Participant1UnreadCount, &conv.Participant2UnreadCount,
		&conv.CreatedAt, &conv.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	conv.ID = conversationID
	conv.Participant1ID, _ = uuid.Parse(p1IDStr)
	conv.Participant2ID, _ = uuid.Parse(p2IDStr)

	// Get messages between participants
	msgQuery := `SELECT id, sender_id, recipient_id, subject, content, is_read, read_at, is_archived, archived_at, metadata, created_at, updated_at
		FROM messages WHERE (sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)
		ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, msgQuery,
		conv.Participant1ID.String(), conv.Participant2ID.String(),
		conv.Participant2ID.String(), conv.Participant1ID.String(),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var m models.Message
		var senderIDStr, recipientIDStr, idStr string
		var subject, metadata sql.NullString
		var readAt, archivedAt sql.NullTime

		err := rows.Scan(&idStr, &senderIDStr, &recipientIDStr, &subject, &m.Content,
			&m.IsRead, &readAt, &m.IsArchived, &archivedAt, &metadata,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}

		m.ID, _ = uuid.Parse(idStr)
		m.SenderID, _ = uuid.Parse(senderIDStr)
		m.RecipientID, _ = uuid.Parse(recipientIDStr)
		if subject.Valid {
			m.Subject = &subject.String
		}
		if readAt.Valid {
			m.ReadAt = &readAt.Time
		}
		if archivedAt.Valid {
			m.ArchivedAt = &archivedAt.Time
		}
		if metadata.Valid {
			m.Metadata = &metadata.String
		}

		messages = append(messages, &m)
	}

	return messages, rows.Err()
}

func (r *messageRepository) GetConversationsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Conversation, error) {
	query := `SELECT id, participant1_id, participant2_id, last_message_id, last_message_at,
		participant1_unread_count, participant2_unread_count, created_at, updated_at
		FROM conversations WHERE participant1_id = ? OR participant2_id = ?
		ORDER BY last_message_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID.String(), userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		var conv models.Conversation
		var p1IDStr, p2IDStr, convIDStr string
		var lastMsgID sql.NullString
		var lastMsgAt sql.NullTime

		err := rows.Scan(&convIDStr, &p1IDStr, &p2IDStr, &lastMsgID, &lastMsgAt,
			&conv.Participant1UnreadCount, &conv.Participant2UnreadCount,
			&conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, err
		}

		conv.ID, _ = uuid.Parse(convIDStr)
		conv.Participant1ID, _ = uuid.Parse(p1IDStr)
		conv.Participant2ID, _ = uuid.Parse(p2IDStr)
		if lastMsgID.Valid {
			id, _ := uuid.Parse(lastMsgID.String)
			conv.LastMessageID = &id
		}
		if lastMsgAt.Valid {
			conv.LastMessageAt = &lastMsgAt.Time
		}

		conversations = append(conversations, &conv)
	}

	return conversations, rows.Err()
}

func (r *messageRepository) GetOrCreateConversation(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Conversation, error) {
	// Try to get existing conversation
	var conv models.Conversation
	var convIDStr, p1IDStr, p2IDStr string
	var lastMsgID sql.NullString
	var lastMsgAt sql.NullTime

	query := `SELECT id, participant1_id, participant2_id, last_message_id, last_message_at,
		participant1_unread_count, participant2_unread_count, created_at, updated_at
		FROM conversations WHERE (participant1_id = ? AND participant2_id = ?) OR (participant1_id = ? AND participant2_id = ?)`
	err := r.db.QueryRowContext(ctx, query,
		user1ID.String(), user2ID.String(),
		user2ID.String(), user1ID.String(),
	).Scan(&convIDStr, &p1IDStr, &p2IDStr, &lastMsgID, &lastMsgAt,
		&conv.Participant1UnreadCount, &conv.Participant2UnreadCount,
		&conv.CreatedAt, &conv.UpdatedAt)

	if err == nil {
		conv.ID, _ = uuid.Parse(convIDStr)
		conv.Participant1ID, _ = uuid.Parse(p1IDStr)
		conv.Participant2ID, _ = uuid.Parse(p2IDStr)
		if lastMsgID.Valid {
			id, _ := uuid.Parse(lastMsgID.String)
			conv.LastMessageID = &id
		}
		if lastMsgAt.Valid {
			conv.LastMessageAt = &lastMsgAt.Time
		}
		return &conv, nil
	}

	// Create new conversation
	conv.ID = uuid.New()
	conv.Participant1ID = user1ID
	conv.Participant2ID = user2ID
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()

	insertQuery := `INSERT INTO conversations (id, participant1_id, participant2_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err = r.db.ExecContext(ctx, insertQuery,
		conv.ID.String(), conv.Participant1ID.String(), conv.Participant2ID.String(),
		conv.CreatedAt, conv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &conv, nil
}

func (r *messageRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE messages SET is_read = 1, read_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *messageRepository) MarkAsArchived(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE messages SET is_archived = 1, archived_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *messageRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages WHERE recipient_id = ? AND is_read = 0`
	err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(&count)
	return count, err
}

