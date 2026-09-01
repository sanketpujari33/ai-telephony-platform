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
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/observability"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/storage"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/transfer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = flag.Bool("version", false, "Print version and exit")
)

const (
	serviceName = "transfer-manager"
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

	logger.Info("Starting Transfer Manager service",
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

	engine := transfer.NewTransferEngine(&cfg.Transfer, logger) // We would inject store here if needed
	grpcHandler := transfer.NewGRPCServer(engine, logger)
	
	lis, err := net.Listen("tcp", cfg.Transfer.Manager.GRPC.Listen)
	if err != nil {
		logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterTransferManagerServer(grpcServer, grpcHandler)

	logger.Info("Transfer Manager service running",
		zap.String("grpc_listen", cfg.Transfer.Manager.GRPC.Listen),
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
	logger.Info("Shutdown signal received, stopping Transfer Manager...")
	grpcServer.GracefulStop()
}
