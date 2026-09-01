# 🔧 Nexora Platform — 10-Gap Resolution Plan

## Overview

Resolves all 10 critical architecture gaps identified in the deep study. Changes span **Go**, **TypeScript/Node.js**, **Python**, and **Flutter/Dart**.

---

## Gap 1 — Redis Call State Persistence (Go Telephony Service)

> **Problem:** `activeCalls map[string]*CallSession` is in-memory and lost on restart.

### Files Changed

#### [MODIFY] [`main.go`](file:///d:/trideocoders/ai-telephony-platform/services/telephony-service/main.go)
- Replace the in-memory `activeCalls` map with Redis-backed persistence
- Add `go-redis/v9` to `go.mod`
- On `POST /api/v1/calls/inbound` and `outbound`: `SET nexora:call:active:{callId}` JSON with 24h TTL
- On `GET /api/v1/calls`: `SCAN nexora:call:active:*` + `MGET`
- Add `DELETE /api/v1/calls/:id` endpoint to terminate + remove call from Redis

#### [MODIFY] [`go.mod`](file:///d:/trideocoders/ai-telephony-platform/services/telephony-service/go.mod)
- Add `github.com/redis/go-redis/v9`

#### [NEW] `pkg/store/redis.go`
- `NewRedisStore()` — connects to `REDIS_URL`
- `SaveCall(ctx, session)` — SET with JSON marshal
- `GetCall(ctx, callId)` — GET + JSON unmarshal
- `ListCalls(ctx)` — SCAN + pipeline MGET
- `DeleteCall(ctx, callId)` — DEL

---

## Gap 2 — Real JWT Auth Middleware (API Gateway)

> **Problem:** `optionalAuthMiddleware` hardcodes `usr_authenticated` — never actually verifies tokens.

### Files Changed

#### [MODIFY] [`auth.middleware.ts`](file:///d:/trideocoders/ai-telephony-platform/services/api-gateway/src/middlewares/auth.middleware.ts)
- Import `jsonwebtoken` and verify the Bearer token using `JWT_SECRET` env var
- Decode `{sub, role, tenantId, email}` from payload
- On invalid/expired token: return `401 Unauthorized`
- Populate `req.user` from verified payload (not hardcoded values)
- Check Redis token blacklist via `@nexora/database` `isTokenBlacklisted()`

#### [NEW] `src/middlewares/ratelimit.middleware.ts`
- See Gap 8 below

---

## Gap 3 — Stripe Billing Integration

> **Problem:** Billing service simulates Stripe checkout with `Math.random()`.

### Files Changed

#### [MODIFY] [`billing-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/billing-service/src/index.ts)
- Replace in-memory `tenantUsage` Map with PostgreSQL (via `@nexora/database`)
- `POST /api/v1/billing/meter` → upsert usage record in `calls` table cost column
- `POST /api/v1/billing/checkout` → call real `stripe.checkout.sessions.create()` with `payment_intent_data`, `line_items`, and `success_url`/`cancel_url`
- `POST /api/v1/billing/webhook` → verify Stripe webhook signature with `stripe.webhooks.constructEvent()`, handle `checkout.session.completed` and `invoice.paid`

#### [MODIFY] `billing-service/package.json`
- Add `stripe` npm package

---

## Gap 4 — Service-to-Service JWT Auth (Go ↔ Python)

> **Problem:** AI Orchestrator (Go) calls Orchestrator Agent (Python) with no authentication.

### Files Changed

#### [NEW] `services/ai-orchestrator/pkg/auth/service_token.go`
- `GenerateServiceToken(secret, issuer, audience, ttl)` — HS256 JWT signed with `SERVICE_JWT_SECRET`
- Signed claims: `{iss: "ai-orchestrator", aud: "orchestrator-agent", iat, exp}`

#### [MODIFY] [`services/ai-orchestrator/main.go`](file:///d:/trideocoders/ai-telephony-platform/services/ai-orchestrator/main.go)
- Generate service JWT before every HTTP call to the Python agent
- Add `Authorization: Bearer <service_token>` header to the `http.Post()` call

#### [MODIFY] [`agents/orchestrator-agent/agent.py`](file:///d:/trideocoders/ai-telephony-platform/agents/orchestrator-agent/agent.py)
- Add `verify_service_token()` FastAPI dependency using `python-jose`
- Validates `iss`, `aud`, `exp` from the incoming `Authorization` header
- Apply as `Depends(verify_service_token)` on the `/api/v1/evaluate` route

#### [MODIFY] `services/ai-orchestrator/go.mod`
- Add `github.com/golang-jwt/jwt/v5`

---

## Gap 5 — Kafka Producers & Consumers

> **Problem:** Kafka container running but no application code publishes or consumes events.

### Files Changed

