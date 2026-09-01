package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/campaign"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/compliance"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/observability"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = flag.Bool("version", false, "Print version and exit")
)

const (
	serviceName = "campaign-manager"
	versionStr  = "1.0.0"
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("%s version %s\n", serviceName, versionStr)
		os.Exit(0)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Campaign Manager service",
		zap.String("service", serviceName),
		zap.String("version", versionStr),
	)

	cfgPath := *configPath
	if cfgPath == "" {
		cfgPath = config.GetConfigPath()
	}

	cfg, err := config.Load(cfgPath, logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize observability
	ctx := context.Background()
	obs, err := observability.NewObservability(ctx, cfg.Observability)
	if err != nil {
		logger.Fatal("Failed to initialize observability", zap.Error(err))
	}
	defer obs.Close(ctx)
	logger = obs.Logger()

	// Initialize storage
	store, err := storage.NewStorage(ctx, cfg.Storage, logger)
	if err != nil {
		logger.Fatal("Failed to initialize storage", zap.Error(err))
	}
	defer store.Close()

	// Initialize compliance checker
	checker := compliance.NewChecker(cfg.Campaign.Manager.Compliance, cfg.Campaign.Manager.CallingWindow, store, logger)

	engine := campaign.NewCampaignEngine(&cfg.Campaign, checker, logger)
	
	grpcHandler := campaign.NewGRPCServer(engine, logger)
	
	lis, err := net.Listen("tcp", cfg.Campaign.Manager.GRPC.Listen)
	if err != nil {
		logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterCampaignManagerServer(grpcServer, grpcHandler)

	logger.Info("Campaign Manager service running",
		zap.String("grpc_listen", cfg.Campaign.Manager.GRPC.Listen),
		zap.Any("engine", engine),
	)
	
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, stopping Campaign Manager...")
	grpcServer.GracefulStop()
}
