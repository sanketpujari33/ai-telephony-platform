package ratelimit

import (
	"net/http"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/auth"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// RateLimiter enforces API request limits per tenant (T12.4)
type RateLimiter interface {
	Middleware(limit int, window time.Duration) func(http.Handler) http.Handler
}

type rateLimiterImpl struct {
	rdb    *redis.Client
	logger *zap.Logger
}

// NewRateLimiter initializes a Redis-backed rate limiter
func NewRateLimiter(rdb *redis.Client, logger *zap.Logger) RateLimiter {
	return &rateLimiterImpl{
		rdb:    rdb,
		logger: logger,
	}
}

func (rl *rateLimiterImpl) Middleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract tenant from context (injected by Auth middleware)
			tenantID, ok := r.Context().Value(auth.TenantContextKey).(string)
			if !ok || tenantID == "" {
				tenantID = r.RemoteAddr // fallback to IP if unauthenticated route
			}

			key := "ratelimit:" + tenantID
			
			// Simple counter-based rate limiting for demonstration
			// In production, use lua scripts for atomic Token Bucket or sliding window
			ctx := r.Context()
			count, err := rl.rdb.Incr(ctx, key).Result()
			if err != nil {
				rl.logger.Error("Rate limiter redis error", zap.Error(err))
				// Fail open if Redis is down
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				rl.rdb.Expire(ctx, key, window)
			}

			if count > int64(limit) {
				rl.logger.Warn("Rate limit exceeded", zap.String("tenant", tenantID))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
