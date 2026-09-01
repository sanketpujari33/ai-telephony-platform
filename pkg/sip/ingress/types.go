package ingress

import (
	"net"
	"time"

	"github.com/emiago/sipgo/sip"
)

// SIPMethod represents SIP methods
type SIPMethod string

const (
	MethodINVITE     SIPMethod = "INVITE"
	MethodACK        SIPMethod = "ACK"
	MethodBYE        SIPMethod = "BYE"
	MethodCANCEL     SIPMethod = "CANCEL"
	MethodREGISTER   SIPMethod = "REGISTER"
	MethodOPTIONS    SIPMethod = "OPTIONS"
	MethodINFO       SIPMethod = "INFO"
	MethodPRACK      SIPMethod = "PRACK"
	MethodUPDATE     SIPMethod = "UPDATE"
	MethodREFER      SIPMethod = "REFER"
	MethodNOTIFY     SIPMethod = "NOTIFY"
	MethodSUBSCRIBE  SIPMethod = "SUBSCRIBE"
	MethodMESSAGE    SIPMethod = "MESSAGE"
	MethodPUBLISH    SIPMethod = "PUBLISH"
)

// SIPTransport represents SIP transport protocols
type SIPTransport string

const (
	TransportUDP SIPTransport = "udp"
	TransportTCP SIPTransport = "tcp"
	TransportTLS SIPTransport = "tls"
	TransportWS  SIPTransport = "ws"
	TransportWSS SIPTransport = "wss"
)

// SIPInvite represents a parsed INVITE request
type SIPInvite struct {
	// Base SIP request
	Request *sip.Request

	// Parsed headers
	CallID       string
	FromURI      *sip.Uri
	ToURI        *sip.Uri
	ContactURI   *sip.Uri
	ViaHeader    *sip.ViaHeader
	CSeq         uint32
	MaxForwards  uint32
	ContentType  string
	ContentLength int

	// SDP body
	SDPBody string

	// Transport info
	Transport SIPTransport
	LocalAddr *net.UDPAddr
	RemoteAddr *net.UDPAddr

	// TLS info
	TLSConnectionState *TLSConnectionState

	// Timestamp
	ReceivedAt time.Time
}

// SIPResponse represents a SIP response
type SIPResponse struct {
	// Base SIP response
	Response *sip.Response

	// Parsed headers
	CallID      string
	CSeq        uint32
	ViaHeader   *sip.ViaHeader
	ContactURI  *sip.Uri
	ContentType string
	ContentLength int

	// SDP body
	SDPBody string

	// Status
	StatusCode int
	ReasonPhrase string

	// Transport info
	Transport SIPTransport
	LocalAddr *net.UDPAddr
	RemoteAddr *net.UDPAddr

	// Timestamp
	SentAt time.Time
}

// TLSConnectionState holds TLS connection information
type TLSConnectionState struct {
	Version          uint16
	CipherSuite      uint16
	ServerName       string
	PeerCertificates [][]byte
	VerifiedChains   [][][]byte
}

// SIPMessage represents a generic SIP message (request or response)
type SIPMessage interface {
	GetCallID() string
	GetCSeq() uint32
	GetMethod() SIPMethod
	GetTransport() SIPTransport
	GetRemoteAddr() *net.UDPAddr
	GetLocalAddr() *net.UDPAddr
	GetTimestamp() time.Time
}

// GetCallID returns the Call-ID header
func (i *SIPInvite) GetCallID() string {
	return i.CallID
}

// GetCSeq returns the CSeq number
func (i *SIPInvite) GetCSeq() uint32 {
	return i.CSeq
}

// GetMethod returns the SIP method
func (i *SIPInvite) GetMethod() SIPMethod {
	return MethodINVITE
}

// GetTransport returns the transport protocol
func (i *SIPInvite) GetTransport() SIPTransport {
	return i.Transport
}

// GetRemoteAddr returns the remote address
func (i *SIPInvite) GetRemoteAddr() *net.UDPAddr {
	return i.RemoteAddr
}

// GetLocalAddr returns the local address
func (i *SIPInvite) GetLocalAddr() *net.UDPAddr {
	return i.LocalAddr
}

// GetTimestamp returns the timestamp
func (i *SIPInvite) GetTimestamp() time.Time {
	return i.ReceivedAt
}

// GetCallID returns the Call-ID header
func (r *SIPResponse) GetCallID() string {
	return r.CallID
}

// GetCSeq returns the CSeq number
func (r *SIPResponse) GetCSeq() uint32 {
	return r.CSeq
}

