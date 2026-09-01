package tts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	metricTTSCacheHit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ai_telephony_tts_cache_hit_total",
		Help: "Total number of TTS cache hits",
	}, []string{"provider"})

	metricTTSCacheMiss = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ai_telephony_tts_cache_miss_total",
		Help: "Total number of TTS cache misses",
	}, []string{"provider"})

	metricTTSLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ai_telephony_tts_first_byte_latency_ms",
		Help:    "TTS time-to-first-byte latency in milliseconds",
		Buckets: []float64{100, 250, 500, 750, 1000, 2000},
	}, []string{"provider"})
)

// Voice represents a TTS voice configuration
type Voice struct {
	ID       string
	Name     string
	Gender   string
	Language string
	Provider string
}

// TTSRequest represents a text-to-speech synthesis request
type TTSRequest struct {
	Text         string
	VoiceID      string
	IsSSML       bool
	LanguageCode string
}

// AudioChunk represents a chunk of PCM audio from TTS
type AudioChunk struct {
	PCMData    []byte
	SampleRate int
	Channels   int
	IsFinal    bool
}

// TTSStream represents an active streaming synthesis session
type TTSStream interface {
	AudioChunks() <-chan AudioChunk
	Close() error
}

// TTSProvider defines the interface for text-to-speech engines
type TTSProvider interface {
	Name() string
	SynthesizeStream(ctx context.Context, req TTSRequest) (TTSStream, error)
	Synthesize(ctx context.Context, req TTSRequest) ([]byte, error)
	SyncVoices(ctx context.Context) ([]Voice, error)
}

// TTSCache provides audio caching in Redis or memory
type TTSCache struct {
	cache  map[string][]byte
	logger *zap.Logger
	mu     sync.RWMutex
}

// NewTTSCache creates a new TTS cache manager
func NewTTSCache(logger *zap.Logger) *TTSCache {
	return &TTSCache{
		cache:  make(map[string][]byte),
		logger: logger,
	}
}

// CacheKey generates a hash key for a prompt text and voice ID
func (c *TTSCache) CacheKey(text, voiceID string) string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s", voiceID, text)))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Get retrieves cached audio bytes
func (c *TTSCache) Get(text, voiceID string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.CacheKey(text, voiceID)
	data, exists := c.cache[key]
	return data, exists
}

// Set stores audio bytes into cache
func (c *TTSCache) Set(text, voiceID string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.CacheKey(text, voiceID)
	c.cache[key] = data
}

// ElevenLabsProvider stub implementation
type ElevenLabsProvider struct {
	name   string
	cfg    config.TTSProviderConfig
	logger *zap.Logger
}

func NewElevenLabsProvider(cfg config.TTSProviderConfig, logger *zap.Logger) *ElevenLabsProvider {
	return &ElevenLabsProvider{
		name:   cfg.Name,
		cfg:    cfg,
		logger: logger,
	}
}

func (p *ElevenLabsProvider) Name() string { return p.name }

func (p *ElevenLabsProvider) SyncVoices(ctx context.Context) ([]Voice, error) {
	// Real implementation would call GET https://api.elevenlabs.io/v1/voices
	p.logger.Info("Synchronizing ElevenLabs voice catalog")
	return []Voice{{ID: "rachel_v2", Name: "Rachel", Gender: "female", Language: "en-US", Provider: p.name}}, nil
}

func (p *ElevenLabsProvider) Synthesize(ctx context.Context, req TTSRequest) ([]byte, error) {
	// Real implementation calls HTTP endpoint and returns full payload
	return make([]byte, 16000), nil
}

type ElevenLabsStream struct {
	ch chan AudioChunk
}
func (s *ElevenLabsStream) AudioChunks() <-chan AudioChunk { return s.ch }
func (s *ElevenLabsStream) Close() error { close(s.ch); return nil }

func (p *ElevenLabsProvider) SynthesizeStream(ctx context.Context, req TTSRequest) (TTSStream, error) {
	stream := &ElevenLabsStream{ch: make(chan AudioChunk, 10)}
	
	go func() {
		defer stream.Close()
		start := time.Now()
		
		// Simulate network latency for first byte
		time.Sleep(300 * time.Millisecond)
		metricTTSLatency.WithLabelValues(p.name).Observe(float64(time.Since(start).Milliseconds()))
		
		stream.ch <- AudioChunk{
			PCMData:    make([]byte, 320),
			SampleRate: 16000,
			Channels:   1,
			IsFinal:    true,
		}
	}()
	return stream, nil
}

// AzureTTSAdapter implements TTSProvider for Azure Neural TTS
type AzureTTSAdapter struct {
	name   string
	cfg    config.TTSProviderConfig
	logger *zap.Logger
}

