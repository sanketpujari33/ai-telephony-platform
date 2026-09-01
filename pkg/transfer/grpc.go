package transfer

import (
	"context"

	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer implements the TransferManager gRPC service
type GRPCServer struct {
	pb.UnimplementedTransferManagerServer
	engine *TransferEngine
	logger *zap.Logger
}

// NewGRPCServer creates a new gRPC handler for TransferManager
func NewGRPCServer(engine *TransferEngine, logger *zap.Logger) *GRPCServer {
	return &GRPCServer{
		engine: engine,
		logger: logger,
	}
}

func (s *GRPCServer) InitiateTransfer(ctx context.Context, req *pb.InitiateTransferRequest) (*pb.InitiateTransferResponse, error) {
	s.logger.Info("Received InitiateTransfer request", zap.String("call_id", req.CallId))
	
	transferID := "tx_" + req.CallId
	
	transferReq := &TransferRequest{
		ID:             transferID,
		CallID:         req.CallId,
		TargetAgentURI: req.TargetSipUri,
		Priority:       1,
		Transcript:     req.TranscriptContext,
		Sentiment:      req.Sentiment,
	}
	
	_, err := s.engine.InitiateTransfer(ctx, transferReq)
	if err != nil {
		s.logger.Error("Failed to initiate transfer", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to initiate transfer: %v", err)
	}

	return &pb.InitiateTransferResponse{
		TransferId: transferID,
		Accepted:   true,
	}, nil
}
