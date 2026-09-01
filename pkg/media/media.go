package media

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/pion/rtp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// Common Codec constants
const (
	PayloadTypePCMU = 0
	PayloadTypePCMA = 8
	PayloadTypeG722 = 9
	PayloadTypeOpus = 111
)

var (
	metricJitter = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ai_telephony_media_jitter_ms",
		Help:    "Jitter in milliseconds",
		Buckets: []float64{0, 10, 30, 50, 100, 200},
	}, []string{"session_id"})

	metricPacketLoss = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ai_telephony_media_packet_loss_total",
		Help: "Total packet loss",
	}, []string{"session_id"})

	metricMOS = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ai_telephony_media_mos",
		Help: "Mean Opinion Score (1.0 - 4.5)",
	}, []string{"session_id"})
)

// CodecSpec represents details for a supported codec
type CodecSpec struct {
	Name        string
	PayloadType uint8
	ClockRate   uint32
	Channels    uint16
}

// MediaSession represents an active media session (audio stream)
type MediaSession struct {
	ID             string
	CallID         string
	Codec          CodecSpec
	LocalRTPAddr   *net.UDPAddr
	RemoteRTPAddr  *net.UDPAddr
	RTPConn        *net.UDPConn
	JitterBuffer   *JitterBuffer
	DTMFEngine     *DTMFEngine
	AMDEngine      *AMDEngine
	Metrics        *MediaQualityMetrics
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
	isActive       bool
}

// JitterBuffer buffers incoming RTP packets to smooth out arrival jitter
type JitterBuffer struct {
	minDelay    time.Duration
	maxDelay    time.Duration
	targetDelay time.Duration
	adaptive    bool
	buffer      map[uint16]*rtp.Packet
	highestSeq  uint16
	mu          sync.Mutex
}

// NewJitterBuffer creates a new JitterBuffer
func NewJitterBuffer(minDelayMs, maxDelayMs, targetDelayMs int, adaptive bool) *JitterBuffer {
	return &JitterBuffer{
		minDelay:    time.Duration(minDelayMs) * time.Millisecond,
		maxDelay:    time.Duration(maxDelayMs) * time.Millisecond,
		targetDelay: time.Duration(targetDelayMs) * time.Millisecond,
		adaptive:    adaptive,
		buffer:      make(map[uint16]*rtp.Packet),
	}
}

// Push adds an RTP packet into the jitter buffer
func (jb *JitterBuffer) Push(pkt *rtp.Packet) {
	jb.mu.Lock()
	defer jb.mu.Unlock()

	jb.buffer[pkt.SequenceNumber] = pkt
	if pkt.SequenceNumber > jb.highestSeq {
		jb.highestSeq = pkt.SequenceNumber
	}
}

// Pop retrieves and removes a packet by sequence number
func (jb *JitterBuffer) Pop(seq uint16) (*rtp.Packet, bool) {
	jb.mu.Lock()
	defer jb.mu.Unlock()

	pkt, exists := jb.buffer[seq]
	if exists {
		delete(jb.buffer, seq)
	}
	return pkt, exists
}

// DTMFEvent represents a detected DTMF digit
type DTMFEvent struct {
	Digit    rune
	Duration time.Duration
	Volume   uint8
	Time     time.Time
}

// DTMFEngine handles RFC 4733 and in-band DTMF detection
type DTMFEngine struct {
	eventsChan chan DTMFEvent
	enabled    bool
}

// NewDTMFEngine creates a new DTMF engine
func NewDTMFEngine(enabled bool) *DTMFEngine {
	return &DTMFEngine{
		eventsChan: make(chan DTMFEvent, 100),
		enabled:    enabled,
	}
}

// ProcessRTPPacket inspects an RTP packet for RFC 4733 telephone-event payload
func (de *DTMFEngine) ProcessRTPPacket(pkt *rtp.Packet, telephoneEventPT uint8) *DTMFEvent {
	if !de.enabled || pkt.PayloadType != telephoneEventPT || len(pkt.Payload) < 4 {
		return nil
	}

	eventCode := pkt.Payload[0]
	endBit := (pkt.Payload[1] & 0x80) != 0
	volume := pkt.Payload[1] & 0x3F
	duration := uint16(pkt.Payload[2])<<8 | uint16(pkt.Payload[3])

	digit := dtmfCodeToDigit(eventCode)
	if digit == 0 {
		return nil
	}

	event := DTMFEvent{
		Digit:    digit,
		Duration: time.Duration(duration) * time.Millisecond,
		Volume:   volume,
		Time:     time.Now(),
	}

	if endBit {
		select {
		case de.eventsChan <- event:
		default:
		}
	}

	return &event
}

