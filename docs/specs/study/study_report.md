# 🧠 Deep Study Report: Nexora AI Telephony & Communication Platform

Based on an extensive architectural review of the repository directories, root monorepo configuration files, and a deep-dive into the core specifications (`docs/specs`), here is a comprehensive breakdown of the platform's structure, technology stack, domain boundaries, and delivery roadmap.

The Nexora platform is an enterprise-grade, polyglot AI Communication Operating System designed around the philosophy that languages and runtimes are partitioned by technical responsibility and latency requirements.

---

## 🏗️ 1. Directory Breakdown

### 📚 `docs/` (Documentation Master Portal)
This directory acts as the central source of truth for the project. 
- **Specifications (`docs/specs/`)**: Contains detailed PRDs/SRDs, system diagrams, delivery roadmaps (spanning 28 sprints over 9 phases), and comprehensive API documentation.
- **Architecture Blueprints**: Defines the multi-model storage layer, runtime partitioning, and real-time audio sequences.

### 📦 `packages/` (The 12 Foundation Packages)
Contains shared business logic and utilities, primarily written in TypeScript (`@nexora/*` scope).
- **`@nexora/ai`**: Multi-LLM provider client (OpenAI, Anthropic, Groq, etc.) with fallback chains.
- **`@nexora/database`**: DB connections (PostgreSQL with `pgvector`, MongoDB, Redis).
- **`@nexora/events`**: Event bus using NATS JetStream and Redis Streams.
- **`@nexora/tools`**: Dynamic tool registry and Model Context Protocol (MCP) manager.
- **Other Packages**: `config`, `logger`, `notifications`, `prompts`, `sdk`, `shared`, `types`, `ui`.

### ⚙️ `services/` (Backend Microservices)
A constellation of **24 Microservices** responsible for core domain logic, utilizing a polyglot approach.
- **Node.js/Express**: Standard API/BFF microservices on ports `:4000` to `:4020` (e.g., `auth-service`, `billing-service`).
- **Go 1.22**: High-throughput components like the `ai-orchestrator` (gRPC :50051) and `telephony-service` (:8085).

### 🤖 `agents/` (Multi-Agent AI Ecosystem)
Contains **20 Specialized Python Agents** (Python 3.11) leveraging LangGraph.
- **`orchestrator-agent`**: The master intent router that directs calls to specific specialists.
- **`phone-agent`**: Real-time voice agent connected via LiveKit WebRTC.
- **Specialists**: `support-agent`, `calendar-agent`, `finance-agent`, `research-agent`, etc.

### 📱 `apps/` (Client Applications)
6 cross-platform client interfaces partitioned by technology:
- **Next.js 14 / React 19**: Web portals including `admin-panel`, `user-dashboard`, `enterprise-portal`, `landing-page`, and `documentation`.
- **Flutter 3.x / Dart 3.x**: Hardware-accelerated applications for `mobile-app` (iOS/Android) and `desktop` (Windows/macOS).

### 🛠️ `scripts/` (Build & Server Operations)
Auto-detecting (Windows, macOS, Ubuntu) shell and Node scripts for multi-platform release packaging and lifecycle management (`setup`, `start:all`).

### 🏗️ `infra/` (Infrastructure & Deployments)
Contains IaC and environment definitions, including Docker, Kubernetes, Monitoring (Prometheus, Grafana, Tempo), Nginx, and Terraform.

---

## 🛠️ 2. Monorepo Configuration & Tooling

The root directory contains the configuration necessary to orchestrate this massive codebase:

- **Monorepo Manager (`pnpm-workspace.yaml`)**: Maps `apps/*`, `services/*`, `packages/*`, and `tests` into a unified PNPM workspace.
- **Build System (`turbo.json` & `.turbo/`)**: Utilizes Turborepo for highly concurrent, cached task execution.
- **Go Workspaces (`go.work`)**: Unifies the Go microservices under a single Go workspace.
- **Multi-Provider Configuration (`.env.example`)**: Highly agnostic setup supporting OpenAI/Anthropic/Groq/Ollama (LLMs), Deepgram/Azure (STT), Cartesia/ElevenLabs (TTS), and Twilio/WhatsApp (Comm).

---

