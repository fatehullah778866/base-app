package models

import (
	"time"

	"github.com/google/uuid"
)

type SearchHistory struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	Query       string    `db:"query" json:"query"`
	SearchType  *string   `db:"search_type" json:"search_type"` // users, dashboard_items, messages, all
	ResultsCount int      `db:"results_count" json:"results_count"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type SearchResult struct {
	Type        string      `json:"type"` // user, dashboard_item, message
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	Data        interface{} `json:"data"`
}