#### [NEW] `services/telephony-service/pkg/kafka/producer.go`
- `NewCallEventProducer(brokers)` using `segmentio/kafka-go`
- `PublishCallStarted(ctx, session)` → topic `nexora.call.started`
- `PublishCallEnded(ctx, callId, duration)` → topic `nexora.call.ended`

#### [MODIFY] `services/telephony-service/main.go`
- Initialize Kafka producer on startup
- Call `producer.PublishCallStarted()` after LiveKit room created
- Call `producer.PublishCallEnded()` on call termination endpoint

#### [MODIFY] `services/telephony-service/go.mod`
- Add `github.com/segmentio/kafka-go`

#### [NEW] `services/analytics-service/src/kafka/consumer.ts`
- `@nexora/events` NATS consumer → bridge to Kafka topics `nexora.call.*`
- Subscribe `nexora.call.started`: increment active call counter metric
- Subscribe `nexora.call.ended`: store call duration analytics, update Prometheus gauge

#### [MODIFY] `services/notification-service/src/index.ts`
- Replace stub `eventBus.subscribe('call.completed')` with real Kafka consumer
- Consume `nexora.call.ended` → trigger real push/SMS notification stub

---

## Gap 6 — Flutter Mobile App Expansion

> **Problem:** `lib/main.dart` is a single 176-line file with only a dialer screen.

### Files Changed

#### [MODIFY] `pubspec.yaml`
- Add dependencies: `go_router`, `provider`, `shared_preferences`, `dio`

#### [MODIFY] `lib/main.dart`
- Add `go_router` routing with named routes
- Routes: `/login`, `/home`, `/dialer`, `/active-call`, `/call-history`

#### [NEW] `lib/screens/login_screen.dart`
- Email + password login form → POST to `$API_BASE_URL/api/v1/auth/login`
- Stores JWT in `shared_preferences`

#### [NEW] `lib/screens/home_screen.dart`
- Bottom nav: Dialer | Calls | AI Chat | Settings
- AI toggle switch (global state via Provider)

#### [NEW] `lib/screens/active_call_screen.dart`
- Live call UI: caller name, duration timer, mute/speaker/hang-up buttons
- AI status indicator (analyzing / speaking / listening)

#### [NEW] `lib/screens/call_history_screen.dart`
- List of past calls with duration, direction (in/out), AI handled badge

#### [NEW] `lib/services/api_service.dart`
- `ApiService` singleton using `dio`
- `login()`, `dispatchCall()`, `getActiveCalls()`, `getCallHistory()`

#### [NEW] `lib/providers/call_provider.dart`
- `ChangeNotifier` provider for active call state
- `startCall(phoneNumber)`, `endCall()`, `toggleAI()`

---

## Gap 7 — 18 Empty Python Agents

> **Problem:** 18 of 20 agent directories are empty (no `agent.py`).

Each agent follows the same **FastAPI + OpenAI structured output** pattern from `orchestrator-agent`.

### Files Created (one per agent)

| Agent | File | Capability |
|---|---|---|
| `calendar-agent` | `agent.py` | Create/list/delete Google Calendar events via tool calls |
| `appointment-agent` | `agent.py` | Slot availability check + booking confirmation |
| `support-agent` | `agent.py` | Ticket creation, warm transfer, FAQ lookup |
| `email-agent` | `agent.py` | Compose/send emails, summarize inbox |
| `whatsapp-agent` | `agent.py` | Send/receive WhatsApp messages via Meta Cloud API |
| `sms-agent` | `agent.py` | Send SMS via Twilio |
| `memory-agent` | `agent.py` | Fact extraction + MongoDB upsert, memory retrieval |
| `reminder-agent` | `agent.py` | Temporal workflow scheduling for reminders |
| `research-agent` | `agent.py` | Web search via Perplexity/Tavily API |
| `browser-agent` | `agent.py` | Playwright headless browser automation |
| `finance-agent` | `agent.py` | Invoice lookup, balance summary, billing FAQ |
| `shopping-agent` | `agent.py` | Product search, cart management stub |
| `travel-agent` | `agent.py` | Flight/hotel search, itinerary creation |
| `coding-agent` | `agent.py` | Code generation, debugging assistance |
| `enterprise-agent` | `agent.py` | Policy lookup, escalation routing |
| `automation-agent` | `agent.py` | Workflow trigger via Temporal |
| `knowledge-agent` | `agent.py` | RAG retrieval from OpenSearch knowledge base |
| `delivery-agent` | `agent.py` | Order status, delivery tracking |

All agents expose `POST /api/v1/handle` and `GET /health` via FastAPI.

---

## Gap 8 — Rate Limiting (Redis Token Bucket wired to API Gateway)

> **Problem:** Redis key namespace `nexora:ratelimit:*` defined but never used.

### Files Changed

