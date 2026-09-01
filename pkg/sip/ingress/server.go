package ingress

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/emiago/sipgo/sip"
	"go.uber.org/zap"
)

// SIPServerImpl implements the SIPServer interface
type SIPServerImpl struct {
	config       *SIPServerConfig
	handler      SIPHandler
	parser       *SIPParser
	trunks       []*SIPTrunk
	workerPool   *WorkerPool
	udpConn      *net.UDPConn
	tcpListener  *net.TCPListener
	tlsListener  *net.Listener
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	logger       *zap.Logger
	mu           sync.RWMutex
}

// NewSIPServer creates a new SIP server
func NewSIPServer(config *SIPServerConfig, logger *zap.Logger) (*SIPServerImpl, error) {
	ctx, cancel := context.WithCancel(context.Background())

	server := &SIPServerImpl{
		config:     config,
		parser:     NewSIPParser(),
		trunks:     make([]*SIPTrunk, 0, len(config.Trunks)),
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
	}

	// Initialize trunks
	for _, trunkConfig := range config.Trunks {
		trunk := &SIPTrunk{
			Config:   trunkConfig,
			Status:   TrunkStatusUnknown,
		}
		server.trunks = append(server.trunks, trunk)
	}

	// Initialize worker pool
	server.workerPool = NewWorkerPool(config.WorkerCount, config.MaxRequestsPerWorker, logger)

	return server, nil
}

// SetHandler sets the SIP message handler
func (s *SIPServerImpl) SetHandler(handler SIPHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}

// Start starts the SIP server
func (s *SIPServerImpl) Start() error {
	s.logger.Info("Starting SIP server")

	// Start UDP listener
	if s.config.UDPListen != "" {
		if err := s.startUDPListener(); err != nil {
			return fmt.Errorf("failed to start UDP listener: %w", err)
		}
	}

	// Start TCP listener
	if s.config.TCPListen != "" {
		if err := s.startTCPListener(); err != nil {
			return fmt.Errorf("failed to start TCP listener: %w", err)
		}
	}

	// Start TLS listener
	if s.config.TLSListen != "" {
		if err := s.startTLSListener(); err != nil {
			return fmt.Errorf("failed to start TLS listener: %w", err)
		}
	}

	// Start worker pool
	s.workerPool.Start()

	// Start keepalive monitoring
	if s.config.KeepaliveEnabled {
		s.wg.Add(1)
		go s.keepaliveMonitor()
	}

	// Start trunk health monitoring
	s.wg.Add(1)
	go s.trunkHealthMonitor()

	s.logger.Info("SIP server started successfully")
	return nil
}

// Stop stops the SIP server
func (s *SIPServerImpl) Stop() error {
	s.logger.Info("Stopping SIP server")

	s.cancel()
	s.workerPool.Stop()

	if s.udpConn != nil {
		s.udpConn.Close()
	}

	if s.tcpListener != nil {
		s.tcpListener.Close()
	}

	if s.tlsListener != nil {
		(*s.tlsListener).Close()
	}

	s.wg.Wait()

	s.logger.Info("SIP server stopped")
	return nil
}

// startUDPListener starts the UDP listener
func (s *SIPServerImpl) startUDPListener() error {
	addr, err := net.ResolveUDPAddr("udp", s.config.UDPListen)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP: %w", err)
	}

	s.udpConn = conn
	s.wg.Add(1)
	go s.udpListenerLoop(conn)

	s.logger.Info("UDP listener started", zap.String("address", s.config.UDPListen))
	return nil
}

// startTCPListener starts the TCP listener
func (s *SIPServerImpl) startTCPListener() error {
	addr, err := net.ResolveTCPAddr("tcp", s.config.TCPListen)
	if err != nil {
		return fmt.Errorf("failed to resolve TCP address: %w", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on TCP: %w", err)
	}

	s.tcpListener = listener
	s.wg.Add(1)
	go s.tcpListenerLoop(listener)

	s.logger.Info("TCP listener started", zap.String("address", s.config.TCPListen))
	return nil
}

// startTLSListener starts the TLS listener
func (s *SIPServerImpl) startTLSListener() error {
	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair(s.config.TLSCertFile, s.config.TLSKeyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}

	if s.config.VerifyClient {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	addr, err := net.ResolveTCPAddr("tcp", s.config.TLSListen)
	if err != nil {
		return fmt.Errorf("failed to resolve TLS address: %w", err)
	}

	listener, err := tls.Listen("tcp", addr.String(), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to listen on TLS: %w", err)
	}

	s.tlsListener = &listener
	s.wg.Add(1)
	go s.tlsListenerLoop(listener)

	s.logger.Info("TLS listener started", zap.String("address", s.config.TLSListen))
	return nil
}

// udpListenerLoop handles incoming UDP packets
func (s *SIPServerImpl) udpListenerLoop(conn *net.UDPConn) {
	defer s.wg.Done()

	buffer := make([]byte, 65535)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if s.ctx.Err() != nil {
					return
				}
				s.logger.Error("UDP read error", zap.Error(err))
				continue
			}

			localAddr := conn.LocalAddr().(*net.UDPAddr)

			// Submit to worker pool
			s.workerPool.Submit(func() {
				s.handleUDPMessage(buffer[:n], remoteAddr, localAddr)
			})
		}
	}
}

