# 🚀 Platform Delivery Phases Master Index

The complete 9-phase engineering roadmap turning the Nexora vision into an operational enterprise platform across 28 two-week sprints.

---

## 📑 Delivery Phases Matrix

| Phase | Title & Focus Area | Sprints | Weeks | Specification Document | Key Deliverables |
| :---: | :--- | :---: | :---: | :--- | :--- |
| **0** | **Foundations & Infrastructure** | 1–2 | 1–4 | [**`phase-0-foundations-and-infrastructure.md`**](./phase-0-foundations-and-infrastructure.md) | 12 `@nexora/*` packages, Docker Compose dev environment, HashiCorp Vault secrets, Auth/User services. |
| **1** | **MVP Core Loop (5 Agents)** | 3–6 | 5–12 | [**`phase-1-mvp-core-loop.md`**](./phase-1-mvp-core-loop.md) | Go telephony, Python voice streaming, Silero VAD, Deepgram STT, Cartesia TTS, 5-agent MVP. |
| **2** | **MVP Client Surfaces** | 7–9 | 13–18 | [**`phase-2-mvp-surfaces.md`**](./phase-2-mvp-surfaces.md) | Flutter Mobile App (Android/iOS), Next.js User Dashboard, Admin Panel v1, 3-min Landing Page demo. |
| **3** | **Omnichannel Expansion** | 10–12 | 19–24 | [**`phase-3-messaging-expansion.md`**](./phase-3-messaging-expansion.md) | WhatsApp Business Cloud API, Twilio/Sinch SMS, Email intelligence, Unified customer timeline. |
| **4** | **Admin, Analytics & Billing** | 13–15 | 25–30 | [**`phase-4-admin-analytics-billing.md`**](./phase-4-admin-analytics-billing.md) | Stripe/Razorpay billing, minute metering, 9 Admin Operations modules, Carrier P&L analytics. |
| **5** | **Specialized AI Agents (20 Total)**| 16–20 | 31–40 | [**`phase-5-specialized-agents.md`**](./phase-5-specialized-agents.md) | Full 20 specialized Python agents, dynamic Tool Registry, Model Context Protocol (MCP) bridge. |
| **6** | **Desktop App & Enterprise Portal**| 21–23 | 41–46 | [**`phase-6-desktop-enterprise-portal.md`**](./phase-6-desktop-enterprise-portal.md) | Flutter Desktop softphone (macOS/Windows), SAML/OIDC SSO, SCIM 2.0, dedicated SIP trunks. |
| **7** | **Security, Compliance & Scale** | 24–26 | 47–52 | [**`phase-7-security-compliance-scale.md`**](./phase-7-security-compliance-scale.md) | SOC 2 / HIPAA compliance, immutable audit logs, PII redaction, 10k concurrent call load test. |
| **8** | **Launch Readiness** | 27–28 | 53–56 | [**`phase-8-launch-readiness.md`**](./phase-8-launch-readiness.md) | Multi-region production cluster, canary deployments, developer docs portal, 99.99% SLA. |

---

## 🎯 Recommended Checkpoint Milestones

- **End of Phase 2 (~Week 18)**: Working MVP — AI dials/answers calls, calendar booking, support CRM on mobile & web.
- **End of Phase 4 (~Week 30)**: Voice + WhatsApp + SMS + Email + Billing maturity. Production-ready SaaS v1.
- **End of Phase 8 (~Week 56)**: Enterprise-scale platform — 20 agents, desktop softphone, SOC 2 compliance, multi-region high availability.



┌────────────────────────────────────────────────────────────────────────────────────────┐
│                         Nexora 24-Microservice Ecosystem                               │
└───────────────────────────────────────────┬────────────────────────────────────────────┘
                                            │
   ┌────────────────────┬───────────────────┼───────────────────┬────────────────────┐
   ▼                    ▼                   ▼                   ▼                    ▼
[ Phase 1: IAM ]    [ Phase 2: Media ]  [ Phase 3: AI ]    [ Phase 4: Comms ]   [ Phase 5: Billing ]
- auth-service (✅)  - telephony (Go)    - ai-orchestrator  - notification-svc   - billing-service
- user-service      - voice-service     - memory-service   - transcript-service - payment-service
- api-gateway       - call-service      - spam-service     - recording-service  - workflow-service
                    - websocket-service                    - contact-service    - scheduler-service
