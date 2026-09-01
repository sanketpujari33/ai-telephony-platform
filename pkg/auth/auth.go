package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type contextKey string

const (
	TenantContextKey contextKey = "tenant_id"
	RoleContextKey   contextKey = "role"
)

// Authenticator handles JWT validation and RBAC
type Authenticator interface {
	Middleware(requiredRole string) func(http.Handler) http.Handler
	VerifyToken(token string) (tenantID string, role string, err error)
}

type authenticatorImpl struct {
	secret string
	logger *zap.Logger
}

// NewAuthenticator initializes a basic JWT/RBAC authenticator
func NewAuthenticator(secret string, logger *zap.Logger) Authenticator {
	return &authenticatorImpl{
		secret: secret,
		logger: logger,
	}
}

func (a *authenticatorImpl) VerifyToken(token string) (string, string, error) {
	// In a production environment, use github.com/golang-jwt/jwt to parse and verify the JWT.
	// For this system, we will mock verification to return predefined roles based on token prefix for demonstration.
	if token == "" {
		return "", "", errors.New("missing token")
	}

	if strings.HasPrefix(token, "admin-") {
		return "tenant-admin", "admin", nil
	} else if strings.HasPrefix(token, "operator-") {
		return "tenant-ops", "operator", nil
	} else if strings.HasPrefix(token, "viewer-") {
		return "tenant-view", "viewer", nil
	}

	return "", "", errors.New("invalid signature")
}

func (a *authenticatorImpl) Middleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: missing Bearer token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			tenantID, role, err := a.VerifyToken(token)
			if err != nil {
				a.logger.Warn("Failed auth", zap.Error(err), zap.String("ip", r.RemoteAddr))
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			if !a.hasAccess(role, requiredRole) {
				a.logger.Warn("RBAC denied", zap.String("role", role), zap.String("required", requiredRole))
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), TenantContextKey, tenantID)
			ctx = context.WithValue(ctx, RoleContextKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (a *authenticatorImpl) hasAccess(currentRole, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"viewer":   1,
		"operator": 2,
		"admin":    3,
	}

	return roleHierarchy[currentRole] >= roleHierarchy[requiredRole]
}
