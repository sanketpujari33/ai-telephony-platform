# 🏃 Sprints 1–2 — Foundations & Infrastructure Plan
*Weeks 1–4 · Phase 0*

---

## Sprint 1: Monorepo, Infra Skeleton & Shared Packages
- **Objective**: Establish monorepo workspace, develop 5 core foundation packages, and start local containerized infrastructure.
- **Tasks**:
  - [x] Scaffold monorepo with Turborepo & pnpm.
  - [x] Implement `@nexora/types` domain contracts.
  - [x] Implement `@nexora/config` with HashiCorp Vault KV v2 client.
  - [x] Implement `@nexora/logger` using Pino with OTel tracing.
  - [x] Implement `@nexora/shared` crypto, UUID, and billing HMAC verifiers.
  - [x] Implement `@nexora/events` with NATS JetStream, Redis Pub/Sub, and in-memory buses.
  - [x] Configure `infra/docker/docker-compose.dev.yml` with 10 database and infrastructure containers.

---

## Sprint 2: Identity, Gateway, Storage & Observability
- **Objective**: Stand up authentication, tenant-aware routing, S3/MinIO media storage, and database adapters.
- **Tasks**:
  - [x] Build `services/auth-service`: JWT generation, bcrypt password hashing, phone OTP.
  - [x] Build `services/user-service`: User account lifecycle, tenant provisioning.
  - [x] Build `services/api-gateway`: Route proxying, auth verification, token bucket rate limits.
  - [x] Build `services/storage-service`: S3 / MinIO upload handlers and presigned URL generators.
  - [x] Build `@nexora/database`: PostgreSQL connection pool, Redis singleton, MongoDB client.
  - [x] Build `@nexora/sdk`: Isomorphic client SDK for web, desktop, and mobile.
  - [x] Set up OpenTelemetry, Prometheus `/metrics`, and Sentry error trackers.
