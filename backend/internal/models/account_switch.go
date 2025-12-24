package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountSwitch struct {
	ID              uuid.UUID `db:"id" json:"id"`
	UserID          uuid.UUID `db:"user_id" json:"user_id"`
	SwitchedToUserID *uuid.UUID `db:"switched_to_user_id" json:"switched_to_user_id"`
	SwitchedToRole  *string    `db:"switched_to_role" json:"switched_to_role"`
	SwitchedFromRole *string   `db:"switched_from_role" json:"switched_from_role"`
	Reason          *string   `db:"reason" json:"reason"`
	IPAddress       *string   `db:"ip_address" json:"ip_address"`
	UserAgent       *string   `db:"user_agent" json:"user_agent"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

