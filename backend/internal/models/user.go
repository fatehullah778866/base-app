package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                     uuid.UUID  `db:"id" json:"id"`
	Email                  string     `db:"email" json:"email"`
	EmailVerified          bool       `db:"email_verified" json:"email_verified"`
	EmailVerificationToken *string    `db:"email_verification_token" json:"-"`
	PasswordHash           string     `db:"password_hash" json:"-"`
	PasswordChangedAt      time.Time  `db:"password_changed_at" json:"password_changed_at"`
	Name                   string     `db:"name" json:"name"`
	FirstName              *string    `db:"first_name" json:"first_name"`
	LastName               *string    `db:"last_name" json:"last_name"`
	PhotoURL               *string    `db:"photo_url" json:"photo_url"`
	Role                   string     `db:"role" json:"role"`
	Phone                  *string    `db:"phone" json:"phone"`
	PhoneVerified          bool       `db:"phone_verified" json:"phone_verified"`
	Status                 string     `db:"status" json:"status"`
	SignupSource           *string    `db:"signup_source" json:"signup_source"`
	CreatedAt              time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at" json:"updated_at"`
	LastLoginAt            *time.Time `db:"last_login_at" json:"last_login_at"`
}
