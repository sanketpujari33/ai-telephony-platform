package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"
	"encoding/json"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/llm"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/stt"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/tts"
	"go.uber.org/zap"
)

// SessionState represents the state of an active AI call session
type SessionState int32

const (
	SessionStateConnecting SessionState = iota
	SessionStateConversing
	SessionStateTransferring
	SessionStateEnding
	SessionStateEnded
)

func (s SessionState) String() string {
	switch s {
	case SessionStateConnecting:
		return "connecting"
	case SessionStateConversing:
		return "conversing"
	case SessionStateTransferring:
		return "transferring"
	case SessionStateEnding:
		return "ending"
	case SessionStateEnded:
		return "ended"
	default:
		return "unknown"
	}
}

// CallSession manages a full-duplex AI conversation pipeline
type CallSession struct {
	ID            string
	CallID        string
	TenantID      string
	State         SessionState
	ContextWindow *llm.ContextWindow
	sttRouter     *stt.STTRouter
	llmRouter     *llm.LLMRouter
	ttsRouter     *tts.TTSRouter
	logger        *zap.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	mu            sync.RWMutex
	bargeInActive bool
	startTime     time.Time
	metadata      CallMetadata
}

// CallMetadata tracks analytics for persistence
type CallMetadata struct {
	TotalCost       float64
	AverageLatency  float64
	Disposition     string
	Escalated       bool
	TransferReason  string
}

// ConversationPersister handles S3/PostgreSQL storage
type ConversationPersister struct {
	logger *zap.Logger
}

func (p *ConversationPersister) SaveTranscriptToPG(sessionID string, messages []llm.ChatMessage, meta CallMetadata) error {
	p.logger.Info("Persisting transcript to PostgreSQL", zap.String("session_id", sessionID), zap.Int("turns", len(messages)))
	// Simulate pgx insertion
	return nil
}

func (p *ConversationPersister) UploadDualChannelToS3(sessionID string) error {
	p.logger.Info("Uploading dual-channel recording to S3", zap.String("session_id", sessionID))
	// Simulate S3 multipart upload
	return nil
}

