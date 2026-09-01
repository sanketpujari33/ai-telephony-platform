package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/sip/ingress"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = flag.Bool("version", false, "Print version and exit")
)

const (
	serviceName = "sip-ingress"
	versionStr  = "1.0.0"
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("%s version %s\n", serviceName, versionStr)
		os.Exit(0)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting SIP Ingress service",
		zap.String("service", serviceName),
		zap.String("version", versionStr),
	)

	// Load configuration
	cfgPath := *configPath
	if cfgPath == "" {
		cfgPath = config.GetConfigPath()
	}

	cfg, err := config.Load(cfgPath, logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded",
		zap.String("config_file", cfgPath),
	)

	// Create SIP server config from platform config
	sipConfig := &ingress.SIPServerConfig{
		UDPListen:            cfg.SIP.Ingress.Listen.UDP,
		TCPListen:            cfg.SIP.Ingress.Listen.TCP,
		TLSListen:            cfg.SIP.Ingress.Listen.TLS,
		WSListen:             cfg.SIP.Ingress.Listen.WS,
		WSSListen:            cfg.SIP.Ingress.Listen.WSS,
		TLSCertFile:          cfg.SIP.Ingress.TLS.CertFile,
		TLSKeyFile:           cfg.SIP.Ingress.TLS.KeyFile,
		TLSCAFile:            cfg.SIP.Ingress.TLS.CAFile,
		VerifyClient:         cfg.SIP.Ingress.TLS.VerifyClient,
		WorkerCount:          cfg.SIP.Ingress.Workers,
		MaxRequestsPerWorker: cfg.SIP.Ingress.MaxRequests,
		Trunks:               make([]ingress.SIPTrunkConfig, len(cfg.SIP.Ingress.Trunks)),
		DNSSRVEnabled:        cfg.SIP.Ingress.DNSSRV.Enabled,
		DNSSRVDomain:         cfg.SIP.Ingress.DNSSRV.Domain,
		DNSSRVTransport:      ingress.SIPTransport(cfg.SIP.Ingress.DNSSRV.Transport),
		KeepaliveEnabled:     cfg.SIP.Ingress.Keepalive.Enabled,
		KeepaliveInterval:    time.Duration(cfg.SIP.Ingress.Keepalive.Interval) * time.Second,
		KeepaliveTimeout:     time.Duration(cfg.SIP.Ingress.Keepalive.Timeout) * time.Second,
		KeepaliveRetryCount:  cfg.SIP.Ingress.Keepalive.RetryCount,
		RateLimitEnabled:     cfg.SIP.Ingress.RateLimit.Enabled,
		RateLimitRPS:         cfg.SIP.Ingress.RateLimit.RequestsPerSecond,
		RateLimitBurst:       cfg.SIP.Ingress.RateLimit.Burst,
	}

	// Convert trunk configs
	for i, trunk := range cfg.SIP.Ingress.Trunks {
		sipConfig.Trunks[i] = ingress.SIPTrunkConfig{
			Name:      trunk.Name,
			Host:      trunk.Host,
			Port:      trunk.Port,
			Transport: ingress.SIPTransport(trunk.Transport),
			Username:  trunk.Username,
			Password:  trunk.Password,
			Realm:     trunk.Realm,
			Priority:  trunk.Priority,
			Weight:    trunk.Weight,
		}
	}

	// Create SIP server
	server, err := ingress.NewSIPServer(sipConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create SIP server", zap.Error(err))
	}

	// Create and set handler
	handler := NewSIPIngressHandler(cfg, logger)
	server.SetHandler(handler)

	// Start server
	if err := server.Start(); err != nil {
		logger.Fatal("Failed to start SIP server", zap.Error(err))
	}

	logger.Info("SIP Ingress service started successfully")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, stopping server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop server
	if err := server.Stop(); err != nil {
		logger.Error("Error stopping server", zap.Error(err))
	}

	// Wait for context cancellation or timeout
	select {
	case <-ctx.Done():
		logger.Info("Shutdown timeout exceeded")
	case <-time.After(100 * time.Millisecond):
	}

	logger.Info("SIP Ingress service stopped")
}

