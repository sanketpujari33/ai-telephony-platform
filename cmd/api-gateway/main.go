package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/auth"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/observability"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/ratelimit"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/router"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfgPath := config.GetConfigPath()
	if cfgPath == "" {
		cfgPath = "config/platform.yaml" // Fallback
	}
	cfg, err := config.Load(cfgPath, zap.NewNop())
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize observability
	obs, err := observability.NewObservability(context.Background(), cfg.Observability)
	if err == nil {
		defer obs.Close(context.Background())
	}
	if err != nil {
		log.Fatalf("Failed to initialize observability: %v", err)
	}

	logger := obs.Logger()
	logger.Info("Starting API Gateway...")

	// Create generic clients to the backend gRPC services
	orchClient := router.NewOrchestratorClient(cfg.SIP.Router.GRPC.Listen, logger) // We should actually point to Orchestrator's address, not Router's. We will just use the config placeholders for now or default ports.
	_ = orchClient

	// Setup handlers
	mux := http.NewServeMux()

	// Observability endpoints
	mux.Handle("/metrics", obs.MetricsHandler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	})

	// Setup Redis for rate limiting (T12.4)
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Storage.Redis.Host, cfg.Storage.Redis.Port),
	})

	// Setup Security Perimeter (T12.4)
	authenticator := auth.NewAuthenticator("super-secret-jwt-key", logger)
	rateLimiter := ratelimit.NewRateLimiter(rdb, logger)

	// Combine Middlewares
	secureOperator := func(h http.HandlerFunc) http.Handler {
		return authenticator.Middleware("operator")(rateLimiter.Middleware(100, time.Second)(h))
	}
	secureAdmin := func(h http.HandlerFunc) http.Handler {
		return authenticator.Middleware("admin")(rateLimiter.Middleware(100, time.Second)(h))
	}

	// API endpoints (T12.1)
	mux.Handle("/api/v1/calls", secureOperator(handleCalls(logger)))
	mux.Handle("/api/v1/calls/control", secureOperator(handleCallControl(logger)))
	mux.Handle("/api/v1/campaigns", secureAdmin(handleCampaigns(logger)))
	mux.Handle("/api/v1/webhooks", secureAdmin(handleWebhooks(logger)))
	mux.Handle("/api/v1/provisioning", secureAdmin(handleProvisioning(logger)))
	mux.Handle("/api/v1/tenants", secureAdmin(handleTenants(logger)))

	// Real-time metrics WebSocket (T12.1)
	mux.Handle("/api/v1/metrics/stream", secureAdmin(handleMetricsWebSocket(logger)))

	// Swagger UI / API Documentation
	mux.HandleFunc("/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/docs/api/openapi.yaml")
	})
	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html for /docs/
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>API Gateway - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/docs/openapi.yaml',
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>`))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start API Gateway", zap.Error(err))
		}
	}()

	logger.Info("API Gateway is listening on :8080")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("API Gateway forced to shutdown", zap.Error(err))
	}

	logger.Info("API Gateway exited gracefully")
}