// NewCallSession creates a new CallSession
func NewCallSession(
	id, callID, tenantID string,
	cfg *config.Config,
	sttRouter *stt.STTRouter,
	llmRouter *llm.LLMRouter,
	ttsRouter *tts.TTSRouter,
	logger *zap.Logger,
) *CallSession {
	ctx, cancel := context.WithCancel(context.Background())

	sysPrompt := "You are a polite, helpful AI telephony assistant."
	cw := llm.NewContextWindow(
		sysPrompt,
		cfg.LLM.Context.SlidingWindow.MaxTurns,
		cfg.LLM.Context.SlidingWindow.MaxTokens,
		cfg.LLM.Context.TokenBudget.TotalBudget,
	)

	return &CallSession{
		ID:            id,
		CallID:        callID,
		TenantID:      tenantID,
		State:         SessionStateConnecting,
		ContextWindow: cw,
		sttRouter:     sttRouter,
		llmRouter:     llmRouter,
		ttsRouter:     ttsRouter,
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Start initiates the audio processing and speech-llm pipeline
func (s *CallSession) Start() error {
	s.mu.Lock()
	if s.State != SessionStateConnecting {
		s.mu.Unlock()
		return fmt.Errorf("cannot start session in state %s", s.State)
	}
	s.State = SessionStateConversing
	s.mu.Unlock()

	s.logger.Info("CallSession started conversing", zap.String("session_id", s.ID), zap.String("call_id", s.CallID))
	return nil
}

// ProcessTranscript handles an incoming user transcript (interim or final)
func (s *CallSession) ProcessTranscript(res stt.TranscriptResult) error {
	s.mu.Lock()
	if s.State != SessionStateConversing {
		s.mu.Unlock()
		return nil
	}

	// Barge-in detection: if user speaks while assistant is generating/playing TTS
	if !res.IsFinal && len(res.Text) > 2 {
		s.bargeInActive = true
		s.logger.Debug("Barge-in detected, interrupting assistant playback", zap.String("session_id", s.ID))
	}

	if !res.IsFinal {
		s.mu.Unlock()
		return nil
	}

	s.bargeInActive = false
	s.mu.Unlock()

	s.logger.Info("User input received", zap.String("session_id", s.ID), zap.String("text", res.Text))
	s.ContextWindow.AddUserTurn(res.Text, len(res.Text)/4)

	// Trigger LLM response generation
	return s.generateAssistantResponse()
}

func (s *CallSession) generateAssistantResponse() error {
	provider, err := s.llmRouter.SelectProvider()
	if err != nil {
		return fmt.Errorf("failed to select LLM provider: %w", err)
	}

	messages := s.ContextWindow.GetMessages()
	stream, err := provider.GenerateStream(s.ctx, messages)
	if err != nil {
		return fmt.Errorf("LLM streaming failed: %w", err)
	}

	go func() {
		var fullResp string
		for resp := range stream {
			if s.bargeInActive {
				s.logger.Info("Barge-in active, aborting LLM generation", zap.String("session_id", s.ID))
				break
			}
			fullResp += resp.Content
		}

		if fullResp != "" && !s.bargeInActive {
			s.ContextWindow.AddAssistantTurn(fullResp, len(fullResp)/4)
			s.logger.Info("Assistant response generated", zap.String("session_id", s.ID), zap.String("text", fullResp))
		}
	}()

	return nil
}

// Stop ends the session gracefully and triggers persistence
func (s *CallSession) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.State = SessionStateEnded
	s.cancel()
	s.logger.Info("CallSession stopped", zap.String("session_id", s.ID))
	
	// Fire persistence in background
	go s.PersistSessionData()
	
	return nil
}

// PersistSessionData flushes the session transcript and dual-channel recordings to storage
func (s *CallSession) PersistSessionData() {
	persister := &ConversationPersister{logger: s.logger}
	
	// 1. Gather Transcript
	messages := s.ContextWindow.GetMessages()
	
	// 2. Gather Metadata
	s.metadata.Disposition = "completed"
	if s.metadata.Escalated {
		s.metadata.Disposition = "transferred"
	}
	
	// 3. Save to PG (JSON format)
	transcriptJSON, _ := json.Marshal(messages)
	s.logger.Debug("Serialized transcript", zap.Int("bytes", len(transcriptJSON)))
	err := persister.SaveTranscriptToPG(s.ID, messages, s.metadata)
	if err != nil {
		s.logger.Error("Failed to save transcript to PostgreSQL", zap.Error(err))
	}
	
	// 4. Upload mixed audio to S3
	err = persister.UploadDualChannelToS3(s.ID)
	if err != nil {
		s.logger.Error("Failed to upload recording to S3", zap.Error(err))
	}
	
	s.logger.Info("CallSession persistence completed successfully", zap.String("session_id", s.ID))
}

// OrchestratorServer manages call sessions
type OrchestratorServer struct {
	config    *config.Config
	logger    *zap.Logger
	sttRouter *stt.STTRouter
	llmRouter *llm.LLMRouter
	ttsRouter *tts.TTSRouter
	sessions  map[string]*CallSession
	mu        sync.RWMutex
}

// NewOrchestratorServer creates an OrchestratorServer instance
func NewOrchestratorServer(cfg *config.Config, logger *zap.Logger) *OrchestratorServer {
	return &OrchestratorServer{
		config:    cfg,
		logger:    logger,
		sttRouter: stt.NewSTTRouter(&cfg.STT, logger),
		llmRouter: llm.NewLLMRouter(&cfg.LLM, logger),
		ttsRouter: tts.NewTTSRouter(&cfg.TTS, logger),
		sessions:  make(map[string]*CallSession),
	}
}

// CreateCallSession creates a new call session
func (os *OrchestratorServer) CreateCallSession(sessionID, callID, tenantID string) (*CallSession, error) {
	os.mu.Lock()
	defer os.mu.Unlock()

	session := NewCallSession(sessionID, callID, tenantID, os.config, os.sttRouter, os.llmRouter, os.ttsRouter, os.logger)
	os.sessions[sessionID] = session
	return session, nil
}
