package stt

import (
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCircuitBreakerStateTransitions(t *testing.T) {
	cb := NewProviderCircuitBreaker(2, 1)

	assert.Equal(t, StateClosed, cb.GetState())
	assert.True(t, cb.Allow())

	cb.RecordFailure()
	assert.Equal(t, StateClosed, cb.GetState())

	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())
	assert.False(t, cb.Allow())

	time.Sleep(1100 * time.Millisecond)

	assert.True(t, cb.Allow()) // Transitions to HalfOpen
	cb.RecordSuccess()
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestSTTRouterGetPrimaryProvider(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.STTConfig{
		Providers: []config.STTProviderConfig{
			{Name: "deepgram", Type: "deepgram", Priority: 1, Enabled: true},
			{Name: "whisper", Type: "whisper", Priority: 2, Enabled: true},
		},
		CircuitBreaker: config.STTCircuitBreakerConfig{
			Enabled:          true,
			FailureThreshold: 2,
			Timeout:          30,
		},
	}

	router := NewSTTRouter(cfg, logger)
	provider, err := router.GetPrimaryProvider()

	assert.NoError(t, err)
	assert.Equal(t, "deepgram", provider.Name())

	// Trip circuit breaker for deepgram
	router.ReportFailure("deepgram")
	router.ReportFailure("deepgram")

	nextProvider, err := router.GetPrimaryProvider()
	assert.NoError(t, err)
	assert.Equal(t, "whisper", nextProvider.Name())
}
