# ⚖️ Framework Choices & Technical Stack Evaluations

This document evaluates the architectural rationale behind framework, language, and runtime selections across all platform tiers.

---

## 1. Comparative Architecture Evaluation

| Layer | Selection | Rationale & Tradeoffs | Verdict |
| :--- | :--- | :--- | :---: |
| **Android App** | Flutter + Dart | Native compilation, high performance, shares 80%+ codebase with Desktop/iOS. | ✅ Selected |
| **Web App & Dashboards** | Next.js 14 + TypeScript | SSR/SSG rendering, SEO readiness, fast load times, rich React ecosystem. | ✅ Selected |
| **Desktop Softphone** | Flutter + Dart | Cross-platform desktop (macOS & Windows) sharing UI and WebRTC code with mobile. | ✅ Selected |
| **Admin Panel** | Next.js 14 + TypeScript | Shares component library (`@nexora/ui`), auth tokens, and types with Web App. | ✅ Selected |
| **API Gateway & Microservices**| Node.js + TypeScript + Express | Fast I/O aggregation, rich npm ecosystem, unified JSON data handling. | ✅ Selected |
| **Real-Time Voice Pipeline** | Python 3.11 + LiveKit Agents | Ecosystem maturity for Silero VAD, Deepgram STT, and Cartesia/ElevenLabs streaming. | ✅ Selected |
| **Telephony & SIP Ingress** | Go 1.22 | Extreme concurrency, sub-millisecond packet processing, low memory overhead. | ✅ Selected |
| **Durable Workflow Automation**| Temporal 1.24 | Crash-resilient state machines across Go, Python, and TypeScript services. | ✅ Selected |
| **Relational & Vector Memory** | PostgreSQL 16 + `pgvector` | Unified ACID transactions and high-dimensional semantic vector search. | ✅ Selected |

---

## 2. Polyglot Language Partitioning Rule

> **"Languages are separated by technical responsibility, not by business feature."**

1. **Go**: Dedicated to ultra-low-latency network packet handling (SIP, RTP, WebRTC audio proxying).
2. **Python**: Dedicated to AI models, VAD frame processing, embeddings, and agent tool execution.
3. **TypeScript / Node.js**: Dedicated to API aggregation, business logic, billing, CRUD persistence, and frontend portals.
4. **Dart / Flutter**: Dedicated to cross-platform mobile and desktop client applications.
