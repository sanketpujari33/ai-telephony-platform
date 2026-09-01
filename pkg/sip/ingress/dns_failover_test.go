package ingress

import (
	"context"
	"net"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestDNSResolver_Disabled(t *testing.T) {
	logger := zap.NewNop()
	cfg := &DNSResolverConfig{
		Enabled: false,
		Domain:  "example.com",
	}

	resolver := NewDNSResolver(cfg, logger)

	_, err := resolver.ResolveSRV(context.Background())
	if err == nil {
		t.Errorf("expected error when disabled, got nil")
	}
}

func TestDNSResolver_Cache(t *testing.T) {
	logger := zap.NewNop()
	cfg := &DNSResolverConfig{
		Enabled:  true,
		Domain:   "example.com",
		CacheTTL: 100 * time.Millisecond,
	}

	resolver := NewDNSResolver(cfg, logger)

	// Inject cache
	cachedRecords := []*net.SRV{
		{Target: "sip1.example.com.", Port: 5060, Priority: 10, Weight: 100},
	}
	
	resolver.mu.Lock()
	resolver.cache[cfg.Domain] = &SRVCacheEntry{
		Records:   cachedRecords,
		ExpiresAt: time.Now().Add(cfg.CacheTTL),
	}
	resolver.mu.Unlock()

	// Should return from cache without doing a real DNS lookup
	records, err := resolver.ResolveSRV(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Target != "sip1.example.com." {
		t.Errorf("expected target sip1.example.com., got %s", records[0].Target)
	}

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Since cache expired, it will attempt a real DNS lookup and fail because example.com doesn't have SIP SRV
	// However, we just want to ensure it doesn't return the cached value.
	// Since example.com might actually have SRV records if we are very unlucky, or the lookup might just timeout,
	// we just check if it's the exact same cache pointer by seeing if the slice is our injected one.
	
	records2, _ := resolver.ResolveSRV(context.Background())
	if len(records2) > 0 && records2[0].Target == "sip1.example.com." {
		// It's possible it returned the same, but highly unlikely unless example.com has that exact SRV
		t.Errorf("expected cache to expire and perform real lookup")
	}
}
