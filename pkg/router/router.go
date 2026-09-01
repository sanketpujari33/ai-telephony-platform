package router

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CallRouterConfig holds configuration for the Call Router
type CallRouterConfig struct {
	GRPCListen           string
	MaxRecvMsgSize       int
	MaxSendMsgSize       int
	TenantLookupCacheTTL int
	TenantLookupCacheSize int
	MediaPortStart       int
	MediaPortEnd         int
	MediaPortsPerCall    int
	MediaServerAddress   string
	OrchestratorAddress  string
}

// CallRouter implements the Call Router service
type CallRouter struct {
	config           *CallRouterConfig
	logger           *zap.Logger
	tenantCache      *TenantCache
	portAllocator    *PortAllocator
	mediaClient      MediaServerClient
	orchestratorClient OrchestratorClient
	grpcServer       *grpc.Server
	mu               sync.RWMutex
	activeCalls      map[string]*CallRouteInfo
}

// TenantCache caches tenant configurations
type TenantCache struct {
	cache map[string]*CachedTenantConfig
	ttl   time.Duration
	mu    sync.RWMutex
}

// CachedTenantConfig holds cached tenant configuration
type CachedTenantConfig struct {
	Config    *config.CampaignConfig
	ExpiresAt time.Time
}

// PortAllocator manages media port allocation
type PortAllocator struct {
	startPort int
	endPort   int
	portsPerCall int
	allocated map[int]bool
	mu        sync.Mutex
	nextPort  int
}

// MediaServerClient is a client for the Media Server
type MediaServerClient interface {
	AllocatePorts(ctx context.Context, count int) ([]int, error)
	ReleasePorts(ctx context.Context, ports []int) error
	HealthCheck(ctx context.Context) error
}

// OrchestratorClient is a client for the Orchestrator
type OrchestratorClient interface {
	StartCall(ctx context.Context, req *StartCallRequest) (*StartCallResponse, error)
	EndCall(ctx context.Context, callID string) error
	HealthCheck(ctx context.Context) error
}

// CallRouteInfo holds information about an active call route
type CallRouteInfo struct {
	CallID           string
	TenantID         string
	MediaPorts       []int
	OrchestratorID   string
	CreatedAt        time.Time
	MediaServerAddr  string
}

// RouteRequest represents a call routing request
type RouteRequest struct {
	CallID     string
	FromURI    string
	ToURI      string
	ContactURI string
	SDPBody    string
	Transport  string
}

// RouteResponse represents a call routing response
type RouteResponse struct {
	MediaPorts      []int
	SDPAnswer       string
	OrchestratorID  string
	MediaServerAddr string
}

// StartCallRequest represents a request to start a call in the orchestrator
type StartCallRequest struct {
	CallID        string
	TenantID      string
	FromURI       string
	ToURI         string
	ContactURI    string
	SDPBody       string
	MediaPorts    []int
	MediaServerAddr string
}

// StartCallResponse represents a response from starting a call
type StartCallResponse struct {
	OrchestratorID string
	SessionID      string
}

// NewCallRouter creates a new Call Router
func NewCallRouter(cfg *CallRouterConfig, logger *zap.Logger) (*CallRouter, error) {
	router := &CallRouter{
		config:      cfg,
		logger:      logger,
		tenantCache: NewTenantCache(cfg.TenantLookupCacheTTL, cfg.TenantLookupCacheSize),
		portAllocator: NewPortAllocator(cfg.MediaPortStart, cfg.MediaPortEnd, cfg.MediaPortsPerCall),
		activeCalls: make(map[string]*CallRouteInfo),
	}

	// Initialize media server client
	router.mediaClient = NewMediaServerClient(cfg.MediaServerAddress, logger)

	// Initialize orchestrator client
	router.orchestratorClient = NewOrchestratorClient(cfg.OrchestratorAddress, logger)

	// Create gRPC server
	router.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
	)

	// Register services
	// pb.RegisterCallRouterServer(router.grpcServer, router)

	return router, nil
}

// Start starts the Call Router gRPC server
func (r *CallRouter) Start() error {
	r.logger.Info("Starting Call Router", zap.String("address", r.config.GRPCListen))

	lis, err := net.Listen("tcp", r.config.GRPCListen)
	if err != nil {
	    return fmt.Errorf("failed to listen: %w", err)
	}
	go r.grpcServer.Serve(lis)

	r.logger.Info("Call Router started")
	return nil
}

// Stop stops the Call Router
func (r *CallRouter) Stop() error {
	r.logger.Info("Stopping Call Router")
	r.grpcServer.GracefulStop()
	return nil
}

