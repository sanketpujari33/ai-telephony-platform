package storage

import (
	"context"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"go.uber.org/zap"
)

// Storage defines the interface for data persistence
type Storage interface {
	Postgres() PostgresStore
	Redis() RedisStore
	S3() S3Store
	Close() error
}

type storageImpl struct {
	postgres PostgresStore
	redis    RedisStore
	s3       S3Store
	logger   *zap.Logger
}

// NewStorage initializes and returns a new Storage instance
func NewStorage(ctx context.Context, cfg config.StorageConfig, logger *zap.Logger) (Storage, error) {
	pgStore, err := NewPostgresStore(ctx, cfg.Postgres, logger)
	if err != nil {
		return nil, err
	}

	redisStore, err := NewRedisStore(ctx, cfg.Redis, logger)
	if err != nil {
		pgStore.Close()
		return nil, err
	}

	s3Store, err := NewS3Store(ctx, cfg.S3, logger)
	if err != nil {
		pgStore.Close()
		redisStore.Close()
		return nil, err
	}

	return &storageImpl{
		postgres: pgStore,
		redis:    redisStore,
		s3:       s3Store,
		logger:   logger,
	}, nil
}

func (s *storageImpl) Postgres() PostgresStore {
	return s.postgres
}

func (s *storageImpl) Redis() RedisStore {
	return s.redis
}

func (s *storageImpl) S3() S3Store {
	return s.s3
}

func (s *storageImpl) Close() error {
	var errs []error
	
	if err := s.postgres.Close(); err != nil {
		errs = append(errs, err)
	}
	
	if err := s.redis.Close(); err != nil {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		return errs[0]
	}
	
	return nil
}
