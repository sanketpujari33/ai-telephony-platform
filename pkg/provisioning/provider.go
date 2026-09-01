package provisioning

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// NumberProvider defines the interface for external trunk provisioning (T12.3)
type NumberProvider interface {
	SearchNumbers(ctx context.Context, areaCode string) ([]string, error)
	PurchaseNumber(ctx context.Context, number string) error
	ConfigureWebhook(ctx context.Context, number, webhookURL string) error
}

type TwilioProvider struct {
	accountSID string
	authToken  string
	logger     *zap.Logger
}

// NewTwilioProvider initializes a Twilio provisioning client
func NewTwilioProvider(sid, token string, logger *zap.Logger) *TwilioProvider {
	return &TwilioProvider{
		accountSID: sid,
		authToken:  token,
		logger:     logger,
	}
}

func (t *TwilioProvider) SearchNumbers(ctx context.Context, areaCode string) ([]string, error) {
	t.logger.Info("Searching for Twilio numbers", zap.String("areaCode", areaCode))
	// In production, HTTP GET to https://api.twilio.com/2010-04-01/Accounts/{AccountSid}/AvailablePhoneNumbers/US/Local.json
	return []string{"+1" + areaCode + "5551234", "+1" + areaCode + "5555678"}, nil
}

func (t *TwilioProvider) PurchaseNumber(ctx context.Context, number string) error {
	t.logger.Info("Purchasing Twilio number", zap.String("number", number))
	// In production, HTTP POST to https://api.twilio.com/2010-04-01/Accounts/{AccountSid}/IncomingPhoneNumbers.json
	return nil
}

func (t *TwilioProvider) ConfigureWebhook(ctx context.Context, number, webhookURL string) error {
	t.logger.Info("Configuring Twilio webhook", zap.String("number", number), zap.String("url", webhookURL))
	// In production, HTTP POST to update VoiceUrl on the IncomingPhoneNumber resource
	return nil
}

type TelnyxProvider struct {
	apiKey string
	logger *zap.Logger
}

// NewTelnyxProvider initializes a Telnyx provisioning client
func NewTelnyxProvider(key string, logger *zap.Logger) *TelnyxProvider {
	return &TelnyxProvider{
		apiKey: key,
		logger: logger,
	}
}

func (t *TelnyxProvider) SearchNumbers(ctx context.Context, areaCode string) ([]string, error) {
	t.logger.Info("Searching for Telnyx numbers", zap.String("areaCode", areaCode))
	// HTTP GET to https://api.telnyx.com/v2/available_phone_numbers
	return []string{"+1" + areaCode + "5559999"}, nil
}

func (t *TelnyxProvider) PurchaseNumber(ctx context.Context, number string) error {
	t.logger.Info("Purchasing Telnyx number", zap.String("number", number))
	// HTTP POST to https://api.telnyx.com/v2/number_orders
	return nil
}

func (t *TelnyxProvider) ConfigureWebhook(ctx context.Context, number, webhookURL string) error {
	t.logger.Info("Configuring Telnyx connection", zap.String("number", number), zap.String("url", webhookURL))
	// Telnyx typically requires mapping the number to a Connection/Application which has the webhook URL
	return nil
}

// ProvisioningManager orchestrates operations across multiple providers
type ProvisioningManager struct {
	providers map[string]NumberProvider
	logger    *zap.Logger
}

func NewProvisioningManager(logger *zap.Logger) *ProvisioningManager {
	return &ProvisioningManager{
		providers: make(map[string]NumberProvider),
		logger:    logger,
	}
}

func (m *ProvisioningManager) RegisterProvider(name string, provider NumberProvider) {
	m.providers[name] = provider
}

func (m *ProvisioningManager) GetProvider(name string) (NumberProvider, error) {
	p, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return p, nil
}
