# 🧠 Deep Study — Nexora AI Communication & Telephony Platform

> **Platform Name:** Nexora AI  
> **Codebase Root:** `d:\trideocoders\ai-telephony-platform`  
> **Study Date:** 2026-08-22  
> **Version:** 0.1.0 (active development)

---

## 1. 🎯 Executive Summary

Nexora is an **enterprise-grade, AI-first telephony operating system** — a polyglot monorepo that unifies:
- Real-time voice AI calls (SIP/WebRTC)
- A multi-agent AI ecosystem (20 specialized agents)
- Omni-channel messaging (WhatsApp, SMS, Email)
- Cross-platform clients (Flutter mobile/desktop + Next.js web)
- A 24-microservice backend

The platform is still in **early active development** (v0.1.0, Sprint 1–6 MVP in progress). Core telephony and AI pipeline are functional; most services are scaffolded with the full MVC structure but not yet production-complete.

---

## 2. 🏗️ Monorepo Structure & Build System

### Build Orchestration
```
Package Manager : pnpm@9.0.0
Build System    : Turborepo v2.3.0
Languages       : TypeScript, Go 1.22, Python 3.11, Dart 3.x
```

### Workspace Layout
```
nexora-ai/
├── apps/          (6 frontend apps)
├── services/      (24 backend microservices)
├── agents/        (20 Python AI agents)
├── packages/      (12 shared libraries)
├── infra/         (Docker, K8s, Terraform, Monitoring, Nginx)
├── docs/          (architecture, sprints, roadmap, specs)
├── tests/         (integration & load tests, mirrors services/)
└── scripts/       (setup, start, stop, seed, healthcheck, build)
```

### Turborepo Task Graph
| Task | Depends On | Inputs |
|---|---|---|
| `build` | `^build` (dependencies first) | `src/**`, `package.json`, `tsconfig.json` |
| `type-check` | `^build` | `src/**`, `tsconfig.json` |
| `lint` | `^lint` | `src/**`, `.eslintrc*` |
| `test` | `^build` | `src/**`, `tests/**` |
| `dev` | — | (persistent, no cache) |

### Go Workspace
```go
// go.work — Go 1.22 multi-module workspace
use (
    ./services/telephony-service
    ./services/ai-orchestrator
    ./tests/services/telephony-service
)
```

---

## 3. 📱 Client Applications (`apps/`)

| App | Tech | Port | Status |
|---|---|---|---|
| `user-dashboard` | Next.js 14 + Tailwind CSS | 3000 | Scaffolded (basic page) |
| `admin-panel` | Next.js 14 | 3001 | Scaffolded |
| `landing-page` | Next.js 14 | 3002 | Scaffolded |
| `documentation` | MDX / Docs | 3003 | Planned |
| `mobile-app` | Flutter 3.x / Dart | Native | Single `main.dart` file |
| `desktop` | Flutter 3.x / Dart | Native | Scaffolded |

### Flutter Mobile (`apps/mobile-app`)
- Full cross-platform setup: Android, iOS, Linux, macOS, Windows, Web
- Single `lib/main.dart` — 5.1 KB entry point
- Dependencies in `pubspec.yaml`
- Likely includes the WebRTC dialer, call UI, and agent chat screen

### Next.js Apps
- Use `@nexora/ui` for shared design system
- Use `@nexora/types` for isomorphic TypeScript types
- All apps share `tsconfig.base.json` from root

---

## 4. ⚙️ Backend Microservices (`services/`)

### Service Inventory (24 services)

