package models

import (
	"time"
)

type AccessRequest struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Title     *string   `db:"title" json:"title"`
	Details   *string   `db:"details" json:"details"`
	Status    string    `db:"status" json:"status"`
	Feedback  *string   `db:"feedback" json:"feedback"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
