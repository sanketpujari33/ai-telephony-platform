# 🏗️ Phase 0 — Foundations & Infrastructure Specification
*Sprints 1–2 · Weeks 1–4*

---

## 1. Phase Goals & Scope

Phase 0 establishes the bedrock foundation of the entire Nexora AI Communication Platform. The objective is to build a reliable, type-safe monorepo containing all shared packages, containerized infrastructure, secrets management, database migration runners, and initial authentication and user services.

---

## 2. Sprint Deliverables

### Sprint 1: Monorepo, Infra Skeleton & Shared Packages
- **Monorepo Architecture**: Setup Turborepo with pnpm workspaces (`apps/`, `services/`, `agents/`, `packages/`, `infra/`, `docs/`, `tests/`, `scripts/`).
- **Foundational Packages**:
  - `@nexora/types`: Universal TypeScript interfaces and domain schemas.
  - `@nexora/config`: Zod-validated environment config and HashiCorp Vault KV v2 client.
  - `@nexora/logger`: High-speed structured Pino logger with OpenTelemetry distributed trace context.
  - `@nexora/events`: NATS JetStream, Redis Pub/Sub, and in-memory event buses with DLQ.
  - `@nexora/shared`: Result monad, AES-256-GCM crypto, E.164 phone normalization, 6 billing HMAC verifiers.
- **Local Dev Infrastructure (`docker-compose.dev.yml`)**:
  - PostgreSQL 16 (`:5432`) with `pgvector`
  - Redis 7.2 (`:6379`)
  - MongoDB 6.0 (`:27017`)
  - NATS JetStream (`:4222`, `:8222`)
  - Apache Kafka (`:9092`)
  - MinIO S3 (`:9000`, `:9001`)
  - OpenSearch 2.13 (`:9200`)
  - Temporal 1.24 (`:7233`, `:8088`)
  - LiveKit WebRTC Server (`:7880`)
  - HashiCorp Vault (`:8200`)

### Sprint 2: Identity, Gateway, Storage & Observability Baseline
- **`services/auth-service`**: JWT access & refresh tokens, password hashing with bcrypt, OTP phone verification lifecycle.
- **`services/user-service`**: User CRUD, tenant provisioning, profile metadata, multi-device tracking.
- **`services/api-gateway`**: Central reverse proxy, JWT authentication middleware, tenant context injection, rate limiting.
- **`services/storage-service`**: MinIO / S3 object upload, presigned URLs for audio recordings and media files.
- **Database & SDK Packages**: `@nexora/database` connection pooling and `@nexora/sdk` TypeScript client library.
- **Observability**: OpenTelemetry tracing hooks, Prometheus `/metrics` endpoints, Sentry error capture.

---

## 3. Exit Criteria & Definition of Done

- [x] All 12 `@nexora/*` packages build with 0 TypeScript errors.
- [x] Docker Compose launches all 10 infrastructure engines locally.
- [x] Secrets are securely retrieved from Vault without hardcoded credentials.
- [x] Centralized tests in `tests/packages/` and `tests/services/` pass 100%.
