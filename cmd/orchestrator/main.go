package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	pb "github.com/ai-telephony-platform/ai-telephony-platform/pkg/api/telephony/v1"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/observability"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/orchestrator"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = flag.Bool("version", false, "Print version and exit")
)

const (
	serviceName = "orchestrator"
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

	logger.Info("Starting AI Orchestrator service",
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
	
	// Initialize proper observability
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

	grpcListenAddr := ":50051"
	server := orchestrator.NewOrchestratorServer(cfg, logger) // In a real app we'd pass store and obs here
	grpcHandler := orchestrator.NewGRPCServer(server, logger)
	
	lis, err := net.Listen("tcp", grpcListenAddr)
	if err != nil {
		logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterAIOrchestratorServer(grpcServer, grpcHandler)

	logger.Info("AI Orchestrator service initialized and running", zap.String("grpc_listen", grpcListenAddr))
	
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, stopping AI Orchestrator...")
	grpcServer.GracefulStop()
}
