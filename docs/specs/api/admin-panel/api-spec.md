# 👨‍💼 Next.js Platform Admin Panel — API Specification with Comprehensive Roles

This specification provides the complete endpoint reference for Platform Administration, segmented by role responsibilities across **SuperAdmin**, **Tenant Ops**, **Marketing**, **Customer Support**, **Bot Support**, and **Developer** domains.

References:
- [**Comprehensive RBAC Matrix & Role Hierarchy**](./rbac-matrix.md)
- [**Roles & Permissions CRUD Management Guide**](./roles-permissions-crud.md)

---

## 📑 Admin API Endpoint Directory

```text
├── 🔐 Authentication & RBAC
│   ├── POST   /api/v1/admin/auth/login                  [Public / Credentials]
│   ├── POST   /api/v1/admin/auth/mfa/verify             [MFA Required]
│   ├── GET    /api/v1/admin/roles                       [View Roles List]
│   ├── POST   /api/v1/admin/roles                       [Add / Create Custom Role]
│   ├── GET    /api/v1/admin/roles/:id                   [View Single Role Details]
│   ├── PATCH  /api/v1/admin/roles/:id                   [Edit / Update Role]
│   ├── DELETE /api/v1/admin/roles/:id                   [Delete Custom Role]
│   ├── GET    /api/v1/admin/permissions                 [View All System Permissions]
│   ├── POST   /api/v1/admin/roles/:id/permissions       [Add Permission to Role]
│   ├── DELETE /api/v1/admin/roles/:id/permissions/:key  [Remove Permission from Role]
│   └── PUT    /api/v1/admin/users/:id/role              [Assign / Edit User Role]
│
├── 🎧 Customer Support & Live Chat Helpdesk
│   ├── GET    /api/v1/admin/support/chats               [SUPPORT_LEAD, CUSTOMER_SUPPORT]
│   ├── POST   /api/v1/admin/support/chats/:id/claim     [SUPPORT_LEAD, CUSTOMER_SUPPORT]
│   ├── POST   /api/v1/admin/support/chats/:id/messages  [SUPPORT_LEAD, CUSTOMER_SUPPORT]
│   ├── POST   /api/v1/admin/support/chats/:id/internal-notes [SUPPORT_LEAD, CUSTOMER_SUPPORT]
│   ├── POST   /api/v1/admin/support/chats/:id/resolve   [SUPPORT_LEAD, CUSTOMER_SUPPORT]
│   ├── GET    /api/v1/admin/support/canned-responses    [SUPPORT_LEAD, CUSTOMER_SUPPORT]
│   ├── POST   /api/v1/admin/support/canned-responses    [SUPPORT_LEAD, SUPERADMIN]
│   ├── GET    /api/v1/admin/support/tickets             [SUPPORT_LEAD, CUSTOMER_SUPPORT]
│   └── GET    /api/v1/admin/support/csat-metrics        [SUPPORT_LEAD, MARKETING_ADMIN]
│
├── 📈 Marketing & Campaign Management
│   ├── GET    /api/v1/admin/marketing/campaigns         [MARKETING_ADMIN, TENANT_ADMIN]
│   ├── POST   /api/v1/admin/marketing/broadcast         [MARKETING_ADMIN, TENANT_ADMIN]
│   └── POST   /api/v1/admin/marketing/templates         [MARKETING_ADMIN, TENANT_ADMIN]
│
├── 🤖 Bot Support & AI Supervision
│   ├── GET    /api/v1/admin/bots/audit-transcripts      [BOT_SUPPORT, SUPPORT_LEAD]
│   ├── POST   /api/v1/admin/bots/:id/tune-prompt        [BOT_SUPPORT, TENANT_ADMIN]
│   ├── PUT    /api/v1/admin/bots/:id/fallbacks          [BOT_SUPPORT, PLATFORM_ADMIN]
│   └── POST   /api/v1/admin/bots/sandbox-test           [BOT_SUPPORT, DEVELOPER]
│
├── 💻 Developer & API Integrations
│   ├── GET    /api/v1/admin/developers/api-keys         [DEVELOPER, TENANT_ADMIN]
│   ├── POST   /api/v1/admin/developers/api-keys         [DEVELOPER, TENANT_ADMIN]
│   ├── POST   /api/v1/admin/developers/webhooks         [DEVELOPER, TENANT_ADMIN]
│   ├── POST   /api/v1/admin/developers/mcp/install      [DEVELOPER, TENANT_ADMIN]
│   └── GET    /api/v1/admin/developers/logs             [DEVELOPER, PLATFORM_ADMIN]
│
├── 📄 Dynamic CMS & Legal Content Management
│   ├── GET    /api/v1/admin/content                     [CONTENT_MANAGER, COMPLIANCE_OFFICER, SUPERADMIN]
│   ├── GET    /api/v1/admin/content/:slug               [CONTENT_MANAGER, COMPLIANCE_OFFICER, SUPERADMIN]
│   ├── PUT    /api/v1/admin/content/:slug               [CONTENT_MANAGER, COMPLIANCE_OFFICER, SUPERADMIN]
│   ├── POST   /api/v1/admin/content/faqs                [CONTENT_MANAGER, SUPERADMIN]
│   ├── PUT    /api/v1/admin/content/faqs/:id            [CONTENT_MANAGER, SUPERADMIN]
│   ├── DELETE /api/v1/admin/content/faqs/:id            [CONTENT_MANAGER, SUPERADMIN]
│   └── PUT    /api/v1/admin/content/faqs/reorder        [CONTENT_MANAGER, SUPERADMIN]
│
├── ⚙️ Remote Config & Feature Flags
│   ├── GET    /api/v1/admin/config/client-settings      [PLATFORM_ADMIN, SUPERADMIN]
│   └── PUT    /api/v1/admin/config/client-settings      [PLATFORM_ADMIN, SUPERADMIN]
│
├── 📶 Carrier Telemetry & MMI Diagnostics
│   ├── GET    /api/v1/admin/telephony/carrier-telemetry [TELECOM_ADMIN, PLATFORM_ADMIN, SUPERADMIN]
│   └── POST   /api/v1/admin/telephony/carriers/mmi-sync [TELECOM_ADMIN, SUPERADMIN]
│
└── 🔍 Platform Health & Audit
    ├── GET    /api/v1/admin/system/health               [PLATFORM_ADMIN, SUPERADMIN]
    └── GET    /api/v1/admin/audit-logs                  [COMPLIANCE_AUDITOR, SUPERADMIN]
```