// handleUDPMessage handles an incoming UDP message
func (s *SIPServerImpl) handleUDPMessage(data []byte, remoteAddr, localAddr *net.UDPAddr) {
	// Parse the request
	invite, err := s.parser.ParseRequest(data, TransportUDP, localAddr, remoteAddr)
	if err != nil {
		s.logger.Error("Failed to parse UDP request", zap.Error(err), zap.String("remote", remoteAddr.String()))
		return
	}

	// Handle the request
	s.handleRequest(invite)
}

// tcpListenerLoop handles incoming TCP connections
func (s *SIPServerImpl) tcpListenerLoop(listener *net.TCPListener) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			listener.SetDeadline(time.Now().Add(1 * time.Second))
			conn, err := listener.AcceptTCP()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if s.ctx.Err() != nil {
					return
				}
				s.logger.Error("TCP accept error", zap.Error(err))
				continue
			}

			s.wg.Add(1)
			go s.handleTCPConnection(conn)
		}
	}
}

// handleTCPConnection handles a TCP connection
func (s *SIPServerImpl) handleTCPConnection(conn *net.TCPConn) {
	defer s.wg.Done()
	defer conn.Close()

	buffer := make([]byte, 65535)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			n, err := conn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				return
			}

			localAddr := conn.LocalAddr().(*net.TCPAddr)
			remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

			// Parse the request
			invite, err := s.parser.ParseRequest(buffer[:n], TransportTCP, 
				&net.UDPAddr{IP: localAddr.IP, Port: localAddr.Port},
				&net.UDPAddr{IP: remoteAddr.IP, Port: remoteAddr.Port})
			if err != nil {
				s.logger.Error("Failed to parse TCP request", zap.Error(err))
				continue
			}

			// Handle the request
			s.handleRequest(invite)
		}
	}
}

// tlsListenerLoop handles incoming TLS connections
func (s *SIPServerImpl) tlsListenerLoop(listener net.Listener) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			listener.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))
			conn, err := listener.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if s.ctx.Err() != nil {
					return
				}
				s.logger.Error("TLS accept error", zap.Error(err))
				continue
			}

			s.wg.Add(1)
			go s.handleTLSConnection(conn)
		}
	}
}

// handleTLSConnection handles a TLS connection
func (s *SIPServerImpl) handleTLSConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	tlsConn := conn.(*tls.Conn)
	buffer := make([]byte, 65535)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			n, err := conn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				return
			}

			localAddr := conn.LocalAddr().(*net.TCPAddr)
			remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

			// Get TLS connection state
			tlsState := tlsConn.ConnectionState()
			tlsConnectionState := &TLSConnectionState{
				Version:          tlsState.Version,
				CipherSuite:      tlsState.CipherSuite,
				ServerName:       tlsState.ServerName,
				PeerCertificates: make([][]byte, len(tlsState.PeerCertificates)),
				VerifiedChains:   make([][][]byte, len(tlsState.VerifiedChains)),
			}
			for i, cert := range tlsState.PeerCertificates {
				tlsConnectionState.PeerCertificates[i] = cert.Raw
			}
			for i, chain := range tlsState.VerifiedChains {
				tlsConnectionState.VerifiedChains[i] = make([][]byte, len(chain))
				for j, cert := range chain {
					tlsConnectionState.VerifiedChains[i][j] = cert.Raw
				}
			}

			// Parse the request
			invite, err := s.parser.ParseRequest(buffer[:n], TransportTLS,
				&net.UDPAddr{IP: localAddr.IP, Port: localAddr.Port},
				&net.UDPAddr{IP: remoteAddr.IP, Port: remoteAddr.Port})
			if err != nil {
				s.logger.Error("Failed to parse TLS request", zap.Error(err))
				continue
			}

			invite.TLSConnectionState = tlsConnectionState

			// Handle the request
			s.handleRequest(invite)
		}
	}
}

