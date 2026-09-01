package campaign

import (
	"context"
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/compliance"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestComplianceValidatorDNCAndCallingWindow(t *testing.T) {
	logger := zap.NewNop()
	compCfg := config.CampaignManagerComplianceConfig{DNCCheckEnabled: false}
	winCfg := config.CampaignManagerCallingWindowConfig{DefaultStart: "08:00", DefaultEnd: "21:00"}

	checker := compliance.NewChecker(compCfg, winCfg, nil, logger)

	// Check clean phone number
	allowed, err := checker.IsAllowed(context.Background(), "+15551234567")
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestPredictiveDialerPacing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.CampaignPredictiveDialerErlangCConfig{
		TargetUtilization: 0.85,
		SafetyMargin:      0.03,
		MaxAbandonRate:    0.03,
	}

	pd := NewPredictiveDialer(cfg, logger)

	// 10 agents, 50% connect rate (0.5)
	pacing := pd.CalculatePacing(10, 0.5, 180)
	assert.True(t, pacing >= 10 && pacing <= 40)

	// 0 available agents -> 0 pacing
	pacingZero := pd.CalculatePacing(0, 0.5, 180)
	assert.Equal(t, 0, pacingZero)
}

func TestAMDDetector(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.CampaignPredictiveDialerAMDConfig{
		HumanThreshold:     0.5,
		VoicemailThreshold: 0.7,
		SilenceDuration:    1.5,
		BeepDetection:      true,
	}

	amd := NewAMDDetector(cfg, logger)

	// Beep detected -> Machine
	resBeep := amd.AnalyzeAudio(100*time.Millisecond, 100*time.Millisecond, true)
	assert.Equal(t, AMDResultMachine, resBeep)

	// Long greeting -> Machine
	resGreeting := amd.AnalyzeAudio(3000*time.Millisecond, 500*time.Millisecond, false)
	assert.Equal(t, AMDResultMachine, resGreeting)

	// Short "Hello?" followed by silence -> Human
	resHuman := amd.AnalyzeAudio(400*time.Millisecond, 1200*time.Millisecond, false)
	assert.Equal(t, AMDResultHuman, resHuman)
}
