package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                    uuid.UUID  `db:"id" json:"id"`
	UserID                uuid.UUID  `db:"user_id" json:"user_id"`
	Token                 string     `db:"token" json:"-"`
	RefreshToken          *string    `db:"refresh_token" json:"-"`
	RefreshTokenExpiresAt *time.Time `db:"refresh_token_expires_at" json:"-"`
	DeviceID              *string    `db:"device_id" json:"device_id"`
	DeviceType            *string    `db:"device_type" json:"device_type"`
	DeviceName            *string    `db:"device_name" json:"device_name"`
	OS                    *string    `db:"os" json:"os"`
	Browser               *string    `db:"browser" json:"browser"`
	IPAddress             *string    `db:"ip_address" json:"ip_address"`
	LocationCountry       *string    `db:"location_country" json:"location_country"`
	LocationCity          *string    `db:"location_city" json:"location_city"`
	IsActive              bool       `db:"is_active" json:"is_active"`
	ExpiresAt             time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	LastUsedAt            time.Time  `db:"last_used_at" json:"last_used_at"`
}