// SIPIngressHandler implements the SIPHandler interface
type SIPIngressHandler struct {
	config *config.Config
	logger *zap.Logger
	routerClient *RouterClient
}

// NewSIPIngressHandler creates a new SIP ingress handler
func NewSIPIngressHandler(cfg *config.Config, logger *zap.Logger) *SIPIngressHandler {
	return &SIPIngressHandler{
		config: cfg,
		logger: logger,
		routerClient: NewRouterClient(cfg.SIP.Router.GRPC.Listen, logger),
	}
}

// HandleInvite handles an incoming INVITE request
func (h *SIPIngressHandler) HandleInvite(invite *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Info("Handling INVITE",
		zap.String("call_id", invite.CallID),
		zap.String("from", invite.FromURI.String()),
		zap.String("to", invite.ToURI.String()),
	)

	// Call the Call Router to get routing information
	routeResp, err := h.routerClient.RouteCall(context.Background(), &RouteRequest{
		CallID:      invite.CallID,
		FromURI:     invite.FromURI.String(),
		ToURI:       invite.ToURI.String(),
		ContactURI:  invite.ContactURI.String(),
		SDPBody:     invite.SDPBody,
		Transport:   string(invite.Transport),
	})
	if err != nil {
		h.logger.Error("Failed to route call", zap.Error(err))
		return h.buildErrorResponse(invite, 500, "Call routing failed"), nil
	}

	// Build response with media information from router
	response := h.buildSuccessResponse(invite, routeResp)
	return response, nil
}

// HandleACK handles an incoming ACK request
func (h *SIPIngressHandler) HandleACK(ack *ingress.SIPInvite) error {
	h.logger.Debug("Handling ACK", zap.String("call_id", ack.CallID))
	// ACK doesn't require a response
	return nil
}

// HandleBYE handles an incoming BYE request
func (h *SIPIngressHandler) HandleBYE(bye *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Info("Handling BYE", zap.String("call_id", bye.CallID))

	// Notify call router about call termination
	if err := h.routerClient.EndCall(context.Background(), bye.CallID); err != nil {
		h.logger.Error("Failed to notify router of call end", zap.Error(err))
	}

	// Build 200 OK response
	response := h.buildSimpleResponse(bye, 200, "OK")
	return response, nil
}

// HandleReInvite handles a re-INVITE request
func (h *SIPIngressHandler) HandleReInvite(invite *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Info("Handling re-INVITE", zap.String("call_id", invite.CallID))

	// Handle re-INVITE (e.g., for hold/resume, codec renegotiation)
	// For now, just accept with current SDP
	response := h.buildSuccessResponse(invite, &RouteResponse{
		MediaPorts: []int{},
		SDPAnswer:  invite.SDPBody, // Echo back for simplicity
	})
	return response, nil
}

// HandleOptions handles an incoming OPTIONS request
func (h *SIPIngressHandler) HandleOptions(options *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Debug("Handling OPTIONS", zap.String("call_id", options.CallID))

	// Build 200 OK response for keepalive
	response := h.buildSimpleResponse(options, 200, "OK")
	return response, nil
}

// HandleCancel handles an incoming CANCEL request
func (h *SIPIngressHandler) HandleCancel(cancel *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Info("Handling CANCEL", zap.String("call_id", cancel.CallID))

	// Build 200 OK response
	response := h.buildSimpleResponse(cancel, 200, "OK")
	return response, nil
}

// HandleRegister handles an incoming REGISTER request
func (h *SIPIngressHandler) HandleRegister(register *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Info("Handling REGISTER", zap.String("call_id", register.CallID))

	// For now, accept all registrations
	// In production, you'd authenticate and store in location service
	response := h.buildSimpleResponse(register, 200, "OK")
	return response, nil
}

// HandleInfo handles an incoming INFO request
func (h *SIPIngressHandler) HandleInfo(info *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Debug("Handling INFO", zap.String("call_id", info.CallID))

	// Handle DTMF via SIP INFO
	response := h.buildSimpleResponse(info, 200, "OK")
	return response, nil
}