func (r *SIPResponse) GetMethod() SIPMethod {
	if r.Response != nil {
		if cseq := r.Response.CSeq(); cseq != nil {
			return SIPMethod(cseq.MethodName)
		}
	}
	return ""
}

// GetTransport returns the transport protocol
func (r *SIPResponse) GetTransport() SIPTransport {
	return r.Transport
}

// GetRemoteAddr returns the remote address
func (r *SIPResponse) GetRemoteAddr() *net.UDPAddr {
	return r.RemoteAddr
}

// GetLocalAddr returns the local address
func (r *SIPResponse) GetLocalAddr() *net.UDPAddr {
	return r.LocalAddr
}

// GetTimestamp returns the timestamp
func (r *SIPResponse) GetTimestamp() time.Time {
	return r.SentAt
}

// SIPHandler defines the interface for handling SIP messages
type SIPHandler interface {
	// HandleInvite handles an incoming INVITE request
	HandleInvite(invite *SIPInvite) (*SIPResponse, error)

	// HandleACK handles an incoming ACK request
	HandleACK(ack *SIPInvite) error

	// HandleBYE handles an incoming BYE request
	HandleBYE(bye *SIPInvite) (*SIPResponse, error)

	// HandleReInvite handles a re-INVITE request
	HandleReInvite(invite *SIPInvite) (*SIPResponse, error)

	// HandleOptions handles an incoming OPTIONS request
	HandleOptions(options *SIPInvite) (*SIPResponse, error)

	// HandleCancel handles an incoming CANCEL request
	HandleCancel(cancel *SIPInvite) (*SIPResponse, error)

	// HandleRegister handles an incoming REGISTER request
	HandleRegister(register *SIPInvite) (*SIPResponse, error)

	// HandleInfo handles an incoming INFO request
	HandleInfo(info *SIPInvite) (*SIPResponse, error)

	// HandleRefer handles an incoming REFER request
	HandleRefer(refer *SIPInvite) (*SIPResponse, error)

	// HandleResponse handles an incoming SIP response
	HandleResponse(response *SIPResponse) error
}

// SIPServerConfig holds configuration for the SIP server
type SIPServerConfig struct {
	// Listen addresses
	UDPListen string
	TCPListen string
	TLSListen string
	WSListen  string
	WSSListen string

	// TLS configuration
	TLSCertFile string
	TLSKeyFile  string
	TLSCAFile   string
	VerifyClient bool

	// Worker pool
	WorkerCount int
	MaxRequestsPerWorker int

	// SIP trunks
	Trunks []SIPTrunkConfig

	// DNS SRV failover
	DNSSRVEnabled bool
	DNSSRVDomain  string
	DNSSRVTransport SIPTransport

	// Keepalive
	KeepaliveEnabled bool
	KeepaliveInterval time.Duration
	KeepaliveTimeout time.Duration
	KeepaliveRetryCount int

	// Rate limiting
	RateLimitEnabled bool
	RateLimitRPS int
	RateLimitBurst int
}

// SIPTrunkConfig holds SIP trunk configuration
type SIPTrunkConfig struct {
	Name     string
	Host     string
	Port     int
	Transport SIPTransport
	Username string
	Password string
	Realm    string
	Priority int
	Weight   int
}

// SIPTrunk represents a SIP trunk connection
type SIPTrunk struct {
	Config SIPTrunkConfig
	Status TrunkStatus
	LastHealthCheck time.Time
	ConsecutiveFailures int
}

// TrunkStatus represents the status of a SIP trunk
type TrunkStatus string

const (
	TrunkStatusHealthy TrunkStatus = "healthy"
	TrunkStatusDegraded TrunkStatus = "degraded"
	TrunkStatusUnhealthy TrunkStatus = "unhealthy"
	TrunkStatusUnknown TrunkStatus = "unknown"
)

// SIPServer defines the interface for a SIP server
type SIPServer interface {
	// Start starts the SIP server
	Start() error

	// Stop stops the SIP server
	Stop() error

	// SetHandler sets the SIP message handler
	SetHandler(handler SIPHandler)

	// SendResponse sends a SIP response
	SendResponse(response *SIPResponse) error

	// SendRequest sends a SIP request
	SendRequest(request *SIPInvite) error

	// GetTrunks returns the list of SIP trunks
	GetTrunks() []*SIPTrunk

	// GetTrunkStatus returns the status of a specific trunk
	GetTrunkStatus(name string) (TrunkStatus, error)
}