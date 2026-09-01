# 🏗️ Nexora System Architecture & Media Topology

---

## 1. High-Level Polyglot Architecture

Nexora is partitioned by **technical execution responsibilities**:
- **Go Engine**: Ultra-low-latency real-time RTP/SIP audio packet streaming and WebRTC media server integration.
- **Python AI Platform**: Multi-agent conversational intelligence, streaming VAD, STT (Deepgram), Multi-LLM routing, and TTS (Cartesia / ElevenLabs).
- **TypeScript / Node.js Microservices**: API routing, authentication, business rules, billing, database access, and event pub/sub.
- **Flutter Client**: Cross-platform mobile (Android, iOS) and desktop (Windows, macOS) communication dialer.
- **Next.js Portals**: Real-time management dashboards, analytics visualizations, and enterprise portals.

```mermaid
graph TD
    subgraph External["🌐 External Ingress & Providers"]
        PSTN["☎️ PSTN / Carrier Network"]
        TwilioSIP["📡 Twilio / Telnyx SIP Trunk"]
        MobileClient["📱 Mobile App (Flutter)"]
        WebPortal["💻 Browser (Next.js Dashboard)"]
        CloudLLM["☁️ LLM Providers (OpenAI, Groq, Claude, DeepSeek)"]
        CloudVoice["🎙️ Voice APIs (Deepgram STT, Cartesia TTS)"]
    end

    subgraph Edge["🛡️ Edge & Gateway"]
        KongEnv["🔒 Reverse Proxy / Kong / Nginx"]
        APIGW["⚡ Nexora API Gateway (:4000)"]
        LiveKit["🎙️ LiveKit SFU WebRTC Server (:7880)"]
        SIPInbound["📞 SIP Ingress Gateway (Go) (:5060)"]
    end

    subgraph ServiceMesh["⚙️ Node.js Microservices Layer"]
        AuthSvc["🔐 auth-service (:4001)"]
        CallSvc["📞 call-service (:4002)"]
        UserSvc["👤 user-service (:4003)"]
        BillingSvc["💳 billing-service (:4004)"]
        NotifSvc["🔔 notification-service (:4005)"]
        ContactSvc["📇 contact-service (:4006)"]
        TransSvc["📝 transcript-service (:4007)"]
        AnalyticsSvc["📊 analytics-service (:4008)"]
    end

    subgraph PythonAI["🤖 Python AI Agent Engine"]
        LiveKitAgentPool["🎙️ LiveKit Voice Agent Pool"]
        VADModule["⚡ Silero VAD (Voice Activity Detection)"]
        STTClient["📝 Deepgram Streaming STT"]
        LLMRouting["🧠 Multi-LLM Model Router"]
        TTSClient["🗣️ Cartesia / ElevenLabs Streaming TTS"]
        ToolsEngine["🔌 Tool & MCP Registry Bridge"]
    end

    subgraph EventAndWorkflow["⚡ Event Bus & Durable Workflows"]
        NATSJetStream["⚡ NATS JetStream / Kafka (:4222 / :9092)"]
        TemporalServer["⏳ Temporal Workflow Orchestrator (:7233)"]
    end

    subgraph Persistence["💾 Persistence Layer"]
        PostgresDB[("🐘 PostgreSQL 16 (:5432)")]
        RedisDB[("🔴 Redis 7.2 (:6379)")]
        MongoDB[("🍃 MongoDB 6.0 (:27017)")]
        MinIODB[("🪣 MinIO S3 (:9000)")]
        OpenSearchDB[("🔍 OpenSearch 2.13 (:9200)")]
    end

    %% Ingress connections
    PSTN --> TwilioSIP
    TwilioSIP --> SIPInbound
    SIPInbound --> LiveKit
    MobileClient --> APIGW
    WebPortal --> APIGW
    MobileClient -.->|WebRTC Audio| LiveKit

    %% Gateway to Microservices
    APIGW --> AuthSvc
    APIGW --> CallSvc
    APIGW --> UserSvc
    APIGW --> BillingSvc
    APIGW --> NotifSvc
    APIGW --> ContactSvc

    %% Telephony & AI Loop
    LiveKit <--> LiveKitAgentPool
    LiveKitAgentPool --> VADModule
    VADModule --> STTClient
    STTClient --> CloudVoice
    STTClient --> LLMRouting
    LLMRouting --> CloudLLM
    LLMRouting --> ToolsEngine
    LLMRouting --> TTSClient
    TTSClient --> CloudVoice
    TTSClient --> LiveKitAgentPool

    %% Event streaming
    CallSvc -.-> NATSJetStream
    LiveKitAgentPool -.-> NATSJetStream
    NATSJetStream --> TransSvc
    NATSJetStream --> AnalyticsSvc
    NATSJetStream --> TemporalServer

    %% Persistence
    ServiceMesh --> PostgresDB
    ServiceMesh --> RedisDB
    TransSvc --> OpenSearchDB
    CallSvc --> MinIODB
    PythonAI --> MongoDB
```

