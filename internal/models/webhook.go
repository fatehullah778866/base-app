package models

import (
	"time"

	"github.com/google/uuid"
)

type WebhookEvent struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	EventType          string     `db:"event_type" json:"event_type"`
	EventVersion       string     `db:"event_version" json:"event_version"`
	EventSource        string     `db:"event_source" json:"event_source"`
	UserID             uuid.UUID  `db:"user_id" json:"user_id"`
	Payload            []byte     `db:"payload" json:"payload"`
	PayloadHash        string     `db:"payload_hash" json:"payload_hash"`
	WebhookURL         string     `db:"webhook_url" json:"webhook_url"`
	WebhookSecret      string     `db:"webhook_secret" json:"-"`
	Status             string     `db:"status" json:"status"`
	DeliveryAttempts   int        `db:"delivery_attempts" json:"delivery_attempts"`
	MaxAttempts        int        `db:"max_attempts" json:"max_attempts"`
	ScheduledAt        time.Time  `db:"scheduled_at" json:"scheduled_at"`
	ProcessedAt        *time.Time `db:"processed_at" json:"processed_at"`
	DeliveredAt        *time.Time `db:"delivered_at" json:"delivered_at"`
	NextRetryAt        *time.Time `db:"next_retry_at" json:"next_retry_at"`
	LastResponseStatus *int       `db:"last_response_status" json:"last_response_status"`
	LastResponseBody   *string    `db:"last_response_body" json:"last_response_body"`
	LastErrorMessage   *string    `db:"last_error_message" json:"last_error_message"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
}

type WebhookSubscription struct {
	ID                    uuid.UUID  `db:"id" json:"id"`
	UserID                *uuid.UUID `db:"user_id" json:"user_id"`
	SubscriptionName      string     `db:"subscription_name" json:"subscription_name"`
	WebhookURL            string     `db:"webhook_url" json:"webhook_url"`
	WebhookSecret         string     `db:"webhook_secret" json:"-"`
	EventTypes            []string   `db:"event_types" json:"event_types"`
	IsActive              bool       `db:"is_active" json:"is_active"`
	IsVerified            bool       `db:"is_verified" json:"is_verified"`
	RateLimitPerMinute    int        `db:"rate_limit_per_minute" json:"rate_limit_per_minute"`
	MaxRetries            int        `db:"max_retries" json:"max_retries"`
	RetryBackoffMultiplier float64   `db:"retry_backoff_multiplier" json:"retry_backoff_multiplier"`
	Description           *string    `db:"description" json:"description"`
	Metadata              []byte     `db:"metadata" json:"metadata"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updated_at"`
}