// ProcessInbandAudio analyzes raw PCM audio for in-band DTMF frequencies (Goertzel algorithm)
func (de *DTMFEngine) ProcessInbandAudio(pcm []int16) {
	if !de.enabled {
		return
	}
	// Stub: In a real implementation, we'd run the Goertzel algorithm over the PCM chunk
	// to detect standard dual-tone frequencies (697/770/852/941 Hz + 1209/1336/1477/1633 Hz).
	// If a digit is continuously detected for >40ms, emit a DTMFEvent.
}

// ProcessSIPInfo handles out-of-band DTMF sent via SIP INFO requests
func (de *DTMFEngine) ProcessSIPInfo(digit rune, duration time.Duration) {
	if !de.enabled {
		return
	}
	
	event := DTMFEvent{
		Digit:    digit,
		Duration: duration,
		Volume:   10,
		Time:     time.Now(),
	}
	
	select {
	case de.eventsChan <- event:
	default:
	}
}

// Events returns the channel for receiving DTMF events
func (de *DTMFEngine) Events() <-chan DTMFEvent {
	return de.eventsChan
}

func dtmfCodeToDigit(code uint8) rune {
	switch code {
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	case 3:
		return '3'
	case 4:
		return '4'
	case 5:
		return '5'
	case 6:
		return '6'
	case 7:
		return '7'
	case 8:
		return '8'
	case 9:
		return '9'
	case 10:
		return '*'
	case 11:
		return '#'
	case 12:
		return 'A'
	case 13:
		return 'B'
	case 14:
		return 'C'
	case 15:
		return 'D'
	default:
		return 0
	}
}

// AMDEngine handles Answering Machine Detection
type AMDEngine struct {
	enabled          bool
	humanThreshold   float64
	voicemailThresh  float64
	silenceDuration  time.Duration
	resultChan       chan string
	speechEnergy     float64
	continuousSpeech time.Duration
}

// NewAMDEngine initializes the AMD engine
func NewAMDEngine(enabled bool, humanThresh, vmThresh float64, silence time.Duration) *AMDEngine {
	return &AMDEngine{
		enabled:         enabled,
		humanThreshold:  humanThresh,
		voicemailThresh: vmThresh,
		silenceDuration: silence,
		resultChan:      make(chan string, 1),
	}
}

// ProcessRTP analyzes RTP payload energy for speech detection
func (a *AMDEngine) ProcessRTP(pkt *rtp.Packet) {
	if !a.enabled || len(pkt.Payload) == 0 {
		return
	}
	
	// A real implementation decodes the payload (e.g., G.711) and calculates RMS energy.
	// We simulate this by checking payload variance or length as a stand-in.
	energy := float64(len(pkt.Payload)) / 160.0 // Very naive placeholder for audio energy
	
	if energy > 0.1 {
		a.speechEnergy += energy
		a.continuousSpeech += 20 * time.Millisecond // Assuming 20ms packets
	} else {
		// Silence detected, evaluate speech duration
		if a.continuousSpeech > a.silenceDuration {
			if a.speechEnergy > a.voicemailThresh {
				select {
				case a.resultChan <- "voicemail":
				default:
				}
			} else if a.speechEnergy > a.humanThreshold {
				select {
				case a.resultChan <- "human":
				default:
				}
			}
		}
		a.speechEnergy = 0
		a.continuousSpeech = 0
	}
}

// Result returns a channel for AMD results
func (a *AMDEngine) Result() <-chan string {
	return a.resultChan
}

// MediaQualityMetrics tracks real-time QoS (jitter, packet loss, MOS) and emits to Prometheus
type MediaQualityMetrics struct {
	SessionID       string
	PacketsReceived uint64
	PacketsLost     uint64
	JitterMs        float64
	MOS             float64
	LastSeq         uint16
	LastTransit     uint32
	mu              sync.Mutex
}

// Update calculates jitter and packet loss according to RFC 3550 and emits metrics
func (m *MediaQualityMetrics) Update(pkt *rtp.Packet, clockRate uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.PacketsReceived++

	if m.PacketsReceived > 1 {
		expectedGap := pkt.SequenceNumber - m.LastSeq
		if expectedGap > 1 && expectedGap < 100 {
			lost := uint64(expectedGap - 1)
			m.PacketsLost += lost
			metricPacketLoss.WithLabelValues(m.SessionID).Add(float64(lost))
		}

		transit := uint32(time.Now().UnixNano()/1e6) - pkt.Timestamp
		d := int32(transit - m.LastTransit)
		if d < 0 {
			d = -d
		}
		m.JitterMs += (float64(d) - m.JitterMs) / 16.0
		metricJitter.WithLabelValues(m.SessionID).Observe(m.JitterMs)
	}

	m.LastSeq = pkt.SequenceNumber
	m.LastTransit = uint32(time.Now().UnixNano() / 1e6)

	m.calculateMOS()
	metricMOS.WithLabelValues(m.SessionID).Set(m.MOS)
}