// RouteCall routes a call and returns media information
func (r *CallRouter) RouteCall(ctx context.Context, req *RouteRequest) (*RouteResponse, error) {
	r.logger.Info("Routing call",
		zap.String("call_id", req.CallID),
		zap.String("from", req.FromURI),
		zap.String("to", req.ToURI),
	)

	// Extract tenant ID from From URI
	tenantID := r.extractTenantID(req.FromURI)
	if tenantID == "" {
		tenantID = "default"
	}

	// Lookup tenant configuration
	_, err := r.getTenantConfig(ctx, tenantID)
	if err != nil {
		r.logger.Error("Failed to get tenant config", zap.Error(err), zap.String("tenant", tenantID))
		return nil, fmt.Errorf("tenant config lookup failed: %w", err)
	}

	// Allocate media ports
	mediaPorts, err := r.portAllocator.Allocate(r.config.MediaPortsPerCall)
	if err != nil {
		r.logger.Error("Failed to allocate media ports", zap.Error(err))
		return nil, fmt.Errorf("media port allocation failed: %w", err)
	}

	// Get media server address
	mediaServerAddr := r.config.MediaServerAddress
	if mediaServerAddr == "" {
		mediaServerAddr = "media-proxy:9091"
	}

	// Start call in orchestrator
	orchestratorResp, err := r.orchestratorClient.StartCall(ctx, &StartCallRequest{
		CallID:          req.CallID,
		TenantID:        tenantID,
		FromURI:         req.FromURI,
		ToURI:           req.ToURI,
		ContactURI:      req.ContactURI,
		SDPBody:         req.SDPBody,
		MediaPorts:      mediaPorts,
		MediaServerAddr: mediaServerAddr,
	})
	if err != nil {
		// Release ports on failure
		r.portAllocator.Release(mediaPorts)
		r.logger.Error("Failed to start call in orchestrator", zap.Error(err))
		return nil, fmt.Errorf("orchestrator start call failed: %w", err)
	}

	// Store call route info
	callInfo := &CallRouteInfo{
		CallID:           req.CallID,
		TenantID:         tenantID,
		MediaPorts:       mediaPorts,
		OrchestratorID:   orchestratorResp.OrchestratorID,
		CreatedAt:        time.Now(),
		MediaServerAddr:  mediaServerAddr,
	}

	r.mu.Lock()
	r.activeCalls[req.CallID] = callInfo
	r.mu.Unlock()

	// Build response
	response := &RouteResponse{
		MediaPorts:      mediaPorts,
		SDPAnswer:       "", // Will be filled by SIP ingress
		OrchestratorID:  orchestratorResp.OrchestratorID,
		MediaServerAddr: mediaServerAddr,
	}

	r.logger.Info("Call routed successfully",
		zap.String("call_id", req.CallID),
		zap.String("orchestrator", orchestratorResp.OrchestratorID),
		zap.Ints("media_ports", mediaPorts),
	)

	return response, nil
}

// EndCall ends a call and releases resources
func (r *CallRouter) EndCall(ctx context.Context, callID string) error {
	r.logger.Info("Ending call", zap.String("call_id", callID))

	r.mu.Lock()
	callInfo, exists := r.activeCalls[callID]
	if exists {
		delete(r.activeCalls, callID)
	}
	r.mu.Unlock()

	if !exists {
		r.logger.Warn("Call not found", zap.String("call_id", callID))
		return nil
	}

	// Release media ports
	if err := r.portAllocator.Release(callInfo.MediaPorts); err != nil {
		r.logger.Error("Failed to release media ports", zap.Error(err), zap.String("call_id", callID))
	}

	// End call in orchestrator
	if err := r.orchestratorClient.EndCall(ctx, callID); err != nil {
		r.logger.Error("Failed to end call in orchestrator", zap.Error(err), zap.String("call_id", callID))
	}

	r.logger.Info("Call ended", zap.String("call_id", callID))
	return nil
}

// GetCallInfo returns information about an active call
func (r *CallRouter) GetCallInfo(callID string) (*CallRouteInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.activeCalls[callID]
	return info, exists
}

// extractTenantID extracts tenant ID from a SIP URI
func (r *CallRouter) extractTenantID(uri string) string {
	// Parse URI and extract tenant from domain or user part
	// For example: sip:user@tenant.example.com -> tenant
	// This is a simplified implementation
	return "default"
}

// getTenantConfig retrieves tenant configuration with caching
func (r *CallRouter) getTenantConfig(ctx context.Context, tenantID string) (*config.CampaignConfig, error) {
	// Check cache first
	if cached := r.tenantCache.Get(tenantID); cached != nil {
		return cached.Config, nil
	}

	// TODO: Fetch from database/config service
	// For now, return default config
	tenantConfig := &config.CampaignConfig{}

	// Cache the result
	r.tenantCache.Set(tenantID, tenantConfig)

	return tenantConfig, nil
}

// NewTenantCache creates a new tenant cache
func NewTenantCache(ttlSeconds int, maxSize int) *TenantCache {
	return &TenantCache{
		cache: make(map[string]*CachedTenantConfig),
		ttl:   time.Duration(ttlSeconds) * time.Second,
	}
}

// Get retrieves a cached tenant config
func (tc *TenantCache) Get(tenantID string) *CachedTenantConfig {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	cached, exists := tc.cache[tenantID]
	if !exists {
		return nil
	}

	if time.Now().After(cached.ExpiresAt) {
		return nil
	}

	return cached
}

