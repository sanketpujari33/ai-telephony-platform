# 🌐 Next.js Landing Page & Public Ingress — Complete API Reference

---

## 📌 1. Public Lead & Contact Inquiries

### `POST /api/v1/public/contact-us`
Submits a public contact or enterprise sales inquiry form.

#### Request Headers:
```http
Content-Type: application/json
```

#### Request Body:
```json
{
  "fullName": "Sarah Connor",
  "email": "sarah@cyberdyne.com",
  "companyName": "Cyberdyne Systems",
  "companySize": "200-500",
  "estimatedMonthlyMinutes": 50000,
  "useCase": "outbound_sales",
  "message": "We want to migrate 200 customer support agents to Nexora voice AI."
}
```

#### Response `201 Created`:
```json
{
  "success": true,
  "statusCode": 201,
  "data": {
    "inquiryId": "inq_99182374921",
    "receivedAt": "2026-08-17T15:50:00Z"
  },
  "message": "Thank you! An enterprise AI specialist will contact you within 24 hours."
}
```

#### Error Responses:
- **`400 Bad Request`**:
  ```json
  {
    "success": false,
    "statusCode": 400,
    "error": "VALIDATION_FAILED",
    "message": "Validation errors on submitted form fields.",
    "details": [
      { "field": "email", "issue": "Must be a valid business email address" }
    ]
  }
  ```

---

## 🎙️ 2. Interactive Browser Voice Playground

### `POST /api/v1/public/playground/session`
Spawns a temporary, rate-limited 3-minute WebRTC voice sandbox for landing page visitors to test real-time AI voice conversation directly in the browser.

#### Request Body:
```json
{
  "selectedPersona": "customer_support",
  "voiceId": "cartesia_sonic_en_1",
  "clientBrowser": "Chrome 128 (macOS)"
}
```

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "livekitUrl": "wss://livekit.nexora.ai",
    "roomName": "room_demo_sandbox_89123",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "maxDurationSeconds": 180,
    "systemPrompt": "You are a friendly customer service demo agent for Nexora AI."
  },
  "message": "Interactive voice sandbox session provisioned"
}
```

---

## 🟢 3. Public Platform Health & Status Ping

### `GET /api/v1/public/status`
Returns public uptime percentages and operational status across all core telephony services.

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "overallStatus": "operational",
    "services": {
      "voiceTelephony": { "status": "operational", "latencyMs": 42 },
      "messagingApis": { "status": "operational", "latencyMs": 18 },
      "aiOrchestrator": { "status": "operational", "latencyMs": 110 },
      "webrtcMediaSFU": { "status": "operational", "latencyMs": 24 }
    },
    "uptime90DaysPercent": 99.98,
    "lastIncidentAt": null,
    "checkedAt": "2026-08-17T15:50:00Z"
  }
}
```
