package campaign

import (
	"context"
	"testing"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
)

func TestGRPCServer_StartCampaign(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.CampaignConfig{}
	
	// nil checker for this simple test since StartCampaign doesn't immediately use it
	engine := NewCampaignEngine(cfg, nil, logger)
	server := NewGRPCServer(engine, logger)

	req := &pb.StartCampaignRequest{
		CampaignId: "camp-123",
		TargetAbandonmentRate: 2.0,
		Contacts: []*pb.ContactMessage{
			{
				Id: "cont-1",
				PhoneNumber: "+1234567890",
				Timezone: "UTC",
			},
		},
	}

	resp, err := server.StartCampaign(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true")
	}

	// Verify engine registered it
	engine.mu.RLock()
	campaign, exists := engine.campaigns["camp-123"]
	engine.mu.RUnlock()

	if !exists {
		t.Fatalf("campaign was not registered in engine")
	}

	if len(campaign.Contacts) != 1 {
		t.Errorf("expected 1 contact, got %d", len(campaign.Contacts))
	}
	
	if campaign.TargetAgents != 2 {
		t.Errorf("expected target agents (mapped from abandonment rate logic) to be 2, got %d", campaign.TargetAgents)
	}
}

func TestGRPCServer_StopCampaign(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.CampaignConfig{}
	
	engine := NewCampaignEngine(cfg, nil, logger)
	server := NewGRPCServer(engine, logger)

	// Register a campaign manually
	engine.RegisterCampaign(&Campaign{
		ID: "camp-456",
		Status: CampaignStatusActive,
	})

	req := &pb.StopCampaignRequest{
		CampaignId: "camp-456",
	}

	resp, err := server.StopCampaign(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success=true")
	}

	// Verify status changed
	engine.mu.RLock()
	campaign := engine.campaigns["camp-456"]
	engine.mu.RUnlock()

	if campaign.Status != CampaignStatusPaused {
		t.Errorf("expected status paused, got %s", campaign.Status)
	}
	
	// Test stopping unknown campaign
	badReq := &pb.StopCampaignRequest{
		CampaignId: "camp-unknown",
	}
	
	_, err = server.StopCampaign(context.Background(), badReq)
	if err == nil {
		t.Errorf("expected error for unknown campaign")
	}
}
