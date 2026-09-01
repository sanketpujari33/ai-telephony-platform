# 🧪 Walkthrough: Unified Testing for Packages & Services

This walkthrough details the centralized testing architecture and verification suites across all `@nexora/*` packages and backend microservices.

---

## 1. Centralized Test Suite Breakdown (`tests/`)

Both packages and backend microservices are tested under [`tests/`](file:///d:/trideocoders/ai-telephony-platform/tests) with **100% pass rate across 15 test files**:

```text
tests/
├── packages/
│   ├── ai/ai.test.ts                      ── 4 tests: Token counting, cost estimation, 13 LLM providers & mocks
│   ├── config/config.test.ts              ── 5 tests: Vault KV v2 client, token auth, AppConfig immutability
│   ├── database/database.test.ts          ── 3 tests: Multitenant PostgreSQL queries & Redis key namespacing
│   ├── events/events.test.ts              ── 2 tests: CloudEvents BaseEvent builder & MemoryEventBus
│   ├── logger/logger.test.ts              ── 5 tests: Pino structured logs, PII & credential redaction
│   ├── notifications/notifications.test.ts ── 3 tests: SMS, WhatsApp, Email, Push router dispatch
│   ├── prompts/prompts.test.ts            ── 7 tests: Template compiler, injection guardrails, agent roles
│   ├── sdk/sdk.test.ts                    ── 1 test:  Isomorphic TypeScript SDK & sub-resources
│   ├── shared/shared.test.ts              ── 8 tests: AES-256-GCM crypto, phone parsing, 6 billing HMAC verifiers
│   ├── tools/tools.test.ts                ── 4 tests: Zod to JSON Schema converter, built-in tools, MCP adapter
│   ├── types/types.test.ts                ── 4 tests: Plan limits, type unions, service & app contract models
│   └── ui/ui.test.ts                      ── 1 test:  Zero-wrapper React component library tree
│
└── services/
    ├── api-gateway/api-gateway.test.ts    ── 2 tests: Outbound call dispatching, CallSessionStore, agent profiles
    ├── auth-service/auth-service.test.ts  ── 5 tests: bcrypt password hashing, JWT generation, registration, OTP lifecycle
    └── call-service/call-service.test.ts  ── 2 tests: Call state transitions (ringing -> in-progress -> completed), CallRegistry
```

---

## 2. Monorepo Verification Matrix

| Check | Result |
| :--- | :--- |
| **Total Test Suites** | **15 Passed (15 Total)** |
| **Total Unit Tests** | **56 Passed (56 Total)** |
| **Type Check (`pnpm run type-check`)** | **28 Tasks Successful (0 errors)** |
| **Turbo Monorepo Build (`pnpm test`)** | **27 Tasks Successful (0 warnings)** |
