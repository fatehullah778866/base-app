package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type webhookRepository struct {
	db *database.DB
}

func NewWebhookRepository(db *database.DB) WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) CreateEvent(ctx context.Context, event *models.WebhookEvent) error {
	query := `
		INSERT INTO webhook_events (
			id, event_type, event_version, event_source, user_id, payload,
			payload_hash, webhook_url, webhook_secret, status, max_attempts,
			scheduled_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		event.ID, event.EventType, event.EventVersion, event.EventSource,
		event.UserID, event.Payload, event.PayloadHash,
		event.WebhookURL, event.WebhookSecret, event.Status, event.MaxAttempts,
		event.ScheduledAt, event.CreatedAt, event.UpdatedAt,
	)

	return err
}

func (r *webhookRepository) GetPendingEvents(ctx context.Context, limit int) ([]*models.WebhookEvent, error) {
	query := `
		SELECT id, event_type, event_version, event_source, user_id, payload,
			payload_hash, webhook_url, webhook_secret, status, delivery_attempts,
			max_attempts, scheduled_at, processed_at, delivered_at, next_retry_at,
			last_response_status, last_response_body, last_error_message,
			created_at, updated_at
		FROM webhook_events
		WHERE status IN ('pending', 'retrying')
			AND (next_retry_at IS NULL OR next_retry_at <= CURRENT_TIMESTAMP)
		ORDER BY scheduled_at ASC
		LIMIT ?
	`

	rows, err := r.db.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.WebhookEvent
	for rows.Next() {
		event := &models.WebhookEvent{}
		err := rows.Scan(
			&event.ID, &event.EventType, &event.EventVersion, &event.EventSource,
			&event.UserID, &event.Payload, &event.PayloadHash,
			&event.WebhookURL, &event.WebhookSecret, &event.Status,
			&event.DeliveryAttempts, &event.MaxAttempts,
			&event.ScheduledAt, &event.ProcessedAt, &event.DeliveredAt,
			&event.NextRetryAt, &event.LastResponseStatus, &event.LastResponseBody,
			&event.LastErrorMessage, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *webhookRepository) UpdateEvent(ctx context.Context, event *models.WebhookEvent) error {
	query := `
		UPDATE webhook_events
			SET status = ?, delivery_attempts = ?, processed_at = ?,
				delivered_at = ?, next_retry_at = ?, last_response_status = ?,
				last_response_body = ?, last_error_message = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		event.Status, event.DeliveryAttempts,
		event.ProcessedAt, event.DeliveredAt, event.NextRetryAt,
		event.LastResponseStatus, event.LastResponseBody, event.LastErrorMessage, event.ID,
	)

	return err
}

