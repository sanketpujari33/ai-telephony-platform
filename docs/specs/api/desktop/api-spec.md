# 💻 Flutter Desktop App (macOS & Windows) — Complete API Reference

---

## 📌 1. Multi-Line Call & Softphone Management

### `GET /api/v1/desktop/lines`
Returns active phone lines, hold states, and concurrent call slots assigned to the desktop agent seat.

#### Request Headers:
```http
Authorization: Bearer <AGENT_ACCESS_TOKEN>
x-tenant-id: tenant_default
```

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "agentSeatId": "seat_desktop_01",
    "lines": [
      {
        "lineIndex": 1,
        "status": "in-call",
        "callId": "call_101",
        "remoteParty": "+15551112233",
        "callerName": "Tech Corp Support",
        "durationSeconds": 142,
        "isMuted": false,
        "isHeld": false
      },
      {
        "lineIndex": 2,
        "status": "on-hold",
        "callId": "call_102",
        "remoteParty": "+15554445566",
        "callerName": "VIP Customer Line",
        "durationSeconds": 310,
        "isMuted": true,
        "isHeld": true
      },
      {
        "lineIndex": 3,
        "status": "idle",
        "callId": null
      }
    ]
  }
}
```

---

### `POST /api/v1/desktop/calls/transfer`
Performs an attended (warm) or blind (cold) call transfer to another agent or external PSTN number.

#### Request Body:
```json
{
  "callId": "call_101",
  "target": "+18005550199",
  "transferType": "warm",
  "notes": "Customer needs level 2 database assistance.",
  "bridgeAudioBeforeTransfer": true
}
```

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "transferId": "xfer_891238472",
    "status": "initiated",
    "target": "+18005550199",
    "timestamp": "2026-08-17T15:45:00Z"
  },
  "message": "Warm call transfer initiated"
}
```

---

## ⚡ 2. Real-Time Live Utterance & Waveform Stream (`WebSocket`)

### `WSS /ws/v1/calls/:id/live-stream`
Full-duplex WebSocket connection streaming live interim transcript tokens, VAD voice activation states, and 60fps PCM audio RMS energy levels.

#### Query Parameters:
- `token` (`string`, required): Signed JWT session token.

#### Incoming Server Stream Messages:

```json
// 1. Live interim transcription token packet
{
  "type": "transcript_interim",
  "speaker": "user",
  "text": "I am inquiring about my enterprise billing plan",
  "isFinal": false,
  "timestampMs": 1420
}

// 2. Finalized utterance packet
{
  "type": "transcript_final",
  "speaker": "user",
  "text": "I am inquiring about my enterprise billing plan.",
  "isFinal": true,
  "confidence": 0.99,
  "sentiment": 0.78,
  "timestampMs": 3100
}

// 3. High-frequency audio RMS energy packet (for visualizer)
{
  "type": "waveform_energy",
  "inputDb": -18.4,
  "outputDb": -22.1,
  "vadActive": true,
  "fps": 60
}
```

---

## 🛠️ 3. Audio Hardware Diagnostics

### `POST /api/v1/desktop/diagnostics/loopback-test`
Initiates a 5-second audio hardware loopback test to measure round-trip audio latency, jitter, and microphone clarity.

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "roundTripLatencyMs": 18.2,
    "jitterMs": 1.4,
    "inputAudioRmsDb": -14.2,
    "packetLossPercent": 0.0,
    "microphoneClarityScore": 0.98,
    "status": "passed"
  }
}
```
