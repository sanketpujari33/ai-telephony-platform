package ratelimit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/auth"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func TestRateLimiter_FailOpen(t *testing.T) {
	logger := zap.NewNop()

	// Use an invalid redis address to simulate redis being down
	rdb := redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:9999", // unlikely to have redis here
		DialTimeout: 10 * time.Millisecond,
		MaxRetries:  1,
	})

	limiter := NewRateLimiter(rdb, logger)
	middleware := limiter.Middleware(2, time.Minute)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	// Inject tenant context
	ctx := context.WithValue(req.Context(), auth.TenantContextKey, "test-tenant")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Since Redis is down, it should fail open and return 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
