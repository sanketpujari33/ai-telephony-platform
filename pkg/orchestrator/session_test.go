package orchestrator

import (
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/llm"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/stt"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/tts"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCallSessionLifecycle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Providers: []config.LLMProviderConfig{
				{Name: "openai", Type: "openai", Enabled: true},
			},
		},
	}

	sttRouter := stt.NewSTTRouter(&cfg.STT, logger)
	llmRouter := llm.NewLLMRouter(&cfg.LLM, logger)
	ttsRouter := tts.NewTTSRouter(&cfg.TTS, logger)

	session := NewCallSession("sess-1", "call-1", "tenant-1", cfg, sttRouter, llmRouter, ttsRouter, logger)
	assert.Equal(t, SessionStateConnecting, session.State)

	err := session.Start()
	assert.NoError(t, err)
	assert.Equal(t, SessionStateConversing, session.State)

	// Process user transcript
	err = session.ProcessTranscript(stt.TranscriptResult{
		Text:      "Hello, I need help with my account.",
		IsFinal:   true,
		Timestamp: time.Now(),
	})
	assert.NoError(t, err)

	err = session.Stop()
	assert.NoError(t, err)
	assert.Equal(t, SessionStateEnded, session.State)
}
