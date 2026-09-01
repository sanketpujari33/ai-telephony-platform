package media

import (
	"testing"
	"time"

	"github.com/pion/rtp"
	"github.com/stretchr/testify/assert"
)

func TestJitterBufferPushAndPop(t *testing.T) {
	jb := NewJitterBuffer(20, 200, 60, true)
	pkt1 := &rtp.Packet{
		Header: rtp.Header{
			SequenceNumber: 100,
			Timestamp:      1000,
		},
		Payload: []byte{0x01, 0x02, 0x03},
	}

	jb.Push(pkt1)
	poppedPkt, exists := jb.Pop(100)
	assert.True(t, exists)
	assert.Equal(t, uint16(100), poppedPkt.SequenceNumber)

	_, exists = jb.Pop(100)
	assert.False(t, exists)
}

func TestDTMFProcessRFC4733(t *testing.T) {
	de := NewDTMFEngine(true)
	// Payload for DTMF '5': event code 5, end bit set, volume 10, duration 160
	payload := []byte{0x05, 0x80 | 0x0A, 0x00, 0xA0}
	pkt := &rtp.Packet{
		Header: rtp.Header{
			PayloadType: 101,
		},
		Payload: payload,
	}

	event := de.ProcessRTPPacket(pkt, 101)
	assert.NotNil(t, event)
	assert.Equal(t, '5', event.Digit)

	select {
	case ev := <-de.Events():
		assert.Equal(t, '5', ev.Digit)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Expected DTMF event in channel")
	}
}

func TestMediaQualityMetrics(t *testing.T) {
	metrics := &MediaQualityMetrics{MOS: 4.5}

	pkt1 := &rtp.Packet{
		Header: rtp.Header{
			SequenceNumber: 1,
			Timestamp:      1000,
		},
	}
	pkt2 := &rtp.Packet{
		Header: rtp.Header{
			SequenceNumber: 3, // Gap of 1 packet
			Timestamp:      1160,
		},
	}

	metrics.Update(pkt1, 8000)
	metrics.Update(pkt2, 8000)

	assert.Equal(t, uint64(2), metrics.PacketsReceived)
	assert.Equal(t, uint64(1), metrics.PacketsLost)
	assert.True(t, metrics.MOS >= 1.0 && metrics.MOS <= 4.5)
}
