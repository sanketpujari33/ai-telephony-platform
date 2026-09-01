package ingress

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

// DNSResolver handles DNS SRV resolution for SIP trunks
type DNSResolver struct {
	config *DNSResolverConfig
	logger *zap.Logger
	cache  map[string]*SRVCacheEntry
	mu     sync.RWMutex
}

// DNSResolverConfig holds DNS resolver configuration
type DNSResolverConfig struct {
	Enabled     bool
	Domain      string
	Transport   SIPTransport
	Timeout     time.Duration
	CacheTTL    time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
}

// SRVCacheEntry holds cached SRV records
type SRVCacheEntry struct {
	Records  []*net.SRV
	ExpiresAt time.Time
}

// NewDNSResolver creates a new DNS resolver
func NewDNSResolver(config *DNSResolverConfig, logger *zap.Logger) *DNSResolver {
	return &DNSResolver{
		config: config,
		logger: logger,
		cache:  make(map[string]*SRVCacheEntry),
	}
}

// ResolveSRV resolves SRV records for a SIP domain
func (r *DNSResolver) ResolveSRV(ctx context.Context) ([]*net.SRV, error) {
	if !r.config.Enabled {
		return nil, fmt.Errorf("DNS SRV resolution is disabled")
	}

	// Check cache first
	r.mu.RLock()
	if cached, exists := r.cache[r.config.Domain]; exists && time.Now().Before(cached.ExpiresAt) {
		r.mu.RUnlock()
		r.logger.Debug("Using cached SRV records", zap.String("domain", r.config.Domain))
		return cached.Records, nil
	}
	r.mu.RUnlock()

	// Resolve SRV records
	records, err := r.resolveSRVWithRetry(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	r.mu.Lock()
	r.cache[r.config.Domain] = &SRVCacheEntry{
		Records:   records,
		ExpiresAt: time.Now().Add(r.config.CacheTTL),
	}
	r.mu.Unlock()

	return records, nil
}

// resolveSRVWithRetry resolves SRV records with retry logic
func (r *DNSResolver) resolveSRVWithRetry(ctx context.Context) ([]*net.SRV, error) {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			r.logger.Debug("Retrying DNS SRV resolution",
				zap.Int("attempt", attempt),
				zap.String("domain", r.config.Domain),
			)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(r.config.RetryDelay):
			}
		}

		records, err := r.resolveSRVOnce(ctx)
		if err == nil {
			return records, nil
		}

		lastErr = err
		r.logger.Warn("DNS SRV resolution failed",
			zap.Error(err),
			zap.Int("attempt", attempt+1),
			zap.String("domain", r.config.Domain),
		)
	}

	return nil, fmt.Errorf("DNS SRV resolution failed after %d retries: %w", r.config.MaxRetries+1, lastErr)
}

// resolveSRVOnce performs a single SRV resolution attempt
func (r *DNSResolver) resolveSRVOnce(ctx context.Context) ([]*net.SRV, error) {
	// Build service name based on transport
	service := r.buildServiceName()

	// Create DNS client
	client := &dns.Client{
		Timeout: r.config.Timeout,
	}

	// Create query
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(service), dns.TypeSRV)
	msg.RecursionDesired = true

	// Send query
	resp, _, err := client.ExchangeContext(ctx, msg, "8.8.8.8:53")
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %w", err)
	}

	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS query returned error: %s", dns.RcodeToString[resp.Rcode])
	}

	// Parse SRV records
	var records []*net.SRV
	for _, rr := range resp.Answer {
		if srv, ok := rr.(*dns.SRV); ok {
			records = append(records, &net.SRV{
				Target:   strings.TrimSuffix(srv.Target, "."),
				Port:     srv.Port,
				Priority: srv.Priority,
				Weight:   srv.Weight,
			})
		}
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no SRV records found for %s", service)
	}

	// Sort by priority and weight (RFC 2782)
	sort.Slice(records, func(i, j int) bool {
		if records[i].Priority != records[j].Priority {
			return records[i].Priority < records[j].Priority
		}
		return records[i].Weight > records[j].Weight
	})

	r.logger.Info("Resolved SRV records",
		zap.String("domain", r.config.Domain),
		zap.Int("count", len(records)),
	)

	return records, nil
}

