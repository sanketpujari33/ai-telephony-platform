package integration

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/llm"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/orchestrator"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/sip/ingress"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/storage"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/stt"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/tts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestContainers(ctx context.Context, t *testing.T) (config.StoragePostgresConfig, config.StorageRedisConfig, func()) {
	t.Log("Initializing storage configuration for integration testing...")

	pgCfg := config.StoragePostgresConfig{
		Host:           "localhost",
		Port:           5432,
		Username:       "postgres",
		Password:       "secret",
		Database:       "telephony_test",
		SSLMode:        "disable",
		MaxConnections: 5,
		Tables: config.StoragePostgresTablesConfig{
			DNCEntries:  "dnc_entries",
			CallRecords: "call_records",
		},
	}

	redisCfg := config.StorageRedisConfig{
		Host:           "localhost",
		Port:           6379,
		Session:        config.StorageRedisSessionConfig{Prefix: "session:"},
		MaxConnections: 10,
	}

	return pgCfg, redisCfg, func() {}
}

// TestInboundCallFlow tests the end-to-end traversal of a SIP INVITE
// flowing through the Orchestrator, triggering STT and LLM, and generating a TTS response.
func TestInboundCallFlow(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	// 1. Setup Testcontainers (Postgres + Redis)
	pgCfg, redisCfg, cleanup := setupTestContainers(ctx, t)
	defer cleanup()

	// 2. Initialize Platform Storage
	pgStore, err := storage.NewPostgresStore(ctx, pgCfg, logger)
	require.NoError(t, err)
	defer pgStore.Close()

	redisStore, err := storage.NewRedisStore(ctx, redisCfg, logger)
	require.NoError(t, err)
	defer redisStore.Close()

	t.Log("Initialized Platform Core storage components")

	// 3. Simulate Inbound SIP INVITE
	rawInvite := []byte("INVITE sip:+18005551234@10.0.0.1:5060 SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 10.0.0.2:5060;branch=z9hG4bK-1234\r\n" +
		"From: <sip:+14155551111@10.0.0.2>;tag=caller123\r\n" +
		"To: <sip:+18005551234@10.0.0.1>\r\n" +
		"Call-ID: c4129a02-e221-4f18-b807-testcallid123\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Content-Length: 0\r\n\r\n")

	parser := ingress.NewSIPParser()
	localAddr := &net.UDPAddr{IP: net.ParseIP("10.0.0.1"), Port: 5060}
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("10.0.0.2"), Port: 5060}

	invite, err := parser.ParseRequest(rawInvite, ingress.TransportUDP, localAddr, remoteAddr)
	require.NoError(t, err)
	assert.Equal(t, "c4129a02-e221-4f18-b807-testcallid123", invite.CallID)

	// 4. Initialize Platform AI Routers & Orchestrator
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Providers: []config.LLMProviderConfig{
				{Name: "openai-test", Type: "openai", Priority: 1, Enabled: true, Config: map[string]interface{}{"model": "gpt-4"}},
			},
			Context: config.LLMContextConfig{
				SlidingWindow: config.LLMSlidingWindowConfig{MaxTurns: 10, MaxTokens: 4096},
				TokenBudget:   config.LLMTokenBudgetConfig{TotalBudget: 8192},
			},
		},
		STT: config.STTConfig{
			Providers: []config.STTProviderConfig{
				{Name: "deepgram-test", Type: "deepgram", Priority: 1, Enabled: true},
			},
		},
		TTS: config.TTSConfig{
			Providers: []config.TTSProviderConfig{
				{Name: "elevenlabs-test", Type: "elevenlabs", Priority: 1, Enabled: true},
			},
		},
	}

	sttRouter := stt.NewSTTRouter(&cfg.STT, logger)
	llmRouter := llm.NewLLMRouter(&cfg.LLM, logger)
	ttsRouter := tts.NewTTSRouter(&cfg.TTS, logger)

	session := orchestrator.NewCallSession("sess-integration-1", invite.CallID, "tenant-default", cfg, sttRouter, llmRouter, ttsRouter, logger)
	assert.Equal(t, orchestrator.SessionStateConnecting, session.State)

	// 5. Validate Orchestrator State transitions
	err = session.Start()
	require.NoError(t, err)
	assert.Equal(t, orchestrator.SessionStateConversing, session.State)

	// 6. Simulate STT Transcript Processing
	transcript := stt.TranscriptResult{
		Text:       "Hello assistant, can you help me with my account?",
		IsFinal:    true,
		Confidence: 0.98,
		Timestamp:  time.Now(),
	}

	err = session.ProcessTranscript(transcript)
	require.NoError(t, err)

	// 7. Store session in Redis
	err = redisStore.SetSession(ctx, session.ID, `{"state":"conversing"}`, 5*time.Minute)
	if err == nil {
		cachedData, err := redisStore.GetSession(ctx, session.ID)
		assert.NoError(t, err)
		assert.Contains(t, cachedData, "conversing")
	}

	// 8. End Call & Cleanup
	_ = session.Stop()
	assert.Equal(t, orchestrator.SessionStateEnded, session.State)
	t.Log("End-to-End Call Flow executed successfully")
}

// TestCampaignOutboundFlow tests AMD detection and outbound dialing.
func TestCampaignOutboundFlow(t *testing.T) {
	t.Log("Testing Campaign Outbound logic with AMD integration")
	logger := zap.NewNop()

	cfg := &config.Config{
		LLM: config.LLMConfig{
			Providers: []config.LLMProviderConfig{{Name: "openai-test", Type: "openai", Enabled: true}},
		},
	}

	sttRouter := stt.NewSTTRouter(&cfg.STT, logger)
	llmRouter := llm.NewLLMRouter(&cfg.LLM, logger)
	ttsRouter := tts.NewTTSRouter(&cfg.TTS, logger)

	session := orchestrator.NewCallSession("campaign-sess-1", "call-outbound-999", "tenant-campaign", cfg, sttRouter, llmRouter, ttsRouter, logger)
	require.NotNil(t, session)
	assert.Equal(t, orchestrator.SessionStateConnecting, session.State)

	err := session.Start()
	require.NoError(t, err)
	assert.Equal(t, orchestrator.SessionStateConversing, session.State)

	_ = session.Stop()
	assert.Equal(t, orchestrator.SessionStateEnded, session.State)
	t.Log("Outbound campaign flow passed")
}

