# 🐳 Infrastructure & 20-Agent Ecosystem Deep Study

This document details the containerized infrastructure topology in `infra/docker/` and the 20 specialized Python AI agents in `agents/`.

---

## 1. Local Infrastructure & Database Topology (`infra/docker/`)

The local development environment (`infra/docker/docker-compose.dev.yml`) provides a complete, production-mirrored infrastructure stack:

```text
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                           LOCAL INFRASTRUCTURE TOPOLOGY                                 │
├──────────────────────┬──────────────────────┬───────────────────────────────────────────┤
│ Service / Engine     │ Port Mapping         │ Purpose & Internal Connection Details     │
├──────────────────────┼──────────────────────┼───────────────────────────────────────────┤
│ 🐘 PostgreSQL 16      │ localhost:5432       │ Relational data, pgvector embeddings, Temporal state (user: nexora) │
│ 🔴 Redis 7.2         │ localhost:6379       │ In-memory caching, tenant rate-limiting, distributed locks │
│ 🍃 MongoDB 6.0       │ localhost:27017      │ Document store for semi-structured AI memory & logs       │
│ ⚡ NATS JetStream     │ localhost:4222, 8222 │ Real-time event bus, Pub/Sub messaging, monitoring        │
│ 📨 Apache Kafka      │ localhost:9092       │ High-throughput stream ingestion for transcripts & events │
│ 🪣 MinIO (S3 API)    │ localhost:9000, 9001 │ Audio recordings, call waveforms, presigned assets        │
│ 🔍 OpenSearch 2.13   │ localhost:9200       │ Full-text search over call transcripts & knowledge bases  │
│ ⏳ Temporal 1.24     │ localhost:7233, 8088 │ Durable multi-step workflow engine & Web UI               │
│ 🎙️ LiveKit 1.7       │ localhost:7880, 7881 │ Real-time WebRTC audio streaming, VoIP, TURN relay        │
│ 🔒 HashiCorp Vault   │ localhost:8200       │ Zero-secret local KV v2 secrets store (token: root)       │
└──────────────────────┴──────────────────────┴───────────────────────────────────────────┘
```

---

## 2. Multi-Database & Redis Client Architecture (`@nexora/database`)

The [`@nexora/database`](file:///d:/trideocoders/ai-telephony-platform/packages/database) package provides unified connection pooling, tenant isolation, and health checking across PostgreSQL, Redis, and MongoDB:

```text
@nexora/database
├── postgres.ts       ── Connection pool, pgvector similarity queries, parameterized query runner
├── redis.ts          ── Singleton ioredis client, tenantKey namespacing, JSON get/set, Redlock single-node locking
├── mongodb.ts        ── MongoClient pooling, collection getters, ping healthchecks
├── repository.ts     ── BasePostgresRepository with mandatory WHERE tenant_id = $x enforcement
├── health.ts         ── Composite multi-database health monitor (Postgres + Redis + Mongo)
└── migrator.ts       ── Automatic schema migrations runner
```

---

## 3. The 20 Specialized Agent Ecosystem (`agents/`)

Every agent in [`agents/`](file:///d:/trideocoders/ai-telephony-platform/agents) is tailored for a specific communication channel and operational domain, powered by LiveKit Agents, Silero VAD, Deepgram STT, OpenAI/Groq/Claude LLMs, and Cartesia/ElevenLabs TTS:

| Agent Directory | Domain Scope | Primary Channel | Voice/STT Config |
| :--- | :--- | :--- | :--- |
| **`orchestrator-agent`** | Intent classification, agent routing, tool handoffs | Omnichannel | N/A (Text / Voice routing) |
| **`phone-agent`** | Live telephone calls, warm transfers, caller verification | Phone / PSTN / WebRTC | Silero VAD + Deepgram + Cartesia |
| **`support-agent`** | Customer service, knowledge base lookup, ticket creation | Voice & Web Chat | Deepgram Nova-2 + Claude 3.5 + Cartesia |
| **`appointment-agent`** | Calendar scheduling, slot booking, conflict checks | Voice & WhatsApp | Deepgram + Groq Llama-3.1 + ElevenLabs |
| **`delivery-agent`** | Order status tracking, ETA notifications, address verification | SMS & WhatsApp | Mistral Codestral / Fast inference |
| **`whatsapp-agent`** | Rich WhatsApp messaging, interactive buttons, quick replies | WhatsApp Cloud API | Meta Cloud API + Claude 3.5 |
| **`sms-agent`** | Concise 160-char transactional SMS, verification codes | SMS | Twilio / Telnyx / Infobip |
| **`email-agent`** | Contextual email inbox triaging, customer reply drafting | Email | SendGrid / Resend + GPT-4o |
| **`calendar-agent`** | Google / Outlook / Cal.com calendar synchronization | Internal Tools | MCP Tool execution |
| **`reminder-agent`** | Automated follow-up alarms, appointment reminders | Voice / SMS | Temporal durable cron workflows |
| **`finance-agent`** | Invoice status, subscription upgrades, payment links | Web & Voice | Stripe / Razorpay / Paddle tools |
| **`shopping-agent`** | E-commerce product catalog search, checkout assistance | Web & WhatsApp | Vector search + Tool Calling |
| **`travel-agent`** | Flight/hotel itineraries, booking confirmations, updates | Voice & WhatsApp | MCP travel integrations |
| **`research-agent`** | Live web search citation, documentation synthesis | Web Chat | Perplexity Sonar / Grok-2 |
| **`browser-agent`** | Headless browser web interaction, web scraping | Automation | Playwright / MCP Browser |
| **`coding-agent`** | Code analysis, webhook payload debugging, SDK snippets | Developer Portal | Claude 3.5 Sonnet / Codestral |
| **`enterprise-agent`** | SSO configuration, SCIM user provisioning, compliance audits | Enterprise Portal | Internal Enterprise Services |
| **`automation-agent`** | Zapier/Make/Webhook trigger processing | Webhooks | Temporal Workflows |
| **`knowledge-agent`** | RAG semantic retrieval over internal company docs | RAG Engine | OpenSearch + pgvector |
| **`memory-agent`** | User preference extraction, conversation long-term memory | Memory Engine | pgvector / MongoDB |
