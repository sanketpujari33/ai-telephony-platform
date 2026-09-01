package llm

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	metricLLMLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ai_telephony_llm_latency_ms",
		Help:    "LLM time-to-first-token in milliseconds",
		Buckets: []float64{100, 250, 500, 1000, 2000, 5000, 10000},
	}, []string{"provider", "model"})

	metricLLMCost = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ai_telephony_llm_cost_usd",
		Help: "Estimated LLM cost in USD",
	}, []string{"provider", "model"})
)

// MessageRole represents the sender role in conversation context
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// ToolSchema represents a function definition for tool calling
type ToolSchema struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

// CostInfo tracks the estimated cost of an LLM request
type CostInfo struct {
	PromptTokens int
	OutputTokens int
	EstimatedUSD float64
}

// ChatMessage represents a single message turn
type ChatMessage struct {
	Role      MessageRole
	Content   string
	Tokens    int
	Timestamp time.Time
}

// LLMResponse represents a generated response chunk or complete completion
type LLMResponse struct {
	Content      string
	IsFinal      bool
	PromptTokens int
	OutputTokens int
	Model        string
	Latency      time.Duration
	Cost         CostInfo
}

// LLMProvider interface for language models
type LLMProvider interface {
	Name() string
	GenerateStream(ctx context.Context, messages []ChatMessage) (<-chan LLMResponse, error)
	CountTokens(text string) int
}

// ContextWindow manages turn history, token budget, and background summarization
type ContextWindow struct {
	systemPrompt string
	turns        []ChatMessage
	maxTurns     int
	maxTokens    int
	tokenBudget  int
	totalTokens  int
	mu           sync.RWMutex
}

// NewContextWindow creates a new context window manager
func NewContextWindow(systemPrompt string, maxTurns, maxTokens, totalBudget int) *ContextWindow {
	cw := &ContextWindow{
		systemPrompt: systemPrompt,
		turns:        make([]ChatMessage, 0),
		maxTurns:     maxTurns,
		maxTokens:    maxTokens,
		tokenBudget:  totalBudget,
	}
	cw.recalculateTokens()
	return cw
}

// AddUserTurn adds a user message turn
func (cw *ContextWindow) AddUserTurn(text string, tokenEst int) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.turns = append(cw.turns, ChatMessage{
		Role:      RoleUser,
		Content:   text,
		Tokens:    tokenEst,
		Timestamp: time.Now(),
	})
	cw.pruneLocked()
}

// AddSystemPrompt updates the system prompt
func (cw *ContextWindow) AddSystemPrompt(prompt string) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.systemPrompt = prompt
}

// AddAssistantTurn adds an assistant response turn
func (cw *ContextWindow) AddAssistantTurn(text string, tokenEst int) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.turns = append(cw.turns, ChatMessage{
		Role:      RoleAssistant,
		Content:   text,
		Tokens:    tokenEst,
		Timestamp: time.Now(),
	})
	cw.pruneLocked()
}

// GetMessages returns formatted payload messages for LLM request
func (cw *ContextWindow) GetMessages() []ChatMessage {
	cw.mu.RLock()
	defer cw.mu.RUnlock()

	msgs := make([]ChatMessage, 0, len(cw.turns)+1)
	if cw.systemPrompt != "" {
		msgs = append(msgs, ChatMessage{
			Role:    RoleSystem,
			Content: cw.systemPrompt,
		})
	}
	msgs = append(msgs, cw.turns...)
	return msgs
}

// ShouldSummarize checks if token count exceeds threshold
func (cw *ContextWindow) ShouldSummarize(thresholdRatio float64) bool {
	cw.mu.RLock()
	defer cw.mu.RUnlock()

	if cw.maxTokens <= 0 {
		return false
	}
	return float64(cw.totalTokens)/float64(cw.maxTokens) >= thresholdRatio
}

// SummarizeBackground runs a background job to compress history when budget exceeds 70%
func (cw *ContextWindow) SummarizeBackground(ctx context.Context, provider LLMProvider) {
	if !cw.ShouldSummarize(0.70) {
		return
	}
	
	// Launch background job
	go func() {
		cw.mu.Lock()
		defer cw.mu.Unlock()
		
		// In a real implementation, this would:
		// 1. Extract the oldest N turns
		// 2. Format a prompt asking gpt-4o-mini to summarize them
		// 3. Call provider.GenerateStream(ctx, sumPrompt)
		// 4. Replace the old turns with a single RoleSystem summary turn
		
		// For now, we simulate a hard prune to keep budget
		if len(cw.turns) > 4 {
			cw.turns = cw.turns[2:] // Drop oldest
			cw.recalculateTokens()
		}
	}()
}

func (cw *ContextWindow) pruneLocked() {
	// Prune by max turns
	if cw.maxTurns > 0 && len(cw.turns) > cw.maxTurns {
		cw.turns = cw.turns[len(cw.turns)-cw.maxTurns:]
	}

	cw.recalculateTokens()

	// Prune by max tokens
	for cw.maxTokens > 0 && cw.totalTokens > cw.maxTokens && len(cw.turns) > 2 {
		cw.turns = cw.turns[2:] // Drop oldest user/assistant pair
		cw.recalculateTokens()
	}
}

func (cw *ContextWindow) recalculateTokens() {
	sum := 0
	for _, t := range cw.turns {
		sum += t.Tokens
	}
	cw.totalTokens = sum
}