---

## 2. Real-Time Voice Audio Pipeline

```mermaid
sequenceDiagram
    autonumber
    actor Caller as ☎️ Caller (PSTN / Mobile)
    participant LK as 🎙️ LiveKit SFU
    participant VAD as ⚡ Silero VAD
    participant STT as 📝 Deepgram STT
    participant LLM as 🧠 LLM Router (Groq/OpenAI)
    participant Tools as 🔌 Tool / MCP Registry
    participant TTS as 🗣️ Cartesia Sonic TTS

    Caller->>LK: 1. Inbound Audio (Opus 48kHz RTP)
    LK->>VAD: 2. Stream Audio Frames (20ms chunks)
    VAD->>VAD: 3. Detect Voice Activity / Speech End
    VAD->>STT: 4. Stream Active Speech Audio
    STT->>STT: 5. Real-time Transcription Stream
    STT-->>LLM: 6. Final Interim Transcript String
    
    rect rgb(240, 248, 255)
        note over LLM,Tools: Context & Tool Execution Phase
        LLM->>Tools: 7. Function Call (e.g. check_availability)
        Tools-->>LLM: 8. Tool Result Data
    end

    LLM-->>TTS: 9. Stream Token Chunks (First token in ~200ms)
    TTS-->>LK: 10. Streaming Audio Stream (PCM 24kHz)
    LK-->>Caller: 11. Real-time Spoken Response (<450ms Total E2E Latency)

    opt Interruption Handling (Barge-in)
        Caller->>LK: 12. Caller starts speaking while Agent speaks
        LK->>VAD: 13. Detect user voice frame
        VAD->>TTS: 14. CANCEL active TTS playback immediately
        VAD->>LLM: 15. Abort current generation & reset turn
    end
```

---

## 3. Network Ports & Internal Service Matrix

| Service | Internal Host / DNS | Exposed Local Port | Protocol | Purpose |
| :--- | :--- | :--- | :--- | :--- |
| **API Gateway** | `nexora-api-gateway` | `4000` | HTTP / REST / WS | Edge routing, JWT validation, rate limiting |
| **Auth Service** | `nexora-auth-service` | `4001` | HTTP / REST | User registration, login, OTP dispatch, session tokens |
| **Call Service** | `nexora-call-service` | `4002` | HTTP / REST | Call lifecycle state machine, dialing triggers |
| **User Service** | `nexora-user-service` | `4003` | HTTP / REST | User profiles, organization teams, permissions |
| **Billing Service** | `nexora-billing-service` | `4004` | HTTP / REST | Invoices, Stripe subscriptions, wallet balances |
| **Notification Service** | `nexora-notification-service` | `4005` | HTTP / REST | SMS, WhatsApp, Email, Push dispatchers |
| **LiveKit Media SFU** | `nexora-livekit` | `7880, 7881` | WebRTC / TCP / UDP | Real-time audio conference bridge & streaming |
| **PostgreSQL** | `nexora-postgres` | `5432` | TCP (SQL) | Primary relational database + `pgvector` |
| **Redis** | `nexora-redis` | `6379` | RESP (In-Memory) | Cache, distributed locks, rate-limit buckets |
| **MongoDB** | `nexora-mongo` | `27017` | BSON / TCP | Conversation memory, semi-structured agent state |
| **NATS JetStream** | `nexora-nats` | `4222, 8222` | TCP / HTTP Monitoring | Event bus and streaming pub/sub |
| **Apache Kafka** | `nexora-kafka` | `9092` | TCP (Binary) | High-throughput telemetry and transcript streaming |
| **OpenSearch** | `nexora-opensearch` | `9200` | HTTP / JSON | Transcript full-text search and RAG embeddings |
| **MinIO (S3)** | `nexora-minio` | `9000, 9001` | HTTP / S3 API | Audio recording WAV/MP3 storage & asset uploads |
| **Temporal Engine** | `nexora-temporal` | `7233, 8088` | gRPC / HTTP Web UI | Durable long-running workflow orchestration |
| **HashiCorp Vault** | `nexora-vault` | `8200` | HTTP / REST | Production secrets management & API keys |
