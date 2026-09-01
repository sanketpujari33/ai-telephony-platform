# 📱 Flutter Mobile App (Android & iOS) — Complete API Specification

Comprehensive API reference for the Flutter Mobile client, covering app bootstrap, phone OTP auth, Dual-SIM MMI forwarding, Live Green Card screening, in-call whispers, 1-tap takeover barge-in, autonomous missions, minute wallets, AI Support Bot, and dynamic CMS content.

---

## 🚀 1. App Startup & Universal Bootstrap

### `GET /api/v1/bootstrap`
Called once immediately upon mobile app launch to obtain feature flags, maintenance status, forced update checks, carrier profiles, payment keys, and STUN/TURN ICE credentials.

#### Request Headers:
```http
Authorization: Bearer <USER_JWT> (Optional if guest)
If-None-Match: "W/\"bootstrap-cfg-20260825-v140\""
Content-Type: application/json
```

#### Query Parameters:
* `platform=android` (or `ios`)
* `version=1.4.0`
* `buildNumber=42`

#### Response `200 OK`:
```json
{
  "success": true,
  "data": {
    "systemStatus": { "isUnderMaintenance": false },
    "appVersionEnforcement": {
      "forceUpdateRequired": false,
      "minimumSupportedVersion": "1.2.0",
      "latestAvailableVersion": "1.4.2"
    },
    "featureFlags": {
      "enableCallScreening": true,
      "enableWhisperGuidance": true,
      "enableOneTapTakeover": true,
      "enableAutonomousOutboundMissions": true,
      "enableDualSimAutoDetection": true,
      "enableTraiSpamShield": true
    },
    "telephonyConfig": {
      "defaultForwardingDid": "+18008490192",
      "livekitWsUrl": "wss://livekit.nexora.ai",
      "websocketGatewayUrl": "wss://api.nexora.ai/ws",
      "stunTurnServers": [
        {
          "urls": ["stun:turn.nexora.ai:3478", "turn:turn.nexora.ai:3478?transport=udp"],
          "username": "nexora_client_ephemeral",
          "credential": "turn_token_99182310"
        }
      ]
    },
    "paymentAndPricing": {
      "activeCurrency": "INR",
      "currencySymbol": "₹",
      "rates": { "inboundScreeningPerMinute": 0.0, "outboundMissionPerMinute": 1.20 },
      "gatewayPublicKeys": { "razorpayKeyId": "rzp_live_9a8b7c6d5e", "stripePublishableKey": "pk_live_..." }
    }
  }
}
```

---

## 🔑 2. Phone Authentication & Device Registration

### `POST /api/v1/auth/otp/send`
Dispatches a 6-digit SMS verification code to the user's mobile number.

```json
{
  "phoneNumber": "+919820188412",
  "channel": "sms"
}
```

### `POST /api/v1/auth/otp/verify`
Validates SMS OTP and returns JWT credentials.

```json
{
  "phoneNumber": "+919820188412",
  "otp": "488488"
}
```

### `POST /api/v1/devices/register`
Registers the device for VoIP and Push Notifications (FCM / APNs).

```json
{
  "deviceId": "dev_android_pix_8pro_9918",
  "platform": "android",
  "fcmToken": "fK9j284jlsdf08-23jksdf...",
  "appVersion": "1.4.0",
  "osVersion": "Android 14 (API 34)",
  "deviceModel": "Google Pixel 8 Pro"
}
```

---

## 📶 3. Dual-SIM Telephony & MMI Forwarding Automation

### `GET /api/v1/telephony/carriers/detect?mccMnc=405854`
Auto-detects the telecom carrier from device SIM telemetry.

```json
{
  "success": true,
  "data": {
    "carrierName": "Reliance Jio 5G",
    "mccMnc": "405854",
    "activationMmi": "*004*+18008490192#",
    "deactivationMmi": "##004#",
    "forwardingDid": "+18008490192"
  }
}
```

### `POST /api/v1/telephony/sim-slots`
Registers user SIM slot preferences (e.g. Unknown number shield, Low battery `<15%` failover, Meeting DND shield).

```json
{
  "deviceId": "dev_android_pix_8pro_9918",
  "slotIndex": 0,
  "carrierName": "Reliance Jio 5G",
  "phoneNumber": "+919820188412",
  "unknownOnlyShield": true,
  "lowBatteryFailover": true,
  "lowBatteryThreshold": 15,
  "dndMeetingShield": true
}
```

### `POST /api/v1/telephony/sim-slots/{id}/verify`
Verifies active forwarding registration with the carrier.

---

## 🟢 4. "Green Card" Live Screening, Whisper & 1-Tap Takeover

### WebSocket Live Stream Connection (`ws://api.nexora.ai/ws?token=<JWT>`)
* **Event `screening:started`**: Displays full-screen Green Card UI with caller identity and intent.
* **Event `screening:transcript_chunk`**: Appends real-time diarized text (`<120ms` chunk).

### `POST /api/v1/calls/{callSessionId}/whisper`
Injects an in-call whisper guidance prompt into the AI agent's speech queue (`<18ms`).

```json
{
  "whisperType": "preset_chip",
  "promptText": "Delivery PIN is 8492. Leave package with guard at Gate 2.",
  "targetLanguage": "hinglish"
}
```

### `POST /api/v1/calls/{callSessionId}/takeover`
Triggers 1-tap live call takeover barge-in (`<15ms`). Instantly mutes the AI agent and bridges the user's WebRTC microphone track to the remote caller.

```json
{
  "deviceId": "dev_android_pix_8pro_9918",
  "muteAiImmediate": true
}
```

---

## 🎯 5. Autonomous Outbound Missions (hiRobin Parity)

### `POST /api/v1/missions`
Creates a goal-oriented autonomous phone calling mission.

```json
{
  "missionType": "restaurant_reservation",
  "targetPhoneNumber": "+918041234567",
  "targetBusinessName": "Olive Bistro",
  "goal": "Book a table for 4 tonight at 8:00 PM (outdoor seating)",
  "constraints": {
    "partySize": 4,
    "preferredTime": "2026-08-25T20:00:00.000Z",
    "seating": "outdoor"
  }
}
```

### `GET /api/v1/missions/{missionId}`
Polls live IVR progress, hold duration, and extracted confirmation details.

---

## 💳 6. Minute Wallets & In-App Payments

### `GET /api/v1/billing/wallet`
Fetches current available call minutes and cash balance.

### `POST /api/v1/payments/intent`
Creates a Razorpay / Stripe payment order to recharge minute wallets.

```json
{
  "gateway": "razorpay",
  "amount": 500.00,
  "currency": "INR",
  "topUpMinutes": 250
}
```

---

## 🤖 7. In-App AI Support Bot & Helpdesk Chat

### `POST /api/v1/support/chat/sessions`
Starts a Tier-1 AI Support Bot chat session for troubleshooting.

```json
{
  "category": "call_forwarding",
  "initialMessage": "My Jio call forwarding failed with MMI error."
}
```

### `POST /api/v1/support/chat/sessions/{id}/messages`
Sends user message to AI bot; escalates to live human support when requested.

### `POST /api/v1/support/chat/sessions/{id}/feedback`
Submits 1–5 star CSAT feedback after resolution.

---

## 📄 8. Dynamic Legal, Privacy Policy & FAQs Content

### `GET /api/v1/content/{slug}`
Fetches dynamic database-driven legal documents (`privacy-policy`, `terms-and-conditions`, `about-us`, `refund-policy`).

### `GET /api/v1/content/faqs?category=all`
Fetches categorized FAQs for instant on-device searching and help center display.
