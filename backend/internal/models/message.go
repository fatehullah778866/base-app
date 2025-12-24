package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	SenderID    uuid.UUID  `db:"sender_id" json:"sender_id"`
	RecipientID uuid.UUID  `db:"recipient_id" json:"recipient_id"`
	Subject     *string    `db:"subject" json:"subject"`
	Content     string     `db:"content" json:"content"`
	IsRead      bool       `db:"is_read" json:"is_read"`
	ReadAt      *time.Time `db:"read_at" json:"read_at"`
	IsArchived  bool       `db:"is_archived" json:"is_archived"`
	ArchivedAt  *time.Time `db:"archived_at" json:"archived_at"`
	Metadata    *string    `db:"metadata" json:"metadata"` // JSON for attachments, etc.
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type Conversation struct {
	ID                      uuid.UUID  `db:"id" json:"id"`
	Participant1ID         uuid.UUID  `db:"participant1_id" json:"participant1_id"`
	Participant2ID         uuid.UUID  `db:"participant2_id" json:"participant2_id"`
	LastMessageID          *uuid.UUID `db:"last_message_id" json:"last_message_id"`
	LastMessageAt          *time.Time `db:"last_message_at" json:"last_message_at"`
	Participant1UnreadCount int       `db:"participant1_unread_count" json:"participant1_unread_count"`
	Participant2UnreadCount int       `db:"participant2_unread_count" json:"participant2_unread_count"`
	CreatedAt              time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at" json:"updated_at"`
}