func (m *MediaQualityMetrics) calculateMOS() {
	if m.PacketsReceived == 0 {
		m.MOS = 4.5
		return
	}

	lossRate := float64(m.PacketsLost) / float64(m.PacketsReceived+m.PacketsLost) * 100.0
	rFactor := 93.2 - (lossRate * 2.5) - (m.JitterMs * 0.1)
	if rFactor < 0 {
		rFactor = 0
	}

	if rFactor > 100 {
		m.MOS = 4.5
	} else {
		m.MOS = 1.0 + (0.035 * rFactor) + (rFactor * (rFactor - 60.0) * (100.0 - rFactor) * 0.000007)
	}

	m.MOS = math.Min(4.5, math.Max(1.0, m.MOS))
}

// MediaProxyServer manages media sessions and RTP allocation
type MediaProxyServer struct {
	config   *config.MediaConfig
	logger   *zap.Logger
	sessions map[string]*MediaSession
	portPool chan int
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewMediaProxyServer creates a new MediaProxyServer
func NewMediaProxyServer(cfg *config.MediaConfig, logger *zap.Logger) (*MediaProxyServer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	startPort := cfg.Proxy.RTPPortRange.Start
	endPort := cfg.Proxy.RTPPortRange.End
	if startPort == 0 || endPort <= startPort {
		startPort = 10000
		endPort = 20000
	}

	portCount := endPort - startPort + 1
	portPool := make(chan int, portCount)
	for p := startPort; p <= endPort; p++ {
		portPool <- p
	}

	return &MediaProxyServer{
		config:   cfg,
		logger:   logger,
		sessions: make(map[string]*MediaSession),
		portPool: portPool,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// CreateSession allocates ports and starts a media proxy session
func (s *MediaProxyServer) CreateSession(sessionID, callID string, codec CodecSpec) (*MediaSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case port := <-s.portPool:
		localAddr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port}
		conn, err := net.ListenUDP("udp", localAddr)
		if err != nil {
			s.portPool <- port
			return nil, fmt.Errorf("failed to listen on RTP port %d: %w", port, err)
		}

		ctx, cancel := context.WithCancel(s.ctx)
		session := &MediaSession{
			ID:           sessionID,
			CallID:       callID,
			Codec:        codec,
			LocalRTPAddr: localAddr,
			RTPConn:      conn,
			JitterBuffer: NewJitterBuffer(
				s.config.Proxy.JitterBuffer.MinDelay,
				s.config.Proxy.JitterBuffer.MaxDelay,
				s.config.Proxy.JitterBuffer.TargetDelay,
				s.config.Proxy.JitterBuffer.Adaptive,
			),
			DTMFEngine: NewDTMFEngine(s.config.Proxy.DTMF.RFC4733Enabled),
			AMDEngine: NewAMDEngine(
				true, // Replace with s.config.Campaign.PredictiveDialer.AMD.Enabled when cross-package config is passed
				0.5,
				0.7,
				1500*time.Millisecond,
			),
			Metrics:    &MediaQualityMetrics{SessionID: sessionID, MOS: 4.5},
			ctx:        ctx,
			cancel:     cancel,
			isActive:   true,
		}

		s.sessions[sessionID] = session
		s.logger.Info("Media session created",
			zap.String("session_id", sessionID),
			zap.String("call_id", callID),
			zap.Int("port", port),
		)
		return session, nil
	default:
		return nil, errors.New("no RTP ports available")
	}
}

// CloseSession terminates a session and returns its allocated port to the pool
func (s *MediaProxyServer) CloseSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	session.mu.Lock()
	if session.isActive {
		session.isActive = false
		session.cancel()
		if session.RTPConn != nil {
			session.RTPConn.Close()
			s.portPool <- session.LocalRTPAddr.Port
		}
	}
	session.mu.Unlock()

	delete(s.sessions, sessionID)
	s.logger.Info("Media session closed", zap.String("session_id", sessionID))
	return nil
}

// Shutdown stops all active media sessions
func (s *MediaProxyServer) Shutdown() {
	s.cancel()
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, session := range s.sessions {
		session.cancel()
		if session.RTPConn != nil {
			session.RTPConn.Close()
		}
		delete(s.sessions, id)
	}
}
