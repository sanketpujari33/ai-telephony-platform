# 🏛️ System Architecture Design Specification

---

## 1. High-Level Polyglot Architecture

Nexora is designed with language partitioning based on **technical responsibility**:

```mermaid
graph TD
    subgraph Clients["📱 Cross-Platform Clients"]
        FLUTTER["Flutter Mobile & Desktop (Dart)"]
        NEXT["Next.js Web Portals (TypeScript)"]
        PSTN["PSTN / Carrier Phone Trunk"]
    end

    subgraph Edge["🚪 Ingress & Gateway Layer"]
        GW["API Gateway (:4000) (Node.js/Express)"]
        SFU["LiveKit WebRTC SFU (:7880) (Go)"]
        SIP["SIP Ingress Gateway (:5060) (Go)"]
    end

    subgraph Microservices["⚙️ Backend Microservices (Node.js/TypeScript)"]
        SERVICES["24 MVC Microservices (Auth, Call, User, Billing, Admin...)"]
    end

    subgraph AIPlatform["🤖 Python AI Engine"]
        AGENT_ORCH["Orchestrator Agent & 20 Specialists (Python 3.11)"]
        VOICE_PIPE["Streaming STT/TTS & Silero VAD (Python 3.11)"]
        MCP["Model Context Protocol (MCP) Bridge"]
    end

    subgraph Storage["💾 Multi-Model Storage Infrastructure"]
        POSTGRES[("PostgreSQL 16 + pgvector")]
        REDIS[("Redis 7.2 Cluster")]
        MONGO[("MongoDB 6.0 Memory")]
        OPENSEARCH[("OpenSearch 2.13")]
        MINIO[("MinIO S3 Audio Storage")]
    end

    FLUTTER --> GW
    NEXT --> GW
    PSTN --> SIP
    SIP --> SFU
    FLUTTER -.->|WebRTC Audio| SFU

    GW --> SERVICES
    SFU <--> VOICE_PIPE
    VOICE_PIPE --> AGENT_ORCH
    AGENT_ORCH --> MCP

    SERVICES --> POSTGRES
    SERVICES --> REDIS
    AGENT_ORCH --> MONGO
    SERVICES --> OPENSEARCH
    SERVICES --> MINIO
```

---

## 2. Ingress & Media Path Separation

- **Control Plane**: HTTP REST and JSON-RPC APIs routed through `api-gateway` on port `4000`.
- **Data Plane (Audio Media)**: Low-latency WebRTC RTP media streams handled directly by `LiveKit` SFU on port `7880` and `telephony-service` SIP gateway on port `5060`.
- **Event Plane**: NATS JetStream and Kafka handle asynchronous event streaming, post-call diarization, and billing ledger updates.
