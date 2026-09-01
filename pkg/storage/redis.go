package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// RedisStore defines the Redis interface
type RedisStore interface {
	Client() *redis.Client
	Close() error
	
	// Sessions
	SetSession(ctx context.Context, id string, data string, ttl time.Duration) error
	GetSession(ctx context.Context, id string) (string, error)
}

type redisStore struct {
	client *redis.Client
	cfg    config.StorageRedisConfig
	logger *zap.Logger
}

// NewRedisStore creates a new Redis connection
func NewRedisStore(ctx context.Context, cfg config.StorageRedisConfig, logger *zap.Logger) (RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.Database,
		PoolSize:     cfg.MaxConnections,
		MinIdleConns: cfg.MinIdleConnections,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("Failed to ping redis, will retry later", zap.Error(err))
	}

	return &redisStore{
		client: client,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (s *redisStore) Client() *redis.Client {
	return s.client
}

func (s *redisStore) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *redisStore) SetSession(ctx context.Context, id string, data string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", s.cfg.Session.Prefix, id)
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *redisStore) GetSession(ctx context.Context, id string) (string, error) {
	key := fmt.Sprintf("%s%s", s.cfg.Session.Prefix, id)
	return s.client.Get(ctx, key).Result()
}
