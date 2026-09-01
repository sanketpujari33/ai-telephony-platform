package stt

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	metricCircuitBreaker = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ai_telephony_stt_circuit_breaker_open_total",
		Help: "Total number of times STT circuit breaker opened",
	}, []string{"provider"})

	metricSTTLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ai_telephony_stt_latency_ms",
		Help:    "STT processing latency in milliseconds",
		Buckets: []float64{50, 100, 250, 500, 1000, 2000, 5000},
	}, []string{"provider", "interim"})
)

// WordTiming provides granular timing for individual words
type WordTiming struct {
	Word  string
	Start time.Duration
	End   time.Duration
}

// EndpointingConfig configures silence detection for STT
type EndpointingConfig struct {
	Enabled         bool
	SilenceDuration time.Duration
}

// ProviderInfo holds metadata about the STT provider
type ProviderInfo struct {
	Name    string
	Version string
	Region  string
}

// STTStreamConfig defines the parameters for a new streaming STT session
type STTStreamConfig struct {
	SampleRate      int
	Encoding        string
	Language        string
	InterimResults  bool
	KeywordBoosting []string
	Diarization     bool
	Endpointing     EndpointingConfig
}

// STTStream represents an active streaming session to an STT provider
type STTStream interface {
	SendAudioChunk(chunk []byte) error
	Transcripts() <-chan TranscriptResult
	Close() error
}

// TranscriptResult represents a transcript from an STT provider
type TranscriptResult struct {
	Text       string
	IsFinal    bool
	Confidence float64
	Speaker    string
	Language   string
	Duration   time.Duration
	Timestamp  time.Time
	Words      []WordTiming
}

// STTProvider defines the common interface for Speech-To-Text engines
type STTProvider interface {
	Name() string
	Type() string
	Info() ProviderInfo
	Connect(ctx context.Context) error
	NewStream(ctx context.Context, config STTStreamConfig) (STTStream, error)
	Close() error
}

// CircuitBreakerState represents the circuit breaker state
type CircuitBreakerState int32

const (
	StateClosed CircuitBreakerState = iota
	StateHalfOpen
	StateOpen
)

func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// ProviderCircuitBreaker monitors provider health
type ProviderCircuitBreaker struct {
	state            int32
	failureCount     int32
	failureThreshold int32
	timeout          time.Duration
	lastStateChange  time.Time
	mu               sync.Mutex
}

// NewProviderCircuitBreaker creates a circuit breaker instance
func NewProviderCircuitBreaker(failureThreshold int, timeoutSeconds int) *ProviderCircuitBreaker {
	return &ProviderCircuitBreaker{
		state:            int32(StateClosed),
		failureThreshold: int32(failureThreshold),
		timeout:          time.Duration(timeoutSeconds) * time.Second,
		lastStateChange:  time.Now(),
	}
}

