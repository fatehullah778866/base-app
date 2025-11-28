package models

import "time"

type ActivityLog struct {
	ID         string     `db:"id" json:"id"`
	ActorID    *string    `db:"actor_id" json:"actor_id"`
	ActorRole  *string    `db:"actor_role" json:"actor_role"`
	Action     string     `db:"action" json:"action"`
	TargetType *string    `db:"target_type" json:"target_type"`
	TargetID   *string    `db:"target_id" json:"target_id"`
	Metadata   *string    `db:"metadata" json:"metadata"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}
