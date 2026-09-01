package media

import (
	"context"

	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer implements the MediaProxy gRPC service
type GRPCServer struct {
	pb.UnimplementedMediaProxyServer
	engine *MediaProxyServer
	logger *zap.Logger
}

// NewGRPCServer creates a new gRPC handler for MediaProxy
func NewGRPCServer(engine *MediaProxyServer, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		engine: engine,
		logger: logger,
	}
}

func (s *GRPCServer) StartMediaSession(ctx context.Context, req *pb.StartMediaRequest) (*pb.StartMediaResponse, error) {
	s.logger.Info("Received StartMediaSession request", zap.String("call_id", req.CallId))
	
	codecSpec := CodecSpec{
		Name:      req.Codec,
		ClockRate: uint32(req.SampleRate),
	}
	
	// Create a unique session ID if not provided
	sessionID := "sess_" + req.CallId
	
	session, err := s.engine.CreateSession(sessionID, req.CallId, codecSpec)
	if err != nil {
		s.logger.Error("Failed to create media session", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	return &pb.StartMediaResponse{
		SessionId:    session.ID,
		LocalRtpPort: int32(session.LocalRTPAddr.Port),
	}, nil
}

func (s *GRPCServer) StopMediaSession(ctx context.Context, req *pb.StopMediaRequest) (*pb.StopMediaResponse, error) {
	s.logger.Info("Received StopMediaSession request", zap.String("session_id", req.SessionId))
	
	err := s.engine.CloseSession(req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "session not found")
	}
	
	return &pb.StopMediaResponse{
		Success: true,
	}, nil
}

func (s *GRPCServer) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	s.engine.mu.RLock()
	session, exists := s.engine.sessions[req.SessionId]
	s.engine.mu.RUnlock()
	
	if !exists {
		return nil, status.Errorf(codes.NotFound, "session not found")
	}
	
	session.Metrics.mu.Lock()
	defer session.Metrics.mu.Unlock()
	
	return &pb.GetMetricsResponse{
		Mos:               float32(session.Metrics.MOS),
		JitterMs:          int32(session.Metrics.JitterMs),
		PacketLossPercent: float32(float64(session.Metrics.PacketsLost) / float64(session.Metrics.PacketsReceived+session.Metrics.PacketsLost) * 100.0),
	}, nil
}
