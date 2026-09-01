package v1

import (
	"context"
	"google.golang.org/grpc"
)

// This is a stubbed package representing the generated Protobuf code 
// for the Telephony Platform. In a real environment, this would be 
// auto-generated using `protoc` from the .proto definitions.

type CampaignManagerServer interface {
	StartCampaign(ctx context.Context, in *StartCampaignRequest) (*StartCampaignResponse, error)
	StopCampaign(ctx context.Context, in *StopCampaignRequest) (*StopCampaignResponse, error)
}
func RegisterCampaignManagerServer(s *grpc.Server, srv CampaignManagerServer) {}
type UnimplementedCampaignManagerServer struct{}

type TransferManagerServer interface {
	InitiateTransfer(ctx context.Context, in *InitiateTransferRequest) (*InitiateTransferResponse, error)
}
func RegisterTransferManagerServer(s *grpc.Server, srv TransferManagerServer) {}
type UnimplementedTransferManagerServer struct{}

type CallRouterServer interface{}
func RegisterCallRouterServer(s *grpc.Server, srv CallRouterServer) {}
type UnimplementedCallRouterServer struct{}

type MediaProxyServer interface{}
func RegisterMediaProxyServer(s *grpc.Server, srv MediaProxyServer) {}
type UnimplementedMediaProxyServer struct{}

type AIOrchestratorServer interface {
	ProcessCall(ctx context.Context, in *ProcessCallRequest) (*ProcessCallResponse, error)
	EndCall(ctx context.Context, in *EndCallRequest) (*EndCallResponse, error)
}
func RegisterAIOrchestratorServer(s *grpc.Server, srv AIOrchestratorServer) {}
type UnimplementedAIOrchestratorServer struct{}

type APIRequest struct{}
type APIResponse struct{}

// Client stubs
type MediaProxyClient interface {
	StartMediaSession(ctx context.Context, in *StartMediaRequest, opts ...grpc.CallOption) (*StartMediaResponse, error)
	StopMediaSession(ctx context.Context, in *StopMediaRequest, opts ...grpc.CallOption) (*StopMediaResponse, error)
	GetMetrics(ctx context.Context, in *GetMetricsRequest, opts ...grpc.CallOption) (*GetMetricsResponse, error)
}
func NewMediaProxyClient(cc grpc.ClientConnInterface) MediaProxyClient { return nil }

type AIOrchestratorClient interface {
	ProcessCall(ctx context.Context, in *ProcessCallRequest, opts ...grpc.CallOption) (*ProcessCallResponse, error)
	EndCall(ctx context.Context, in *EndCallRequest, opts ...grpc.CallOption) (*EndCallResponse, error)
}
func NewAIOrchestratorClient(cc grpc.ClientConnInterface) AIOrchestratorClient { return nil }

// Request / Response stubs
type StartMediaRequest struct {
	CallId string
	Codec string
	SampleRate int32
}
type StartMediaResponse struct {
	LocalRtpPort int32
	SessionId string
}
type StopMediaRequest struct {
	SessionId string
}
type StopMediaResponse struct {
	Success bool
}
type GetMetricsRequest struct {
	SessionId string
}
type GetMetricsResponse struct {
	Mos float32
	JitterMs int32
	PacketLossPercent float32
}

type ProcessCallRequest struct {
	CallId string
	InitialPrompt string
}
type ProcessCallResponse struct {
	OrchestratorSessionId string
	Accepted bool
}
type EndCallRequest struct {
	OrchestratorSessionId string
}
type EndCallResponse struct {
	TranscriptSummary string
}

type InitiateTransferRequest struct {
	CallId string
	TargetSipUri string
	TranscriptContext string
	Sentiment string
}
type InitiateTransferResponse struct {
	TransferId string
	Accepted bool
}

type ContactMessage struct {
	Id string
	PhoneNumber string
	Timezone string
}

type StartCampaignRequest struct {
	CampaignId string
	Contacts []*ContactMessage
	TargetAbandonmentRate float32
}
type StartCampaignResponse struct {
	Success bool
}
type StopCampaignRequest struct {
	CampaignId string
}
type StopCampaignResponse struct {
	Success bool
}
