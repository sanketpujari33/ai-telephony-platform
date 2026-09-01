package ingress

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/emiago/sipgo/sip"
)

// SIPParser handles parsing of SIP messages
type SIPParser struct{}

// NewSIPParser creates a new SIP parser
func NewSIPParser() *SIPParser {
	return &SIPParser{}
}

// ParseRequest parses a SIP request from raw bytes
func (p *SIPParser) ParseRequest(data []byte, transport SIPTransport, localAddr, remoteAddr *net.UDPAddr) (*SIPInvite, error) {
	request, err := sip.ParseMessage(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SIP request: %w", err)
	}
	req, ok := request.(*sip.Request)
	if !ok {
		return nil, fmt.Errorf("parsed message is not a request")
	}

	invite := &SIPInvite{
		Request:     req,
		Transport:   transport,
		LocalAddr:   localAddr,
		RemoteAddr:  remoteAddr,
		ReceivedAt:  time.Now(),
	}

	// Parse headers
	if err := p.parseRequestHeaders(invite); err != nil {
		return nil, fmt.Errorf("failed to parse request headers: %w", err)
	}

	// Parse body
	if body := request.Body(); len(body) > 0 {
		invite.SDPBody = string(body)
		invite.ContentLength = len(body)
	}

	return invite, nil
}

// ParseResponse parses a SIP response from raw bytes
func (p *SIPParser) ParseResponse(data []byte, transport SIPTransport, localAddr, remoteAddr *net.UDPAddr) (*SIPResponse, error) {
	response, err := sip.ParseMessage(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SIP response: %w", err)
	}
	respObj, ok := response.(*sip.Response)
	if !ok {
		return nil, fmt.Errorf("parsed message is not a response")
	}

	resp := &SIPResponse{
		Response:    respObj,
		Transport:   transport,
		LocalAddr:   localAddr,
		RemoteAddr:  remoteAddr,
		SentAt:      time.Now(),
	}

	// Parse headers
	if err := p.parseResponseHeaders(resp); err != nil {
		return nil, fmt.Errorf("failed to parse response headers: %w", err)
	}

	// Parse body
	if body := response.Body(); len(body) > 0 {
		resp.SDPBody = string(body)
		resp.ContentLength = len(body)
	}

	return resp, nil
}

// parseRequestHeaders parses headers from a SIP request
func (p *SIPParser) parseRequestHeaders(invite *SIPInvite) error {
	request := invite.Request

	// Call-ID
	if callID := request.CallID(); callID != nil {
		invite.CallID = callID.Value()
	}

	// From header
	if from := request.From(); from != nil {
		invite.FromURI = &from.Address
	}

	// To header
	if to := request.To(); to != nil {
		invite.ToURI = &to.Address
	}

	// Contact header
	if contact := request.Contact(); contact != nil {
		invite.ContactURI = &contact.Address
	}

	// Via header
	if via := request.Via(); via != nil {
		invite.ViaHeader = via
	}

	// CSeq
	if cseq := request.CSeq(); cseq != nil {
		invite.CSeq = cseq.SeqNo
	}

	// Max-Forwards
	if maxForwards := request.MaxForwards(); maxForwards != nil {
		if val, err := strconv.Atoi(maxForwards.Value()); err == nil {
			invite.MaxForwards = uint32(val)
		}
	}

	// Content-Type
	if contentType := request.ContentType(); contentType != nil {
		invite.ContentType = contentType.Value()
	}

	// Content-Length
	if contentLength := request.ContentLength(); contentLength != nil {
		if val, err := strconv.Atoi(contentLength.Value()); err == nil {
			invite.ContentLength = val
		}
	}

	return nil
}

// parseResponseHeaders parses headers from a SIP response
func (p *SIPParser) parseResponseHeaders(resp *SIPResponse) error {
	response := resp.Response

	// Call-ID
	if callID := response.CallID(); callID != nil {
		resp.CallID = callID.Value()
	}

	// CSeq
	if cseq := response.CSeq(); cseq != nil {
		resp.CSeq = cseq.SeqNo
	}

	// Via header
	if via := response.Via(); via != nil {
		resp.ViaHeader = via
	}

	// Contact header
	if contact := response.Contact(); contact != nil {
		resp.ContactURI = &contact.Address
	}

	// Content-Type
	if contentType := response.ContentType(); contentType != nil {
		resp.ContentType = contentType.Value()
	}

	// Content-Length
	if contentLength := response.ContentLength(); contentLength != nil {
		if val, err := strconv.Atoi(contentLength.Value()); err == nil {
			resp.ContentLength = val
		}
	}

	// Status code and reason phrase
	resp.StatusCode = response.StatusCode
	resp.ReasonPhrase = response.Reason

	return nil
}

