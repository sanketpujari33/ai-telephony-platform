# 🛡️ Admin Model, RBAC Matrix & Permission Architecture Deep Study

This document details the architectural design of the Nexora Admin Model, 12-Role RBAC hierarchy, recursive sub-module permissions, two-tier Redis permission caching, and real-time token invalidation.

---

## 1. System Roles & RBAC Hierarchy

Nexora enforces a multi-tiered 12-role Role-Based Access Control (RBAC) model across the platform:

```text
┌────────────────────────────────────────────────────────────────────────┐
│                         SUPER_ADMIN (Level 100)                        │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
        ┌───────────────────────────┴───────────────────────────┐
        │                                                       │
┌───────▼────────┐                                      ┌───────▼────────┐
│ SYSTEM_ADMIN   │                                      │ SECURITY_ADMIN │
│ (Level 90)     │                                      │ (Level 90)     │
└───────┬────────┘                                      └───────┬────────┘
        │                                                       │
        ├──────────────────────┬──────────────────────┐         │
┌───────▼────────┐    ┌────────▼───────┐    ┌─────────▼──────┐  │
│ BILLING_ADMIN  │    │ COMPLIANCE_OFF │    │ DEVELOPER      │  │
│ (Level 80)     │    │ (Level 80)     │    │ (Level 70)     │  │
└────────────────┘    └────────────────┘    └────────────────┘  │
                                                                │
┌───────────────────────────────────────────────────────────────┴────────┐
│                         TENANT_ADMIN (Level 60)                        │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
        ┌───────────────────────────┼───────────────────────────┐
┌───────▼────────┐          ┌───────▼────────┐          ┌───────▼────────┐
│ AGENT_MANAGER  │          │ AGENT_OPERATOR │          │ VIEWER         │
│ (Level 50)     │          │ (Level 40)     │          │ (Level 20)     │
└────────────────┘          └────────────────┘          └────────────────┘
```

---

## 2. Granular Permission Scopes

Each role maps to explicit string scopes following the `resource:action` convention:

| Domain | Permission Scope | Description |
| :--- | :--- | :--- |
| **Banners** | `admin:banners:read`, `admin:banners:write`, `admin:banners:delete` | Manage marketing and system announcement banners |
| **Promoters** | `admin:promoters:read`, `admin:promoters:payout`, `admin:promoters:approve` | Affiliate / promoter lifecycle and revenue share |
| **Funds** | `admin:funds:receive`, `admin:funds:refund`, `admin:funds:withdraw` | Manual ledger top-ups, dispute refunds, withdrawals |
| **P&L Margins**| `admin:pnl:read`, `admin:pnl:export`, `admin:pnl:configure_rates` | Telephony carrier cost vs retail margin analytics |
| **Dead Reports**| `admin:dead_reports:read`, `admin:dead_reports:resolve` | Inactive campaigns, disconnected SIP trunks, dead-letter jobs |
| **Data Monitor**| `admin:datamonitor:read`, `admin:datamonitor:keywords_manage` | Real-time spam keywords, PII redaction rules, compliance flags |
| **Feedback** | `admin:feedback:read`, `admin:feedback:reply` | End-user NPS, CSAT ratings, agent call reviews |
| **General Tabs**| `admin:generaltabs:read`, `admin:generaltabs:write` | System settings, portal navigation links, feature flags |

---

## 3. High-Performance Two-Tier Permission Caching

To prevent database bottlenecking on every HTTP/WebSocket request, user permissions and token blacklists are cached across two tiers:

```text
┌─────────────────────────────────────────────────────────┐
│              Incoming HTTP Request + JWT                │
└────────────────────────────┬────────────────────────────┘
                             │
                             ▼
              ┌─────────────────────────────┐
              │ Tier 1: Local In-Memory LRU │  Hit (<1ms) ──► Allow / Deny
              │ (Node.js lru-cache / 60s)   │
              └──────────────┬──────────────┘
                             │ Miss
                             ▼
              ┌─────────────────────────────┐
              │ Tier 2: Redis Cluster 7.2   │  Hit (<3ms) ──► Populate L1
              │ Key: `rbac:perms:{userId}`  │
              └──────────────┬──────────────┘
                             │ Miss
                             ▼
              ┌─────────────────────────────┐
              │ Tier 3: PostgreSQL 16 DB    │  Hit (~15ms) ─► Populate L2 + L1
              │ `user_roles` JOIN `roles`   │
              └─────────────────────────────┘
```

---

## 4. Real-time Permission Revocation & Invalidation

When an administrator updates a role or revokes user permissions:
1. `admin-service` updates PostgreSQL tables within an ACID transaction.
2. An invalidation event is emitted to NATS JetStream: `auth.rbac.invalidated`.
3. All service instances delete local in-memory LRU keys and flush `rbac:perms:{userId}` in Redis.
4. Active WebSocket connections for the affected user immediately re-evaluate permissions or terminate session gracefully.