// HandleRefer handles an incoming REFER request
func (h *SIPIngressHandler) HandleRefer(refer *ingress.SIPInvite) (*ingress.SIPResponse, error) {
	h.logger.Info("Handling REFER", zap.String("call_id", refer.CallID))

	// Handle call transfer
	// For now, accept
	response := h.buildSimpleResponse(refer, 202, "Accepted")
	return response, nil
}

// HandleResponse handles an incoming SIP response
func (h *SIPIngressHandler) HandleResponse(response *ingress.SIPResponse) error {
	h.logger.Debug("Handling SIP response",
		zap.String("call_id", response.CallID),
		zap.Int("status_code", response.StatusCode),
	)
	return nil
}

// buildSimpleResponse builds a simple SIP response
func (h *SIPIngressHandler) buildSimpleResponse(invite *ingress.SIPInvite, statusCode int, reason string) *ingress.SIPResponse {
	parser := ingress.NewSIPParser()
	response, err := parser.BuildResponse(invite.Request, statusCode, reason, nil, "")
	if err != nil {
		h.logger.Error("Failed to build response", zap.Error(err))
		return nil
	}

	return &ingress.SIPResponse{
		Response:     response,
		Transport:    invite.Transport,
		LocalAddr:    invite.LocalAddr,
		RemoteAddr:   invite.RemoteAddr,
		StatusCode:   statusCode,
		ReasonPhrase: reason,
		SentAt:       time.Now(),
	}
}

// buildErrorResponse builds an error SIP response
func (h *SIPIngressHandler) buildErrorResponse(invite *ingress.SIPInvite, statusCode int, reason string) *ingress.SIPResponse {
	return h.buildSimpleResponse(invite, statusCode, reason)
}

// buildSuccessResponse builds a successful SIP response with SDP
func (h *SIPIngressHandler) buildSuccessResponse(invite *ingress.SIPInvite, routeResp *RouteResponse) *ingress.SIPResponse {
	parser := ingress.NewSIPParser()

	// Use SDP from router if available, otherwise create default
	sdpBody := routeResp.SDPAnswer
	if sdpBody == "" {
		// Create default SDP answer
		localIP := invite.LocalAddr.IP.String()
		localPort := routeResp.MediaPorts[0] // Use first media port
		if localPort == 0 {
			localPort = 10000
		}
		sdpBody = parser.CreateSDPResponse(localIP, localPort, []string{"111", "0", "8", "9"})
	}

	response, err := parser.BuildResponse(invite.Request, 200, "OK", nil, sdpBody)
	if err != nil {
		h.logger.Error("Failed to build success response", zap.Error(err))
		return h.buildErrorResponse(invite, 500, "Failed to build response")
	}

	return &ingress.SIPResponse{
		Response:     response,
		Transport:    invite.Transport,
		LocalAddr:    invite.LocalAddr,
		RemoteAddr:   invite.RemoteAddr,
		StatusCode:   200,
		ReasonPhrase: "OK",
		SDPBody:      sdpBody,
		SentAt:       time.Now(),
	}
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
	MediaPorts []int
	SDPAnswer  string
	OrchestratorID string
}

// RouterClient is a gRPC client for the Call Router
type RouterClient struct {
	address string
	logger  *zap.Logger
	// conn    *grpc.ClientConn
	// client  pb.CallRouterClient
}

// NewRouterClient creates a new router client
func NewRouterClient(address string, logger *zap.Logger) *RouterClient {
	return &RouterClient{
		address: address,
		logger:  logger,
	}
}

// RouteCall routes a call via the Call Router
func (c *RouterClient) RouteCall(ctx context.Context, req *RouteRequest) (*RouteResponse, error) {
	c.logger.Debug("Routing call", zap.String("call_id", req.CallID))

	// In a full deployment, this uses grpc.DialContext(c.address)
	// For standalone compilation, we return the local route allocation:
	return &RouteResponse{
		MediaPorts:     []int{10000, 10002},
		SDPAnswer:      "",
		OrchestratorID: "orchestrator-1",
	}, nil
}

// EndCall notifies the router that a call has ended
func (c *RouterClient) EndCall(ctx context.Context, callID string) error {
	c.logger.Debug("Ending call", zap.String("call_id", callID))
	// In a full deployment, this issues the EndCall gRPC command.
	return nil
}