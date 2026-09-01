package media

import (
	"context"
	"testing"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"go.uber.org/zap"
)

func TestGRPCServer_StartStopMediaSession(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.MediaConfig{
		Proxy: config.MediaProxyConfig{
			RTPPortRange: config.MediaProxyRTPPortRangeConfig{
				Start: 20000,
				End:   20010,
			},
		},
	}
	
	engine, err := NewMediaProxyServer(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create media proxy server: %v", err)
	}
	
	server := NewGRPCServer(engine, logger)

	req := &pb.StartMediaRequest{
		CallId: "call-123",
		Codec: "PCMU",
		SampleRate: 8000,
	}

	resp, err := server.StartMediaSession(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error starting session: %v", err)
	}

	if resp.SessionId == "" {
		t.Errorf("expected session ID to be generated")
	}

	if resp.LocalRtpPort < 20000 || resp.LocalRtpPort > 20010 {
		t.Errorf("expected RTP port between 20000 and 20010, got %d", resp.LocalRtpPort)
	}

	// Verify metrics
	metricReq := &pb.GetMetricsRequest{
		SessionId: resp.SessionId,
	}
	metricResp, err := server.GetMetrics(context.Background(), metricReq)
	if err != nil {
		t.Fatalf("unexpected error getting metrics: %v", err)
	}
	
	// Initial metrics should be 4.5/default
	if metricResp.Mos != 4.5 {
		t.Errorf("expected initial MOS 4.5, got %f", metricResp.Mos)
	}

	// Stop session
	stopReq := &pb.StopMediaRequest{
		SessionId: resp.SessionId,
	}
	
	stopResp, err := server.StopMediaSession(context.Background(), stopReq)
	if err != nil {
		t.Fatalf("unexpected error stopping session: %v", err)
	}
	
	if !stopResp.Success {
		t.Errorf("expected stop success=true")
	}

	// Metrics should fail now
	_, err = server.GetMetrics(context.Background(), metricReq)
	if err == nil {
		t.Errorf("expected error getting metrics for stopped session")
	}
}