| # | Service | Language | Port | Responsibility |
|---|---|---|---|---|
| 1 | `api-gateway` | Node.js/TS | 8080 | Central entry point, routing, auth middleware |
| 2 | `auth-service` | Node.js/TS | 8082 | JWT, OTP, OAuth, refresh tokens |
| 3 | `user-service` | Node.js/TS | 4003 | User profiles, RBAC |
| 4 | `call-service` | Node.js/TS | 4002 | Call records, session management |
| 5 | `telephony-service` | **Go 1.22** | 8085 | SIP trunk, LiveKit room mgmt, inbound/outbound calls |
| 6 | `ai-orchestrator` | **Go 1.22** | 8086 | Routes user intent to Python agent ecosystem |
| 7 | `voice-service` | **Python 3.11** | — | STT/TTS pipeline, LiveKit agent sessions |
| 8 | `memory-service` | Node.js/TS | 4013 | User context, facts storage |
| 9 | `billing-service` | Node.js/TS | 8089 | Usage metering, Stripe checkout |
| 10 | `payment-service` | Node.js/TS | 4010 | Payment processing |
| 11 | `notification-service` | Node.js/TS | 4005 | Push/SMS/email notifications |
| 12 | `analytics-service` | Node.js/TS | 4008 | Metrics, call analytics |
| 13 | `contact-service` | Node.js/TS | 4006 | Address book, contact management |
| 14 | `spam-service` | Node.js/TS | 4012 | Call spam detection, Truecaller integration |
| 15 | `scheduler-service` | Node.js/TS | 4018 | Appointment scheduling |
| 16 | `recording-service` | Node.js/TS | 4017 | Call recording storage → MinIO |
| 17 | `transcript-service` | Node.js/TS | 4007 | STT results → OpenSearch |
| 18 | `search-service` | Node.js/TS | 4019 | Full-text search over transcripts |
| 19 | `storage-service` | Node.js/TS | 4011 | S3/MinIO abstraction |
| 20 | `audit-service` | Node.js/TS | 4015 | Audit log trail |
| 21 | `admin-service` | Node.js/TS | 4014 | Admin panel APIs |
| 22 | `enterprise-service` | Node.js/TS | 4016 | Enterprise tenant features |
| 23 | `websocket-service` | Node.js/TS | 4009 | Real-time WS push to clients |
| 24 | `workflow-service` | Node.js/TS | 4020 | Temporal workflow client |

### Node.js Services — Common Architecture Pattern
All Node.js services follow a consistent **MVC pattern**:
```
src/
├── index.ts          (Express app bootstrap, PORT config)
├── routes/           (Route definitions)
│   ├── index.ts      (Route aggregator)
│   ├── *.routes.ts   (Feature route files)
│   └── health.routes.ts
├── controllers/      (Request handlers)
├── services/         (Business logic)
├── models/           (DB models / schemas)
├── middlewares/      (auth, logging, error)
└── utils/            (Helpers)
```

### API Gateway — Key Details
- Express + CORS + JSON body parser
- `@nexora/logger` for structured logging
- `optionalAuthMiddleware` (JWT validation)
- `requestLoggerMiddleware` (request tracing)
- `errorHandlerMiddleware` (global error handling)
- Routes: `/api/v1/calls`, `/api/v1/agents`, `/health`

### Auth Service — Key Details
- Routes: `/auth/*` (login, register, refresh), `/otp/*` (send, verify), `/health`
- Controllers: `auth.controller.ts`, `otp.controller.ts`, `health.controller.ts`
- Port: `8082`

### Telephony Service (Go) — Key Details
The **most complete** Go service:
- `pkg/livekit/client.go` — LiveKit SDK client wrapper
- On startup: creates SIP Trunk + SIP Dispatch Rule if env vars set
- `POST /api/v1/calls/inbound` — Creates LiveKit room, issues SIP+Agent JWT tokens
- `POST /api/v1/calls/outbound` — Same flow for outbound calls
- `GET /api/v1/calls` — Lists active in-memory call sessions
- Active calls stored in-memory `map[string]*CallSession` (not yet persisted)
- Uses `gin`, `livekit/protocol`, `google/uuid`, `rs/zerolog`

### AI Orchestrator (Go) — Key Details
- `POST /api/v1/orchestrate` — Receives `{phone_number, user_input, session_id}`
- Calls Python orchestrator agent at `ORCHESTRATOR_AGENT_URL` (default: `http://localhost:8087/api/v1/evaluate`)
- Returns `{selected_agent, agent_response, tools_executed, confidence, memory_updated}`
- Port: `8086`

