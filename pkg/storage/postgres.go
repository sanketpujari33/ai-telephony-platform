package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// PostgresStore defines the PostgreSQL interface
type PostgresStore interface {
	DB() *sql.DB
	Close() error
	
	// Compliance checks
	IsDNC(ctx context.Context, phoneNumber string) (bool, error)
	OptOut(ctx context.Context, phoneNumber string) error
	
	// Retention (T10.3)
	CleanupTranscripts(ctx context.Context, retentionYears int) error
}

type postgresStore struct {
	db     *sql.DB
	cfg    config.StoragePostgresConfig
	logger *zap.Logger
}

// NewPostgresStore creates a new PostgreSQL connection
func NewPostgresStore(ctx context.Context, cfg config.StoragePostgresConfig, logger *zap.Logger) (PostgresStore, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Failed to open postgres connection", zap.Error(err))
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)

	if err := db.PingContext(ctx); err != nil {
		logger.Warn("Failed to ping postgres, will retry later", zap.Error(err))
		// Continue anyway to allow robust startup
	}

	return &postgresStore{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (s *postgresStore) DB() *sql.DB {
	return s.db
}

func (s *postgresStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *postgresStore) IsDNC(ctx context.Context, phoneNumber string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE phone_number = $1 LIMIT 1", s.cfg.Tables.DNCEntries)
	var exists int
	err := s.db.QueryRowContext(ctx, query, phoneNumber).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *postgresStore) OptOut(ctx context.Context, phoneNumber string) error {
	s.logger.Info("Adding number to DNC list (Opt-out)", zap.String("phone", phoneNumber))
	query := fmt.Sprintf("INSERT INTO %s (phone_number, created_at, source) VALUES ($1, NOW(), 'opt_out') ON CONFLICT DO NOTHING", s.cfg.Tables.DNCEntries)
	_, err := s.db.ExecContext(ctx, query, phoneNumber)
	return err
}

func (s *postgresStore) CleanupTranscripts(ctx context.Context, retentionYears int) error {
	s.logger.Info("Running transcript retention cleanup", zap.Int("retention_years", retentionYears))
	query := fmt.Sprintf("DELETE FROM %s WHERE created_at < NOW() - INTERVAL '%d years'", s.cfg.Tables.CallRecords, retentionYears)
	_, err := s.db.ExecContext(ctx, query)
	return err
}
