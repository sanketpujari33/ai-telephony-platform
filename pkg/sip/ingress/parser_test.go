package ingress

import (
	"net"
	"strings"
	"testing"
)

func TestSIPParser_ParseRequest(t *testing.T) {
	parser := NewSIPParser()
	localAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5060}
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("192.168.1.100"), Port: 5060}

	rawSIP := []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds\r\n" +
		"Max-Forwards: 70\r\n" +
		"To: Bob <sip:bob@biloxi.com>\r\n" +
		"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
		"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
		"CSeq: 314159 INVITE\r\n" +
		"Contact: <sip:alice@pc33.atlanta.com>\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 129\r\n" +
		"\r\n" +
		"v=0\r\n" +
		"o=alice 2890844526 2890844526 IN IP4 10.1.1.1\r\n" +
		"s=\r\n" +
		"c=IN IP4 10.1.1.1\r\n" +
		"t=0 0\r\n" +
		"m=audio 49170 RTP/AVP 0\r\n" +
		"a=rtpmap:0 PCMU/8000\r\n")

	invite, err := parser.ParseRequest(rawSIP, TransportUDP, localAddr, remoteAddr)
	if err != nil {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	if invite.CallID != "a84b4c76e66710@pc33.atlanta.com" {
		t.Errorf("expected CallID a84b4c76e66710@pc33.atlanta.com, got %s", invite.CallID)
	}
	
	if string(invite.Request.Method) != "INVITE" {
		t.Errorf("expected Method INVITE, got %s", string(invite.Request.Method))
	}

	if !strings.Contains(invite.SDPBody, "m=audio 49170 RTP/AVP 0") {
		t.Errorf("missing SDP body content")
	}

	if invite.ContentLength != 129 {
		t.Errorf("expected ContentLength 129, got %d", invite.ContentLength)
	}
}

func TestSIPParser_ParseResponse(t *testing.T) {
	parser := NewSIPParser()
	localAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5060}
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("192.168.1.100"), Port: 5060}

	rawSIP := []byte("SIP/2.0 200 OK\r\n" +
		"Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds\r\n" +
		"To: Bob <sip:bob@biloxi.com>;tag=a6c85cf\r\n" +
		"From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n" +
		"Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n" +
		"CSeq: 314159 INVITE\r\n" +
		"Contact: <sip:bob@192.0.2.4>\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 127\r\n" +
		"\r\n" +
		"v=0\r\n" +
		"o=bob 2890844527 2890844527 IN IP4 192.0.2.4\r\n" +
		"s=\r\n" +
		"c=IN IP4 192.0.2.4\r\n" +
		"t=0 0\r\n" +
		"m=audio 3456 RTP/AVP 0\r\n" +
		"a=rtpmap:0 PCMU/8000\r\n")

	resp, err := parser.ParseResponse(rawSIP, TransportUDP, localAddr, remoteAddr)
	if err != nil {
		t.Fatalf("unexpected error parsing response: %v", err)
	}

	if resp.CallID != "a84b4c76e66710@pc33.atlanta.com" {
		t.Errorf("expected CallID a84b4c76e66710@pc33.atlanta.com, got %s", resp.CallID)
	}
	
	if resp.StatusCode != 200 {
		t.Errorf("expected StatusCode 200, got %d", resp.StatusCode)
	}

	if !strings.Contains(resp.SDPBody, "m=audio 3456 RTP/AVP 0") {
		t.Errorf("missing SDP body content")
	}
}

func TestSIPParser_InvalidData(t *testing.T) {
	parser := NewSIPParser()
	localAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5060}
	remoteAddr := &net.UDPAddr{IP: net.ParseIP("192.168.1.100"), Port: 5060}

	rawSIP := []byte("INVALID_SIP_MESSAGE")

	_, err := parser.ParseRequest(rawSIP, TransportUDP, localAddr, remoteAddr)
	if err == nil {
		t.Errorf("expected error when parsing invalid request")
	}

	_, err = parser.ParseResponse(rawSIP, TransportUDP, localAddr, remoteAddr)
	if err == nil {
		t.Errorf("expected error when parsing invalid response")
	}
}