### Voice Service (Python) — Key Details
- Full Python service with structured layout: `api/`, `core/`, `services/`, `repository/`, `hook/`, `utils/`, `config/`
- `core/agent_session.py` — LiveKit AgentSession wrapper
- `agent.py` — Entry point for LiveKit voice agent worker
- Requirements: FastAPI, livekit-agents, deepgram, cartesia, openai

### Billing Service — Key Details (Functional MVP)
- In-memory `Map<tenantId, UsageRecord>`
- `POST /api/v1/billing/meter` — Meters call seconds + tokens at `$0.02/min + $0.001/1k tokens`
- `GET /api/v1/billing/usage/:tenantId` — Returns usage summary
- `POST /api/v1/billing/checkout` — Returns mock Stripe checkout URL
- Port: `8089`

---

## 5. 🤖 AI Agent Ecosystem (`agents/`)

### 20 Specialized Agents (all Python)
```
orchestrator-agent  phone-agent         calendar-agent
appointment-agent   support-agent       email-agent
whatsapp-agent      sms-agent           memory-agent
reminder-agent      research-agent      browser-agent
finance-agent       shopping-agent      travel-agent
coding-agent        enterprise-agent    automation-agent
knowledge-agent     delivery-agent
```

### Orchestrator Agent (`agents/orchestrator-agent/agent.py`)
- **FastAPI** server on port `8087`
- Uses **OpenAI Structured Outputs** with Pydantic schema `OrchestratorDecision`
- Model: `gpt-4o` with `temperature=0.0`
- Routes to: `calendar-agent`, `support-agent`, `phone-agent`, `finance-agent`
- Mock fallback when `OPENAI_API_KEY=mock-key` (keyword-based routing)
- Endpoint: `POST /api/v1/evaluate` → `{selected_agent, confidence, required_tools, agent_response}`

### Phone Agent (`agents/phone-agent/agent.py`)
- **LiveKit Agents SDK** — runs as a LiveKit worker
- **STT**: Deepgram `nova-2-general` with `interim_results=True`, 500ms endpointing
- **LLM**: OpenAI `gpt-4o`
- **TTS**: Cartesia `sonic-english` (WebSocket streaming for low TTFAB)
- **VAD**: Silero (`min_silence_duration=0.1s`, `min_speech_duration=0.05s`)
- **Metrics**: `MetricsCollector` for TTFAB tracking
- System Prompt: "Handle incoming/outgoing calls, identify caller, resolve inquiries, schedule appointments, warm transfers"

### AI Data Flow (Call Path)
```
Caller → SIP/WebRTC → Telephony Service (Go:8085)
    → LiveKit Room Created (ws://localhost:7880)
    → Phone Agent dispatched (Python LiveKit Worker)
    → Silero VAD → Deepgram STT → LLM (gpt-4o/groq)
    → Cartesia TTS → LiveKit → Caller Audio
    → [if complex intent] → AI Orchestrator (Go:8086)
    → Orchestrator Agent (Python:8087) → Specialized Agent
```

---

## 6. 📦 Shared Packages (`packages/`)

| Package | Name | Key Exports |
|---|---|---|
| `ai` | `@nexora/ai` | `ModelRouter`, `AIClient`, provider adapters (13 LLMs) |
| `config` | `@nexora/config` | Env + Vault config loader |
| `database` | `@nexora/database` | PostgreSQL, MongoDB, Redis clients; SQL migrator |
| `events` | `@nexora/events` | NATS JetStream + Redis pub/sub event bus |
| `logger` | `@nexora/logger` | Structured JSON logger with PII redaction |
| `notifications` | `@nexora/notifications` | Push/SMS/email notification sender |
| `prompts` | `@nexora/prompts` | System prompt templates |
| `sdk` | `@nexora/sdk` | Public API SDK for third-party developers |
| `shared` | `@nexora/shared` | Utility functions (crypto, billing math, phone, time) |
| `tools` | `@nexora/tools` | Tool registry for agents |
| `types` | `@nexora/types` | Shared TypeScript interfaces & enums |
| `ui` | `@nexora/ui` | React component library |