// buildServiceName builds the SRV service name based on transport
func (r *DNSResolver) buildServiceName() string {
	var service string
	switch r.config.Transport {
	case TransportUDP:
		service = "_sip._udp"
	case TransportTCP:
		service = "_sip._tcp"
	case TransportTLS:
		service = "_sips._tcp"
	default:
		service = "_sip._udp"
	}
	return fmt.Sprintf("%s.%s", service, r.config.Domain)
}

// GetTrunkEndpoints returns trunk endpoints from SRV records
func (r *DNSResolver) GetTrunkEndpoints(ctx context.Context) ([]TrunkEndpoint, error) {
	records, err := r.ResolveSRV(ctx)
	if err != nil {
		return nil, err
	}

	var endpoints []TrunkEndpoint
	for _, record := range records {
		endpoints = append(endpoints, TrunkEndpoint{
			Host:     record.Target,
			Port:     int(record.Port),
			Priority: int(record.Priority),
			Weight:   int(record.Weight),
			Transport: r.config.Transport,
		})
	}

	return endpoints, nil
}

// TrunkEndpoint represents a SIP trunk endpoint from SRV records
type TrunkEndpoint struct {
	Host     string
	Port     int
	Priority int
	Weight   int
	Transport SIPTransport
}

// FailoverManager manages SIP trunk failover
type FailoverManager struct {
	config       *FailoverConfig
	logger       *zap.Logger
	trunks       []*SIPTrunk
	dnsResolver  *DNSResolver
	currentIndex int
	sipClient    *sipgo.Client
	mu           sync.RWMutex
}

// FailoverConfig holds failover configuration
type FailoverConfig struct {
	Enabled          bool
	DNSSRVEnabled    bool
	DNSSRVDomain     string
	DNSSRVTransport  SIPTransport
	HealthCheckInterval time.Duration
	FailureThreshold int
	RecoveryThreshold int
}

// NewFailoverManager creates a new failover manager
func NewFailoverManager(config *FailoverConfig, logger *zap.Logger) *FailoverManager {
	fm := &FailoverManager{
		config:  config,
		logger:  logger,
		trunks:  make([]*SIPTrunk, 0),
	}

	if config.DNSSRVEnabled {
		dnsConfig := &DNSResolverConfig{
			Enabled:    true,
			Domain:     config.DNSSRVDomain,
			Transport:  config.DNSSRVTransport,
			Timeout:    5 * time.Second,
			CacheTTL:   5 * time.Minute,
			MaxRetries: 3,
			RetryDelay: 1 * time.Second,
		}
		fm.dnsResolver = NewDNSResolver(dnsConfig, logger)
	}

	// Initialize SIP client for health checks
	ua, err := sipgo.NewUA()
	if err == nil {
		client, err := sipgo.NewClient(ua)
		if err == nil {
			fm.sipClient = client
		} else {
			logger.Error("Failed to create sipgo Client for health checks", zap.Error(err))
		}
	} else {
		logger.Error("Failed to create sipgo UA for health checks", zap.Error(err))
	}

	return fm
}

// AddTrunk adds a static trunk
func (fm *FailoverManager) AddTrunk(trunk *SIPTrunk) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.trunks = append(fm.trunks, trunk)
}

// GetActiveTrunk returns the currently active trunk
func (fm *FailoverManager) GetActiveTrunk() *SIPTrunk {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if len(fm.trunks) == 0 {
		return nil
	}

	// Find first healthy trunk starting from current index
	for i := 0; i < len(fm.trunks); i++ {
		idx := (fm.currentIndex + i) % len(fm.trunks)
		if fm.trunks[idx].Status == TrunkStatusHealthy {
			return fm.trunks[idx]
		}
	}

	// If no healthy trunk, return the current one (degraded)
	return fm.trunks[fm.currentIndex]
}