#### [NEW] `services/api-gateway/src/middlewares/ratelimit.middleware.ts`
- Redis sliding window rate limiter using `ioredis` + Lua script
- Key: `nexora:ratelimit:{tenantId}:{ip}` — 100 req/min default
- Returns `429 Too Many Requests` with `Retry-After` header on breach
- Configurable per-route limits: `/api/v1/calls/*` → 20/min, `/api/v1/orchestrate` → 10/min

#### [MODIFY] `services/api-gateway/src/index.ts`
- Import and mount `rateLimitMiddleware` before route handlers
- Initialize Redis client from `@nexora/database` on startup

---

## Gap 9 — Prometheus Metrics Exported from All Services

> **Problem:** `prometheus.yml` has only commented scrape targets. No service exports `/metrics`.

### Files Changed

#### [MODIFY] `services/api-gateway/src/index.ts`
- Mount `prom-client` default registry on `GET /metrics`
- Custom counters: `nexora_http_requests_total{method,route,status}`, `nexora_http_duration_ms`

#### [MODIFY] `services/telephony-service/main.go`
- Already imports `prometheus/client_golang` (in `go.mod`)
- Add `/metrics` handler using `promhttp.Handler()`
- Custom gauges: `nexora_active_calls_total`, `nexora_call_duration_seconds`

#### [MODIFY] `agents/orchestrator-agent/agent.py`
- Add `prometheus-fastapi-instrumentator` and mount `/metrics` endpoint
- Custom counters: `nexora_agent_requests_total{agent}`, `nexora_llm_latency_seconds{provider}`

#### [MODIFY] `infra/monitoring/prometheus/prometheus.yml`
- Uncomment and add all service scrape jobs:
  - `api-gateway:8080`, `auth-service:8082`, `telephony-service:8085`
  - `ai-orchestrator:8086`, `orchestrator-agent:8087`, `billing-service:8089`
  - `notification-service:8090`

---

## Gap 10 — Twilio & WhatsApp Webhook Verification

> **Problem:** No webhook signature validation for Twilio or WhatsApp callbacks.

### Files Changed

#### [NEW] `services/telephony-service/pkg/webhook/twilio.go`
- `ValidateTwilioSignature(authToken, url, params, signature)` using Twilio's HMAC-SHA1 scheme
- Returns `bool` — used as middleware guard on `POST /api/v1/webhooks/twilio`

#### [MODIFY] `services/telephony-service/main.go`
- Add `POST /api/v1/webhooks/twilio` route with Twilio signature validation middleware
- Handle: `CallStatus` events (update call record in Redis), `RecordingComplete` (notify recording-service)

#### [NEW] `services/notification-service/src/webhooks/whatsapp.webhook.ts`
- `verifyWhatsAppWebhook()` — validate `X-Hub-Signature-256` header using `crypto.createHmac('sha256', WHATSAPP_APP_SECRET)`
- `GET /webhooks/whatsapp` — respond to Meta's verification challenge (`hub.challenge`)
- `POST /webhooks/whatsapp` — parse inbound messages, emit to NATS `nexora.whatsapp.inbound`

#### [MODIFY] `services/notification-service/src/index.ts`
- Mount WhatsApp webhook router
- Add `WHATSAPP_APP_SECRET` to required env vars

---

## Verification Plan

### Automated Tests
```bash
# Go tests (telephony Redis persistence)
cd tests/services/telephony-service && go test ./...

# Node.js tests (rate limiting, auth middleware)
pnpm run test --filter=@nexora/api-gateway

# Python agent tests
cd agents/orchestrator-agent && pytest tests/

# Full integration healthcheck
pnpm run healthcheck
```

### Manual Verification
1. Restart telephony-service → active calls survive restart (Redis persisted)
2. API Gateway: call without token → 401, with valid JWT → 200, exceed rate limit → 429
3. Billing: POST `/checkout` returns real Stripe URL with `cs_` prefix
4. Trigger inbound call → check Kafka topic `nexora.call.started` for event
5. Flutter app: compile + run on device/emulator → login + dial screens render
6. `curl http://localhost:8080/metrics` → Prometheus output
7. `curl http://localhost:8085/metrics` → Go Prometheus output
8. Prometheus UI → all 7 scrape targets show UP status

---

## Open Questions

> [!IMPORTANT]
> **Stripe mode**: Should billing use Stripe **Test mode** (`sk_test_*`) keys, or stub calls behind a feature flag for now?

> [!IMPORTANT]
> **Flutter routing**: Should the mobile app use `go_router` (recommended, URL-based) or Flutter's built-in `Navigator.push`?

> [!IMPORTANT]
> **Agent scope**: Should the 18 agents be **full implementations** (with real API calls) or **production-ready stubs** (correct structure, mock responses, ready to fill in real API keys)?

> [!WARNING]
> **Kafka broker address**: The current docker-compose uses `kafka:29092` (internal) and `localhost:9092` (host). Go code should use `kafka:29092` when running in Docker. Confirm deployment target.
