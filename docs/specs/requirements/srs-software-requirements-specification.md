# 📜 Software Requirements Specification (SRS) — Nexora AI Platform
*Conforming to ISO/IEC/IEEE 29148:2018 Standards*

---

## 1. Introduction & Purpose

This Software Requirements Specification (SRS) establishes the formal functional, behavioral, and interface requirements for the **Nexora AI Telephony Platform**, including the core backend microservices, Go SIP engine, LiveKit WebRTC SFU, Python AI agent runtime, and the Flutter cross-platform mobile client.

---

## 2. Functional Requirements Specification

### SRS-01: Dual-SIM Telephony & Supplementary Service Management
- **SRS-01.1**: The mobile client shall interface with Android `TelephonyManager` / iOS `CoreTelephony` to detect active SIM cards, slot indices, MCC/MNC carrier codes, and IMS registration status.
- **SRS-01.2**: The system shall generate standard 3GPP Supplementary Service MMI strings (`*004*<DID>#`, `*21*<DID>#`, `##004#`, `##002#`) formatted specifically for the detected carrier.
- **SRS-01.3**: The client shall execute sequential MMI execution across SIM 1 and SIM 2 upon 1-click master activation.

### SRS-02: Real-Time Inbound Voice Screening & WebRTC Streaming
- **SRS-02.1**: The Go SIP Trunk Engine (`:5060`) shall accept incoming conditional forwarded calls and bridge SIP audio to LiveKit WebRTC SFU (`:7880`) within 100ms.
- **SRS-02.2**: The Python Voice Agent shall stream synthesized audio and receive caller audio over bidirectional Opus 48kHz WebRTC tracks.
- **SRS-02.3**: Deepgram STT shall stream speech transcripts to the client over WebSocket (`:4009`) with Time-To-First-Word < 120ms.

### SRS-03: In-Call Whisper Guidance & 1-Tap Live Take-Over
- **SRS-03.1**: The client shall send whisper prompt injection payloads over WebSocket to the AI Orchestrator (`:8086`), modifying the LLM conversational context dynamically during an active call.
- **SRS-03.2**: When the user triggers live call take-over, LiveKit SFU shall mute the AI voice track and bridge the mobile client's audio track directly to the PSTN caller with audio transition latency < 15ms.

### SRS-04: Outbound Missions & IVR Navigation
- **SRS-04.1**: The AI Orchestrator shall execute multi-step outbound missions with configurable user goals, contact numbers, and memory constraints.
- **SRS-04.2**: The voice pipeline shall parse audio prompts for multi-level IVRs and transmit RFC 4733 DTMF tones.
- **SRS-04.3**: The FFT audio classifier shall identify hold music loops and trigger mobile push notifications upon human voice onset.

### SRS-05: Multilingual Dialect Processing
- **SRS-05.1**: The voice agent shall classify caller language within the first 3 seconds of speech across 9+ Indian dialects (Hindi, Hinglish, Marathi, Tamil, Telugu, Kannada, Bengali, Gujarati, Punjabi) and international languages.
- **SRS-05.2**: The LLM prompt and TTS synthesizer voice model shall switch dynamically without disconnecting the SIP session.

---

## 3. Non-Functional & Quality Attributes

| Attribute | Metric & Acceptance Criteria |
| :--- | :--- |
| **Acoustic Latency** | Median round-trip voice latency < 80ms over WebRTC. |
| **Throughput** | Support 10,000 concurrent active calls per regional cluster. |
| **Availability** | 99.99% monthly service uptime. |
| **Security & Privacy** | TLS 1.3 encryption, SRTP media encryption, zero-retention ephemeral audio mode. |
| **Drop Rate** | Carrier forwarding failure rate < 0.05% via redundant STUN/TURN media relay. |
