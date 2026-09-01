package router

import (
	"testing"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
)

func TestPortAllocator(t *testing.T) {
	allocator := NewPortAllocator(10000, 10005, 2)

	// Should be able to allocate 3 pairs of ports
	p1, err := allocator.Allocate(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p1) != 2 || p1[0] != 10000 || p1[1] != 10001 {
		t.Errorf("unexpected ports: %v", p1)
	}

	p2, err := allocator.Allocate(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p2) != 2 || p2[0] != 10002 || p2[1] != 10003 {
		t.Errorf("unexpected ports: %v", p2)
	}

	p3, err := allocator.Allocate(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p3) != 2 || p3[0] != 10004 || p3[1] != 10005 {
		t.Errorf("unexpected ports: %v", p3)
	}

	// Should fail because no ports left
	_, err = allocator.Allocate(2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// Release p2 and it should be available again
	_ = allocator.Release(p2)

	p4, err := allocator.Allocate(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p4) != 2 {
		t.Errorf("expected 2 ports, got %d", len(p4))
	}
}

func TestTenantCache(t *testing.T) {
	cache := NewTenantCache(1, 10)
	cfg := &config.CampaignConfig{}

	// Test Set and Get
	cache.Set("tenant-1", cfg)

	val := cache.Get("tenant-1")
	if val == nil {
		t.Fatalf("expected to find tenant config in cache")
	}

	// Expired key check
	cacheShort := &TenantCache{
		cache: make(map[string]*CachedTenantConfig),
		ttl:   1 * time.Millisecond,
	}
	cacheShort.Set("tenant-2", cfg)
	time.Sleep(5 * time.Millisecond)

	valExpired := cacheShort.Get("tenant-2")
	if valExpired != nil {
		t.Errorf("expected cache to expire")
	}
}