### `@nexora/ai` — Multi-Provider LLM Router
- **13 LLM provider adapters**: OpenAI, Anthropic, Gemini, Groq, Mistral, DeepSeek, Cohere, Cerebras, Fireworks, Together, Ollama, Perplexity, xAI
- **`ModelRouter`**: Executes completions with automatic fallback chain
  - Default provider: `groq`
  - Fallback chain: `[primary, openai, anthropic, gemini]`
  - Retries next provider on any error
- **`AIClient`**: Unified interface regardless of provider

### `@nexora/database` — Database Abstraction Layer
- **PostgreSQL** (`postgres.ts`): pg-based client, connection pooling
- **MongoDB** (`mongodb.ts`): Mongoose-based client
- **Redis** (`redis.ts`): ioredis client, 6.7KB — rate limiting, sessions, OTPs, call state, Redlock
- **Repository** (`repository.ts`): Base CRUD repository pattern
- **Migrator** (`migrator.ts`): SQL migration runner
- **Health** (`health.ts`): Health check pings for all DB connections

### `@nexora/events` — Event Bus
- **NATS JetStream** (`nats.ts`): Primary event bus for high-throughput events
- **Redis Pub/Sub** (`redis.ts`): Lightweight in-process events
- **DLQ** (`dlq.ts`): Dead-letter queue for failed events
- **Middleware** (`middleware.ts`): Event interceptors
- **Memory bus** (`memory.ts`): In-process event bus for testing

### `@nexora/logger` — Structured Logger
- 9.9 KB main file — full-featured pino-based logger
- `redact.ts` — PII redaction (credit cards, SSNs, passwords, JWTs)
- `serializers.ts` — Request/response serialization
- `middleware.ts` — Express middleware integration

### `@nexora/shared` — Utility Library
- `crypto.ts` — Hashing, token generation
- `billing.ts` — Billing math helpers
- `phone.ts` — Phone number validation/formatting (E.164)
- `time.ts` — Time zone utilities
- `strings.ts` — String manipulation
- `collections.ts` — Array/map utilities
- `result.ts` — Result/Either type for error handling
- `async.ts` — Async utilities, retry logic

### `@nexora/types` — Type System
Service-specific TypeScript interfaces:
- `auth.ts` — JWT payloads, OTP types, auth responses
- `call.ts` — Call direction, status, session types
- `telephony.ts` — SIP trunk, LiveKit room types
- `voice.ts` — STT/TTS provider types
- `orchestrator.ts` — Agent routing types
- `gateway.ts` — API Gateway request/response types
- `analytics.ts` — Analytics event types

---

## 7. 💾 Database Architecture

### PostgreSQL 16 + pgvector (`:5432`)
Primary relational database with **multi-tenancy from day 1**:

```sql
-- Every table has tenant_id UUID NOT NULL
tenants    → plan, billing_email, settings (JSONB)
users      → RBAC roles, phone_number, email (UNIQUE per tenant)
agents     → LLM config, system_prompt, voice_config (JSONB), tools_allowed[]
calls      → direction, status, room_name, spam_score, duration, recording_url
memories   → user facts, tags[], source_channel, confidence, expires_at
embeddings → VECTOR(1536) with HNSW index (cosine similarity)
```

Extensions: `uuid-ossp`, `vector` (pgvector >= 0.5.0)

### Redis 7.2 (`:6379`)
Key-space design:
```
nexora:ratelimit:{tenantId}:{ip}      (Token bucket)
nexora:session:{userId}:{sessionId}   (JWT sessions)
nexora:otp:{phoneNumber}              (5-min OTP codes)
nexora:call:active:{callId}           (Real-time call state)
nexora:lock:call:{callId}             (Distributed mutex / Redlock)
```

### MongoDB 6.0 (`:27017`)
Semi-structured AI memory:
```
conversation_memories   (context, sentiment history per session)
knowledge_chunks        (knowledge base documents for RAG)
```

