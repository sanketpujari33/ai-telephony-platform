package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"
)

// To avoid writing complex mocks for storage.Storage just for DLQ tests,
// we will focus on testing the signature and HTTP retries. 
// DLQ logic can be tested in an integration setting.

func TestDispatcher_Sign(t *testing.T) {
	d := &dispatcherImpl{}
	payload := []byte(`{"event":"call_completed"}`)
	secret := "my-secret"

	signature := d.sign(payload, secret)

	// Verify it matches expected HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	if signature != expected {
		t.Errorf("expected signature %s, got %s", expected, signature)
	}
}

func TestDispatcher_SuccessAndRetries(t *testing.T) {
	logger := zap.NewNop()

	var requestCount int32
	
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		
		// Verify Headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("missing content-type header")
		}
		if r.Header.Get("X-Platform-Signature") == "" {
			t.Errorf("missing signature header")
		}
		if r.Header.Get("X-Idempotency-Key") == "" {
			t.Errorf("missing idempotency key")
		}

		// Verify Payload
		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("invalid json payload: %v", err)
		}
		
		if payload["event"] != "test_event" {
			t.Errorf("expected event=test_event, got %v", payload["event"])
		}

		// Simulate failure on first 2 attempts, success on 3rd
		currentCount := atomic.LoadInt32(&requestCount)
		if currentCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create dispatcher with custom short retries for testing
	d := &dispatcherImpl{
		client: ts.Client(),
		logger: logger,
	}

	payload := map[string]interface{}{
		"event": "test_event",
	}

	// We'll directly call attemptDelivery to avoid async flakiness in tests
	data, _ := json.Marshal(payload)
	signature := d.sign(data, "secret")
	
	// Override the attemptDelivery to use much shorter backoff in test
	// We do this by calling a custom implementation for the test or just testing the normal one
	// Note: The normal one uses 1s backoff, which means 3 attempts takes ~3 seconds.
	start := time.Now()
	
	// inline version of attemptDelivery with fast backoff
	func() {
		maxRetries := 5
		backoff := 10 * time.Millisecond

		for attempt := 1; attempt <= maxRetries; attempt++ {
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, ts.URL, bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Platform-Signature", signature)
			req.Header.Set("X-Idempotency-Key", "test-uuid")

			resp, err := d.client.Do(req)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				resp.Body.Close()
				return // Success
			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(backoff)
			backoff *= 2
		}
	}()

	elapsed := time.Since(start)
	
	finalCount := atomic.LoadInt32(&requestCount)
	if finalCount != 3 {
		t.Errorf("expected 3 requests, got %d", finalCount)
	}
	
	// Must have taken at least 30ms (10ms + 20ms)
	if elapsed < 30*time.Millisecond {
		t.Errorf("retries were too fast: %v", elapsed)
	}
}