## 💾 3. Multi-Model Storage Infrastructure
The platform avoids a monolithic database, instead utilizing a highly specialized data layer:
- **PostgreSQL 16 + pgvector**: Relational multi-tenant isolation, RBAC, billing, and HNSW embeddings.
- **Redis 7.2**: High-speed token buckets for rate limiting, active call session locks, and JWT auth.
- **MongoDB 6.0**: Semi-structured document storage for conversation memories and facts.
- **OpenSearch 2.13**: Full-text RAG search and transcript indexing.
- **MinIO (S3)**: Blob storage for audio recordings and waveform data.
- **HashiCorp Vault**: Zero-secret KV v2 credential storage.

---

## 🎯 4. Product Requirements & KPIs (from PRD)

The platform targets Enterprise Operations Leaders, Contact Center Agents, Developers, and Tenant Admins. Its core KPIs are extremely strict:
- **First-Call Resolution (FCR)**: >80% automated resolution without human transfer.
- **Voice Latency (TTFAB - Time to First Audio Byte)**: <800ms median across global PSTN calls.
- **Natural Barge-In**: Silero VAD detects user speech during agent playback and interrupts audio within **20ms**.
- **System Availability**: 99.99% monthly SLA.

---

## 🎙️ 5. Real-Time Voice Pipeline Design

To achieve the sub-800ms TTFAB requirement, the voice engine operates as a highly optimized continuous streaming frame processor:

1. **Ingest**: PSTN caller RTP packets hit the **Go SIP Gateway** and are bridged into a **LiveKit SFU** WebRTC room.
2. **Audio Processing**: The Python `phone-agent` receives 20ms audio frame chunks. **Silero VAD** detects human voice.
3. **Synthesis Pipeline**: **Deepgram Nova-2** streams interim transcripts to the LLM (**Groq / Claude**), which streams response tokens directly into the TTS engine (**Cartesia Sonic**). 
4. **Chunked Playback**: Audio chunks are pushed back to LiveKit as soon as the first 50ms of audio is synthesized. 
5. **Barge-In**: If the user interrupts, the VAD issues a cancellation signal that halts TTS generation, flushes the LiveKit jitter buffers, and aborts the LLM context, all in under 20ms.

---

## 🤖 6. Multi-Agent Orchestration & Tool Registry

The platform uses a hub-and-spoke multi-agent topology orchestrated by the `ai-orchestrator`:
- **4 Agent Clusters**: VoiceCore (Phone, Support), Messaging (WhatsApp, SMS, Email), Business (Finance, Shopping, Travel), and Intelligence (Research, Knowledge, Coding, Browser).
- **Tool Registry (`@nexora/tools`)**: Instead of agents integrating APIs directly, tools (Calendar, CRM, Payment) are exposed to the AI ecosystem via an **MCP (Model Context Protocol)** server. Agents dynamically connect to remote MCP servers via JSON-RPC.
- **Shared Memory**: All agents use the same Memory Service (backed by MongoDB and OpenSearch), meaning the `phone-agent` remembers what was texted to the `whatsapp-agent`.

---

## 🚀 7. Delivery Roadmap (13 Months / 28 Sprints)

The project is structured into a massive 9-Phase, 28-Sprint delivery plan:

- **Phase 0 (Sprints 1-2)**: Foundations & Infrastructure (Monorepo, Docker, Security baseline).
- **Phase 1 (Sprints 3-6)**: MVP Core Loop (Realtime Voice, Memory, Orchestrator, Calendar, and Support Agents).
- **Phase 2 (Sprints 7-9)**: MVP Surfaces (Flutter Mobile App, Next.js Dashboard, Admin Panel).
- **Phase 3 (Sprints 10-12)**: Omnichannel Messaging (WhatsApp, SMS, Email).
- **Phase 4 (Sprints 13-15)**: Admin, Analytics & Billing Maturity (Stripe payments, RBAC, P&L analytics).
- **Phase 5 (Sprints 16-20)**: Specialized AI Agents (Rolling out the remaining 15 agents including Browser, Finance, Coding, and RAG).
- **Phase 6 (Sprints 21-23)**: Desktop App & Enterprise Portal (SSO, Windows/macOS softphone).
- **Phase 7 (Sprints 24-26)**: Security, Compliance & Scale (SOC 2, 10k concurrent calls, chaos testing).
- **Phase 8 (Sprints 27-28)**: Launch Readiness & GA.