// handleRequest routes the request to the appropriate handler
func (s *SIPServerImpl) handleRequest(invite *SIPInvite) {
	s.mu.RLock()
	handler := s.handler
	s.mu.RUnlock()

	if handler == nil {
		s.logger.Warn("No handler set, sending 500")
		s.sendErrorResponse(invite, 500, "Internal Server Error")
		return
	}

	var response *SIPResponse
	var err error

	switch invite.GetMethod() {
	case MethodINVITE:
		if s.parser.IsReInvite(invite.Request, invite.CallID, invite.CSeq) {
			response, err = handler.HandleReInvite(invite)
		} else {
			response, err = handler.HandleInvite(invite)
		}
	case MethodACK:
		err = handler.HandleACK(invite)
	case MethodBYE:
		response, err = handler.HandleBYE(invite)
	case MethodCANCEL:
		response, err = handler.HandleCancel(invite)
	case MethodOPTIONS:
		response, err = handler.HandleOptions(invite)
	case MethodREGISTER:
		response, err = handler.HandleRegister(invite)
	case MethodINFO:
		response, err = handler.HandleInfo(invite)
	case MethodREFER:
		response, err = handler.HandleRefer(invite)
	default:
		s.logger.Warn("Unsupported method", zap.String("method", string(invite.GetMethod())))
		s.sendErrorResponse(invite, 501, "Not Implemented")
		return
	}

	if err != nil {
		s.logger.Error("Handler error", zap.Error(err), zap.String("method", string(invite.GetMethod())))
		s.sendErrorResponse(invite, 500, "Internal Server Error")
		return
	}

	if response != nil {
		if err := s.SendResponse(response); err != nil {
			s.logger.Error("Failed to send response", zap.Error(err))
		}
	}
}

// sendErrorResponse sends an error response
func (s *SIPServerImpl) sendErrorResponse(invite *SIPInvite, statusCode int, reason string) {
	response, err := s.parser.BuildResponse(invite.Request, statusCode, reason, nil, "")
	if err != nil {
		s.logger.Error("Failed to build error response", zap.Error(err))
		return
	}

	sipResponse := &SIPResponse{
		Response:    response,
		Transport:   invite.Transport,
		LocalAddr:   invite.LocalAddr,
		RemoteAddr:  invite.RemoteAddr,
		StatusCode:  statusCode,
		ReasonPhrase: reason,
		SentAt:      time.Now(),
	}

	if err := s.SendResponse(sipResponse); err != nil {
		s.logger.Error("Failed to send error response", zap.Error(err))
	}
}

// SendResponse sends a SIP response
func (s *SIPServerImpl) SendResponse(response *SIPResponse) error {
	data := []byte(response.Response.String())

	var err error
	switch response.Transport {
	case TransportUDP:
		if s.udpConn != nil {
			_, err = s.udpConn.WriteToUDP(data, response.RemoteAddr)
		}
	case TransportTCP, TransportTLS:
		// For TCP/TLS, we need to maintain connection state
		// This is a simplified implementation
		s.logger.Debug("TCP/TLS response sending not fully implemented")
	default:
		return fmt.Errorf("unsupported transport: %s", response.Transport)
	}

	return err
}

// SendRequest sends a SIP request
func (s *SIPServerImpl) SendRequest(request *SIPInvite) error {
	data := []byte(request.Request.String())

	var err error
	switch request.Transport {
	case TransportUDP:
		if s.udpConn != nil {
			_, err = s.udpConn.WriteToUDP(data, request.RemoteAddr)
		}
	case TransportTCP, TransportTLS:
		s.logger.Debug("TCP/TLS request sending not fully implemented")
	default:
		return fmt.Errorf("unsupported transport: %s", request.Transport)
	}

	return err
}

// GetTrunks returns the list of SIP trunks
func (s *SIPServerImpl) GetTrunks() []*SIPTrunk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trunks := make([]*SIPTrunk, len(s.trunks))
	copy(trunks, s.trunks)
	return trunks
}

// GetTrunkStatus returns the status of a specific trunk
func (s *SIPServerImpl) GetTrunkStatus(name string) (TrunkStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, trunk := range s.trunks {
		if trunk.Config.Name == name {
			return trunk.Status, nil
		}
	}

	return TrunkStatusUnknown, fmt.Errorf("trunk not found: %s", name)
}

