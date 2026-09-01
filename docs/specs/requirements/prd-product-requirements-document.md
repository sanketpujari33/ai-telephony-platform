# 📄 Product Requirements Document (PRD) — Nexora AI Platform

---

## 1. Executive Summary & Product Vision

Nexora is an enterprise-grade, omnichannel AI Communication Operating System uniting personal telephony intelligence (**Equal AI** & **hiRobin**) with scalable enterprise communication infrastructure. The platform empowers users to deploy autonomous AI agents capable of screening incoming calls with sub-80ms streaming transcripts, executing 1-tap live call takeovers, dispatching autonomous outbound phone missions, navigating IVRs and hold music, managing Dual-SIM carrier forwarding with zero drop rates, and communicating across 9+ Indian dialects and international languages.

---

## 2. Target User Personas & JTBD

| Persona | Role | Key Needs & Jobs to be Done (JTBD) |
| :--- | :--- | :--- |
| **Everyday Smartphone User** | Mobile Subscriber (Android / iOS) | Never be interrupted by spam or unknown calls; auto-screen delivery drivers and security guards; 1-tap live take-over. |
| **Busy Professional / Executive**| Founder / Manager / Consultant | Delegate tedious phone calls (appointments, reservations, customer support hold lines); whisper prompts to AI mid-call. |
| **Enterprise Support Leader** | VP Operations / Support Director | Eliminate hold times, automate multi-tier phone queues, enforce TCPA/TRAI DND compliance, reduce cost per contact by 70%. |
| **Contact Center Operator** | Live Support Specialist | Softphone dialer, live waveform transcripts, smart suggested replies, automated wrap-up summaries. |
| **Platform Developer** | AI & Telephony Engineer | WebSocket streaming APIs, MCP tool integration, customizable LLM prompts, and low-latency voice pipelines. |

---

## 3. Flagship Telephony Requirements (Equal AI + hiRobin Integration)

### 3.1 Inbound AI Call Screening ("Green Card" System)
- **FR-01: Unknown Callers Only Shield**: Saved device contacts ring phone directly; unknown numbers are intercepted and screened by AI.
- **FR-02: Real-Time Live Waveform & Streaming Transcript**: Stream caller audio waveform and speech-to-text text with sub-80ms latency.
- **FR-03: 1-Tap Live Call Take-Over (Barge-In)**: Allow users to instantly jump into the active call, swapping audio to their mic/earpiece in <15ms.
- **FR-04: In-Call Whisper Guidance (Prompt Injection)**: Provide quick prompt chips (*"Leave at door"*, *"Call back in 30m"*) and custom text input for mid-call AI prompt injection.
- **FR-05: Post-Call Action Cards**: Automatically extract calendar appointments, delivery gate codes, OTPs, and actionable tasks.

### 3.2 Dual-SIM Carrier Forwarding & Auto-Detection Engine
- **FR-06: Automatic Dual-SIM Discovery**: Auto-detect SIM 1 & SIM 2 / eSIM carriers (Jio 5G, Airtel 4G, Vi, BSNL, AT&T, Verizon, T-Mobile).
- **FR-07: 1-Click Master Activation**: Sequential MMI shortcode triggering (`*004*+DID#`) with zero manual dialing.
- **FR-08: Background Automation Rules**: Automatic forwarding during DND / Calendar Meetings, on Low Battery (<15%), and after 3 unanswered rings (15s).
- **FR-09: Live Carrier Handshake Diagnostic**: In-app acoustic ping measuring latency and verifying 0.0% failure rate.

### 3.3 Autonomous Outbound Missions (hiRobin Parity)
- **FR-10: Goal-Driven Calling**: User submits natural language task (e.g. *"Book table for 4 at Olive Bistro tonight at 8 PM"*).
- **FR-11: Multi-Level IVR & DTMF Auto-Navigation**: Deepgram STT parses voice menus; RFC 4733 transmits required DTMF tones autonomously.
- **FR-12: Hold Music Spectral Detection**: FFT audio classifier detects hold music loops and alerts user when human agent joins.
- **FR-13: Debriefing & Rich Notifications**: Send push notifications with confirmation codes, audio snippets, and calendar sync.

### 3.4 Multilingual Dialect Engine
- **FR-14: 9+ Indian Dialects**: Native support for Hindi, Hinglish, Marathi, Gujarati, Tamil, Telugu, Kannada, Bengali, Punjabi, and English.
- **FR-15: Dynamic 3-Second Auto-Detect**: Fast Fourier spectral classifier switches prompt and TTS voice in the first 3 seconds of caller speech.

### 3.5 Compliance & Privacy
- **FR-16: Zero-Retention Mode**: Purge audio recordings immediately upon transcript synthesis.
- **FR-17: Automated PII Redaction**: Mask credit cards, OTPs, Aadhaar, and SSNs in persistent storage.
- **FR-18: Mandatory Disclosures**: Compliance announcements conforming to TCPA (US) and TRAI TCCCPR (India).

---

## 4. Key Performance Indicators (KPIs) & SLAs

- **Acoustic Latency**: <80ms median WebRTC streaming latency.
- **Call Drop Rate**: <0.05% carrier forwarding audio failure rate (vs. 42.6% industry baseline).
- **Barge-In Handover**: <15ms audio track transition on live take-over.
- **System Uptime**: 99.99% monthly SLA across global SIP trunk nodes.
