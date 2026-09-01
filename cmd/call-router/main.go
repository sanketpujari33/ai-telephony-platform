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
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/router"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = flag.Bool("version", false, "Print version and exit")
)

const (
	serviceName = "call-router"
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

	logger.Info("Starting Call Router service",
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

	// Create Call Router config from platform config
	routerConfig := &router.CallRouterConfig{
		GRPCListen:            cfg.SIP.Router.GRPC.Listen,
		MaxRecvMsgSize:        cfg.SIP.Router.GRPC.MaxRecvMsgSize,
		MaxSendMsgSize:        cfg.SIP.Router.GRPC.MaxSendMsgSize,
		TenantLookupCacheTTL:  cfg.SIP.Router.TenantLookup.CacheTTL,
		TenantLookupCacheSize: cfg.SIP.Router.TenantLookup.CacheSize,
		MediaPortStart:        cfg.SIP.Router.MediaPorts.Start,
		MediaPortEnd:          cfg.SIP.Router.MediaPorts.End,
		MediaPortsPerCall:     cfg.SIP.Router.MediaPorts.PerCall,
		MediaServerAddress:    "media-proxy:9091",
		OrchestratorAddress:   "orchestrator:9092",
	}

	// Create Call Router
	callRouter, err := router.NewCallRouter(routerConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create Call Router", zap.Error(err))
	}

	// Start Call Router
	if err := callRouter.Start(); err != nil {
		logger.Fatal("Failed to start Call Router", zap.Error(err))
	}

	logger.Info("Call Router service started successfully")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, stopping server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop Call Router
	if err := callRouter.Stop(); err != nil {
		logger.Error("Error stopping Call Router", zap.Error(err))
	}

	// Wait for context cancellation or timeout
	select {
	case <-ctx.Done():
		logger.Info("Shutdown timeout exceeded")
	case <-time.After(100 * time.Millisecond):
	}

	logger.Info("Call Router service stopped")
}