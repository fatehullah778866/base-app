package webhooks

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type Dispatcher struct {
	webhookRepo repositories.WebhookRepository
	httpClient  *http.Client
	logger      *zap.Logger
	maxRetries  int
	backoffMult float64
	secret      string // Default webhook secret
}

func NewDispatcher(
	webhookRepo repositories.WebhookRepository,
	secret string,
	maxRetries int,
	backoffMult float64,
	logger *zap.Logger,
) *Dispatcher {
	return &Dispatcher{
		webhookRepo: webhookRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:      logger,
		maxRetries:  maxRetries,
		backoffMult: backoffMult,
		secret:      secret,
	}
}

func (d *Dispatcher) ProcessPendingEvents(ctx context.Context) error {
	events, err := d.webhookRepo.GetPendingEvents(ctx, 100)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := d.deliverEvent(ctx, event); err != nil {
			d.logger.Error("Failed to deliver webhook event", zap.Error(err), zap.String("event_id", event.ID.String()))
		}
	}

	return nil
}

func (d *Dispatcher) deliverEvent(ctx context.Context, event *models.WebhookEvent) error {
	// Update status to processing
	now := time.Now()
	event.Status = "processing"
	event.ProcessedAt = &now
	d.webhookRepo.UpdateEvent(ctx, event)

	// Get subscription to retrieve webhook secret
	subscription, err := d.webhookRepo.GetSubscriptionByURL(ctx, event.WebhookURL)
	webhookSecret := d.secret // Default secret
	if err == nil && subscription != nil && subscription.WebhookSecret != "" {
		webhookSecret = subscription.WebhookSecret
	}

	// Prepare payload
	payload := map[string]interface{}{
		"event_id":      event.ID.String(),
		"event_type":    event.EventType,
		"event_version": event.EventVersion,
		"event_source":  event.EventSource,
		"timestamp":     event.CreatedAt.Format(time.RFC3339),
		"user_id":       event.UserID.String(),
		"payload":       json.RawMessage(event.Payload),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	signature := d.signPayload(timestamp, payloadBytes, webhookSecret)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", event.WebhookURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", fmt.Sprintf("sha256=%s", signature))
	req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Webhook-Event-ID", event.ID.String())
	req.Header.Set("X-Webhook-Event-Type", event.EventType)

	// Send request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return d.handleDeliveryError(ctx, event, err)
	}
	defer resp.Body.Close()

	// Read response body for logging
	bodyBytes := make([]byte, 1024)
	n, _ := resp.Body.Read(bodyBytes)
	responseBody := string(bodyBytes[:n])

	// Check response
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Success
		deliveredAt := time.Now()
		event.Status = "delivered"
		event.DeliveredAt = &deliveredAt
		event.LastResponseStatus = &resp.StatusCode
		event.LastResponseBody = &responseBody
		d.webhookRepo.UpdateEvent(ctx, event)
		return nil
	}

	// Failure
	return d.handleDeliveryError(ctx, event, fmt.Errorf("unexpected status code: %d", resp.StatusCode))
}

func (d *Dispatcher) handleDeliveryError(ctx context.Context, event *models.WebhookEvent, err error) error {
	event.DeliveryAttempts++
	errMsg := err.Error()
	event.LastErrorMessage = &errMsg

	if event.DeliveryAttempts >= d.maxRetries {
		event.Status = "failed"
		d.webhookRepo.UpdateEvent(ctx, event)
		return fmt.Errorf("max retries exceeded: %w", err)
	}

	// Calculate next retry time with exponential backoff
	backoff := time.Duration(float64(time.Second*60) * d.backoffMult * float64(event.DeliveryAttempts))
	nextRetry := time.Now().Add(backoff)
	event.Status = "retrying"
	event.NextRetryAt = &nextRetry
	d.webhookRepo.UpdateEvent(ctx, event)

	return nil
}

func (d *Dispatcher) signPayload(timestamp int64, payload []byte, secret string) string {
	message := fmt.Sprintf("%d.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

