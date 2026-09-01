package transfer

import (
	"context"
	"testing"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSIPReferBuilderContextHeaders(t *testing.T) {
	cfg := &config.TransferManagerContextHeadersConfig{
		Enabled:              true,
		HeaderPrefix:         "X-Conversation-",
		IncludeTranscript:    true,
		IncludeExtractedData: true,
		IncludeSentiment:     true,
	}

	builder := NewSIPReferBuilder(cfg)
	req := &TransferRequest{
		ID:             "t1",
		CallID:         "call-123",
		TargetAgentURI: "sip:agent01@domain.com",
		Transcript:     "Customer wants refund.",
		Sentiment:      "neutral",
		ExtractedData:  map[string]string{"Intent": "Billing"},
	}

	msg := builder.BuildReferMessage(req)
	assert.Equal(t, "call-123", msg.CallID)
	assert.Contains(t, msg.ReferTo, "sip:agent01@domain.com")
	assert.Contains(t, msg.ReferTo, "Replaces=call-123")
	assert.Equal(t, "Customer wants refund.", msg.Headers["X-Conversation-Transcript"])
	assert.Equal(t, "neutral", msg.Headers["X-Conversation-Sentiment"])
	assert.Equal(t, "Billing", msg.Headers["X-Conversation-Data-Intent"])
}

func TestTransferEngineLifecycle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.TransferConfig{
		Manager: config.TransferManagerConfig{
			ContextHeaders: config.TransferManagerContextHeadersConfig{
				HeaderPrefix: "X-Conversation-",
			},
		},
	}

	engine := NewTransferEngine(cfg, logger)

	req := &TransferRequest{
		ID:             "tr-1",
		CallID:         "call-99",
		TargetAgentURI: "sip:agent02@domain.com",
		Priority:       1,
	}

	referMsg, err := engine.InitiateTransfer(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, referMsg)

	stored, exists := engine.GetTransfer("tr-1")
	assert.True(t, exists)
	assert.Equal(t, TransferStateQueued, stored.State)

	err = engine.UpdateState("tr-1", TransferStateConnected)
	assert.NoError(t, err)

	updated, _ := engine.GetTransfer("tr-1")
	assert.Equal(t, TransferStateConnected, updated.State)
}
