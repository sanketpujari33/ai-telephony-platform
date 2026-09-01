package transfer

import (
	"context"
	"testing"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
)

func TestGRPCServer_InitiateTransfer(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.TransferConfig{}
	
	engine := NewTransferEngine(cfg, logger)
	server := NewGRPCServer(engine, logger)

	req := &pb.InitiateTransferRequest{
		CallId:            "call-555",
		TargetSipUri:      "sip:agent@callcenter.com",
		TranscriptContext: "User needs help with billing.",
		Sentiment:         "frustrated",
	}

	resp, err := server.InitiateTransfer(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Accepted {
		t.Errorf("expected Accepted=true")
	}

	if resp.TransferId == "" {
		t.Errorf("expected TransferId to be set")
	}

	// Verify engine recorded it
	engine.mu.RLock()
	reqObj, exists := engine.requests[resp.TransferId]
	engine.mu.RUnlock()

	if !exists {
		t.Fatalf("transfer request not found in engine")
	}

	if reqObj.CallID != "call-555" {
		t.Errorf("expected CallID call-555, got %s", reqObj.CallID)
	}
	
	if reqObj.State != TransferStateQueued {
		t.Errorf("expected initial state Queued, got %s", reqObj.State)
	}
}

func TestGRPCServer_InitiateTransfer_MissingCallID(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.TransferConfig{}
	
	engine := NewTransferEngine(cfg, logger)
	server := NewGRPCServer(engine, logger)

	req := &pb.InitiateTransferRequest{
		// Missing CallId
		TargetSipUri: "sip:agent@callcenter.com",
	}

	_, err := server.InitiateTransfer(context.Background(), req)
	if err == nil {
		t.Errorf("expected error for missing CallId")
	}
}
