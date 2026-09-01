package compliance

import (
	"context"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/storage"
	"go.uber.org/zap"
)

// Checker defines the interface for compliance checks
type Checker interface {
	// IsAllowed checks if a phone number is allowed to be called right now
	IsAllowed(ctx context.Context, phoneNumber string) (bool, error)
	// HandleOptOut immediately adds a user to the DNC registry (T10.2)
	HandleOptOut(ctx context.Context, phoneNumber string) error
}

type checkerImpl struct {
	cfg     config.CampaignManagerComplianceConfig
	windows config.CampaignManagerCallingWindowConfig
	store   storage.Storage
	logger  *zap.Logger
}

// NewChecker initializes a new compliance checker
func NewChecker(cfg config.CampaignManagerComplianceConfig, windows config.CampaignManagerCallingWindowConfig, store storage.Storage, logger *zap.Logger) Checker {
	return &checkerImpl{
		cfg:     cfg,
		windows: windows,
		store:   store,
		logger:  logger,
	}
}

func (c *checkerImpl) IsAllowed(ctx context.Context, phoneNumber string) (bool, error) {
	// Check calling window first (cheaper)
	if c.windows.Enabled {
		allowed, err := c.checkCallingWindow()
		if err != nil {
			c.logger.Error("Failed to check calling window", zap.Error(err))
			return false, err
		}
		if !allowed {
			return false, nil
		}
	}

	// Check DNC database
	if c.cfg.DNCCheckEnabled {
		isDNC, err := c.store.Postgres().IsDNC(ctx, phoneNumber)
		if err != nil {
			c.logger.Error("Failed to check DNC registry", zap.Error(err))
			return false, err
		}
		if isDNC {
			return false, nil
		}
	}

	return true, nil
}

func (c *checkerImpl) checkCallingWindow() (bool, error) {
	loc, err := time.LoadLocation(c.windows.TimezoneDB)
	if err != nil {
		// Fallback to local time if timezone is invalid
		loc = time.Local
	}

	now := time.Now().In(loc)
	currentTime := now.Format("15:04")

	if currentTime >= c.windows.DefaultStart && currentTime <= c.windows.DefaultEnd {
		return true, nil
	}

	return false, nil
}

func (c *checkerImpl) HandleOptOut(ctx context.Context, phoneNumber string) error {
	c.logger.Info("Processing compliance opt-out request", zap.String("phone", phoneNumber))
	return c.store.Postgres().OptOut(ctx, phoneNumber)
}
