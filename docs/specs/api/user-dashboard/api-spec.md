# 💻 Next.js User Dashboard — Complete API Reference

Comprehensive API reference for the Next.js User Dashboard web portal, covering Authentication, WebRTC softphone dialing, call logs & transcript inspection, AI Agent customization, Outbound Mission delegation, Minute Wallets, Contact CRM & VIP whitelist, AI Support Bot, and Dynamic CMS content.

---

## 📌 1. Authentication & Session APIs

### `POST /api/v1/auth/login`
Authenticates a user with email and password.

#### Request Body:
```json
{
  "email": "sarah.admin@acme.com",
  "password": "SecurePassword123!"
}
```

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "user": {
      "id": "usr_H47BMt2P6_uh4P24IYcgA",
      "email": "sarah.admin@acme.com",
      "name": "Sarah Connor",
      "role": "tenant_admin",
      "phoneNumber": "+15557778899",
      "isVerified": true,
      "createdAt": "2026-01-10T08:30:00.000Z"
    },
    "tokens": {
      "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expiresInSeconds": 86400
    },
    "tenant": {
      "id": "tenant_default",
      "name": "Acme Communications",
      "slug": "acme-comm",
      "plan": "pro"
    }
  }
}
```

---

## 📞 2. WebRTC Softphone & Call Dispatch APIs

### `POST /api/v1/calls/dispatch`
Initiates an outbound AI call to a destination phone number.

```json
{
  "agentId": "agent_phone_default",
  "to": "+15552345678",
  "from": "+18005559000",
  "maxDurationSeconds": 1800,
  "recordCall": true,
  "variables": {
    "customerName": "Alex Johnson",
    "appointmentDate": "2026-08-20T10:00:00Z"
  }
}
```

### `GET /api/v1/calls`
Retrieves a paginated list of calls with filtering by date, status, direction, and agent.

### `GET /api/v1/calls/{id}/transcript`
Fetches the full utterance transcript and sentiment analysis for a completed call.

---

## 🤖 3. AI Agent Persona Customization APIs

### `GET /api/v1/agents`
Lists all conversational agent profiles configured in the tenant workspace.

### `PATCH /api/v1/agents/{id}`
Updates agent parameters, prompt instructions, tools, and voice settings (Cartesia / Deepgram).

```json
{
  "name": "Nexora Advanced Phone Pro",
  "temperature": 0.5,
  "systemPrompt": "You are Nexora, an elite voice assistant. Speak concisely and professionally.",
  "tools": ["calendar_check_availability", "appointment_create", "sms_send"]
}
```

---

## 🎯 4. Autonomous Outbound Missions Management (hiRobin Parity)

### `POST /api/v1/missions`
Creates an autonomous phone calling mission.

```json
{
  "missionType": "restaurant_reservation",
  "targetPhoneNumber": "+918041234567",
  "targetBusinessName": "Olive Bistro",
  "goal": "Book a table for 4 tonight at 8:00 PM (outdoor seating)",
  "constraints": {
    "partySize": 4,
    "preferredTime": "2026-08-25T20:00:00.000Z"
  }
}
```

### `GET /api/v1/missions`
Lists all active and completed missions with IVR progress, hold time, and extracted confirmation details.

---

## 💳 5. Minute Wallets & Billing Ledger

### `GET /api/v1/billing/wallet`
Retrieves current available balance and minute usage counters.

### `GET /api/v1/billing/ledger`
Fetches itemized debit/credit records of all calls and mission consumption.

### `POST /api/v1/payments/intent`
Creates a Stripe / Razorpay payment order to top up minute wallets.

---

## 📇 6. Contact CRM & VIP Screening Rules

### `GET /api/v1/crm/contacts`
Lists contacts with screening policies (`direct_pass`, `ai_screen`, `block`).

### `POST /api/v1/crm/contacts`
Creates a new contact with customized greeting prompts and direct-pass whitelist rules.

```json
{
  "name": "Dr. Sarah Jenkins",
  "phoneNumber": "+15559876543",
  "category": "family_vip",
  "screeningRule": "direct_pass",
  "customAiGreeting": "Hi Dr. Jenkins, Sanket is currently in a meeting, but let me notify him immediately."
}
```

---

## 💬 7. AI Support Bot & Helpdesk Chat

### `POST /api/v1/support/chat/sessions`
Starts a support chat session with Tier-1 AI Support Bot.

### `POST /api/v1/support/chat/sessions/{id}/messages`
Sends user message, searches knowledge base, and connects to human support when needed.

---

## 📄 8. Dynamic CMS & Help Center Content

### `GET /api/v1/content/{slug}`
Fetches dynamic Privacy Policy, Terms of Service, or About Us document.

### `GET /api/v1/content/faqs`
Fetches categorized FAQs for user help center display.
