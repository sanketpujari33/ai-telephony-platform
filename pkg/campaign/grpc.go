package campaign

import (
	"context"

	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer implements the CampaignManager gRPC service
type GRPCServer struct {
	pb.UnimplementedCampaignManagerServer
	engine *CampaignEngine
	logger *zap.Logger
}

// NewGRPCServer creates a new gRPC handler for CampaignManager
func NewGRPCServer(engine *CampaignEngine, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		engine: engine,
		logger: logger,
	}
}

func (s *GRPCServer) StartCampaign(ctx context.Context, req *pb.StartCampaignRequest) (*pb.StartCampaignResponse, error) {
	s.logger.Info("Received StartCampaign request", zap.String("campaign_id", req.CampaignId))
	
	contacts := make([]Contact, 0, len(req.Contacts))
	for _, c := range req.Contacts {
		contacts = append(contacts, Contact{
			ID:              c.Id,
			PhoneNumber:     c.PhoneNumber,
			Timezone:        c.Timezone,
			ConsentVerified: true, // Assuming true for this demo based on the protobuf definition
		})
	}
	
	campaign := &Campaign{
		ID:           req.CampaignId,
		TenantID:     "default_tenant",
		Name:         "Campaign " + req.CampaignId,
		Status:       CampaignStatusActive,
		Contacts:     contacts,
		TargetAgents: int(req.TargetAbandonmentRate),
	}
	
	s.engine.RegisterCampaign(campaign)
	
	return &pb.StartCampaignResponse{
		Success: true,
	}, nil
}

func (s *GRPCServer) StopCampaign(ctx context.Context, req *pb.StopCampaignRequest) (*pb.StopCampaignResponse, error) {
	s.logger.Info("Received StopCampaign request", zap.String("campaign_id", req.CampaignId))
	
	s.engine.mu.Lock()
	defer s.engine.mu.Unlock()
	
	campaign, exists := s.engine.campaigns[req.CampaignId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "campaign not found")
	}
	
	campaign.Status = CampaignStatusPaused
	
	return &pb.StopCampaignResponse{
		Success: true,
	}, nil
}
