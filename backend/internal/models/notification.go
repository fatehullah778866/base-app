package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	Type      string     `db:"type" json:"type"` // message, alert, promotion, security, system
	Title     string     `db:"title" json:"title"`
	Message   string     `db:"message" json:"message"`
	Link      *string    `db:"link" json:"link"`
	IsRead    bool       `db:"is_read" json:"is_read"`
	ReadAt    *time.Time `db:"read_at" json:"read_at"`
	Metadata  *string    `db:"metadata" json:"metadata"` // JSON
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

