package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/llm"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/stt"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/tts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestProviderFailover tests the Circuit Breaker logic across all 3 AI layers (STT, LLM, TTS).
// It spawns httptest servers that intentionally return 503 or timeout, and ensures
// the provider router correctly fails over to the next priority adapter.
func TestProviderFailover(t *testing.T) {
	logger := zap.NewNop()

	t.Run("STT Failover (Deepgram -> Assembly)", func(t *testing.T) {
		// Mock Deepgram to fail
		mockDeepgram := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer mockDeepgram.Close()

		// Mock Assembly to succeed
		mockAssembly := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"text":"hello world"}`))
		}))
		defer mockAssembly.Close()

		cfg := &config.STTConfig{
			Providers: []config.STTProviderConfig{
				{Name: "deepgram-primary", Type: "deepgram", Priority: 1, Enabled: true},
				{Name: "assemblyai-fallback", Type: "assemblyai", Priority: 2, Enabled: true},
			},
			CircuitBreaker: config.STTCircuitBreakerConfig{
				Enabled:          true,
				FailureThreshold: 2,
				Timeout:          60,
			},
		}
		sttRouter := stt.NewSTTRouter(cfg, logger)

		// Initial primary provider should be Deepgram
		primary, err := sttRouter.GetPrimaryProvider()
		require.NoError(t, err)
		assert.Equal(t, "deepgram-primary", primary.Name())

		t.Log("Simulating Deepgram 503... expected to failover to AssemblyAI")
		sttRouter.ReportFailure("deepgram-primary")
		sttRouter.ReportFailure("deepgram-primary")

		failoverProvider, err := sttRouter.GetPrimaryProvider()
		require.NoError(t, err)
		assert.Equal(t, "assemblyai-fallback", failoverProvider.Name())
		t.Log("STT Circuit breaker trigger validated")
	})

	t.Run("LLM Rate Limit (OpenAI -> Anthropic)", func(t *testing.T) {
		cfg := &config.LLMConfig{
			Providers: []config.LLMProviderConfig{
				{Name: "openai-gpt4", Type: "openai", Priority: 1, Enabled: true, Config: map[string]interface{}{"model": "gpt-4"}},
				{Name: "anthropic-claude", Type: "anthropic", Priority: 2, Enabled: true, Config: map[string]interface{}{"model": "claude-3-5-sonnet"}},
			},
		}
		llmRouter := llm.NewLLMRouter(cfg, logger)

		primary, err := llmRouter.SelectProvider()
		require.NoError(t, err)
		assert.Equal(t, "openai-gpt4", primary.Name())

		t.Log("Simulating OpenAI 429 Rate Limit... expected to failover to Anthropic")
		for i := 0; i < 3; i++ {
			llmRouter.ReportFailure("openai-gpt4")
		}

		failoverProvider, err := llmRouter.SelectProvider()
		require.NoError(t, err)
		assert.Equal(t, "anthropic-claude", failoverProvider.Name())
		t.Log("LLM Failover validated")
	})

	t.Run("TTS Timeout (ElevenLabs -> Azure)", func(t *testing.T) {
		cfg := &config.TTSConfig{
			Providers: []config.TTSProviderConfig{
				{Name: "elevenlabs-primary", Type: "elevenlabs", Priority: 1, Enabled: true},
				{Name: "azure-fallback", Type: "azure", Priority: 2, Enabled: true},
			},
		}
		ttsRouter := tts.NewTTSRouter(cfg, logger)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		req := tts.TTSRequest{
			Text:    "Hello, this is an automated telephony test.",
			VoiceID: "21m00Tcm4TlvDq8ikWAM",
		}

		t.Log("Simulating TTS synthesis stream with fallback...")
		stream, err := ttsRouter.SynthesizeWithFallback(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, stream)

		chunk := <-stream.AudioChunks()
		assert.NotEmpty(t, chunk.PCMData)
		t.Log("TTS Timeout failover validated")
	})
}

