# 🏗️ Study 08: 24-Microservice Master Architecture, Implementation Roadmap & Interconnection Blueprint

> **Classification:** Engineering Specification & Implementation Study  
> **Status:** Active Reference & Standard  
> **Scope:** 24 Backend Microservices + 14 Foundation Packages + 12 Infrastructure Engines + 21 Autonomous AI Agents  
> **Target Ports:** `:8080` through `:8106`, `:5060`, `:50051`  

---

## 🧭 Executive Summary

The **Nexora AI Telephony & Communication Platform** is architected as a decoupled, resilient, polyglot microservice mesh composed of **24 specialized microservices** across Node.js/TypeScript, Go, and Python. Every service implements clean architecture (`controllers`, `middlewares`, `models`, `routes`, `services`, `utils`), consumes shared foundation packages (`@nexora/*` in TypeScript, `packages/go/nexora` in Go, `packages/python/nexora_core` in Python), and interfaces with polyglot database engines (PostgreSQL 16, MongoDB 6.0, Redis 7.2, MinIO S3, NATS JetStream, Kafka, OpenSearch, Temporal).

---

## 📑 Master 24-Microservice Catalog & Port Matrix

| Phase | # | Microservice Directory | Port | Language | Primary Role & Specialization | Data Layer & Infrastructure | Shared Packages |
| :---: | :---: | :--- | :---: | :---: | :--- | :--- | :--- |
| **P1** | 01 | [`services/auth-service`](file:///d:/trideocoders/ai-telephony-platform/services/auth-service) | `:8081` | Node.js | JWT token issuance, Bcrypt hashing, Phone OTP, API keys, RBAC | Redis TTL + PostgreSQL | `@nexora/shared`, `@nexora/logger`, `@nexora/types` |
| **P1** | 02 | [`services/user-service`](file:///d:/trideocoders/ai-telephony-platform/services/user-service) | `:8082` | Node.js | Tenant user profiles, DID assignments, team roles & caller preferences | PostgreSQL 16 (`users`, `tenants`) | `@nexora/database`, `@nexora/types` |
| **P1** | 03 | [`services/api-gateway`](file:///d:/trideocoders/ai-telephony-platform/services/api-gateway) | `:8080` | Node.js | Universal reverse proxy, Spec 11 `/bootstrap`, rate limiting | Redis Token Bucket | `@nexora/config`, `@nexora/shared` |
| **P2** | 04 | [`services/telephony-service`](file:///d:/trideocoders/ai-telephony-platform/services/telephony-service) | `:8085` / `:5060` | Go | SIP trunking, LiveKit SIP ingress, dual-SIM carrier MMI forwarding | LiveKit SFU (`:7880`) | `packages/go/nexora` |
| **P2** | 05 | [`services/voice-service`](file:///d:/trideocoders/ai-telephony-platform/services/voice-service) | `:8092` | Python | Deepgram Nova-2 STT, Cartesia TTS, WebRTC bidirectional audio pipeline | LiveKit SFU (`:7880`) | `packages/python/nexora_core` |
| **P2** | 06 | [`services/call-service`](file:///d:/trideocoders/ai-telephony-platform/services/call-service) | `:8091` | Node.js | Call session state machine, outbound dialer queue, call routing | PostgreSQL + Redis Locks | `@nexora/events`, `@nexora/database` |
| **P2** | 07 | [`services/websocket-service`](file:///d:/trideocoders/ai-telephony-platform/services/websocket-service) | `:8093` | Node.js | Live softphone WebSocket feed, whisper audio queue, barge-in signals | Redis PubSub (`:6379`) | `@nexora/shared`, `@nexora/events` |
| **P3** | 08 | [`services/ai-orchestrator`](file:///d:/trideocoders/ai-telephony-platform/services/ai-orchestrator) | `:8086` / `:50051` | Go | gRPC multi-agent router, prompt compilation, MCP tool sandbox | HashiCorp Vault (`:8200`) | `@nexora/prompts`, `packages/go/nexora` |
| **P3** | 09 | [`services/memory-service`](file:///d:/trideocoders/ai-telephony-platform/services/memory-service) | `:8094` | Node.js | `pgvector` HNSW 1536-dim vector store, long-term caller episodic facts | PostgreSQL (`pgvector`) | `@nexora/ai`, `@nexora/database` |
| **P3** | 10 | [`services/spam-service`](file:///d:/trideocoders/ai-telephony-platform/services/spam-service) | `:8097` | Node.js | TRAI DND registry checks, STIR/SHAKEN caller verification, risk scoring | PostgreSQL + Redis Cache | `@nexora/shared`, `@nexora/types` |
| **P4** | 11 | [`services/transcript-service`](file:///d:/trideocoders/ai-telephony-platform/services/transcript-service) | `:8095` | Node.js | Diarized caller/agent transcripts, sentiment analysis, AI call summaries | MongoDB 6.0 (`transcripts`) | `@nexora/ai`, `@nexora/database` |
| **P4** | 12 | [`services/recording-service`](file:///d:/trideocoders/ai-telephony-platform/services/recording-service) | `:8098` | Node.js | Dual-channel call recording streamer, presigned playback URLs | MinIO S3 (`:9000`) | `@nexora/shared`, `@nexora/database` |
| **P4** | 13 | [`services/contact-service`](file:///d:/trideocoders/ai-telephony-platform/services/contact-service) | `:8096` | Node.js | Address book, contact grouping, caller ID enrichment, CRM sync | PostgreSQL 16 (`contacts`) | `@nexora/database`, `@nexora/types` |
| **P4** | 14 | [`services/storage-service`](file:///d:/trideocoders/ai-telephony-platform/services/storage-service) | `:8103` | Node.js | MinIO/S3 bucket lifecycle, presigned upload/download token generator | MinIO S3 API (`:9000`) | `@nexora/config`, `@nexora/shared` |
| **P5** | 15 | [`services/notification-service`](file:///d:/trideocoders/ai-telephony-platform/services/notification-service) | `:8090` | Node.js | Multi-channel dispatch (SMS 10DLC, WhatsApp Cloud API, FCM Push, APNs) | NATS JetStream (`:4222`) | `@nexora/notifications` |
| **P5** | 16 | [`services/workflow-service`](file:///d:/trideocoders/ai-telephony-platform/services/workflow-service) | `:8099` | Node.js | Temporal durable workflows, multi-step post-call triggers, Zapier hooks | Temporal Engine (`:7233`) | `@nexora/events`, `@nexora/shared` |
| **P5** | 17 | [`services/scheduler-service`](file:///d:/trideocoders/ai-telephony-platform/services/scheduler-service) | `:8100` | Node.js | BullMQ scheduled jobs, automated recurring reminders, dialer batches | Redis BullMQ (`:6379`) | `@nexora/events`, `@nexora/database` |
| **P6** | 18 | [`services/billing-service`](file:///d:/trideocoders/ai-telephony-platform/services/billing-service) | `:8083` | Node.js | Subscription plans, metered voice minutes usage, balance deductions | PostgreSQL 16 Ledgers | `@nexora/database`, `@nexora/types` |
| **P6** | 19 | [`services/payment-service`](file:///d:/trideocoders/ai-telephony-platform/services/payment-service) | `:8084` | Node.js | Stripe & Razorpay checkout sessions, webhook handlers, auto top-ups | HashiCorp Vault (`:8200`) | `@nexora/config`, `@nexora/shared` |
| **P7** | 20 | [`services/admin-service`](file:///d:/trideocoders/ai-telephony-platform/services/admin-service) | `:8101` | Node.js | Platform operations center, CMS legal docs, Live admin support chat | PostgreSQL (Migration `022`) | `@nexora/database`, `@nexora/shared` |
| **P7** | 21 | [`services/audit-service`](file:///d:/trideocoders/ai-telephony-platform/services/audit-service) | `:8102` | Node.js | Tamper-evident audit logs, SHA-256 hash chaining, compliance trails | PostgreSQL Immutable Logs | `@nexora/database`, `@nexora/types` |
| **P7** | 22 | [`services/enterprise-service`](file:///d:/trideocoders/ai-telephony-platform/services/enterprise-service) | `:8104` | Node.js | SAML 2.0 / OIDC enterprise SSO, SCIM 2.0 user sync, SLA enforcer | PostgreSQL + Vault | `@nexora/database`, `@nexora/types` |
| **P8** | 23 | [`services/search-service`](file:///d:/trideocoders/ai-telephony-platform/services/search-service) | `:8105` | Node.js | OpenSearch full-text & vector indexing for calls, contacts, knowledge | OpenSearch Engine (`:9200`) | `@nexora/shared`, `@nexora/logger` |
| **P8** | 24 | [`services/analytics-service`](file:///d:/trideocoders/ai-telephony-platform/services/analytics-service) | `:8106` | Node.js | Kafka stream consumer, call duration KPIs, time-series aggregations | Kafka KRaft Cluster (`:9092`) | `@nexora/events`, `@nexora/logger` |

---

## 🏛️ Standard Microservice Implementation Pattern

Every microservice in the monorepo adheres strictly to the **Nexora Clean Architecture Pattern**:

```text
services/<service-name>/
├── src/
│   ├── controllers/      ← HTTP & RPC Request Handlers (Request DTO validation, status code mapping)
│   ├── middlewares/      ← Interceptors (requireAuth, requireRole, requireScope, rateLimiter, requestLogger)
│   ├── models/           ← Polyglot Persistence Repositories (PostgreSQL / MongoDB / Redis / MinIO)
│   ├── routes/           ← Versioned REST API Route Definitions (/api/v1/...)
│   ├── services/         ← Core Domain Business Logic & External Client Adapters
│   ├── utils/            ← Envelopes, Hashers, Parsers, Sanitizers
│   └── index.ts          ← Lightweight bootstrap engine (< 20 lines) via @nexora/shared createMicroservice()
├── package.json          ← Scoped workspace package (@nexora/<service-name>)
└── tsconfig.json         ← Strict TypeScript compilation configuration
```

### Standard Response Envelope Format
```typescript
{
  "success": boolean,
  "statusCode": number,
  "data": any,
  "message": string,
  "pagination"?: {
    "page": number,
    "pageSize": number,
    "totalPages": number,
    "totalCount": number,
    "hasNextPage": boolean,
    "hasPrevPage": boolean
  }
}
```

---

## 🔄 Phased Execution Roadmap

### 🟢 Phase 1: IAM, Core User Management & Gateway ✅ *Phase Completed*
1. **`auth-service` (`:8081`)** ✅ *Completed*
   - MVC architecture implemented with controllers, middlewares, models, routes, services, and utils.
   - Dual authentication: Bearer JWT (`15m` / `7d`) + `x-api-key` header with SHA-256 verification.
   - Role-Based Access Control (`requireRole`) & Fine-Grained API Scoping (`requireScope`).
   - 6-digit Phone OTP with 5-minute Redis TTL and automatic cooldown.
   - 15/15 unit and integration tests passing.
2. **`user-service` (`:8082`)** ✅ *Completed*
   - MVC architecture implemented with controllers, middlewares, models, routes, services, and utils.
   - Tenant user lifecycle, account profile management, MinIO S3 presigned avatar uploads.
   - 4-mode screening preferences (`ai_answer_all`, `vip_priority_only`, `screen_unknown_only`, `do_not_disturb`).
   - Dual-SIM carrier forwarding engine (Jio, Airtel, Vi, BSNL, AT&T, Verizon).
   - Platform Admin model & CRUD APIs (`/api/v1/admins/*`) with department and RBAC role assignments.
   - 19/19 unit and integration tests passing.
3. **`api-gateway` (`:8080`)** ✅ *Completed*
   - Reverse proxy routing for all 24 downstream services with graceful fallbacks.
   - Redis 7.2 sliding-window rate limiting with atomic Lua script execution.
   - Universal Client Bootstrap & Remote Config (Spec 11 `/api/v1/bootstrap`).
   - Interactive OpenAPI 3.0 Swagger UI documentation for Telephony, Mobile, Admin, Desktop, and Enterprise.
   - Prometheus metrics instrumentation (`httpRequestsTotal`, `httpDurationMs`).
   - 20/20 unit and integration tests passing.

### 🔵 Phase 2: Telephony, SIP & Real-Time Audio Streaming
4. **`telephony-service` (`:8085` / `:5060`)** ✅ *Completed*
   - Go SIP stack, LiveKit SIP ingress, dual-SIM carrier MMI forwarding engine (`*004*`, `*67*`, `*61*`, `*62*`).
   - Redis active call store, Kafka `call.started`/`call.ended` producer, Twilio/Telnyx/Plivo signature validators.
   - 9/9 Go test suites passing.
5. **`voice-service` (`:8092`)** ✅ *Completed*
   - Python LiveKit WebRTC pipeline, Deepgram Nova-2 STT (<120ms), Cartesia Sonic TTS (<80ms), Barge-in interruption.
   - Sub-18ms in-call whisper injection queue (`WhisperService`) & Sub-15ms live call takeover barge-in mutex (`TakeoverService`).
   - Full modular architecture: `api/`, `config/`, `core/`, `hook/`, `plugins/`, `repository/`, `services/`, `utils/`.
   - 10/10 Python unit and integration tests passing.
6. **`call-service` (`:8091`)** ⏳ *Next Step*
   - Active call session state machine (`ringing`, `in-progress`, `barge-in`, `whisper`, `completed`), outbound campaign dialer.
7. **`websocket-service` (`:8093`)**
   - WebSocket softphone gateway, real-time audio waveform stream, in-call whisper queue (<18ms).

### 🟣 Phase 3: AI Orchestration, Context Memory & Telephony Protection
8. **`ai-orchestrator` (`:8086` / `:50051`)**
   - Go gRPC multi-agent brain, prompt compilation from `@nexora/prompts`, MCP tool execution sandbox.
9. **`memory-service` (`:8094`)**
   - PostgreSQL `pgvector` HNSW 1536-dim vector store, long-term caller episodic facts, memory retrieval.
10. **`spam-service` (`:8097`)**
    - India TRAI DND registry lookup, US STIR/SHAKEN caller verification, real-time spam risk scoring.

### 🟡 Phase 4: Call Artifacts, Recordings, Transcripts & Contacts
11. **`transcript-service` (`:8095`)**
    - Diarized transcripts (Caller vs AI Agent), sentiment scoring, AI action item generation.
12. **`recording-service` (`:8098`)**
    - Dual-channel call recording streamer, MinIO S3 storage, presigned audio playback tokens.
13. **`contact-service` (`:8096`)**
    - Address book, contact grouping, caller ID enrichment, CRM sync.
14. **`storage-service` (`:8103`)**
    - S3 bucket lifecycle management, presigned upload/download token generator.

### 🟠 Phase 5: Multi-Channel Messaging, Automations & Scheduling
15. **`notification-service` (`:8090`)**
    - Multi-channel notification SDK (Twilio SMS, Meta WhatsApp Cloud API, FCM Push, APNs).
16. **`workflow-service` (`:8099`)**
    - Temporal durable execution, multi-step post-call automations, Webhook triggers.
17. **`scheduler-service` (`:8100`)**
    - BullMQ job scheduler, automated voice reminders, scheduled campaign dialers.

### 🔴 Phase 6: Commerce, Billing & Payment Processing
18. **`billing-service` (`:8083`)**
    - Subscription plans, metered voice minutes usage, wallet balances, immutable accounting ledgers.
19. **`payment-service` (`:8084`)**
    - Stripe & Razorpay checkout sessions, webhook handlers, automatic wallet top-ups.

### 🛡️ Phase 7: Platform Administration, Security & Enterprise
20. **`admin-service` (`:8101`)**
    - Platform operations center, CMS legal documents, Live admin support chat.
21. **`audit-service` (`:8102`)**
    - Tamper-evident audit logs, SHA-256 hash chaining, compliance trails.
22. **`enterprise-service` (`:8104`)**
    - Enterprise SAML 2.0 / OIDC SSO, SCIM 2.0 user sync, SLA policy monitoring.

### 📊 Phase 8: Search, Real-Time Analytics & Event Streaming
23. **`search-service` (`:8105`)**
    - OpenSearch full-text & vector indexing for calls, contacts, knowledge base.
24. **`analytics-service` (`:8106`)**
    - Kafka stream consumer, call duration KPIs, time-series aggregations.

---

## 🧪 Unified Quality & Test Standards

For every service implemented:
1. Must compile with `0` errors via `pnpm run type-check`.
2. Must include a dedicated test file under `tests/services/<service-name>/`.
3. Must pass 100% of unit and HTTP integration tests via `pnpm run test:all`.
4. Must expose `/health`, `/live`, `/ready`, and Prometheus `/metrics` endpoints.
