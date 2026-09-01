# 🏗️ Nexora Platform Entrypoints, Build Targets & Release Manifest

This document provides the definitive directory mapping of all application and microservice entrypoints, runtime configurations, port matrix, environment variables, build outputs, and multi-platform distribution targets in [`.github/release/`](file:///d:/trideocoders/ai-telephony-platform/.github/release).

---

## 🎯 1. Platform Master Entrypoints Matrix

### A. Applications Layer (`apps/`)

| Application | Entrypoint File | Port | Runtime | Output Target |
| :--- | :--- | :---: | :---: | :--- |
| **`apps/user-dashboard`** | [`apps/user-dashboard/app/page.tsx`](file:///d:/trideocoders/ai-telephony-platform/apps/user-dashboard/app/page.tsx) | `3000` | Next.js 14 / React 19 | `.github/release/web/user-dashboard` |
| **`apps/admin-panel`** | [`apps/admin-panel/app/page.tsx`](file:///d:/trideocoders/ai-telephony-platform/apps/admin-panel/app/page.tsx) | `3001` | Next.js 14 / React 19 | `.github/release/web/admin-panel` |
| **`apps/landing-page`** | [`apps/landing-page/app/page.tsx`](file:///d:/trideocoders/ai-telephony-platform/apps/landing-page/app/page.tsx) | `3002` | Next.js 14 / React 19 | `.github/release/web/landing-page` |
| **`apps/documentation`** | [`apps/documentation/app/page.tsx`](file:///d:/trideocoders/ai-telephony-platform/apps/documentation/app/page.tsx) | `3003` | Next.js 14 / React 19 | `.github/release/web/documentation` |
| **`apps/mobile-app`** | [`apps/mobile-app/lib/main.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/main.dart) | - | Flutter 3.x / Dart | `.github/release/android/` & `.github/release/ios/` |
| **`apps/desktop`** | [`apps/desktop/lib/main.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/desktop/lib/main.dart) | - | Flutter 3.x / Dart | `.github/release/mac/` & `.github/release/win/` |

---

### B. Core Backend Microservices Layer (`services/` — 24 Microservices)

| Microservice | Main Entrypoint | Port | Tech Stack | Primary Responsibility |
| :--- | :--- | :---: | :---: | :--- |
| **`telephony-service`** | [`services/telephony-service/main.go`](file:///d:/trideocoders/ai-telephony-platform/services/telephony-service/main.go) | `8085` / `5060` | Go 1.22 / Gin | SIP Protocol Trunking, Inbound/Outbound Calls, LiveKit Bridge |
| **`ai-orchestrator`** | [`services/ai-orchestrator/main.go`](file:///d:/trideocoders/ai-telephony-platform/services/ai-orchestrator/main.go) | `8086` / `50051`| Go 1.22 / Gin / gRPC | Real-Time Intent Classification & Sub-Agent Dispatch |
| **`voice-service`** | [`services/voice-service/agent.py`](file:///d:/trideocoders/ai-telephony-platform/services/voice-service/agent.py) | - | Python 3.11 / LiveKit | LiveKit WebRTC Voice Worker, Silero VAD, Deepgram STT, Cartesia TTS |
| **`api-gateway`** | [`services/api-gateway/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/api-gateway/src/index.ts) | `4000` | Node.js / Express | Route Proxy, JWT Verification, Token Bucket Rate Limiting |
| **`auth-service`** | [`services/auth-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/auth-service/src/index.ts) | `4001` | Node.js / Express | JWT Issuance, OTP Phone Verification, Password Hashing, RBAC |
| **`call-service`** | [`services/call-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/call-service/src/index.ts) | `4002` | Node.js / Express | Call Lifecycle State Machine, LiveKit Room Allocation |
| **`user-service`** | [`services/user-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/user-service/src/index.ts) | `4003` | Node.js / Express | User Profiles, Workspace Multi-Tenancy, Team Roles |
| **`billing-service`** | [`services/billing-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/billing-service/src/index.ts) | `4004` | Node.js / Express | Minute Usage Metering, Tiered Subscriptions, Invoices |
| **`notification-service`**| [`services/notification-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/notification-service/src/index.ts) | `4005` | Node.js / Express | Multi-Channel Dispatcher (Push, SMS, Email, WhatsApp) |
| **`contact-service`** | [`services/contact-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/contact-service/src/index.ts) | `4006` | Node.js / Express | Multi-Tenant CRM Contact Directory & In-Call Notes |
| **`transcript-service`**| [`services/transcript-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/transcript-service/src/index.ts) | `4007` | Node.js / Express | Audio Transcript Archival & OpenSearch Indexing |
| **`analytics-service`** | [`services/analytics-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/analytics-service/src/index.ts) | `4008` | Node.js / Express | Carrier P&L Margins, MOS Audio Quality Metrics |
| **`websocket-service`** | [`services/websocket-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/websocket-service/src/index.ts) | `4009` | Node.js / WS | Real-time 60fps Audio Waveform & Telemetry Streaming |
| **`payment-service`** | [`services/payment-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/payment-service/src/index.ts) | `4010` | Node.js / Express | Stripe/Razorpay Webhooks & Pre-paid Credit Top-Ups |
| **`storage-service`** | [`services/storage-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/storage-service/src/index.ts) | `4011` | Node.js / Express | S3 / MinIO Presigned URL Dispatcher & Asset Retention |
| **`spam-service`** | [`services/spam-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/spam-service/src/index.ts) | `4012` | Node.js / Express | Truecaller Reputation Signals & Spam Risk Scoring |
| **`memory-service`** | [`services/memory-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/memory-service/src/index.ts) | `4013` | Node.js / Express | Semantic Fact Extraction & MongoDB Vector Memory |
| **`admin-service`** | [`services/admin-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/admin-service/src/index.ts) | `4014` | Node.js / Express | 12-Role RBAC Matrix, Banners, Promoters, General Settings |
| **`audit-service`** | [`services/audit-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/audit-service/src/index.ts) | `4015` | Node.js / Express | Immutable SOC 2 / HIPAA Audit Logs & Tamper Trails |
| **`enterprise-service`**| [`services/enterprise-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/enterprise-service/src/index.ts) | `4016` | Node.js / Express | SAML 2.0 / OIDC SSO & SCIM 2.0 Directory Sync |
| **`recording-service`** | [`services/recording-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/recording-service/src/index.ts) | `4017` | Node.js / Express | Call Audio Ingestion, MinIO Uploads & Waveform Generation |
| **`scheduler-service`** | [`services/scheduler-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/scheduler-service/src/index.ts) | `4018` | Node.js / Express | Scheduled Outbound Dialer Campaigns & Cron Triggers |
| **`search-service`** | [`services/search-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/search-service/src/index.ts) | `4019` | Node.js / Express | OpenSearch BM25 Full-Text & Hybrid Vector Retrieval |
| **`workflow-service`** | [`services/workflow-service/src/index.ts`](file:///d:/trideocoders/ai-telephony-platform/services/workflow-service/src/index.ts) | `4020` | Node.js / Express | Temporal Durable Workflows & Multi-Step Agent Graphs |

---

### C. 20 Specialized Python AI Agents (`agents/`)

```text
agents/
├── phone-agent/          -> agents/phone-agent/agent.py         [Real-time Voice STT/TTS]
├── orchestrator-agent/   -> agents/orchestrator-agent/agent.py  [Intent Classifier & Router]
├── support-agent/        -> agents/support-agent/agent.py       [Customer Service & CRM]
├── calendar-agent/       -> agents/calendar-agent/agent.py      [Scheduling & Free/Busy]
├── appointment-agent/    -> agents/appointment-agent/agent.py   [Slot Reservation]
├── whatsapp-agent/       -> agents/whatsapp-agent/agent.py      [WhatsApp Cloud API]
├── sms-agent/            -> agents/sms-agent/agent.py           [Twilio / Sinch SMS]
├── email-agent/          -> agents/email-agent/agent.py         [Email Summarization & Drafts]
├── finance-agent/        -> agents/finance-agent/agent.py       [Invoices & Balances]
├── enterprise-agent/     -> agents/enterprise-agent/agent.py    [Enterprise Policies & Audit]
├── knowledge-agent/      -> agents/knowledge-agent/agent.py     [Vector RAG QA Engine]
├── memory-agent/         -> agents/memory-agent/agent.py        [Long-term Memory Facts]
├── research-agent/       -> agents/research-agent/agent.py      [Web Research & Scraping]
├── reminder-agent/       -> agents/reminder-agent/agent.py      [Scheduled Follow-ups]
├── travel-agent/         -> agents/travel-agent/agent.py        [Flights & Hospitality]
├── shopping-agent/       -> agents/shopping-agent/agent.py      [Product Recommendations]
├── delivery-agent/       -> agents/delivery-agent/agent.py      [Shipment Tracking]
├── automation-agent/     -> agents/automation-agent/agent.py    [Temporal & Webhook Tasks]
├── coding-agent/         -> agents/coding-agent/agent.py        [Technical API Generator]
└── browser-agent/        -> agents/browser-agent/agent.py       [Headless Web Navigation]
```

---

## 📦 2. Multi-Platform Release Artifacts (`.github/release/`)

```text
.github/release/
├── release-manifest.json             [Master multi-platform manifest]
├── android/
│   └── release-info.json             [Flutter Android APK/AAB build target]
├── ios/
│   └── release-info.json             [Flutter iOS IPA build target]
├── mac/
│   └── release-info.json             [Flutter macOS DMG build target]
├── ubuntu/
│   └── release-info.json             [Ubuntu Linux AppImage/Docker build target]
├── web/
│   └── release-info.json             [Next.js 14 SSR & Static build target]
└── win/
    └── release-info.json             [Windows 10/11 Softphone EXE/MSI target]
```
