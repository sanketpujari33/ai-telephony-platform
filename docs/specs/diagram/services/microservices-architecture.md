# ⚙️ Node.js Microservices Architecture & MVC Topology

All backend microservices in [`services/`](file:///d:/trideocoders/ai-telephony-platform/services) adhere to a strict **Model-View-Controller (MVC)** design pattern, powered by Express, TypeScript ESM, and `@nexora/*` foundation packages.

---

## 1. Microservices Interaction Topology

```mermaid
graph TD
    Client["📱 Client Applications (Web / Mobile / Dialer)"]
    
    subgraph Gateway["🚪 API Gateway (:4000)"]
        Router["⚡ Reverse Proxy & Route Matcher"]
        AuthGuard["🛡️ Auth & JWT Middleware"]
        RateLimiter["⏱️ Redis Token-Bucket Limiter"]
    end

    subgraph Microservices["⚙️ Core Microservices"]
        AuthService["🔐 auth-service (:4001)<br/>JWT, OAuth, OTP verification"]
        CallService["📞 call-service (:4002)<br/>Call lifecycle & state machine"]
        UserService["👤 user-service (:4003)<br/>RBAC, Teams & Profiles"]
        BillingService["💳 billing-service (:4004)<br/>Stripe, Invoicing & Usage"]
        NotificationService["🔔 notification-service (:4005)<br/>SMS, WhatsApp, Email"]
        ContactService["📇 contact-service (:4006)<br/>CRM customer profiles & notes"]
        TranscriptService["📝 transcript-service (:4007)<br/>Utterance indexing & sentiment"]
        AnalyticsService["📊 analytics-service (:4008)<br/>Real-time telemetry & aggregation"]
    end

    subgraph EventMesh["⚡ Event Bus (NATS JetStream)"]
        TopicCalls["📢 Topic: call.*"]
        TopicAuth["📢 Topic: auth.*"]
        TopicBilling["📢 Topic: billing.*"]
        TopicNotif["📢 Topic: notification.*"]
    end

    subgraph DataMesh["💾 Data Mesh"]
        PG[("🐘 PostgreSQL (Relational)")]
        Redis[("🔴 Redis (State & Cache)")]
        Mongo[("🍃 MongoDB (Documents)")]
    end

    Client --> Router
    Router --> AuthGuard
    AuthGuard --> RateLimiter

    RateLimiter -->|/api/v1/auth/*| AuthService
    RateLimiter -->|/api/v1/calls/*| CallService
    RateLimiter -->|/api/v1/users/*| UserService
    RateLimiter -->|/api/v1/billing/*| BillingService
    RateLimiter -->|/api/v1/notifications/*| NotificationService
    RateLimiter -->|/api/v1/contacts/*| ContactService

    %% Event Publications
    CallService -.->|pub: call.initiated, call.completed| TopicCalls
    AuthService -.->|pub: auth.registered, auth.login| TopicAuth
    BillingService -.->|pub: invoice.paid| TopicBilling

    %% Event Subscriptions
    TopicCalls --> TranscriptService
    TopicCalls --> AnalyticsService
    TopicCalls --> NotificationService
    TopicBilling --> NotificationService

    %% DB Persistence
    AuthService --> PG
    CallService --> PG
    CallService --> Redis
    UserService --> PG
    BillingService --> PG
    NotificationService --> Mongo
    ContactService --> PG
```

---

## 2. Strict Service MVC Directory Blueprint

Every microservice in `services/` is structured consistently:

```text
services/<microservice-name>/
├── src/
│   ├── controllers/       ── HTTP Request/Response handlers, status codes & validation
│   ├── routes/            ── Express Router definitions & URL paths
│   ├── services/          ── Pure business logic, external API calls & domain rules
│   ├── models/            ── Data access objects, in-memory stores & DB repositories
│   ├── middlewares/       ── Authentication, rate limiting, error handling & structured logging
│   ├── utils/             ── Helper methods, response formatters, crypto & token utilities
│   └── index.ts           ── Express application bootstrap & port listener
├── package.json           ── Workspace package metadata & dependencies
├── tsconfig.json          ── TypeScript compiler configuration (NodeNext)
└── dist/                  ── Compiled JavaScript output
```

---

## 3. Standard API Response Structure (`utils/response.ts`)

Every HTTP endpoint returns a unified standard envelope:

```typescript
// Success Response (HTTP 200 / 201 / 202)
export interface ApiResponse<T> {
  success: true;
  data: T;
  message?: string;
  statusCode: number;
}

// Error Response (HTTP 400 / 401 / 403 / 404 / 500)
export interface ApiErrorResponse {
  success: false;
  error: string;
  statusCode: number;
  details?: unknown;
}
```

---

## 4. NATS JetStream Event Taxonomy

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                              EVENT BUS TOPOLOGY (NATS)                                 │
├──────────────────────┬──────────────────────┬──────────────────────────────────────────┤
│ Event Subject        │ Source Service       │ Target Consumers                         │
├──────────────────────┼──────────────────────┼──────────────────────────────────────────┤
│ `call.initiated`     │ api-gateway          │ call-service, livekit-agent-pool         │
│ `call.ringing`       │ call-service         │ analytics-service, websocket-service     │
│ `call.in-progress`   │ call-service         │ analytics-service, transcript-service    │
│ `call.completed`     │ call-service         │ memory-agent, transcript-service, billing│
│ `call.failed`        │ call-service         │ notification-service, analytics-service  │
│ `auth.registered`    │ auth-service         │ notification-service (Welcome email)     │
│ `auth.otp.sent`      │ auth-service         │ sms-agent, audit-service                 │
│ `billing.invoice.paid`│ billing-service     │ notification-service, user-service       │
│ `mcp.tool.executed`  │ mcp-manager          │ audit-service, analytics-service         │
└──────────────────────┴──────────────────────┴──────────────────────────────────────────┘
```
