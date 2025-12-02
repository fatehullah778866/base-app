package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	UserID          uuid.UUID  `db:"user_id" json:"user_id"`
	DeviceID        string     `db:"device_id" json:"device_id"`
	DeviceName      *string    `db:"device_name" json:"device_name"`
	DeviceType      *string    `db:"device_type" json:"device_type"`
	OS              *string    `db:"os" json:"os"`
	Browser         *string    `db:"browser" json:"browser"`
	IPAddress       *string    `db:"ip_address" json:"ip_address"`
	LocationCountry *string    `db:"location_country" json:"location_country"`
	LocationCity    *string    `db:"location_city" json:"location_city"`
	IsTrusted       bool       `db:"is_trusted" json:"is_trusted"`
	TrustedAt       *time.Time `db:"trusted_at" json:"trusted_at"`
	LastUsedAt      time.Time  `db:"last_used_at" json:"last_used_at"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}
