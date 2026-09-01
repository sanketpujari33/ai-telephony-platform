package orchestrator

import (
	"context"
	"fmt"

	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer implements the AIOrchestrator gRPC service
type GRPCServer struct {
	pb.UnimplementedAIOrchestratorServer
	engine *OrchestratorServer
	logger *zap.Logger
}

// NewGRPCServer creates a new gRPC handler for AIOrchestrator
func NewGRPCServer(engine *OrchestratorServer, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		engine: engine,
		logger: logger,
	}
}

func (s *GRPCServer) ProcessCall(ctx context.Context, req *pb.ProcessCallRequest) (*pb.ProcessCallResponse, error) {
	s.logger.Info("Received ProcessCall request", zap.String("call_id", req.CallId))
	
	sessionID := "orch_" + req.CallId
	tenantID := "default_tenant" // In a real app this would come from the context/auth
	
	session, err := s.engine.CreateCallSession(sessionID, req.CallId, tenantID)
	if err != nil {
		s.logger.Error("Failed to create call session", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	if req.InitialPrompt != "" {
		session.ContextWindow.AddSystemPrompt(req.InitialPrompt)
	}

	if err := session.Start(); err != nil {
		s.logger.Error("Failed to start session", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to start session: %v", err)
	}

	return &pb.ProcessCallResponse{
		OrchestratorSessionId: session.ID,
		Accepted:              true,
	}, nil
}

func (s *GRPCServer) EndCall(ctx context.Context, req *pb.EndCallRequest) (*pb.EndCallResponse, error) {
	s.logger.Info("Received EndCall request", zap.String("session_id", req.OrchestratorSessionId))
	
	s.engine.mu.RLock()
	session, exists := s.engine.sessions[req.OrchestratorSessionId]
	s.engine.mu.RUnlock()
	
	if !exists {
		return nil, status.Errorf(codes.NotFound, "session not found")
	}
	
	if err := session.Stop(); err != nil {
		s.logger.Error("Failed to stop session", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to stop session: %v", err)
	}
	
	// Create a simple transcript summary
	summary := fmt.Sprintf("Conversation ended with %d turns.", len(session.ContextWindow.GetMessages()))
	
	return &pb.EndCallResponse{
		TranscriptSummary: summary,
	}, nil
}