// keepaliveMonitor monitors trunk health using SIP OPTIONS
func (s *SIPServerImpl) keepaliveMonitor() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.config.KeepaliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.sendKeepaliveToTrunks()
		}
	}
}

// sendKeepaliveToTrunks sends OPTIONS to all trunks
func (s *SIPServerImpl) sendKeepaliveToTrunks() {
	s.mu.RLock()
	trunks := make([]*SIPTrunk, len(s.trunks))
	copy(trunks, s.trunks)
	s.mu.RUnlock()

	for _, trunk := range trunks {
		go s.sendKeepaliveToTrunk(trunk)
	}
}

// sendKeepaliveToTrunk sends OPTIONS to a specific trunk
func (s *SIPServerImpl) sendKeepaliveToTrunk(trunk *SIPTrunk) {
	// Build OPTIONS request
	ruri := sip.Uri{
		User: "",
		Host: trunk.Config.Host,
		Port: trunk.Config.Port,
	}

	request := sip.NewRequest(sip.RequestMethod(MethodOPTIONS), ruri)
	request.AppendHeader(sip.NewHeader("User-Agent", "AI-Telephony-Platform/1.0"))
	request.AppendHeader(sip.NewHeader("Max-Forwards", "70"))

	// Send request
	invite := &SIPInvite{
		Request:    request,
		Transport:  trunk.Config.Transport,
		RemoteAddr: &net.UDPAddr{IP: net.ParseIP(trunk.Config.Host), Port: trunk.Config.Port},
	}

	if err := s.SendRequest(invite); err != nil {
		s.logger.Error("Failed to send keepalive", zap.Error(err), zap.String("trunk", trunk.Config.Name))
		s.updateTrunkStatus(trunk, TrunkStatusUnhealthy)
		return
	}

	// Wait for response with timeout
	// This is simplified - in production, you'd use a transaction layer
	time.Sleep(s.config.KeepaliveTimeout)
	
	// For now, assume success if no error
	s.updateTrunkStatus(trunk, TrunkStatusHealthy)
}

// updateTrunkStatus updates the status of a trunk
func (s *SIPServerImpl) updateTrunkStatus(trunk *SIPTrunk, status TrunkStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	trunk.Status = status
	trunk.LastHealthCheck = time.Now()

	if status == TrunkStatusUnhealthy {
		trunk.ConsecutiveFailures++
	} else {
		trunk.ConsecutiveFailures = 0
	}
}

// trunkHealthMonitor monitors trunk health and handles failover
func (s *SIPServerImpl) trunkHealthMonitor() {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkTrunkHealth()
		}
	}
}

// checkTrunkHealth checks the health of all trunks
func (s *SIPServerImpl) checkTrunkHealth() {
	s.mu.RLock()
	trunks := make([]*SIPTrunk, len(s.trunks))
	copy(trunks, s.trunks)
	s.mu.RUnlock()

	for _, trunk := range trunks {
		if trunk.ConsecutiveFailures >= 3 {
			s.logger.Warn("Trunk unhealthy, marking for failover", zap.String("trunk", trunk.Config.Name))
			s.updateTrunkStatus(trunk, TrunkStatusUnhealthy)
		}
	}
}

// WorkerPool manages a pool of workers for handling SIP messages
type WorkerPool struct {
	workerCount       int
	maxRequestsPerWorker int
	jobQueue          chan func()
	wg                sync.WaitGroup
	ctx               context.Context
	cancel            context.CancelFunc
	logger            *zap.Logger
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount, maxRequestsPerWorker int, logger *zap.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workerCount:         workerCount,
		maxRequestsPerWorker: maxRequestsPerWorker,
		jobQueue:            make(chan func(), workerCount*10),
		ctx:                 ctx,
		cancel:              cancel,
		logger:              logger,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	wp.logger.Info("Worker pool started", zap.Int("workers", wp.workerCount))
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
	wp.logger.Info("Worker pool stopped")
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job func()) {
	select {
	case wp.jobQueue <- job:
	case <-wp.ctx.Done():
		wp.logger.Warn("Worker pool stopped, dropping job")
	}
}

// worker processes jobs from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	requestCount := 0

	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}

			job()
			requestCount++

			// Check if we need to recycle the worker
			if wp.maxRequestsPerWorker > 0 && requestCount >= wp.maxRequestsPerWorker {
				wp.logger.Debug("Worker recycled", zap.Int("worker", id), zap.Int("requests", requestCount))
				requestCount = 0
			}
		}
	}
}