// OpenAIProvider implements LLMProvider for OpenAI API
type OpenAIProvider struct {
	name   string
	cfg    config.LLMProviderConfig
	logger *zap.Logger
}

func NewOpenAIProvider(cfg config.LLMProviderConfig, logger *zap.Logger) *OpenAIProvider {
	return &OpenAIProvider{
		name:   cfg.Name,
		cfg:    cfg,
		logger: logger,
	}
}

func (p *OpenAIProvider) Name() string { return p.name }

func (p *OpenAIProvider) CountTokens(text string) int {
	return len(text) / 4 // Approximate 4 chars per token (tiktoken is heavier)
}

func (p *OpenAIProvider) GenerateStream(ctx context.Context, messages []ChatMessage) (<-chan LLMResponse, error) {
	ch := make(chan LLMResponse, 100)
	
	// Real implementation initiates a POST request to https://api.openai.com/v1/chat/completions
	// with "stream": true and parses Server-Sent Events (SSE)
	p.logger.Info("Initiating OpenAI completion stream", zap.Int("messages", len(messages)))
	
	go func() {
		defer close(ch)
		// We mock the SSE streaming parser for compilation without failing without keys
		ch <- LLMResponse{
			Content:      "This is a connected OpenAI response stream.",
			IsFinal:      true,
			PromptTokens: 50,
			OutputTokens: 10,
			Model:        p.cfg.Config["model"].(string),
			Latency:      150 * time.Millisecond,
		}
	}()
	
	return ch, nil
}

// AnthropicProvider implements LLMProvider for Anthropic API
type AnthropicProvider struct {
	name   string
	cfg    config.LLMProviderConfig
	logger *zap.Logger
}

func NewAnthropicProvider(cfg config.LLMProviderConfig, logger *zap.Logger) *AnthropicProvider {
	return &AnthropicProvider{
		name:   cfg.Name,
		cfg:    cfg,
		logger: logger,
	}
}

func (p *AnthropicProvider) Name() string { return p.name }

func (p *AnthropicProvider) CountTokens(text string) int {
	return len(text) / 4 // Claude has its own tokenizer
}

func (p *AnthropicProvider) GenerateStream(ctx context.Context, messages []ChatMessage) (<-chan LLMResponse, error) {
	ch := make(chan LLMResponse, 100)
	
	// Real implementation initiates a POST request to https://api.anthropic.com/v1/messages
	// with "stream": true and parses Anthropic SSE format
	p.logger.Info("Initiating Anthropic completion stream", zap.Int("messages", len(messages)))
	
	go func() {
		defer close(ch)
		ch <- LLMResponse{
			Content:      "This is a connected Anthropic Claude response stream.",
			IsFinal:      true,
			PromptTokens: 50,
			OutputTokens: 11,
			Model:        p.cfg.Config["model"].(string),
			Latency:      120 * time.Millisecond,
		}
	}()
	
	return ch, nil
}

// LLMCircuitBreaker manages LLM provider health state
type LLMCircuitBreaker struct {
	failures int
	isOpen   bool
	mu       sync.RWMutex
}

func (cb *LLMCircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= 3 {
		cb.isOpen = true
	}
}

func (cb *LLMCircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return !cb.isOpen
}

// LLMRouter handles multi-provider routing and fallback
type LLMRouter struct {
	config    *config.LLMConfig
	logger    *zap.Logger
	providers []LLMProvider
	breakers  map[string]*LLMCircuitBreaker
	cache     map[string]string // Simple response cache for deterministic prompts
	mu        sync.RWMutex
}

// NewLLMRouter creates a new LLMRouter instance
func NewLLMRouter(cfg *config.LLMConfig, logger *zap.Logger) *LLMRouter {
	router := &LLMRouter{
		config:    cfg,
		logger:    logger,
		providers: make([]LLMProvider, 0),
		breakers:  make(map[string]*LLMCircuitBreaker),
		cache:     make(map[string]string),
	}

	for _, pCfg := range cfg.Providers {
		if pCfg.Enabled {
			var p LLMProvider
			switch pCfg.Type {
			case "anthropic":
				p = NewAnthropicProvider(pCfg, logger)
			default:
				p = NewOpenAIProvider(pCfg, logger)
			}
			router.providers = append(router.providers, p)
			router.breakers[p.Name()] = &LLMCircuitBreaker{}
		}
	}

	return router
}

// ReportFailure records a failure for a named provider to trigger circuit breaker
func (r *LLMRouter) ReportFailure(providerName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if cb, exists := r.breakers[providerName]; exists {
		cb.RecordFailure()
	}
}

// SelectProvider returns the active primary LLM provider respecting priority and circuit breakers
func (r *LLMRouter) SelectProvider() (LLMProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.providers) == 0 {
		return nil, errors.New("no active LLM provider configured")
	}

	for _, p := range r.providers {
		if r.breakers[p.Name()].Allow() {
			return p, nil
		}
	}

	return nil, errors.New("all LLM providers are currently unavailable (circuit breakers open)")
}

// CheckCache returns a cached response if available for deterministic requests
func (r *LLMRouter) CheckCache(promptHash string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res, ok := r.cache[promptHash]
	return res, ok
}

// RecordCost logs cost info to Prometheus
func (r *LLMRouter) RecordCost(provider, model string, cost float64) {
	metricLLMCost.WithLabelValues(provider, model).Add(cost)
}
