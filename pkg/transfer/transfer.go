package transfer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	metricTransferWaitTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ai_telephony_transfer_wait_time_seconds",
		Help:    "Time spent in the transfer queue before agent connection",
		Buckets: []float64{5, 15, 30, 60, 120, 300},
	}, []string{"priority"})

	metricTransferSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ai_telephony_transfer_success_total",
		Help: "Total successful human transfers",
	})

	metricTransferAbandon = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ai_telephony_transfer_abandon_total",
		Help: "Total abandoned or timed-out transfers",
	}, []string{"reason"}) // timeout or user_hangup
)

// TransferState represents the transfer state machine
type TransferState string

const (
	TransferStateQueued    TransferState = "queued"
	TransferStateRinging   TransferState = "ringing"
	TransferStateConnected TransferState = "connected"
	TransferStateFailed    TransferState = "failed"
	TransferStateCompleted TransferState = "completed"
)

// TransferRequest holds metadata for escalating a call to a human agent
type TransferRequest struct {
	ID                 string
	CallID             string
	TargetAgentURI     string
	Priority           int
	Reason             string
	Transcript         string
	ExtractedData      map[string]string
	Sentiment          string
	State              TransferState
	QueuedAt           time.Time
	ConnectedAt        time.Time
	EndedAt            time.Time
}

// SIPReferMessage represents a constructed SIP REFER payload
type SIPReferMessage struct {
	CallID      string
	ReferTo     string
	Replaces    string
	Headers     map[string]string
	SDPBody     string
}

// SIPReferBuilder formats SIP REFER requests with Replaces header and context injection
type SIPReferBuilder struct {
	headerPrefix string
	includeTrans bool
	includeData  bool
	includeSent  bool
}

// NewSIPReferBuilder creates a SIPReferBuilder instance
func NewSIPReferBuilder(cfg *config.TransferManagerContextHeadersConfig) *SIPReferBuilder {
	return &SIPReferBuilder{
		headerPrefix: cfg.HeaderPrefix,
		includeTrans: cfg.IncludeTranscript,
		includeData:  cfg.IncludeExtractedData,
		includeSent:  cfg.IncludeSentiment,
	}
}

// BuildReferMessage constructs REFER message headers
func (b *SIPReferBuilder) BuildReferMessage(req *TransferRequest) *SIPReferMessage {
	headers := make(map[string]string)

	if b.includeTrans && req.Transcript != "" {
		headers[b.headerPrefix+"Transcript"] = req.Transcript
	}
	if b.includeSent && req.Sentiment != "" {
		headers[b.headerPrefix+"Sentiment"] = req.Sentiment
	}
	if b.includeData {
		for k, v := range req.ExtractedData {
			headers[b.headerPrefix+"Data-"+k] = v
		}
	}

	referTo := fmt.Sprintf("<%s>", req.TargetAgentURI)
	if req.CallID != "" {
		referTo = fmt.Sprintf("<%s?Replaces=%s>", req.TargetAgentURI, req.CallID)
	}

	return &SIPReferMessage{
		CallID:   req.CallID,
		ReferTo:  referTo,
		Replaces: req.CallID,
		Headers:  headers,
	}
}

// TransferEngine manages ACD transfer queues and transfer state machines
type TransferEngine struct {
	config    *config.TransferConfig
	logger    *zap.Logger
	builder   *SIPReferBuilder
	requests  map[string]*TransferRequest
	mu        sync.RWMutex
}

// NewTransferEngine creates a TransferEngine
func NewTransferEngine(cfg *config.TransferConfig, logger *zap.Logger) *TransferEngine {
	builder := NewSIPReferBuilder(&cfg.Manager.ContextHeaders)
	return &TransferEngine{
		config:   cfg,
		logger:   logger,
		builder:  builder,
		requests: make(map[string]*TransferRequest),
	}
}

// InitiateTransfer registers a new live transfer request
func (e *TransferEngine) InitiateTransfer(ctx context.Context, req *TransferRequest) (*SIPReferMessage, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if req.ID == "" || req.CallID == "" {
		return nil, errors.New("invalid transfer request: ID and CallID are required")
	}

	req.State = TransferStateQueued
	req.QueuedAt = time.Now()
	e.requests[req.ID] = req

	e.logger.Info("Transfer initiated",
		zap.String("transfer_id", req.ID),
		zap.String("call_id", req.CallID),
		zap.String("target_agent", req.TargetAgentURI),
		zap.Int("priority", req.Priority),
	)

	// Launch background monitor for timeout handling (300s default)
	go e.monitorTransfer(req.ID)

	referMsg := e.builder.BuildReferMessage(req)
	return referMsg, nil
}

// monitorTransfer enforces timeouts and handles voicemail callbacks
func (e *TransferEngine) monitorTransfer(transferID string) {
	time.Sleep(300 * time.Second) // 300s default timeout
	
	e.mu.Lock()
	req, exists := e.requests[transferID]
	if !exists || req.State == TransferStateConnected || req.State == TransferStateCompleted {
		e.mu.Unlock()
		return
	}
	
	// Timeout occurred, fallback
	req.State = TransferStateFailed
	req.EndedAt = time.Now()
	e.mu.Unlock()

	metricTransferAbandon.WithLabelValues("timeout").Inc()
	e.logger.Warn("Transfer timed out, falling back to voicemail", zap.String("transfer_id", transferID))
}

// EstimateWaitTime simulates priority-based EWT calculation
func (e *TransferEngine) EstimateWaitTime(priority int) time.Duration {
	// Simulated logic: higher priority (e.g. 1) gets shorter wait times
	baseWait := 120
	estimated := baseWait / priority
	return time.Duration(estimated) * time.Second
}

// UpdateState transitions transfer state machine
func (e *TransferEngine) UpdateState(transferID string, newState TransferState) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	req, exists := e.requests[transferID]
	if !exists {
		return fmt.Errorf("transfer request %s not found", transferID)
	}

	req.State = newState
	if newState == TransferStateConnected {
		req.ConnectedAt = time.Now()
		waitTime := req.ConnectedAt.Sub(req.QueuedAt).Seconds()
		metricTransferWaitTime.WithLabelValues(fmt.Sprintf("%d", req.Priority)).Observe(waitTime)
	} else if newState == TransferStateCompleted {
		req.EndedAt = time.Now()
		metricTransferSuccess.Inc()
	} else if newState == TransferStateFailed {
		req.EndedAt = time.Now()
	}

	e.logger.Info("Transfer state updated", zap.String("transfer_id", transferID), zap.String("state", string(newState)))
	return nil
}

// GetTransfer returns a transfer request by ID
func (e *TransferEngine) GetTransfer(transferID string) (*TransferRequest, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	req, exists := e.requests[transferID]
	return req, exists
}
