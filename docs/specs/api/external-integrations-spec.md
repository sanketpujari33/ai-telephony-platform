# 🔌 External Integrations & Carrier Gateways Specification

This document details the external API integrations, carrier webhook handling, and third-party AI interfaces supported by Nexora.

---

## 1. Carrier & Telephony Integrations

### 1.1 Twilio & Telnyx SIP Trunking
- **Inbound Webhook Verification**: Validates `X-Twilio-Signature` using SHA-1 HMAC with the auth token.
- **Outbound Dispatch**: REST API invocation via `@nexora/shared` circuit-breaker with automatic failover to secondary trunk.
- **Status Callbacks**: `initiated`, `ringing`, `answered`, `completed`, `busy`, `failed`, `no-answer`.

### 1.2 Meta WhatsApp Cloud API
- **Webhook Verification**: Challenge-response handshake on `GET /webhook/whatsapp` (`hub.verify_token` + `hub.challenge`).
- **Signature Verification**: Validates `X-Hub-Signature-256` HMAC-SHA256 with the app secret.
- **Interactive Messages**: Supports Quick Reply Buttons, List Messages, and Media Uploads.

---

## 2. Speech & Voice AI Gateways

### 2.1 Deepgram STT (Streaming & Prerecorded)
- **Model**: `nova-2` with Smart Format, interim results, punctuation, and keyword boosting.
- **Protocol**: WebSocket binary PCM audio stream (`wss://api.deepgram.com/v1/listen?model=nova-2&punctuate=true&interim_results=true`).

### 2.2 Cartesia & ElevenLabs TTS
- **Cartesia Sonic**: Ultra-low latency (<100ms TTFB) WebRTC streaming TTS.
- **ElevenLabs Flash v2.5**: High expressiveness voice synthesis for warm customer service interactions.

---

## 3. Payment Gateway Webhook Security

All billing webhooks are cryptographically validated in `@nexora/shared`:
- **Stripe**: `stripe-signature` timestamp + HMAC-SHA256.
- **Razorpay**: `x-razorpay-signature` HMAC-SHA256.
- **LemonSqueezy**: `x-signature` raw payload hash.
- **Paddle**: `Paddle-Signature` timestamped HMAC.
- **PayPal**: RSA signature + CRC32 verification.
- **Coinbase Commerce**: `X-CC-Webhook-Signature` HMAC.
