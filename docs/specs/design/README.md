# 🏛️ Platform Subsystem Designs Master Index

Technical architecture design specifications detailing polyglot subsystem boundaries, streaming audio pipelines, multi-agent coordination, database clustering, and security controls.

---

## 📑 Subsystem Design Documents

| Specification Document | Focus Area | Key Architectural Content |
| :--- | :--- | :--- |
| [**1. System Architecture Design**](./system-architecture-design.md) | Polyglot Ingress & Boundaries | Go SIP/WebRTC ingress, Node.js API Gateway, Python AI platform, multi-database persistence. |
| [**2. Real-Time Voice Pipeline Design**](./realtime-voice-pipeline-design.md) | Voice Audio Processing Loop | Silero VAD, Deepgram STT, Cartesia TTS, Sub-20ms barge-in interruption flow. |
| [**3. Multi-Agent Orchestration Design**](./multi-agent-orchestration-design.md) | Agent Graph & MCP Tools | Master intent router, 20 specialized Python agents, dynamic Tool Registry, Model Context Protocol. |
| [**4. Data Storage & Caching Design**](./data-storage-and-caching-design.md) | Polyglot Persistence Layer | PostgreSQL 16 + pgvector, MongoDB memory, Redis 7.2 caching, OpenSearch 2.13, MinIO S3. |
| [**5. Security & Compliance Design**](./security-and-compliance-design.md) | Cryptography & Zero-Secrets | HashiCorp Vault KV v2 zero-secret architecture, AES-256-GCM encryption, PII redaction. |
