package models

import (
	"time"

	"github.com/google/uuid"
)

type DashboardItem struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	Title       string     `db:"title" json:"title"`
	Description *string    `db:"description" json:"description"`
	Category    *string    `db:"category" json:"category"`
	Status      string     `db:"status" json:"status"` // active, archived, deleted
	Priority    int        `db:"priority" json:"priority"`
	Metadata    *string    `db:"metadata" json:"metadata"` // JSON for extensibility
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