### OpenSearch 2.13 (`:9200`)
Full-text search:
```
call-transcripts-v1   (BM25 search + nested utterance timestamps)
```

### MinIO S3 (`:9000`)
Object storage:
```
/recordings/{tenantId}/{callId}.wav    (Raw audio recordings)
/waveforms/{callId}.json               (Pre-computed audio waveforms)
```

### HashiCorp Vault (`:8200`)
Zero-secret management:
```
secret/data/nexora/{env}    (All API keys, DB passwords, tokens)
```

---

## 8. 🏛️ Infrastructure (`infra/`)

### Docker Compose — Dev Stack (`infra/docker/docker-compose.dev.yml`)
10 containers managed as `nexora-dev-net` bridge network:
| Container | Image | Purpose |
|---|---|---|
| `nexora-postgres` | `pgvector/pgvector:pg16` | Primary DB + vector search |
| `nexora-redis` | `redis:7.2-alpine` | Cache + sessions |
| `nexora-minio` | `minio/minio:latest` | S3-compatible storage |
| `nexora-opensearch` | `opensearchproject/opensearch:2.13.0` | Full-text search |
| `nexora-nats` | `nats:2.10-alpine` | Event bus (JetStream) |
| `nexora-temporal` | `temporalio/auto-setup:1.24` | Durable workflow engine |
| `nexora-temporal-ui` | `temporalio/ui:latest` | Workflow dashboard |
| `nexora-livekit` | `livekit/livekit-server:v1.7` | WebRTC SFU |
| `nexora-vault` | `hashicorp/vault:1.16` | Secrets management |
| `nexora-mongodb` | `mongo:6.0` | Document DB |
| `nexora-kafka` | `apache/kafka:latest` | Event streaming |

### Kubernetes (`infra/kubernetes/`)
```
namespaces.yaml         (Namespace definitions)
base/
├── deployment.yaml     (Main deployment spec)
├── service.yaml        (Service definitions)
├── ingress.yaml        (Ingress routing)
├── kafka.yaml          (Kafka StatefulSet)
├── mongodb.yaml        (MongoDB StatefulSet)
└── kustomization.yaml  (Kustomize base)
overlays/               (Dev/staging/prod overlays)
```

### Terraform (`infra/terraform/`)
```
main.tf         (GCP/AWS provider, VPC, GKE/EKS cluster)
iam.tf          (IAM roles & policies)
kafka.tf        (Managed Kafka config)
mongodb.tf      (MongoDB Atlas config)
secrets.tf      (Vault/SM secrets)
variables.tf    (Environment variables)
outputs.tf      (Exported values)
environments/   (Dev/staging/prod tfvars)
```

### Monitoring (`infra/monitoring/`)
```
prometheus/     (prometheus.yml, alert rules)
grafana/        (dashboards, datasources)
loki/           (log aggregation config)
tempo/          (distributed tracing config)
otel/           (OpenTelemetry collector config)
```

---

## 9. 🔄 CI/CD Pipeline (`.github/workflows/`)

### `ci.yml` — Continuous Integration
- Runs on: PR to main/develop
- Steps: Install → Type-check → Lint → Test (Node + Go + Python)
- Parallelized via Turborepo

### `cd.yml` — Continuous Deployment
- Runs on: Push to main
- Steps: Build → Docker build → Push to registry → Deploy to K8s

### `release.yml` — Release Pipeline
- Triggered manually or on tags
- Runs `build-release.mjs` → generates release manifests in `.github/release/`

---

## 10. 📜 Scripts System (`scripts/`)

Cross-platform scripts (auto-detect Windows/macOS/Linux):

| Script | Function |
|---|---|
| `scripts/server/auto/setup.mjs` | Initialize Docker DBs, install deps, build packages, seed DB |
| `scripts/server/auto/start.mjs` | Launch all containers + microservices + web portals |
| `scripts/server/auto/stop.mjs` | Graceful shutdown of all processes |
| `scripts/server/auto/seed.mjs` | Seed: 1 Tenant, 1 SuperAdmin, 12 RBAC Roles, 20 AI Agent profiles |
| `scripts/server/health/healthcheck.mjs` | TCP/HTTP ping for all 17 ports |
| `scripts/build/build-release.mjs` | Compile + write release manifests |