// Failover switches to the next available trunk
func (fm *FailoverManager) Failover() (*SIPTrunk, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if len(fm.trunks) <= 1 {
		return nil, fmt.Errorf("no backup trunks available")
	}

	oldIndex := fm.currentIndex
	for i := 1; i < len(fm.trunks); i++ {
		fm.currentIndex = (fm.currentIndex + 1) % len(fm.trunks)
		if fm.trunks[fm.currentIndex].Status == TrunkStatusHealthy {
			fm.logger.Info("Failover to backup trunk",
				zap.String("from", fm.trunks[oldIndex].Config.Name),
				zap.String("to", fm.trunks[fm.currentIndex].Config.Name),
			)
			return fm.trunks[fm.currentIndex], nil
		}
	}

	// No healthy backup, stay on current
	fm.currentIndex = oldIndex
	return nil, fmt.Errorf("no healthy backup trunks available")
}

// RefreshFromDNS refreshes trunk list from DNS SRV records
func (fm *FailoverManager) RefreshFromDNS(ctx context.Context) error {
	if fm.dnsResolver == nil {
		return fmt.Errorf("DNS resolver not configured")
	}

	endpoints, err := fm.dnsResolver.GetTrunkEndpoints(ctx)
	if err != nil {
		return fmt.Errorf("failed to get trunk endpoints from DNS: %w", err)
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Convert endpoints to trunks
	newTrunks := make([]*SIPTrunk, 0, len(endpoints))
	for i, ep := range endpoints {
		trunk := &SIPTrunk{
			Config: SIPTrunkConfig{
				Name:      fmt.Sprintf("dns-trunk-%d", i),
				Host:      ep.Host,
				Port:      ep.Port,
				Transport: ep.Transport,
				Priority:  ep.Priority,
				Weight:    ep.Weight,
			},
			Status: TrunkStatusUnknown,
		}
		newTrunks = append(newTrunks, trunk)
	}

	// Preserve current index if possible
	if fm.currentIndex >= len(newTrunks) {
		fm.currentIndex = 0
	}

	fm.trunks = newTrunks
	fm.logger.Info("Refreshed trunks from DNS SRV", zap.Int("count", len(newTrunks)))

	return nil
}

// StartHealthChecks starts periodic health checks
func (fm *FailoverManager) StartHealthChecks(ctx context.Context) {
	if !fm.config.Enabled {
		return
	}

	ticker := time.NewTicker(fm.config.HealthCheckInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				fm.performHealthChecks(ctx)
			}
		}
	}()
}

// performHealthChecks checks health of all trunks
func (fm *FailoverManager) performHealthChecks(ctx context.Context) {
	fm.mu.RLock()
	trunks := make([]*SIPTrunk, len(fm.trunks))
	copy(trunks, fm.trunks)
	fm.mu.RUnlock()

	for _, trunk := range trunks {
		go fm.checkTrunkHealth(ctx, trunk)
	}
}

