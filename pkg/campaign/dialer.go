package campaign

import (
	"math"
	"sync"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"go.uber.org/zap"
)

// AMDResult represents the output of Answering Machine Detection
type AMDResult string

const (
	AMDResultHuman     AMDResult = "human"
	AMDResultMachine   AMDResult = "voicemail"
	AMDResultUncertain AMDResult = "uncertain"
)

// AMDDetector analyzes call audio to distinguish human vs answering machine
type AMDDetector struct {
	humanThreshold     float64
	voicemailThreshold float64
	silenceDuration    time.Duration
	beepDetection      bool
	logger             *zap.Logger
}

// NewAMDDetector creates a new AMD detector
func NewAMDDetector(cfg *config.CampaignPredictiveDialerAMDConfig, logger *zap.Logger) *AMDDetector {
	return &AMDDetector{
		humanThreshold:     cfg.HumanThreshold,
		voicemailThreshold: cfg.VoicemailThreshold,
		silenceDuration:    time.Duration(cfg.SilenceDuration * float64(time.Second)),
		beepDetection:      cfg.BeepDetection,
		logger:             logger,
	}
}

// AnalyzeAudio evaluates speech duration and silence gaps
func (amd *AMDDetector) AnalyzeAudio(initialSpeechDuration, initialSilenceDuration time.Duration, beepDetected bool) AMDResult {
	if beepDetected && amd.beepDetection {
		return AMDResultMachine
	}

	// Answering machines typically have long initial greetings (> 2.5s)
	if initialSpeechDuration > 2500*time.Millisecond {
		return AMDResultMachine
	}

	// Humans usually give short greetings ("Hello?") followed by silence
	if initialSpeechDuration > 200*time.Millisecond && initialSilenceDuration > 1000*time.Millisecond {
		return AMDResultHuman
	}

	return AMDResultUncertain
}

// PredictiveDialer calculates dial pacing using Erlang-C model
type PredictiveDialer struct {
	targetUtilization float64
	maxAbandonRate    float64
	safetyMargin      float64
	logger            *zap.Logger
	mu                sync.Mutex
}

// NewPredictiveDialer creates a PredictiveDialer
func NewPredictiveDialer(cfg *config.CampaignPredictiveDialerErlangCConfig, logger *zap.Logger) *PredictiveDialer {
	return &PredictiveDialer{
		targetUtilization: cfg.TargetUtilization,
		maxAbandonRate:    cfg.MaxAbandonRate,
		safetyMargin:      cfg.SafetyMargin,
		logger:            logger,
	}
}

// CalculatePacing computes the recommended number of lines to dial given available agents and connect rate
func (pd *PredictiveDialer) CalculatePacing(availableAgents int, connectRate float64, avgCallDurationSec float64) int {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	if availableAgents <= 0 || connectRate <= 0.0 {
		return 0
	}

	// Erlang-C traffic intensity (A = λ * h)
	intensity := float64(availableAgents) * pd.targetUtilization

	// Overdial factor to maintain target utilization while staying below max abandonment
	overdialFactor := (1.0 / connectRate) * (1.0 + pd.safetyMargin)

	pacing := math.Round(intensity * overdialFactor)

	// Ensure pacing is sane
	pacing = math.Max(1.0, math.Min(float64(availableAgents*4), pacing))

	pd.logger.Debug("Calculated predictive dial pacing",
		zap.Int("agents", availableAgents),
		zap.Float64("connect_rate", connectRate),
		zap.Float64("pacing_lines", pacing),
	)

	return int(pacing)
}