OS-specific scripts also exist under `scripts/server/win/`, `scripts/server/mac/`, `scripts/server/ubuntu/`.

---

## 11. 📋 RBAC System

12-role matrix:
```
SUPERADMIN          → Full platform access
PLATFORM_ADMIN      → Platform management
MARKETING_ADMIN     → Marketing content
SUPPORT_LEAD        → Team lead for support
CUSTOMER_SUPPORT    → Handle customer tickets
BOT_SUPPORT         → Manage AI bots
DEVELOPER           → API access, testing
BILLING_ADMIN       → Billing management
COMPLIANCE_AUDITOR  → Audit logs only
TENANT_ADMIN        → Manage their tenant
AGENT_OPERATOR      → Configure agents
VIEWER              → Read-only access
```

---

## 12. 🔌 External Integrations

### LLM Providers (13)
OpenAI, Anthropic (Claude), Google Gemini, Groq, Mistral, DeepSeek, Cohere, Cerebras, Fireworks AI, Together AI, Ollama (self-hosted), Perplexity, xAI (Grok)

### STT Providers
Deepgram Nova-2 (primary), Azure Speech, Whisper, Sarvam AI

### TTS Providers
Cartesia Sonic (primary, WebSocket streaming), ElevenLabs, Azure Speech

### Telephony Providers
Twilio, LiveKit (WebRTC/SIP), planned: Exotel, Knowlarity, Plivo, Vonage

### Spam Detection
Truecaller API (server-side, not SDK)

### Payment
Stripe (checkout sessions, subscription billing)

### WhatsApp
WhatsApp Business API (Meta Cloud API)

---

## 13. 🌊 Real-Time Architecture

### Call Flow (Step by Step)
1. Caller dials → PSTN/SIP → **Telephony Service** (Go:8085)
2. Go creates LiveKit room → issues `sip_token` + `agent_token`
3. **Phone Agent** (Python) dispatches as LiveKit worker, joins room
4. Silero VAD detects speech → Deepgram streams audio → Final transcript
5. Agent queries **Memory Service** (MongoDB) for caller context
6. LLM (Groq/GPT-4o) generates response (streamed)
7. Cartesia TTS streams audio back to LiveKit room → Caller hears response
8. Caller interrupts → Silero VAD fires → TTS stream cancelled in <30ms
9. [Complex intent] → **AI Orchestrator** (Go:8086) routes to specialized Python agent
10. Tool execution → Memory update → Call ends → Recording + Transcript stored

### Latency Targets (Sprint 6)
- **TTFAB** (Time to First Audio Byte): < 800ms
- **Barge-in latency** (interruption): < 30ms
- **VAD speech detection**: < 50ms (`min_speech_duration=0.05s`)

---

## 14. 📅 Development Roadmap (Sprint Plan)

| Sprint Range | Focus |
|---|---|
| Sprint 01–02 | Foundations (monorepo, CI/CD, Docker infra, shared packages) |
| Sprint 03–06 | **MVP Core Loop** (telephony, AI pipeline, voice agent, auth) |
| Sprint 07–09 | MVP Surfaces (user dashboard, admin panel, mobile app) |
| Sprint 10–12 | Messaging (WhatsApp, SMS, Email agents) |
| Sprint 13–15 | Admin, Analytics & Billing |
| Sprint 16–20 | Specialized Agents (all 20) |
| Sprint 21–23 | Desktop App & Enterprise Portal |
| Sprint 24–26 | Security Hardening (pen-test, compliance) |
| Sprint 27–28 | Launch Readiness |

---

## 15. 🔑 Environment Variables (Key Categories)

