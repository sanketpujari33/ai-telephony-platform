package observability

import (
	"context"
	"net/http"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

	// Observability encapsulates logging, metrics, and tracing
type Observability interface {
	Logger() *zap.Logger
	MetricsHandler() http.Handler
	Close(ctx context.Context) error
}

type observabilityImpl struct {
	logger *zap.Logger
	tp     *sdktrace.TracerProvider
	cfg    config.ObservabilityConfig
}

// NewObservability initializes the observability stack
func NewObservability(ctx context.Context, cfg config.ObservabilityConfig) (Observability, error) {
	logger, err := InitLogger(cfg.Logging)
	if err != nil {
		return nil, err
	}
	
	if cfg.Metrics.Enabled {
		InitMetrics(cfg.Metrics)
	}

	tp, err := InitTracing(ctx, cfg.Tracing)
	if err != nil {
		return nil, err
	}

	return &observabilityImpl{
		logger: logger,
		tp:     tp,
		cfg:    cfg,
	}, nil
}

func (o *observabilityImpl) Logger() *zap.Logger {
	return o.logger
}

func (o *observabilityImpl) MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func (o *observabilityImpl) Close(ctx context.Context) error {
	if o.logger != nil {
		_ = o.logger.Sync()
	}
	if o.tp != nil {
		return o.tp.Shutdown(ctx)
	}
	return nil
}