---

## 👥 1. Role & Permission CRUD APIs

*(Detailed schemas, request/response models, and safety rules available in [**`roles-permissions-crud.md`**](./roles-permissions-crud.md))*

### Endpoints:
* `GET /api/v1/admin/roles` — Paginated list of all system and custom roles.
* `POST /api/v1/admin/roles` — Create new custom role with assigned permission scopes.
* `GET /api/v1/admin/roles/:id` — View role details, assigned users, and permissions.
* `PATCH /api/v1/admin/roles/:id` — Modify role title, description, or permissions.
* `DELETE /api/v1/admin/roles/:id` — Delete custom role with user reassignment option.
* `GET /api/v1/admin/permissions` — Catalog of all available system permission keys.
* `POST /api/v1/admin/roles/:id/permissions` — Attach permission scope to role.
* `DELETE /api/v1/admin/roles/:id/permissions/:key` — Revoke permission scope from role.
* `PUT /api/v1/admin/users/:id/role` — Assign or change a user's role.

---

## 🎧 2. Customer Support & Live Chat Helpdesk APIs

### `GET /api/v1/admin/support/chats`
Fetches all active, queued, or escalated customer support chat sessions.

#### Response `200 OK`:
```json
{
  "success": true,
  "data": [
    {
      "id": "chat_session_88192",
      "userId": "usr_991823",
      "userName": "Sanket Sharma",
      "status": "agent_queue",
      "priority": "high",
      "category": "call_forwarding",
      "escalationReason": "Jio MMI error on SIM slot 0",
      "unreadMessages": 2,
      "createdAt": "2026-08-25T14:40:00Z"
    }
  ]
}
```

### `POST /api/v1/admin/support/chats/{id}/claim`
Human support agent claims the chat and enters the room.

### `POST /api/v1/admin/support/chats/{id}/messages`
Sends live human support response or canned snippet.

### `POST /api/v1/admin/support/chats/{id}/internal-notes`
Adds staff-only private note for audit and handoff.

### `POST /api/v1/admin/support/chats/{id}/resolve`
Resolves ticket and triggers 5-star CSAT survey on the client device.

### `GET /api/v1/admin/support/canned-responses`
Fetches quick reply templates (e.g. `/jio-mmi`, `/refund-policy`, `/takeover-help`).

---

## 📄 3. Dynamic CMS & Legal Content Management APIs

### `GET /api/v1/admin/content/:slug`
Fetches full editable sections and raw markdown for any dynamic legal document (`privacy-policy`, `terms-and-conditions`, `about-us`).

### `PUT /api/v1/admin/content/:slug`
Saves and publishes new version of legal documents. Invalidate Redis cache `<2ms` instantly.

### `POST /api/v1/admin/content/faqs`
Creates a new categorized FAQ with tags and display ordering.

---

## ⚙️ 4. Universal Client Remote Config APIs

### `GET /api/v1/admin/config/client-settings`
Fetches global feature flags, maintenance windows, forced update rules, and payment gateways.

### `PUT /api/v1/admin/config/client-settings`
Updates client settings and broadcasts zero-restart `config:updated` WebSocket event.
