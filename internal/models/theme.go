package models

import (
	"time"

	"github.com/google/uuid"
)

type ThemePreferences struct {
	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	Theme         string    `db:"kompassui_theme" json:"theme"`
	Contrast      string    `db:"kompassui_contrast" json:"contrast"`
	TextDirection string    `db:"kompassui_text_direction" json:"text_direction"`
	Brand         *string   `db:"kompassui_brand" json:"brand"`
	SyncedAt      time.Time `db:"theme_synced_at" json:"synced_at"`
	SyncEnabled   bool      `db:"theme_sync_enabled" json:"sync_enabled"`
}

type ProductThemeOverride struct {
	ID            uuid.UUID `db:"id" json:"id"`
	UserID        uuid.UUID `db:"user_id" json:"user_id"`
	ProductName   string    `db:"product_name" json:"product_name"`
	Theme         string    `db:"theme" json:"theme"`
	Contrast      string    `db:"contrast" json:"contrast"`
	TextDirection string    `db:"text_direction" json:"text_direction"`
	Brand         *string   `db:"brand" json:"brand"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

