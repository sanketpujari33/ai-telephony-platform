package campaign

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/compliance"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"go.uber.org/zap"
)

// ContactDisposition represents the outcome of a dial attempt
type ContactDisposition string

const (
	DispositionVoicemailLeft ContactDisposition = "voicemail_left"
	DispositionAnswered      ContactDisposition = "answered"
	DispositionBusy          ContactDisposition = "busy"
	DispositionNoAnswer      ContactDisposition = "no_answer"
	DispositionPending       ContactDisposition = "pending"
)

// Contact represents a campaign target contact
type Contact struct {
	ID              string
	PhoneNumber     string
	Name            string
	Timezone        string
	ConsentVerified bool
	LastCalled      time.Time
	Disposition     ContactDisposition
	RetryCount      int
}

// CampaignStatus represents the status of an outbound campaign
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCompleted CampaignStatus = "completed"
)

// Campaign represents an outbound calling campaign
type Campaign struct {
	ID           string
	TenantID     string
	Name         string
	Status       CampaignStatus
	Contacts     []Contact
	TargetAgents int
	CreatedAt    time.Time
}

// CampaignEngine manages campaign loading and dispatching
type CampaignEngine struct {
	config     *config.CampaignConfig
	logger     *zap.Logger
	checker    compliance.Checker
	campaigns  map[string]*Campaign
	mu         sync.RWMutex
}

// NewCampaignEngine creates a CampaignEngine instance
func NewCampaignEngine(cfg *config.CampaignConfig, checker compliance.Checker, logger *zap.Logger) *CampaignEngine {
	return &CampaignEngine{
		config:     cfg,
		logger:     logger,
		checker:    checker,
		campaigns:  make(map[string]*Campaign),
	}
}

// RegisterCampaign registers a new campaign
func (e *CampaignEngine) RegisterCampaign(c *Campaign) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.campaigns[c.ID] = c
	e.logger.Info("Campaign registered", zap.String("campaign_id", c.ID), zap.Int("contacts", len(c.Contacts)))
}

// GetCallableContacts returns compliant contacts for a campaign
func (e *CampaignEngine) GetCallableContacts(ctx context.Context, campaignID string) ([]Contact, error) {
	e.mu.RLock()
	campaign, exists := e.campaigns[campaignID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("campaign %s not found", campaignID)
	}

	callable := make([]Contact, 0)
	for _, c := range campaign.Contacts {
		if !c.ConsentVerified {
			e.logger.Debug("Contact filtered by compliance", zap.String("phone", c.PhoneNumber), zap.String("reason", "missing explicit consent"))
			continue
		}
		
		ok, err := e.checker.IsAllowed(ctx, c.PhoneNumber)
		if err != nil {
			e.logger.Error("Failed compliance check", zap.Error(err))
			continue
		}
		
		if ok {
			callable = append(callable, c)
		} else {
			e.logger.Debug("Contact filtered by compliance", zap.String("phone", c.PhoneNumber), zap.String("reason", "failed DNC or calling window"))
		}
	}

	return callable, nil
}

// Checker returns the internal compliance checker
func (e *CampaignEngine) Checker() compliance.Checker {
	return e.checker
}

// ProcessOutboundResult handles the AMD outcome of a dialed contact
func (e *CampaignEngine) ProcessOutboundResult(campaignID, contactID string, amd AMDResult, isBusy, isNoAnswer bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	campaign, exists := e.campaigns[campaignID]
	if !exists {
		return
	}

	for i := range campaign.Contacts {
		if campaign.Contacts[i].ID == contactID {
			c := &campaign.Contacts[i]
			c.LastCalled = time.Now()

			// Disposition mapping based on telephony/AMD signals
			if isBusy {
				c.Disposition = DispositionBusy
				e.scheduleRetry(c)
			} else if isNoAnswer {
				c.Disposition = DispositionNoAnswer
				e.scheduleRetry(c)
			} else {
				// Answered
				if amd == AMDResultMachine {
					c.Disposition = DispositionVoicemailLeft
					e.logger.Info("AMD detected Voicemail, triggering voicemail playback branch", zap.String("contact_id", c.ID))
					// In a real app, this triggers Media Proxy to play the voicemail track
				} else {
					c.Disposition = DispositionAnswered
					e.logger.Info("AMD detected Human, routing to Orchestrator", zap.String("contact_id", c.ID))
					// In a real app, this spawns an orchestrator session
				}
			}
			
			e.logger.Debug("Updated contact disposition", zap.String("contact_id", c.ID), zap.String("disposition", string(c.Disposition)))
			break
		}
	}
}

func (e *CampaignEngine) scheduleRetry(c *Contact) {
	if c.RetryCount < e.config.PredictiveDialer.Retry.MaxAttempts { // Just reusing a config stub for max retries
		c.RetryCount++
		c.Disposition = DispositionPending
		e.logger.Info("Contact scheduled for retry", zap.String("contact_id", c.ID), zap.Int("attempt", c.RetryCount))
	} else {
		e.logger.Warn("Contact exhausted retries", zap.String("contact_id", c.ID))
	}
}