// IsInvite checks if the request is an INVITE
func (p *SIPParser) IsInvite(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodINVITE))
}

// IsAck checks if the request is an ACK
func (p *SIPParser) IsAck(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodACK))
}

// IsBye checks if the request is a BYE
func (p *SIPParser) IsBye(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodBYE))
}

// IsCancel checks if the request is a CANCEL
func (p *SIPParser) IsCancel(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodCANCEL))
}

// IsOptions checks if the request is an OPTIONS
func (p *SIPParser) IsOptions(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodOPTIONS))
}

// IsRegister checks if the request is a REGISTER
func (p *SIPParser) IsRegister(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodREGISTER))
}

// IsInfo checks if the request is an INFO
func (p *SIPParser) IsInfo(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodINFO))
}

// IsRefer checks if the request is a REFER
func (p *SIPParser) IsRefer(request *sip.Request) bool {
	return strings.EqualFold(string(request.Method), string(MethodREFER))
}

// IsReInvite checks if the request is a re-INVITE (INVITE with same Call-ID and higher CSeq)
func (p *SIPParser) IsReInvite(request *sip.Request, existingCallID string, existingCSeq uint32) bool {
	if !p.IsInvite(request) {
		return false
	}

	callID := request.CallID()
	if callID == nil {
		return false
	}

	if callID.Value() != existingCallID {
		return false
	}

	cseq := request.CSeq()
	if cseq == nil {
		return false
	}

	return cseq.SeqNo > existingCSeq
}

// GetMethodFromRequest extracts the SIP method from a request
func (p *SIPParser) GetMethodFromRequest(request *sip.Request) SIPMethod {
	return SIPMethod(strings.ToUpper(string(request.Method)))
}

// BuildResponse builds a SIP response
func (p *SIPParser) BuildResponse(request *sip.Request, statusCode int, reasonPhrase string, headers []sip.Header, body string) (*sip.Response, error) {
	response := sip.NewResponseFromRequest(request, statusCode, reasonPhrase, []byte(body))

	// Add headers
	for _, h := range headers {
		response.AppendHeader(h)
	}

	return response, nil
}

// BuildRequest builds a SIP request
func (p *SIPParser) BuildRequest(method SIPMethod, ruri *sip.Uri, headers []sip.Header, body string) (*sip.Request, error) {
	if ruri == nil {
		ruri = &sip.Uri{}
	}
	request := sip.NewRequest(sip.RequestMethod(method), *ruri)

	// Add headers
	for _, h := range headers {
		request.AppendHeader(h)
	}

	// Add body if provided
	if body != "" {
		request.SetBody([]byte(body))
		contentType := sip.NewHeader("Content-Type", "application/sdp")
		request.AppendHeader(contentType)
		contentLength := sip.NewHeader("Content-Length", fmt.Sprintf("%d", len(body)))
		request.AppendHeader(contentLength)
	}

	return request, nil
}

// CreateSDPResponse creates an SDP response for an INVITE
func (p *SIPParser) CreateSDPResponse(localIP string, localPort int, codecs []string) string {
	var sdp strings.Builder
	sdp.WriteString("v=0\r\n")
	sdp.WriteString(fmt.Sprintf("o=- %d %d IN IP4 %s\r\n", time.Now().Unix(), time.Now().Unix(), localIP))
	sdp.WriteString("s=AI Telephony Platform\r\n")
	sdp.WriteString("c=IN IP4 " + localIP + "\r\n")
	sdp.WriteString("t=0 0\r\n")
	sdp.WriteString(fmt.Sprintf("m=audio %d RTP/AVP", localPort))

	for _, codec := range codecs {
		sdp.WriteString(" " + codec)
	}
	sdp.WriteString("\r\n")

	for _, codec := range codecs {
		switch codec {
		case "0": // PCMU
			sdp.WriteString("a=rtpmap:0 PCMU/8000\r\n")
		case "8": // PCMA
			sdp.WriteString("a=rtpmap:8 PCMA/8000\r\n")
		case "9": // G722
			sdp.WriteString("a=rtpmap:9 G722/16000\r\n")
		case "111": // Opus
			sdp.WriteString("a=rtpmap:111 opus/48000/2\r\n")
			sdp.WriteString("a=fmtp:111 useinbandfec=1; usedtx=1\r\n")
		}
	}

	sdp.WriteString("a=sendrecv\r\n")
	sdp.WriteString("a=rtcp-mux\r\n")

	return sdp.String()
}