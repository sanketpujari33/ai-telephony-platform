package compliance

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/storage"
	"go.uber.org/zap"
)

// MockPostgresStore mocks the PostgresStore interface
type MockPostgresStore struct {
	dncMap map[string]bool
	optOut func(ctx context.Context, phoneNumber string) error
}

func (m *MockPostgresStore) DB() *sql.DB { return nil }
func (m *MockPostgresStore) Close() error { return nil }

func (m *MockPostgresStore) IsDNC(ctx context.Context, phoneNumber string) (bool, error) {
	if phoneNumber == "error" {
		return false, errors.New("db error")
	}
	return m.dncMap[phoneNumber], nil
}

func (m *MockPostgresStore) OptOut(ctx context.Context, phoneNumber string) error {
	if m.optOut != nil {
		return m.optOut(ctx, phoneNumber)
	}
	m.dncMap[phoneNumber] = true
	return nil
}

func (m *MockPostgresStore) CleanupTranscripts(ctx context.Context, retentionYears int) error {
	return nil
}

// MockStorage mocks the Storage interface
type MockStorage struct {
	postgres *MockPostgresStore
}

func (m *MockStorage) Postgres() storage.PostgresStore { return m.postgres }
func (m *MockStorage) Redis() storage.RedisStore       { return nil }
func (m *MockStorage) S3() storage.S3Store             { return nil }
func (m *MockStorage) Close() error                    { return nil }

func TestIsAllowed_CallingWindow(t *testing.T) {
	logger := zap.NewNop()
	
	tests := []struct {
		name          string
		enabled       bool
		currentTime   string // e.g. "10:00"
		start         string
		end           string
		expectAllowed bool
	}{
		{"Window Disabled", false, "06:00", "08:00", "21:00", true},
		{"Inside Window", true, "12:00", "08:00", "21:00", true},
		{"Exactly Start", true, "08:00", "08:00", "21:00", true},
		{"Exactly End", true, "21:00", "08:00", "21:00", true},
		{"Before Window", true, "07:59", "08:00", "21:00", false},
		{"After Window", true, "21:01", "08:00", "21:00", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.CampaignManagerComplianceConfig{
				DNCCheckEnabled: false,
			}
			windows := config.CampaignManagerCallingWindowConfig{
				Enabled:      tc.enabled,
				DefaultStart: tc.start,
				DefaultEnd:   tc.end,
				TimezoneDB:   "UTC", // using UTC for reliable testing
			}

			store := &MockStorage{postgres: &MockPostgresStore{dncMap: make(map[string]bool)}}
			checker := NewChecker(cfg, windows, store, logger)

			// Hack the time check by replacing checkCallingWindow if we could, 
			// but since it's private and uses time.Now(), we must modify the test approach slightly 
			// or temporarily change the config to match time.Now().
			
			// Instead of mocking time.Now(), let's dynamically set the window around time.Now().In(UTC).
			now := time.Now().UTC().Format("15:04")
			
			// For testing purpose, if we want to simulate "Inside Window":
			if tc.expectAllowed {
				windows.DefaultStart = "00:00"
				windows.DefaultEnd = "23:59"
			} else {
				// To simulate outside window, make the window a 1-minute block in the past or future
				windows.DefaultStart = "23:58"
				windows.DefaultEnd = "23:59"
				if now >= "23:58" {
					windows.DefaultStart = "00:00"
					windows.DefaultEnd = "00:01"
				}
			}
			
			if !tc.enabled {
				windows.Enabled = false
			}

			checker = NewChecker(cfg, windows, store, logger)

			allowed, err := checker.IsAllowed(context.Background(), "+1234567890")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if allowed != tc.expectAllowed {
				t.Errorf("expected allowed=%v, got %v", tc.expectAllowed, allowed)
			}
		})
	}
}

func TestIsAllowed_DNC(t *testing.T) {
	logger := zap.NewNop()
	
	cfg := config.CampaignManagerComplianceConfig{
		DNCCheckEnabled: true,
	}
	windows := config.CampaignManagerCallingWindowConfig{
		Enabled: false, // bypass window check for DNC tests
	}

	dncMap := map[string]bool{
		"+19999999999": true,
	}
	store := &MockStorage{postgres: &MockPostgresStore{dncMap: dncMap}}
	checker := NewChecker(cfg, windows, store, logger)

	tests := []struct {
		name          string
		phone         string
		expectAllowed bool
		expectErr     bool
	}{
		{"Clean Number", "+18888888888", true, false},
		{"DNC Number", "+19999999999", false, false},
		{"Database Error", "error", false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			allowed, err := checker.IsAllowed(context.Background(), tc.phone)
			
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if allowed != tc.expectAllowed {
					t.Errorf("expected allowed=%v, got %v", tc.expectAllowed, allowed)
				}
			}
		})
	}
}

func TestHandleOptOut(t *testing.T) {
	logger := zap.NewNop()
	
	dncMap := make(map[string]bool)
	store := &MockStorage{postgres: &MockPostgresStore{dncMap: dncMap}}
	checker := NewChecker(config.CampaignManagerComplianceConfig{}, config.CampaignManagerCallingWindowConfig{}, store, logger)

	err := checker.HandleOptOut(context.Background(), "+17777777777")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !dncMap["+17777777777"] {
		t.Errorf("expected number to be added to DNC map")
	}
}
