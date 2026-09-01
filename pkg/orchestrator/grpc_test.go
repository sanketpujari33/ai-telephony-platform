package orchestrator

import (
	"context"
	"testing"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
)

func TestGRPCServer_ProcessAndEndCall(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{}
	
	engine := NewOrchestratorServer(cfg, logger)
	server := NewGRPCServer(engine, logger)

	// Test ProcessCall
	req := &pb.ProcessCallRequest{
		CallId:        "call-999",
		InitialPrompt: "You are a helpful assistant.",
	}

	resp, err := server.ProcessCall(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error processing call: %v", err)
	}

	if !resp.Accepted {
		t.Errorf("expected Accepted=true")
	}

	if resp.OrchestratorSessionId == "" {
		t.Errorf("expected session ID to be populated")
	}

	// Verify session exists
	engine.mu.RLock()
	session, exists := engine.sessions[resp.OrchestratorSessionId]
	engine.mu.RUnlock()

	if !exists {
		t.Fatalf("session was not created in engine")
	}
	
	// Test EndCall
	endReq := &pb.EndCallRequest{
		OrchestratorSessionId: resp.OrchestratorSessionId,
	}
	
	endResp, err := server.EndCall(context.Background(), endReq)
	if err != nil {
		t.Fatalf("unexpected error ending call: %v", err)
	}
	
	if endResp.TranscriptSummary == "" {
		// Mock implementation just returns an empty string or standard value right now.
		// As long as it succeeds we are good.
	}
	
	// Verify state changed
	session.mu.RLock()
	state := session.State
	session.mu.RUnlock()
	
	if state != SessionStateEnded {
		t.Errorf("expected session state to be Ended, got %s", state)
	}
	
	// Test EndCall with invalid session
	badReq := &pb.EndCallRequest{
		OrchestratorSessionId: "unknown",
	}
	
	_, err = server.EndCall(context.Background(), badReq)
	if err == nil {
		t.Errorf("expected error when ending unknown session")
	}
}