// checkTrunkHealth checks health of a single trunk
func (fm *FailoverManager) checkTrunkHealth(ctx context.Context, trunk *SIPTrunk) {
	// Send OPTIONS request
	healthy := fm.sendHealthCheck(ctx, trunk)

	fm.mu.Lock()
	defer fm.mu.Unlock()

	if healthy {
		trunk.ConsecutiveFailures = 0
		if trunk.Status != TrunkStatusHealthy {
			trunk.Status = TrunkStatusHealthy
			fm.logger.Info("Trunk recovered", zap.String("trunk", trunk.Config.Name))
		}
	} else {
		trunk.ConsecutiveFailures++
		if trunk.ConsecutiveFailures >= fm.config.FailureThreshold {
			if trunk.Status != TrunkStatusUnhealthy {
				trunk.Status = TrunkStatusUnhealthy
				fm.logger.Warn("Trunk marked unhealthy",
					zap.String("trunk", trunk.Config.Name),
					zap.Int("failures", trunk.ConsecutiveFailures),
				)

				// Trigger failover if this is the active trunk
				if fm.currentIndex < len(fm.trunks) && fm.trunks[fm.currentIndex] == trunk {
					fm.failoverLocked()
				}
			}
		} else if trunk.ConsecutiveFailures >= fm.config.FailureThreshold/2 {
			if trunk.Status == TrunkStatusHealthy {
				trunk.Status = TrunkStatusDegraded
				fm.logger.Warn("Trunk degraded", zap.String("trunk", trunk.Config.Name))
			}
		}
	}
	trunk.LastHealthCheck = time.Now()
}

// failoverLocked performs failover (must hold lock)
func (fm *FailoverManager) failoverLocked() {
	if len(fm.trunks) <= 1 {
		return
	}

	oldIndex := fm.currentIndex
	for i := 1; i < len(fm.trunks); i++ {
		fm.currentIndex = (fm.currentIndex + 1) % len(fm.trunks)
		if fm.trunks[fm.currentIndex].Status == TrunkStatusHealthy {
			fm.logger.Info("Automatic failover to backup trunk",
				zap.String("from", fm.trunks[oldIndex].Config.Name),
				zap.String("to", fm.trunks[fm.currentIndex].Config.Name),
			)
			return
		}
	}

	// No healthy backup
	fm.currentIndex = oldIndex
	fm.logger.Error("No healthy backup trunks available for failover")
}

// sendHealthCheck sends a SIP OPTIONS request to check trunk health
func (fm *FailoverManager) sendHealthCheck(ctx context.Context, trunk *SIPTrunk) bool {
	if fm.sipClient == nil {
		fm.logger.Warn("SIP client not initialized, skipping real OPTIONS check")
		return false
	}

	// Determine scheme based on transport
	scheme := "sip"
	if trunk.Config.Transport == TransportTLS {
		scheme = "sips"
	}

	// Construct recipient URI (e.g. sip:trunk.example.com:5060)
	target := fmt.Sprintf("%s:%s:%d", scheme, trunk.Config.Host, trunk.Config.Port)
	recipient := &sip.Uri{}
	err := sip.ParseUri(target, recipient)
	if err != nil {
		fm.logger.Error("Failed to parse trunk URI for OPTIONS check", zap.Error(err))
		return false
	}

	// Create OPTIONS request
	req := sip.NewRequest(sip.OPTIONS, *recipient)
	
	// Ensure we set necessary headers
	req.AppendHeader(sip.NewHeader("Max-Forwards", "70"))
	req.AppendHeader(sip.NewHeader("User-Agent", "ai-telephony-platform/1.0"))
	
	// Create context with timeout for the transaction
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Send request and wait for response
	tx, err := fm.sipClient.TransactionRequest(checkCtx, req)
	if err != nil {
		fm.logger.Debug("Trunk OPTIONS check failed", 
			zap.String("trunk", trunk.Config.Name), 
			zap.Error(err))
		return false
	}
	defer tx.Terminate()

	select {
	case <-checkCtx.Done():
		fm.logger.Debug("Trunk OPTIONS check timed out", zap.String("trunk", trunk.Config.Name))
		return false
	case res := <-tx.Responses():
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			// 200 OK means trunk is healthy
			return true
		}
		// Some gateways might return 4xx (e.g., 405 Method Not Allowed) but it still implies they are alive
		// However, typically 200 OK is expected for OPTIONS keepalives.
		if res.StatusCode == 405 || res.StatusCode == 404 {
			return true 
		}
		fm.logger.Debug("Trunk returned non-success code to OPTIONS", 
			zap.String("trunk", trunk.Config.Name), 
			zap.Int("code", int(res.StatusCode)))
		return false
	}
}