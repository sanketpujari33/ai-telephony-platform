package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestVerifyToken(t *testing.T) {
	logger := zap.NewNop()
	auth := NewAuthenticator("secret", logger)

	tests := []struct {
		name          string
		token         string
		expectTenant  string
		expectRole    string
		expectErr     bool
	}{
		{"Admin Token", "admin-123", "tenant-admin", "admin", false},
		{"Operator Token", "operator-123", "tenant-ops", "operator", false},
		{"Viewer Token", "viewer-123", "tenant-view", "viewer", false},
		{"Invalid Signature", "fake-token", "", "", true},
		{"Empty Token", "", "", "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tenant, role, err := auth.VerifyToken(tc.token)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tenant != tc.expectTenant {
					t.Errorf("expected tenant %s, got %s", tc.expectTenant, tenant)
				}
				if role != tc.expectRole {
					t.Errorf("expected role %s, got %s", tc.expectRole, role)
				}
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	logger := zap.NewNop()
	auth := NewAuthenticator("secret", logger)

	tests := []struct {
		name           string
		authHeader     string
		requiredRole   string
		expectedStatus int
	}{
		{"Valid Admin Accessing Admin", "Bearer admin-123", "admin", http.StatusOK},
		{"Valid Operator Accessing Operator", "Bearer operator-123", "operator", http.StatusOK},
		{"Admin Accessing Operator", "Bearer admin-123", "operator", http.StatusOK},
		{"Operator Accessing Admin", "Bearer operator-123", "admin", http.StatusForbidden},
		{"Viewer Accessing Operator", "Bearer viewer-123", "operator", http.StatusForbidden},
		{"Missing Bearer Prefix", "admin-123", "admin", http.StatusUnauthorized},
		{"Invalid Token", "Bearer fake", "viewer", http.StatusUnauthorized},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := auth.Middleware(tc.requiredRole)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}