// Set stores a tenant config in cache
func (tc *TenantCache) Set(tenantID string, config *config.CampaignConfig) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.cache[tenantID] = &CachedTenantConfig{
		Config:    config,
		ExpiresAt: time.Now().Add(tc.ttl),
	}
}

// NewPortAllocator creates a new port allocator
func NewPortAllocator(startPort, endPort, portsPerCall int) *PortAllocator {
	return &PortAllocator{
		startPort:     startPort,
		endPort:       endPort,
		portsPerCall:  portsPerCall,
		allocated:     make(map[int]bool),
		nextPort:      startPort,
	}
}

// Allocate allocates a block of ports
func (pa *PortAllocator) Allocate(count int) ([]int, error) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	if count <= 0 {
		count = pa.portsPerCall
	}

	ports := make([]int, 0, count)
	attempts := 0
	maxAttempts := pa.endPort - pa.startPort + 1

	for len(ports) < count && attempts < maxAttempts {
		port := pa.nextPort
		pa.nextPort++
		if pa.nextPort > pa.endPort {
			pa.nextPort = pa.startPort
		}

		if !pa.allocated[port] {
			pa.allocated[port] = true
			ports = append(ports, port)
		}
		attempts++
	}

	if len(ports) < count {
		// Release any ports we allocated
		for _, p := range ports {
			delete(pa.allocated, p)
		}
		return nil, fmt.Errorf("insufficient ports available")
	}

	return ports, nil
}

// Release releases allocated ports
func (pa *PortAllocator) Release(ports []int) error {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	for _, port := range ports {
		delete(pa.allocated, port)
	}

	return nil
}

// MediaServerClientImpl implements MediaServerClient
type MediaServerClientImpl struct {
	address string
	logger  *zap.Logger
	conn    *grpc.ClientConn
	client  pb.MediaProxyClient
}

// NewMediaServerClient creates a new media server client
func NewMediaServerClient(address string, logger *zap.Logger) MediaServerClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to Media Proxy", zap.String("address", address), zap.Error(err))
	}
	
	var client pb.MediaProxyClient
	if conn != nil {
		client = pb.NewMediaProxyClient(conn)
	}
	
	return &MediaServerClientImpl{
		address: address,
		logger:  logger,
		conn:    conn,
		client:  client,
	}
}

func (c *MediaServerClientImpl) AllocatePorts(ctx context.Context, count int) ([]int, error) {
	if c.client == nil {
		return []int{10000, 10002}, nil // Fallback
	}
	
	req := &pb.StartMediaRequest{
		CallId:     uuid.New().String(),
		Codec:      "PCMU",
		SampleRate: 8000,
	}
	
	resp, err := c.client.StartMediaSession(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return []int{int(resp.LocalRtpPort)}, nil
}

func (c *MediaServerClientImpl) ReleasePorts(ctx context.Context, ports []int) error {
	// Not implemented perfectly here because we don't track session ID for ports in this demo stub
	return nil
}

func (c *MediaServerClientImpl) HealthCheck(ctx context.Context) error {
	// TODO: Implement health check
	return nil
}

// OrchestratorClientImpl implements OrchestratorClient
type OrchestratorClientImpl struct {
	address string
	logger  *zap.Logger
	conn    *grpc.ClientConn
	client  pb.AIOrchestratorClient
}

// NewOrchestratorClient creates a new orchestrator client
func NewOrchestratorClient(address string, logger *zap.Logger) OrchestratorClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to AI Orchestrator", zap.String("address", address), zap.Error(err))
	}
	
	var client pb.AIOrchestratorClient
	if conn != nil {
		client = pb.NewAIOrchestratorClient(conn)
	}

	return &OrchestratorClientImpl{
		address: address,
		logger:  logger,
		conn:    conn,
		client:  client,
	}
}

func (c *OrchestratorClientImpl) StartCall(ctx context.Context, req *StartCallRequest) (*StartCallResponse, error) {
	if c.client == nil {
		orchestratorID := fmt.Sprintf("orchestrator-%s", uuid.New().String()[:8])
		return &StartCallResponse{
			OrchestratorID: orchestratorID,
			SessionID:      uuid.New().String(),
		}, nil
	}
	
	resp, err := c.client.ProcessCall(ctx, &pb.ProcessCallRequest{
		CallId:         req.CallID,
		InitialPrompt:  "Hello, how can I help you?",
	})
	
	if err != nil {
		return nil, err
	}

	return &StartCallResponse{
		OrchestratorID: resp.OrchestratorSessionId,
		SessionID:      resp.OrchestratorSessionId,
	}, nil
}

func (c *OrchestratorClientImpl) EndCall(ctx context.Context, callID string) error {
	if c.client == nil {
		return nil
	}
	
	_, err := c.client.EndCall(ctx, &pb.EndCallRequest{
		OrchestratorSessionId: "orch_" + callID,
	})
	
	return err
}

func (c *OrchestratorClientImpl) HealthCheck(ctx context.Context) error {
	return nil
}