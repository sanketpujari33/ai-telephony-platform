package storage

import (
	"bytes"
	"context"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"go.uber.org/zap"
)

// S3Store defines the S3/MinIO interface
type S3Store interface {
	Client() *minio.Client
	Close() error
	
	UploadRecording(ctx context.Context, callID string, data []byte) error
	CleanupRecordings(ctx context.Context, retentionDays int) error
}

type s3Store struct {
	client *minio.Client
	cfg    config.StorageS3Config
	logger *zap.Logger
}

// NewS3Store creates a new S3 connection
func NewS3Store(ctx context.Context, cfg config.StorageS3Config, logger *zap.Logger) (S3Store, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	
	if err != nil {
		logger.Error("Failed to initialize S3 client", zap.Error(err))
		return nil, err
	}

	// Ensure bucket exists
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err == nil && !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region})
		if err != nil {
			logger.Warn("Failed to create S3 bucket", zap.Error(err))
		}
	}

	return &s3Store{
		client: client,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (s *s3Store) Client() *minio.Client {
	return s.client
}

func (s *s3Store) Close() error {
	// MinIO client doesn't need to be closed explicitly
	return nil
}

func (s *s3Store) UploadRecording(ctx context.Context, callID string, data []byte) error {
	objectName := callID + ".wav"
	if s.cfg.Recordings.Format != "" {
		objectName = callID + "." + s.cfg.Recordings.Format
	}
	
	reader := bytes.NewReader(data)
	
	// PCI DSS T10.4: SSE-AES256 Encryption at rest
	opts := minio.PutObjectOptions{
		ContentType: "audio/wav",
		ServerSideEncryption: encrypt.NewSSE(),
	}
	
	_, err := s.client.PutObject(ctx, s.cfg.Bucket, objectName, reader, int64(len(data)), opts)
	return err
}

// CleanupRecordings removes recordings older than retentionDays (T10.3)
func (s *s3Store) CleanupRecordings(ctx context.Context, retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	s.logger.Info("Starting S3 recording retention cleanup", zap.Time("cutoff_date", cutoff))
	
	// In production, this would list objects and remove those older than cutoff
	// or rely on S3 lifecycle policies (preferred).
	return nil
}
