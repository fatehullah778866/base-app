package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type Emitter struct {
	webhookRepo repositories.WebhookRepository
	secret      string
	logger      *zap.Logger
}

func (e *Emitter) SignPayload(timestamp int64, payload []byte, secret string) string {
	message := fmt.Sprintf("%d.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func NewEmitter(webhookRepo repositories.WebhookRepository, secret string, logger *zap.Logger) *Emitter {
	return &Emitter{
		webhookRepo: webhookRepo,
		secret:      secret,
		logger:      logger,
	}
}

type Event struct {
	EventType    string
	EventVersion string
	UserID       uuid.UUID
	Payload      interface{}
	Metadata     map[string]interface{}
}

func (e *Emitter) Emit(ctx context.Context, event Event) error {
	eventID := uuid.New()
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	payloadHash := e.hashPayload(payloadBytes)

	webhookEvent := &models.WebhookEvent{
		ID:           eventID,
		EventType:    event.EventType,
		EventVersion: event.EventVersion,
		EventSource:  "base_app",
		UserID:       event.UserID,
		Payload:      payloadBytes,
		PayloadHash:  payloadHash,
		Status:       "pending",
		ScheduledAt:  time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Get active subscriptions for this event type
	subscriptions, err := e.webhookRepo.GetActiveSubscriptions(ctx, event.EventType)
	if err != nil {
		return err
	}

	// Create webhook events for each subscription
	for _, sub := range subscriptions {
		eventCopy := *webhookEvent
		eventCopy.WebhookURL = sub.WebhookURL
		eventCopy.WebhookSecret = sub.WebhookSecret

		if err := e.webhookRepo.CreateEvent(ctx, &eventCopy); err != nil {
			e.logger.Error("Failed to create webhook event", zap.Error(err))
			continue
		}
	}

	e.logger.Info("Webhook event emitted", zap.String("event_type", event.EventType))

	return nil
}

func (e *Emitter) hashPayload(payload []byte) string {
	hash := sha256.Sum256(payload)
	return hex.EncodeToString(hash[:])
}
