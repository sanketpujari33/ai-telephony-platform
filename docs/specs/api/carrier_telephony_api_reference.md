# 🌐 Carrier Telephony, Live Take-Over, & Missions API Reference

## 📑 Overview

This API reference defines the core HTTP/REST and real-time WebSocket endpoints powering Dual-SIM carrier forwarding management, live call screening, whisper prompt injection, 1-tap takeover, and outbound mission execution across the Nexora AI Telephony Platform.

---

## 📡 1. Carrier & Dual-SIM Management Endpoints

### 1.1 Discover Active SIMs & Carrier Profiles
- **Endpoint:** `POST /api/v1/telephony/carrier/discover`
- **Headers:** `Authorization: Bearer <jwt_token>`
- **Request Body:**
```json
{
  "device_id": "dev_99812_android",
  "sim_slots": [
    {
      "slot_index": 0,
      "iccid": "899185400291029384",
      "carrier_name": "Jio 5G",
      "mcc_mnc": "405854",
      "phone_number": "+919820188412"
    },
    {
      "slot_index": 1,
      "iccid": "899187000182736451",
      "carrier_name": "Airtel",
      "mcc_mnc": "404045",
      "phone_number": "+919871044921"
    }
  ]
}
```
- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "virtual_did": "+18008490192",
    "regional_did_in": "+918047360000",
    "configurations": [
      {
        "slot_index": 0,
        "carrier": "Jio 5G",
        "activation_mmi": "*004*+918047360000#",
        "deactivation_mmi": "##004#",
        "status_check_mmi": "*#67#",
        "status": "ready"
      },
      {
        "slot_index": 1,
        "carrier": "Airtel",
        "activation_mmi": "*004*+918047360000#",
        "deactivation_mmi": "##004#",
        "status_check_mmi": "*#62#",
        "status": "ready"
      }
    ]
  }
}
```

### 1.2 Run Carrier Handshake Diagnostic Ping
- **Endpoint:** `POST /api/v1/telephony/carrier/diagnostic-ping`
- **Request Body:**
```json
{
  "slot_index": 0,
  "carrier": "Jio 5G",
  "test_mode": "acoustic_latency"
}
```
- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "round_trip_latency_ms": 78,
    "jitter_ms": 4.2,
    "packet_loss_percent": 0.0,
    "status": "HANDSHAKE_VERIFIED",
    "failure_rate": "0.0%"
  }
}
```

---

## 🎙️ 2. Live Screening, Whisper, & Take-Over Endpoints

### 2.1 Live Call WebSocket Stream
- **WebSocket URL:** `wss://api.nexora.ai/ws/calls/stream?call_id=<call_id>&token=<jwt_token>`
- **Inbound Server Message (Real-Time Transcript Chunk):**
```json
{
  "event": "transcript_chunk",
  "call_id": "call_8829103",
  "speaker": "caller",
  "text": "Hi, I am outside with your Swiggy delivery.",
  "confidence": 0.98,
  "timestamp_ms": 1420,
  "intent_category": "delivery"
}
```

### 2.2 Inject In-Call Whisper Guidance
- **Endpoint:** `POST /api/v1/calls/:call_id/whisper`
- **Request Body:**
```json
{
  "whisper_type": "preset_chip",
  "prompt_text": "Please leave the parcel with security guard at Gate 2."
}
```
- **Response (200 OK):**
```json
{
  "success": true,
  "status": "whisper_injected",
  "ai_acknowledgement": "Prompt queued for immediate speech output"
}
```

### 2.3 1-Tap Live Call Take-Over (Barge-In)
- **Endpoint:** `POST /api/v1/calls/:call_id/takeover`
- **Request Body:**
```json
{
  "device_id": "dev_99812_android",
  "audio_channel": "webrtc"
}
```
- **Response (200 OK):**
```json
{
  "success": true,
  "status": "takeover_active",
  "transition_latency_ms": 12,
  "livekit_room": "room_call_8829103",
  "token": "livekit_jwt_token_..."
}
```

---

## 🚀 3. Outbound Mission Endpoints (hiRobin Parity)

### 3.1 Start Autonomous Outbound Mission
- **Endpoint:** `POST /api/v1/missions/start`
- **Request Body:**
```json
{
  "target_phone_number": "+918041234567",
  "mission_type": "appointment_booking",
  "goal": "Book a table for 4 at Olive Bistro tonight at 8:00 PM",
  "constraints": {
    "max_acceptable_delay_mins": 30,
    "seating_preference": "outdoor",
    "name_for_reservation": "Sanket"
  }
}
```
- **Response (202 Accepted):**
```json
{
  "success": true,
  "data": {
    "mission_id": "msn_7739102",
    "status": "dialing",
    "created_at": "2026-08-25T19:00:00Z"
  }
}
```

### 3.2 Get Mission Status & Debriefing
- **Endpoint:** `GET /api/v1/missions/:mission_id`
- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "mission_id": "msn_7739102",
    "status": "completed",
    "call_duration_seconds": 142,
    "hold_duration_seconds": 38,
    "result": {
      "confirmed": true,
      "confirmation_code": "Table 14 (8:15 PM)",
      "summary": "Reserved outdoor table for 4 at Olive Bistro for 8:15 PM under Sanket.",
      "audio_recording_url": "https://s3.nexora.ai/recordings/msn_7739102.mp3"
    }
  }
}
```
