package observability

import (
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes the global zap logger
func InitLogger(cfg config.ObservabilityLoggingConfig) (*zap.Logger, error) {
	var zapCfg zap.Config
	
	if cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	if cfg.Output == "stdout" {
		zapCfg.OutputPaths = []string{"stdout"}
	} else if cfg.Output != "" {
		zapCfg.OutputPaths = []string{cfg.Output}
	}

	logger, err := zapCfg.Build(zap.Fields(zap.String("service", cfg.ServiceName)))
	if err != nil {
		return nil, err
	}

	// Replace the global logger so zap.L() works anywhere
	zap.ReplaceGlobals(logger)

	return logger, nil
}