```env
# Infrastructure
DATABASE_URL     = postgresql://nexora:nexorapassword@localhost:5432/nexoradb
REDIS_URL        = redis://localhost:6379
NATS_URL         = nats://localhost:4222
TEMPORAL_ADDRESS = localhost:7233
LIVEKIT_URL      = ws://localhost:7880
LIVEKIT_API_KEY  = devkey
VAULT_ADDR       = http://localhost:8200

# AI/LLM
OPENAI_API_KEY   = (required for real calls)
GROQ_API_KEY     = (fastest, default)
ANTHROPIC_API_KEY, GEMINI_API_KEY, etc.

# STT/TTS
DEEPGRAM_API_KEY, CARTESIA_API_KEY, ELEVENLABS_API_KEY

# Telephony
TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_PHONE_NUMBER

# Feature Flags
FEATURE_WHATSAPP_ENABLED=false
FEATURE_SMS_ENABLED=false
FEATURE_EMAIL_ENABLED=false
FEATURE_SPAM_CHECK_ENABLED=true
FEATURE_RECORDING_ENABLED=true
```

---

## 16. ⚠️ Current Development State Assessment

### ✅ Working / Complete
- Monorepo setup (pnpm + Turborepo)
- Docker dev infrastructure (all 11 containers)
- Go telephony service (SIP + LiveKit integration)
- Go AI orchestrator (intent routing bridge)
- Python phone agent (full voice pipeline: VAD+STT+LLM+TTS)
- Python orchestrator agent (LLM-based routing with structured outputs)
- Shared packages: `@nexora/ai`, `@nexora/database`, `@nexora/logger`, `@nexora/events`, `@nexora/shared`, `@nexora/types`
- Initial PostgreSQL schema (multi-tenant, pgvector)
- CI/CD pipeline (ci.yml, cd.yml, release.yml)
- Billing service (in-memory MVP, usage metering)
- Scripts system (cross-platform setup/start/stop/seed/healthcheck)

### 🔨 In Progress / Scaffolded
- Auth service (structure exists, controllers scaffolded)
- User service (src/index.ts only)
- All other Node.js services (MVC structure exists, business logic pending)
- Flutter mobile app (`lib/main.dart` only)
- Next.js apps (basic page.tsx, no feature UI)
- Kubernetes manifests (base + kustomize structure)
- Terraform configs (structure exists)

### 📋 Planned / Empty
- Most specialized agents (directory exists, no code)
- `@nexora/ui` component library
- `@nexora/sdk` public SDK
- `@nexora/tools` MCP tool registry
- `@nexora/prompts` prompt management
- Desktop app (Flutter shell only)
- All messaging agents (WhatsApp, SMS, Email)

---

## 17. 🏆 Architectural Strengths

1. **Polyglot by design** — each language chosen for its domain strengths
2. **Multi-tenancy from day 1** — `tenant_id` on every DB table
3. **AI provider abstraction** — 13 providers with automatic fallback chain
4. **Zero-secret commits** — HashiCorp Vault for all credentials
5. **PII-safe logging** — auto-redaction of sensitive data
6. **Latency-optimized voice pipeline** — sub-800ms TTFAB, <30ms barge-in
7. **Durable workflows** — Temporal for multi-step, stateful AI operations
8. **Dual event buses** — NATS JetStream (high-throughput) + Redis pub/sub (lightweight)
9. **pgvector for RAG** — HNSW index on 1536-dim embeddings directly in Postgres

---

## 18. 🚨 Risks & Gaps to Address

1. **In-memory call state** — `activeCalls` map in Go telephony service is not persisted (lost on restart)
2. **No distributed session persistence** — call state should move to Redis
3. **Mock billing** — Stripe is simulated, not integrated
4. **Agent auth missing** — no service-to-service authentication between Go ↔ Python
5. **Kafka underutilized** — declared in Docker but no producer/consumer code yet
6. **Flutter app minimal** — mobile app has only an entry point file
7. **Empty agents** — 18 of 20 agents are empty directories
8. **No rate limiting implemented** — Redis key-space defined but not wired to middleware
9. **Monitoring dashboards** — Grafana/Prometheus configs exist but no service metrics exported yet
10. **Webhook handling** — No Twilio/WhatsApp webhook verification implemented yet