func NewAzureTTSAdapter(cfg config.TTSProviderConfig, logger *zap.Logger) *AzureTTSAdapter {
	return &AzureTTSAdapter{
		name:   cfg.Name,
		cfg:    cfg,
		logger: logger,
	}
}

func (p *AzureTTSAdapter) Name() string { return p.name }

func (p *AzureTTSAdapter) SyncVoices(ctx context.Context) ([]Voice, error) {
	p.logger.Info("Synchronizing Azure voice catalog")
	return []Voice{{ID: "en-US-JennyNeural", Name: "Jenny", Gender: "female", Language: "en-US", Provider: p.name}}, nil
}

func (p *AzureTTSAdapter) Synthesize(ctx context.Context, req TTSRequest) ([]byte, error) {
	return make([]byte, 16000), nil
}

type AzureStream struct {
	ch chan AudioChunk
}
func (s *AzureStream) AudioChunks() <-chan AudioChunk { return s.ch }
func (s *AzureStream) Close() error { close(s.ch); return nil }

func (p *AzureTTSAdapter) SynthesizeStream(ctx context.Context, req TTSRequest) (TTSStream, error) {
	stream := &AzureStream{ch: make(chan AudioChunk, 10)}
	go func() {
		defer stream.Close()
		start := time.Now()
		time.Sleep(250 * time.Millisecond) // Simulated latency
		metricTTSLatency.WithLabelValues(p.name).Observe(float64(time.Since(start).Milliseconds()))
		
		stream.ch <- AudioChunk{
			PCMData:    make([]byte, 320),
			SampleRate: 16000,
			Channels:   1,
			IsFinal:    true,
		}
	}()
	return stream, nil
}

// TTSRouter manages active TTS providers and fallback logic
type TTSRouter struct {
	config    *config.TTSConfig
	logger    *zap.Logger
	cache     *TTSCache
	providers []TTSProvider
	mu        sync.RWMutex
}

// NewTTSRouter creates a new TTSRouter
func NewTTSRouter(cfg *config.TTSConfig, logger *zap.Logger) *TTSRouter {
	router := &TTSRouter{
		config:    cfg,
		logger:    logger,
		cache:     NewTTSCache(logger),
		providers: make([]TTSProvider, 0),
	}

	for _, pCfg := range cfg.Providers {
		if pCfg.Enabled {
			var p TTSProvider
			switch pCfg.Type {
			case "azure":
				p = NewAzureTTSAdapter(pCfg, logger)
			default:
				p = NewElevenLabsProvider(pCfg, logger)
			}
			router.providers = append(router.providers, p)
		}
	}

	// Warm cache with common static prompts if configured
	if cfg.Cache.Enabled && cfg.Cache.Warming.Enabled {
		router.warmCache(cfg.Cache.Warming.Prompts)
	}

	return router
}

// SynthesizeWithFallback synthesizes audio using primary provider and falls back if latency > timeout
func (r *TTSRouter) SynthesizeWithFallback(ctx context.Context, req TTSRequest) (TTSStream, error) {
	// Check Cache
	if data, exists := r.cache.Get(req.Text, req.VoiceID); exists {
		metricTTSCacheHit.WithLabelValues("router").Inc()
		// Wrap cached data in a simple stream
		stream := &ElevenLabsStream{ch: make(chan AudioChunk, 1)}
		stream.ch <- AudioChunk{PCMData: data, SampleRate: 16000, Channels: 1, IsFinal: true}
		close(stream.ch)
		return stream, nil
	}
	metricTTSCacheMiss.WithLabelValues("router").Inc()

	r.mu.RLock()
	providers := make([]TTSProvider, len(r.providers))
	copy(providers, r.providers)
	r.mu.RUnlock()

	if len(providers) == 0 {
		return nil, errors.New("no active TTS provider configured")
	}

	// Try providers sequentially with timeout fallback
	for _, p := range providers {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
		
		stream, err := p.SynthesizeStream(ctxWithTimeout, req)
		if err == nil {
			// Found a working provider
			// We defer cancel here so it remains active for the stream lifetime, 
			// though in a real implementation we'd pass a detached context to the stream worker.
			go func() { <-ctx.Done(); cancel() }()
			return stream, nil
		}
		
		cancel()
		r.logger.Warn("TTS provider failed, falling back", zap.String("failed_provider", p.Name()), zap.Error(err))
	}

	return nil, errors.New("all TTS providers failed or timed out")
}

func (r *TTSRouter) warmCache(prompts []string) {
	for _, prompt := range prompts {
		r.cache.Set(prompt, "default", make([]byte, 640))
	}
	r.logger.Info("Pre-warmed TTS prompt cache", zap.Int("prompts_count", len(prompts)))
}

// Cache returns the internal cache engine
func (r *TTSRouter) Cache() *TTSCache {
	return r.cache
}
