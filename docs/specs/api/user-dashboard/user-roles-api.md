# 👤 User, Team Member & Agent Operator — Complete API Reference

This document specifies the complete API contracts for **Tenant Admins**, **Agent Operators (Softphone)**, **Team Leads**, **Viewers**, and **End-User Customer Self-Service** surfaces.

---

## 📑 User API Directory by Role

```text
├── 👔 Tenant Admin & Team Management
│   ├── POST   /api/v1/user/team/invite              [Invite new team member]
│   ├── GET    /api/v1/user/team/members             [List workspace team members]
│   ├── PATCH  /api/v1/user/team/members/:id/role    [Update member role & seat]
│   └── DELETE /api/v1/user/team/members/:id         [Revoke team member seat]
│
├── 📞 Agent Operator Softphone & In-Call APIs
│   ├── POST   /api/v1/user/dialer/call              [Initiate outbound softphone call]
│   ├── POST   /api/v1/user/dialer/calls/:id/hold    [Mute / Hold active call]
│   ├── POST   /api/v1/user/dialer/calls/:id/transfer[Initiate warm/cold transfer]
│   ├── POST   /api/v1/user/dialer/calls/:id/dtmf    [Send DTMF keypad digits]
│   ├── POST   /api/v1/user/dialer/calls/:id/hangup  [Terminate active call]
│   └── POST   /api/v1/user/contacts/:id/notes       [Add in-call CRM notes]
│
├── 📊 Viewer & Reporting APIs
│   ├── GET    /api/v1/user/analytics/daily-calls    [Daily call volume & duration]
│   ├── GET    /api/v1/user/analytics/sentiment      [Sentiment & CSAT breakdown]
│   └── GET    /api/v1/user/billing/statements       [Download monthly invoices]
│
└── 🌐 End-User Customer Self-Service APIs
    ├── GET    /api/v1/user/profile/preferences      [View communication channels]
    ├── PUT    /api/v1/user/profile/preferences      [Update preferred call window]
    └── POST   /api/v1/user/profile/dnd-toggle       [Toggle DND / marketing consent]
```

---

## 👔 1. Tenant Admin & Team Management APIs

### `POST /api/v1/user/team/invite`
Invites a new team member or agent operator to the tenant workspace.
- **Required Roles**: `TENANT_ADMIN`
- **Required Scope**: `users:invite`

#### Request Body:
```json
{
  "email": "sarah.operator@acme.com",
  "name": "Sarah Jenkins",
  "role": "agent_operator",
  "assignedDID": "+18005559001",
  "dailyOutboundLimit": 500
}
```

#### Response `201 Created`:
```json
{
  "success": true,
  "statusCode": 201,
  "data": {
    "invitationId": "inv_891238472",
    "email": "sarah.operator@acme.com",
    "role": "agent_operator",
    "status": "pending_acceptance",
    "expiresAt": "2026-08-24T15:00:00Z"
  },
  "message": "Team member invitation email dispatched"
}
```

---

### `GET /api/v1/user/team/members`
Lists all active team members, softphone status, and assigned roles.
- **Required Roles**: `TENANT_ADMIN`, `TEAM_LEAD`, `VIEWER`

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": [
    {
      "userId": "usr_991823749",
      "name": "Sarah Jenkins",
      "email": "sarah.operator@acme.com",
      "role": "agent_operator",
      "softphoneStatus": "available",
      "assignedDID": "+18005559001",
      "callsHandledToday": 28,
      "joinedAt": "2026-08-10T09:00:00Z"
    }
  ],
  "pagination": { "page": 1, "pageSize": 20, "totalPages": 1, "totalCount": 1 }
}
```

---

## 📞 2. Agent Operator Softphone & In-Call APIs

### `POST /api/v1/user/dialer/call`
Initiates a softphone dialer call with LiveKit WebRTC audio connection.
- **Required Roles**: `AGENT_OPERATOR`, `TENANT_ADMIN`
- **Required Scope**: `calls:dial`

#### Request Body:
```json
{
  "to": "+15552345678",
  "from": "+18005559001",
  "enableAIAdvisor": true,
  "recordCall": true
}
```

#### Response `201 Created`:
```json
{
  "success": true,
  "statusCode": 201,
  "data": {
    "callId": "call_op_991823749",
    "roomName": "room_call_op_991823749",
    "livekitToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "status": "ringing",
    "startedAt": "2026-08-17T15:30:00Z"
  }
}
```

---

### `POST /api/v1/user/dialer/calls/:id/transfer`
Performs a warm or cold transfer to another team operator or external number.

#### Request Body:
```json
{
  "targetUserId": "usr_support_supervisor_02",
  "transferType": "warm",
  "whisperNotes": "Customer needs technical database assistance."
}
```

---

### `POST /api/v1/user/contacts/:id/notes`
Attaches in-call notes and action items to a customer CRM record.

#### Request Body:
```json
{
  "note": "Customer confirmed they will sign the enterprise annual contract by Friday.",
  "linkedCallId": "call_op_991823749",
  "tags": ["contract", "enterprise"]
}
```

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "noteId": "nte_89123847",
    "contactId": "cnt_99182374",
    "note": "Customer confirmed they will sign the enterprise annual contract by Friday.",
    "createdAt": "2026-08-17T15:35:00Z"
  },
  "message": "Contact note saved successfully"
}
```

---

## 🌐 3. End-User Customer Self-Service APIs

### `GET /api/v1/user/profile/preferences`
Retrieves customer communication channel preferences and preferred calling windows.

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "preferredChannel": "whatsapp",
    "allowedChannels": ["whatsapp", "sms", "voice"],
    "preferredCallingWindow": {
      "start": "09:00",
      "end": "12:00",
      "timezone": "America/New_York"
    },
    "marketingConsent": true,
    "dndEnabled": false
  }
}
```

### `PUT /api/v1/user/profile/preferences`
Updates customer preferred channels and contact window.

#### Request Body:
```json
{
  "preferredChannel": "voice",
  "preferredCallingWindow": {
    "start": "14:00",
    "end": "17:00",
    "timezone": "America/New_York"
  }
}
```
