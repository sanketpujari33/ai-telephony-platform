# 🌐 Nexora Platform Master API Catalog & Specification

Comprehensive REST, WebSocket, and WebRTC API specifications across all applications, microservices, and external platform integrations.

---

## 🚀 0. Interactive Swagger UI & OpenAPI 3.0 Portals (Live Testing)

The Nexora API Gateway (:8080) exposes live, interactive Swagger UI portals with "Try it out" live execution support:

| Documentation Portal | Swagger UI Interactive URL | OpenAPI 3.0 Schema JSON | Target Audience / Surface |
| :--- | :--- | :--- | :--- |
| **🌟 Master A–Z API Explorer** | [`http://localhost:8080/docs`](http://localhost:8080/docs) | [`/api-docs/openapi.json`](http://localhost:8080/api-docs/openapi.json) | Complete Unified Platform API Catalog |
| **👔 Platform Admin Operations** | [`http://localhost:8080/docs/admin`](http://localhost:8080/docs/admin) | [`/api-docs/admin.json`](http://localhost:8080/api-docs/admin.json) | Operations Center, RBAC, CMS, Helpdesk |
| **📱 User Mobile Application** | [`http://localhost:8080/docs/user-mobile`](http://localhost:8080/docs/user-mobile) | [`/api-docs/mobile.json`](http://localhost:8080/api-docs/mobile.json) | Flutter Mobile (Android / iOS) Softphone |
| **💻 User Desktop Softphone** | [`http://localhost:8080/docs/user-desktop`](http://localhost:8080/docs/user-desktop) | [`/api-docs/desktop.json`](http://localhost:8080/api-docs/desktop.json) | Desktop 4-Line Softphone & Waveforms |
| **🌐 User Web Dashboard** | [`http://localhost:8080/docs/user-dashboard`](http://localhost:8080/docs/user-dashboard) | [`/api-docs/dashboard.json`](http://localhost:8080/api-docs/dashboard.json) | Next.js 14 User Portal & Dialers |
| **🏢 Enterprise & Governance** | [`http://localhost:8080/docs/enterprise`](http://localhost:8080/docs/enterprise) | [`/api-docs/enterprise.json`](http://localhost:8080/api-docs/enterprise.json) | SAML 2.0 / Okta SSO, SCIM 2.0, SOC 2 |

---

## 📌 1. Universal Communication Protocols & Base URLs

| Environment | Base URL | Protocol | Usage |
| :--- | :--- | :--- | :--- |
| **API Gateway (Dev)** | `http://localhost:8080` | HTTP / REST | Local development & unit tests |
| **WebSocket Gateway (Dev)** | `ws://localhost:8093/ws` | WebSocket JSON | Real-time live screening & whispers |
| **LiveKit WebRTC SFU (Dev)** | `ws://localhost:7880` | WebRTC / WSS | Real-time low-latency audio streams |
| **LiveKit SIP Gateway (Dev)** | `localhost:5060` | SIP (UDP/TCP) | Telecom carrier PSTN call ingress |
| **Production API Cluster** | `https://api.nexora.ai` | HTTPS / WSS | Global production cluster |
| **Production WebRTC SFU** | `wss://livekit.nexora.ai` | WebRTC / WSS | Production audio streams |


---

## 🔐 2. Authentication & Header Standards

All authenticated endpoints require standard headers:

```http
# User Session Requests (Web / Mobile / Desktop)
Authorization: Bearer <JWT_ACCESS_TOKEN>
x-tenant-id: <TENANT_UUID>
Content-Type: application/json
x-correlation-id: req_98fbc182-42ab-41c8-a921-99c019284729

# Service / Machine-to-Machine API Requests
x-api-key: nx_live_9a8b7c6d5e4f3a2b1c0987654321
x-tenant-id: <TENANT_UUID>
Content-Type: application/json
x-correlation-id: req_98fbc182-42ab-41c8-a921-99c019284729
```

---

## 📑 3. Master Platform API Specifications (A to Z)

| Spec # | Specification Document | Core Domains & Endpoints Covered |
| :--- | :--- | :--- |
| **01** | [**`01_master_api_reference_and_gateway_routing.md`**](./01_master_api_reference_and_gateway_routing.md) | Global API gateway topology, service port catalog (:8081-:8105), standard envelopes, and route mappings. |
| **02** | [**`02_telephony_dual_sim_and_carrier_forwarding_api.md`**](./02_telephony_dual_sim_and_carrier_forwarding_api.md) | Dual-SIM management, carrier detection (Jio, Airtel, Vi, BSNL, US carriers), MMI codes, and failure telemetry. |
| **03** | [**`03_inbound_screening_whisper_and_takeover_api.md`**](./03_inbound_screening_whisper_and_takeover_api.md) | "Green Card" live call screening, <18ms whisper prompt injection, and <15ms 1-tap takeover barge-in. |
| **04** | [**`04_autonomous_outbound_missions_api.md`**](./04_autonomous_outbound_missions_api.md) | hiRobin parity outbound phone missions, automated IVR navigation, hold music detection, and confirmation parsing. |
| **05** | [**`05_ai_agents_personas_and_voice_pipeline_api.md`**](./05_ai_agents_personas_and_voice_pipeline_api.md) | 20 AI Agent profiles, persona customization, Cartesia Sonic / Deepgram Nova-2 / LiveKit tokens, and multilingual dialects. |
| **06** | [**`06_memory_knowledge_and_rag_vectors_api.md`**](./06_memory_knowledge_and_rag_vectors_api.md) | Episodic AI memory facts (delivery PINs, gate codes), PostgreSQL `pgvector` HNSW semantic search, and RAG document ingestion. |
| **07** | [**`07_billing_wallets_subscriptions_and_metering_api.md`**](./07_billing_wallets_subscriptions_and_metering_api.md) | Subscriptions, minute wallets, auto-recharge rules, Stripe/Razorpay payment intents, and per-call metering ledgers. |
| **08** | [**`08_crm_contacts_whitelist_and_spam_shield_api.md`**](./08_crm_contacts_whitelist_and_spam_shield_api.md) | Contact CRM, VIP bypass rules, custom contact greetings, and TRAI 140/160 series telemarketer reputation checks. |
| **09** | [**`09_enterprise_governance_sso_and_compliance_api.md`**](./09_enterprise_governance_sso_and_compliance_api.md) | SAML 2.0 / OIDC SSO configuration, SCIM 2.0 user directory provisioning, 12 RBAC role assignments, and audit logs. |
| **10** | [**`10_websocket_and_webrtc_streaming_protocol.md`**](./10_websocket_and_webrtc_streaming_protocol.md) | WebSocket framing protocol, real-time event subscriptions, and LiveKit WebRTC audio room streaming. |
| **11** | [**`11_universal_client_bootstrap_and_remote_config_api.md`**](./11_universal_client_bootstrap_and_remote_config_api.md) | Single-call client bootstrap endpoint (`GET /api/v1/bootstrap`): Feature flags, app version force upgrade, maintenance mode, payment gateways, and carrier presets. |
| **12** | [**`12_dynamic_cms_legal_faqs_and_compliance_content_api.md`**](./12_dynamic_cms_legal_faqs_and_compliance_content_api.md) | Dynamic CMS & legal documentation API (`GET /api/v1/content/:slug`, `GET /api/v1/content/faqs`): DB-driven Privacy Policy, Terms, About Us, FAQs, and role-based Admin Panel CRUD editing. |
| **13** | [**`13_ai_support_bot_and_live_admin_chat_api.md`**](./13_ai_support_bot_and_live_admin_chat_api.md) | AI Support Bot & Live Admin Chat Helpdesk API (`POST /api/v1/support/chat/sessions`, `GET /api/v1/admin/support/chats`): Tier-1 RAG Bot, bot-to-human escalation, canned quick replies, and CSAT ratings. |

---

## 📱 4. Client Application Specific API Specifications

| Application Surface | Specifications & Documentation | Core Responsibilities |
| :--- | :--- | :--- |
| **Mobile App (Flutter)** | [**`mobile-app/api-spec.md`**](./mobile-app/api-spec.md) | Flutter Mobile WebRTC live voice sessions, FCM/APNs push registrations, OTP auth, Dual-SIM dialing. |
| **User Dashboard (Next.js)** | [**`user-dashboard/api-spec.md`**](./user-dashboard/api-spec.md)<br/>[**`user-dashboard/user-roles-api.md`**](./user-dashboard/user-roles-api.md) | Voice agent management, softphone dialer, call transcripts, recordings, contact CRM. |
| **Admin Panel (Next.js)** | [**`admin-panel/api-spec.md`**](./admin-panel/api-spec.md)<br/>[**`admin-panel/rbac-matrix.md`**](./admin-panel/rbac-matrix.md)<br/>[**`admin-panel/roles-permissions-crud.md`**](./admin-panel/roles-permissions-crud.md)<br/>[**`admin-panel/admin-operations-api.md`**](./admin-panel/admin-operations-api.md)<br/>[**`admin-panel/admin-model-study.md`**](./admin-panel/admin-model-study.md) | 12-Role RBAC hierarchy, dynamic permission CRUD, Banners, Promoters, Funds, P&L, Dead Reports, Keywords. |
| **Enterprise Portal** | [**`enterprise/api-spec.md`**](./enterprise/api-spec.md) | SAML 2.0 / OIDC SSO, SCIM provisioning, dedicated SIP trunks, audit exports. |
| **Desktop App (Electron)** | [**`desktop/api-spec.md`**](./desktop/api-spec.md) | Multi-line softphone, global keyboard shortcuts, local audio loopback tests. |
| **Landing Page (Next.js)** | [**`landing-page/api-spec.md`**](./landing-page/api-spec.md) | Lead capture, 3-minute interactive WebRTC sandbox, platform status health. |
| **External Integrations** | [**`external-integrations-spec.md`**](./external-integrations-spec.md) | Telecom carrier SIP endpoints, Deepgram STT, Cartesia TTS, WhatsApp Cloud API. |