func (r *webhookRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error) {
	query := `
		SELECT id, event_type, event_version, event_source, user_id, payload,
			payload_hash, webhook_url, webhook_secret, status, delivery_attempts,
			max_attempts, scheduled_at, processed_at, delivered_at, next_retry_at,
			last_response_status, last_response_body, last_error_message,
			created_at, updated_at
		FROM webhook_events
		WHERE id = ?
	`

	event := &models.WebhookEvent{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.EventType, &event.EventVersion, &event.EventSource,
		&event.UserID, &event.Payload, &event.PayloadHash,
		&event.WebhookURL, &event.WebhookSecret, &event.Status,
		&event.DeliveryAttempts, &event.MaxAttempts,
		&event.ScheduledAt, &event.ProcessedAt, &event.DeliveredAt,
		&event.NextRetryAt, &event.LastResponseStatus, &event.LastResponseBody,
		&event.LastErrorMessage, &event.CreatedAt, &event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("webhook event not found")
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *webhookRepository) GetActiveSubscriptions(ctx context.Context, eventType string) ([]*models.WebhookSubscription, error) {
	query := `
		SELECT id, user_id, subscription_name, webhook_url, webhook_secret,
			event_types, is_active, is_verified, rate_limit_per_minute,
			max_retries, retry_backoff_multiplier, description, metadata,
			created_at, updated_at
		FROM webhook_subscriptions
		WHERE is_active = 1
	`

	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.WebhookSubscription
	for rows.Next() {
		sub := &models.WebhookSubscription{}
		var eventTypes string
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.SubscriptionName, &sub.WebhookURL,
			&sub.WebhookSecret, &eventTypes, &sub.IsActive, &sub.IsVerified,
			&sub.RateLimitPerMinute, &sub.MaxRetries, &sub.RetryBackoffMultiplier,
			&sub.Description, &sub.Metadata, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sub.EventTypes = splitEventTypes(eventTypes)
		if eventType != "" && !containsEvent(sub.EventTypes, eventType) {
			continue
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *webhookRepository) GetSubscriptionByURL(ctx context.Context, url string) (*models.WebhookSubscription, error) {
	query := `
		SELECT id, user_id, subscription_name, webhook_url, webhook_secret,
			event_types, is_active, is_verified, rate_limit_per_minute,
			max_retries, retry_backoff_multiplier, description, metadata,
			created_at, updated_at
		FROM webhook_subscriptions
		WHERE webhook_url = ? AND is_active = 1
		LIMIT 1
	`

	sub := &models.WebhookSubscription{}
	var eventTypes string
	err := r.db.DB.QueryRowContext(ctx, query, url).Scan(
		&sub.ID, &sub.UserID, &sub.SubscriptionName, &sub.WebhookURL,
		&sub.WebhookSecret, &eventTypes, &sub.IsActive, &sub.IsVerified,
		&sub.RateLimitPerMinute, &sub.MaxRetries, &sub.RetryBackoffMultiplier,
		&sub.Description, &sub.Metadata, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found, return nil
	}
	if err != nil {
		return nil, err
	}

	sub.EventTypes = splitEventTypes(eventTypes)
	return sub, nil
}

func (r *webhookRepository) CreateSubscription(ctx context.Context, sub *models.WebhookSubscription) error {
	query := `
		INSERT INTO webhook_subscriptions (
			id, user_id, subscription_name, webhook_url, webhook_secret,
			event_types, is_active, is_verified, rate_limit_per_minute,
			max_retries, retry_backoff_multiplier, description, metadata,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		sub.ID, sub.UserID, sub.SubscriptionName, sub.WebhookURL, sub.WebhookSecret,
		joinEventTypes(sub.EventTypes), sub.IsActive, sub.IsVerified,
		sub.RateLimitPerMinute, sub.MaxRetries, sub.RetryBackoffMultiplier,
		sub.Description, sub.Metadata, sub.CreatedAt, sub.UpdatedAt,
	)

	return err
}

func (r *webhookRepository) UpdateSubscription(ctx context.Context, sub *models.WebhookSubscription) error {
	query := `
		UPDATE webhook_subscriptions
		SET subscription_name = ?, webhook_url = ?, webhook_secret = ?,
			event_types = ?, is_active = ?, is_verified = ?,
			rate_limit_per_minute = ?, max_retries = ?,
			retry_backoff_multiplier = ?, description = ?,
			metadata = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		sub.SubscriptionName, sub.WebhookURL, sub.WebhookSecret,
		joinEventTypes(sub.EventTypes), sub.IsActive, sub.IsVerified,
		sub.RateLimitPerMinute, sub.MaxRetries, sub.RetryBackoffMultiplier,
		sub.Description, sub.Metadata, sub.ID,
	)

	return err
}

func (r *webhookRepository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM webhook_subscriptions WHERE id = ?`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	return err
}

func joinEventTypes(types []string) string {
	return strings.Join(types, ",")
}

func splitEventTypes(types string) []string {
	if types == "" {
		return []string{}
	}
	parts := strings.Split(types, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func containsEvent(types []string, event string) bool {
	for _, t := range types {
		if t == event {
			return true
		}
	}
	return false
}
