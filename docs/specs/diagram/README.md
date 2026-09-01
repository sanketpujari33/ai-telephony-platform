# 🗺️ Nexora Platform Master Architecture & Diagrams Blueprint

Comprehensive architectural blueprints, entity-relationship diagrams (ERDs), call sequence flows, and multi-agent coordination topologies for the Nexora AI Communication Platform.

---

## 📑 Visual Architecture Index

| Document | Focus Area | Key Diagrams & Data Structures |
| :--- | :--- | :--- |
| [**1. Diagram Inventory & Specs**](./diagram-inventory-and-specs.md) | Complete Catalog | Master inventory of all diagrams, notations, and diagram sources. |
| [**2. System Architecture**](./system-architecture.md) | Polyglot System Topology & Ingress | Master system map, WebRTC media path, API Gateway routing, NATS/Kafka event bus. |
| [**3. Database & Data Structures ERD**](./db/database-schemas-erd.md) | Multi-Database Storage Layer | PostgreSQL relational ERD, MongoDB memory/fact schemas, Redis key-space & TTLs, OpenSearch indices. |
| [**4. Microservices Architecture**](./services/microservices-architecture.md) | Node.js / Express Microservices | MVC service layers, Inter-service RPC & event payloads, API Gateway route map. |
| [**5. Real-time Call Flow Sequences**](./call-flow-sequence.md) | Telephony & Voice Audio Pipelines | Inbound/outbound call lifecycle, LiveKit room allocation, streaming STT/TTS, tool execution. |
| [**6. Multi-Agent & MCP Coordination**](./multi-agent-coordination.md) | AI Engine & MCP Ecosystem | 20 specialized agent graph, intent routing, dynamic tool registry, remote MCP server bridge. |
| [**7. Admin Operations & Platform Modules**](./services/admin-operations-flow.md) | 9 Core Operations Pipelines | Banner delivery, Promoter/affiliate payouts, Receive funds, Profit/Loss margins, Notifications with DND, Dead call reports, Data monitoring & keywords, Feedback CSAT, General tabs. |
| [**8. User & Agent Operator Flow**](./services/user-operations-flow.md) | Workspace & Softphone Calling | Softphone outbound/inbound dialer flow, live in-call assistance AI, team invites, and customer self-service. |

---

## 🏛️ Platform Core Architecture Overview

```mermaid
graph TB
    subgraph Clients["📱 Client Layer (Cross-Platform)"]
        FlutterApp["📱 Flutter Mobile & Desktop<br/>(Android / iOS / macOS / Windows)"]
        NextWeb["🌐 Next.js Portals<br/>(User Dashboard / Admin Panel / Enterprise)"]
        SIPPhone["☎️ Hardphone / PSTN Caller<br/>(Carrier Trunking / VoIP)"]
    end

    subgraph Gateway["🚪 API Gateway & Telephony Ingress"]
        APIGateway["⚡ API Gateway (Node.js/Express)<br/>Port: 4000 (Auth, Rate Limit, Proxy)"]
        LiveKitServer["🎙️ LiveKit Media Server<br/>Port: 7880 (WebRTC Audio Bridge)"]
        SIPTrunk["📞 SIP Trunk Ingress (Go)<br/>(PSTN / FreeSWITCH / Asterisk)"]
    end

    subgraph CoreServices["⚙️ Backend Microservices (Node.js/TypeScript)"]
        AuthSvc["🔐 auth-service (:4001)"]
        CallSvc["📞 call-service (:4002)"]
        UserSvc["👤 user-service (:4003)"]
        BillingSvc["💳 billing-service (:4004)"]
        NotifSvc["🔔 notification-service (:4005)"]
        ContactSvc["📇 contact-service (:4006)"]
        TransSvc["📝 transcript-service (:4007)"]
        AnalyticsSvc["📊 analytics-service (:4008)"]
    end

    subgraph AIEngine["🤖 Multi-Agent AI Engine (Python / LiveKit Agents)"]
        OrchestratorAgent["🧠 Orchestrator Agent<br/>(Intent Classifier & Router)"]
        SpecializedAgents["🤖 20 Specialized Agents<br/>(Phone, Support, WhatsApp, SMS, Finance, RAG, etc.)"]
        MCPManagerBridge["🔌 Tool & MCP Manager<br/>(Local Tools & Remote MCP Servers)"]
    end

    subgraph OperationsMesh["💼 Admin & User Operations Engine"]
        BannerMgr["🖼️ Banner Manager"]
        PromoterMgr["🤝 Promoter & Affiliate Manager"]
        FundMgr["💰 Receive Fund Manager"]
        ProfitLossMgr["🧮 Profit Loss Manager"]
        DeadReportMgr["💀 Dead Report Manager"]
        DataMonitor["🛡️ Data Monitoring & Keywords"]
        FeedbackMgr["⭐ Feedback & CSAT Review"]
        GeneralTabs["⚙️ General Tabs & Settings"]
        UserSoftphone["📞 Agent Softphone & Team Workspace"]
    end

    subgraph StorageLayer["💾 Multi-Model Storage Infrastructure"]
        PostgresDB[("🐘 PostgreSQL 16 (:5432)")]
        RedisCache[("🔴 Redis 7.2 (:6379)")]
        MongoStore[("🍃 MongoDB 6.0 (:27017)")]
        OpenSearchIdx[("🔍 OpenSearch 2.13 (:9200)")]
        MinIOStore[("🪣 MinIO / S3 (:9000)")]
    end

    FlutterApp --> APIGateway
    NextWeb --> APIGateway
    SIPPhone --> SIPTrunk
    SIPTrunk --> LiveKitServer
    FlutterApp -.->|WebRTC Audio| LiveKitServer

    APIGateway --> CoreServices
    APIGateway --> OperationsMesh
    LiveKitServer <--> AIEngine
    OrchestratorAgent --> SpecializedAgents
    SpecializedAgents --> MCPManagerBridge

    CoreServices --> PostgresDB
    CoreServices --> RedisCache
    OperationsMesh --> PostgresDB
    OperationsMesh --> RedisCache
    AIEngine --> MongoStore
    TransSvc --> OpenSearchIdx
    CallSvc --> MinIOStore
```
