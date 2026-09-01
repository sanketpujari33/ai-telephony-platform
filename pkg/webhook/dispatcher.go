package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Dispatcher handles reliable event delivery to CRMs (T12.2)
type Dispatcher interface {
	Dispatch(ctx context.Context, tenantID, endpoint string, payload interface{}, secret string) error
}

type dispatcherImpl struct {
	client *http.Client
	store  storage.Storage
	logger *zap.Logger
}

// NewDispatcher initializes a webhook dispatcher
func NewDispatcher(store storage.Storage, logger *zap.Logger) Dispatcher {
	return &dispatcherImpl{
		client: &http.Client{Timeout: 10 * time.Second},
		store:  store,
		logger: logger,
	}
}

func (d *dispatcherImpl) Dispatch(ctx context.Context, tenantID, endpoint string, payload interface{}, secret string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	idempotencyKey := uuid.New().String()
	signature := d.sign(data, secret)

	// Run in background for retries
	go d.attemptDelivery(context.Background(), tenantID, endpoint, data, signature, idempotencyKey)

	return nil
}

func (d *dispatcherImpl) attemptDelivery(ctx context.Context, tenantID, endpoint string, data []byte, signature, idempotencyKey string) {
	maxRetries := 5
	backoff := 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
		if err != nil {
			d.logger.Error("Failed to create webhook request", zap.Error(err))
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Platform-Signature", signature)
		req.Header.Set("X-Idempotency-Key", idempotencyKey)

		resp, err := d.client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			d.logger.Info("Webhook delivered successfully", zap.String("endpoint", endpoint))
			return
		}

		if resp != nil {
			resp.Body.Close()
		}

		d.logger.Warn("Webhook delivery failed, retrying...", zap.Int("attempt", attempt), zap.String("endpoint", endpoint))
		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff (1s, 2s, 4s, 8s, 16s)
	}

	// Exhausted retries, send to Dead Letter Queue (DLQ) backed by Postgres (T12.2)
	d.sendToDLQ(ctx, tenantID, endpoint, data)
}

func (d *dispatcherImpl) sendToDLQ(ctx context.Context, tenantID, endpoint string, data []byte) {
	d.logger.Error("Webhook exhausted retries, moving to DLQ", zap.String("tenant", tenantID), zap.String("endpoint", endpoint))
	
	// Assuming a `webhook_dlq` table exists in Postgres
	query := `INSERT INTO webhook_dlq (tenant_id, endpoint, payload, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := d.store.Postgres().DB().ExecContext(ctx, query, tenantID, endpoint, string(data))
	if err != nil {
		d.logger.Error("Failed to write to Webhook DLQ", zap.Error(err))
	}
}

func (d *dispatcherImpl) sign(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