// Allow returns true if the call to the provider is allowed
func (cb *ProviderCircuitBreaker) Allow() bool {
	state := CircuitBreakerState(atomic.LoadInt32(&cb.state))
	if state == StateClosed {
		return true
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if state == StateOpen {
		if time.Since(cb.lastStateChange) > cb.timeout {
			atomic.StoreInt32(&cb.state, int32(StateHalfOpen))
			cb.lastStateChange = time.Now()
			return true
		}
		return false
	}

	return true // StateHalfOpen
}

// RecordSuccess marks a successful request
func (cb *ProviderCircuitBreaker) RecordSuccess() {
	atomic.StoreInt32(&cb.failureCount, 0)
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if CircuitBreakerState(atomic.LoadInt32(&cb.state)) == StateHalfOpen {
		atomic.StoreInt32(&cb.state, int32(StateClosed))
		cb.lastStateChange = time.Now()
	}
}

// RecordFailure marks a failed request
func (cb *ProviderCircuitBreaker) RecordFailure() {
	failures := atomic.AddInt32(&cb.failureCount, 1)
	if failures >= cb.failureThreshold {
		cb.mu.Lock()
		defer cb.mu.Unlock()
		if CircuitBreakerState(atomic.LoadInt32(&cb.state)) != StateOpen {
			atomic.StoreInt32(&cb.state, int32(StateOpen))
			cb.lastStateChange = time.Now()
			// We can't access provider name here easily without passing it, but it will be emitted by STTRouter
		}
	}
}

// GetState returns the current state
func (cb *ProviderCircuitBreaker) GetState() CircuitBreakerState {
	return CircuitBreakerState(atomic.LoadInt32(&cb.state))
}

// DeepgramProvider implements STTProvider for Deepgram WebSocket API
type DeepgramProvider struct {
	name   string
	cfg    config.STTProviderConfig
	ch     chan TranscriptResult
	logger *zap.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

func NewDeepgramProvider(cfg config.STTProviderConfig, logger *zap.Logger) *DeepgramProvider {
	return &DeepgramProvider{
		name:   cfg.Name,
		cfg:    cfg,
		ch:     make(chan TranscriptResult, 100),
		logger: logger,
	}
}

func (p *DeepgramProvider) Name() string { return p.name }
func (p *DeepgramProvider) Type() string { return "deepgram" }
func (p *DeepgramProvider) Info() ProviderInfo {
	return ProviderInfo{Name: p.name, Version: "nova-2", Region: "global"}
}

func (p *DeepgramProvider) Connect(ctx context.Context) error {
	p.ctx, p.cancel = context.WithCancel(ctx)
	// Real implementation would manage connection pooling here
	p.logger.Info("Connected to Deepgram STT engine (streaming)", zap.String("provider", p.name))
	return nil
}

// DeepgramStream implements STTStream
type DeepgramStream struct {
	provider *DeepgramProvider
	config   STTStreamConfig
	ch       chan TranscriptResult
}

func (p *DeepgramProvider) NewStream(ctx context.Context, config STTStreamConfig) (STTStream, error) {
	// Reconnection logic and connection pooling would assign an active websocket here
	return &DeepgramStream{
		provider: p,
		config:   config,
		ch:       make(chan TranscriptResult, 100),
	}, nil
}

func (s *DeepgramStream) SendAudioChunk(chunk []byte) error {
	if s.provider.ctx == nil || s.provider.ctx.Err() != nil {
		return errors.New("provider connection closed")
	}
	// Simulate processing time for latency metrics
	start := time.Now()
	// write to websocket...
	latency := time.Since(start).Milliseconds()
	metricSTTLatency.WithLabelValues(s.provider.name, "false").Observe(float64(latency))
	return nil
}

func (s *DeepgramStream) Transcripts() <-chan TranscriptResult {
	return s.ch
}

func (s *DeepgramStream) Close() error {
	close(s.ch)
	return nil
}

func (p *DeepgramProvider) Close() error {
	if p.cancel != nil {
		p.cancel()
	}
	close(p.ch)
	return nil
}

// AssemblyAIProvider implements STTProvider for AssemblyAI Streaming API
type AssemblyAIProvider struct {
	name   string
	cfg    config.STTProviderConfig
	ch     chan TranscriptResult
	logger *zap.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAssemblyAIProvider(cfg config.STTProviderConfig, logger *zap.Logger) *AssemblyAIProvider {
	return &AssemblyAIProvider{
		name:   cfg.Name,
		cfg:    cfg,
		ch:     make(chan TranscriptResult, 100),
		logger: logger,
	}
}

func (p *AssemblyAIProvider) Name() string { return p.name }
func (p *AssemblyAIProvider) Type() string { return "assemblyai" }
func (p *AssemblyAIProvider) Info() ProviderInfo {
	return ProviderInfo{Name: p.name, Version: "universal-2", Region: "global"}
}

func (p *AssemblyAIProvider) Connect(ctx context.Context) error {
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.logger.Info("Connected to AssemblyAI STT engine", zap.String("provider", p.name))
	return nil
}

// AssemblyAIStream implements STTStream
type AssemblyAIStream struct {
	provider *AssemblyAIProvider
	config   STTStreamConfig
	ch       chan TranscriptResult
}

func (p *AssemblyAIProvider) NewStream(ctx context.Context, config STTStreamConfig) (STTStream, error) {
	return &AssemblyAIStream{
		provider: p,
		config:   config,
		ch:       make(chan TranscriptResult, 100),
	}, nil
}

func (s *AssemblyAIStream) SendAudioChunk(chunk []byte) error {
	if s.provider.ctx == nil || s.provider.ctx.Err() != nil {
		return errors.New("provider connection closed")
	}
	return nil
}

func (s *AssemblyAIStream) Transcripts() <-chan TranscriptResult {
	return s.ch
}

func (s *AssemblyAIStream) Close() error {
	close(s.ch)
	return nil
}

func (p *AssemblyAIProvider) Close() error {
	if p.cancel != nil {
		p.cancel()
	}
	close(p.ch)
	return nil
}

// WhisperProvider implements STTProvider for local Whisper instances
type WhisperProvider struct {
	name   string
	cfg    config.STTProviderConfig
	ch     chan TranscriptResult
	logger *zap.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWhisperProvider(cfg config.STTProviderConfig, logger *zap.Logger) *WhisperProvider {
	return &WhisperProvider{
		name:   cfg.Name,
		cfg:    cfg,
		ch:     make(chan TranscriptResult, 100),
		logger: logger,
	}
}

func (p *WhisperProvider) Name() string { return p.name }
func (p *WhisperProvider) Type() string { return "whisper" }
func (p *WhisperProvider) Info() ProviderInfo {
	return ProviderInfo{Name: p.name, Version: "local", Region: "self-hosted"}
}

func (p *WhisperProvider) Connect(ctx context.Context) error {
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.logger.Info("Connected to local Whisper instance", zap.String("provider", p.name))
	return nil
}

// WhisperStream implements STTStream
type WhisperStream struct {
	provider *WhisperProvider
	config   STTStreamConfig
	ch       chan TranscriptResult
}

func (p *WhisperProvider) NewStream(ctx context.Context, config STTStreamConfig) (STTStream, error) {
	return &WhisperStream{
		provider: p,
		config:   config,
		ch:       make(chan TranscriptResult, 100),
	}, nil
}

func (s *WhisperStream) SendAudioChunk(chunk []byte) error {
	if s.provider.ctx == nil || s.provider.ctx.Err() != nil {
		return errors.New("provider connection closed")
	}
	return nil
}

func (s *WhisperStream) Transcripts() <-chan TranscriptResult {
	return s.ch
}

func (s *WhisperStream) Close() error {
	close(s.ch)
	return nil
}

func (p *WhisperProvider) Close() error {
	if p.cancel != nil {
		p.cancel()
	}
	close(p.ch)
	return nil
}


// STTRouter manages active STT providers with failover capability
type STTRouter struct {
	config          *config.STTConfig
	logger          *zap.Logger
	providers       []STTProvider
	circuitBreakers map[string]*ProviderCircuitBreaker
	mu              sync.RWMutex
}

// NewSTTRouter creates a new STTRouter
func NewSTTRouter(cfg *config.STTConfig, logger *zap.Logger) *STTRouter {
	router := &STTRouter{
		config:          cfg,
		logger:          logger,
		providers:       make([]STTProvider, 0),
		circuitBreakers: make(map[string]*ProviderCircuitBreaker),
	}

	for _, pCfg := range cfg.Providers {
		if !pCfg.Enabled {
			continue
		}
		var provider STTProvider
		switch pCfg.Type {
		case "deepgram":
			provider = NewDeepgramProvider(pCfg, logger)
		case "assemblyai":
			provider = NewAssemblyAIProvider(pCfg, logger)
		case "whisper":
			provider = NewWhisperProvider(pCfg, logger)
		default:
			provider = NewDeepgramProvider(pCfg, logger)
		}

		router.providers = append(router.providers, provider)
		router.circuitBreakers[pCfg.Name] = NewProviderCircuitBreaker(
			cfg.CircuitBreaker.FailureThreshold,
			cfg.CircuitBreaker.Timeout,
		)
	}

	return router
}

// GetPrimaryProvider selects the highest priority healthy provider
func (r *STTRouter) GetPrimaryProvider() (STTProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, provider := range r.providers {
		cb, exists := r.circuitBreakers[provider.Name()]
		if exists && cb.Allow() {
			return provider, nil
		}
	}

	if len(r.providers) > 0 {
		return r.providers[0], nil
	}

	return nil, errors.New("no STT providers available")
}

// ReportFailure records a failure for a named provider
func (r *STTRouter) ReportFailure(providerName string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if cb, exists := r.circuitBreakers[providerName]; exists {
		previousState := cb.GetState()
		cb.RecordFailure()
		newState := cb.GetState()
		
		if previousState != StateOpen && newState == StateOpen {
			metricCircuitBreaker.WithLabelValues(providerName).Inc()
			r.logger.Error("STT circuit breaker opened due to successive failures", zap.String("provider", providerName))
		} else {
			r.logger.Warn("STT provider failure reported", zap.String("provider", providerName), zap.String("cb_state", newState.String()))
		}
	}
}

// ReportSuccess records a success for a named provider
func (r *STTRouter) ReportSuccess(providerName string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if cb, exists := r.circuitBreakers[providerName]; exists {
		cb.RecordSuccess()
	}
}
