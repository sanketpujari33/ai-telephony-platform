# 🧠 Nexora Architecture & Packages Deep Study

This document provides a comprehensive technical breakdown of the core platform architecture and the 12 foundational `@nexora/*` packages.

---

## 1. Executive Summary & Vision

The **Nexora AI Communication Platform** is an enterprise-grade, polyglot AI Communication Operating System. It powers conversational agents across **real-time voice telephony (WebRTC, SIP, PSTN)**, **omnichannel messaging (WhatsApp, SMS, Email, Push)**, **durable multi-agent workflows**, and **multi-tenant business intelligence**.

### Core Architecture Philosophy
1. **Language Split by Technical Responsibility**:
   - ⚡ **Go (`telephony-service`, `voice-platform`)**: Ultra-low-latency WebRTC media pipelines, SIP trunking, LiveKit integration, and real-time audio packet routing.
   - 🤖 **Python (`ai-orchestrator`, `ai-platform`)**: Multi-agent coordination, streaming STT/TTS pipelines (Pipecat, Deepgram, Cartesia, ElevenLabs), memory embeddings, and LangGraph/Temporal durable execution.
   - 🚀 **TypeScript / Node.js (`api-gateway`, backend microservices, packages)**: API aggregation, auth, billing, webhooks, database persistence, and shared runtime libraries.
   - 📱 **Flutter / Dart (`mobile-app`, `desktop`)**: Cross-platform Android, iOS, and Desktop client applications sharing 80%+ codebase.
   - 🌐 **Next.js (`user-dashboard`, `admin-panel`, `landing-page`, `documentation`)**: SSR/SSG web interfaces, tenant control centers, and developer portal.

---

## 2. Deep Dive: The 12 `@nexora/*` Core Packages

```text
packages/
├── types/          ── Zero-dependency compile-time single source of truth (18 domain modules)
├── config/         ── Zod-validated env + HashiCorp Vault secrets integration
├── logger/         ── High-performance Pino logger with OTel tracing & PII redaction
├── shared/         ── Async primitives, 6 billing HMAC verifiers, phone parsing, Result monad
├── events/         ── Multi-transport Event Bus (Redis, NATS JetStream, Memory) + DLQ
├── database/       ── Polyglot DB pooling (PostgreSQL + pgvector, MongoDB, Redis) & repositories
├── ai/             ── Multi-LLM Engine with 13 providers + ModelRouter fallback chains
├── notifications/  ── Omnichannel Router (SMS, WhatsApp, Email, Push) across 22+ providers
├── prompts/        ── Versioned template compiler & prompt-injection security guardrails
├── tools/          ── Agent tool registry (OpenAI / Anthropic spec) + MCP server adapter
├── sdk/            ── Isomorphic TypeScript SDK for web, mobile, and backend microservices
└── ui/             ── Premium React component library & design system
```

---

### 2.1 `@nexora/types` — Centralized Domain Contracts

- **Zero Runtime Overhead**: Pure compile-time interfaces, types, and enums.
- **Module Structure**:
  - `agent.ts`: `AgentConfig`, `AgentSession`, `LLMProvider`, `TTSProvider`, `STTProvider`.
  - `ai.ts`: `AIMessage`, `AIOptions`, `AIResponse`, `AIStreamChunk`, `EmbeddingResult`, `IAIProvider`.
  - `billing.ts`: `Plan`, `BillingPeriod`, `BillingProvider`, `Subscription`, `UsageMeter`, `PaymentIntent`.
  - `call.ts`: `Call`, `CallDirection`, `CallStatus`, `CallRecording`, `TranscriptUtterance`.
  - `config.ts`: `AppConfig`, `SecretsMap`, `VaultOptions`, `DatabaseSecrets`, `AISecrets`, etc.
  - `events.ts`: `BaseEvent`, `EventChannel`, `IEventBus`, `EventHandler`, `EventSubscription`.
  - `notification.ts`: `ISMSProvider`, `IWhatsAppProvider`, `IEmailProvider`, `IPushProvider`, `DeliveryReceipt`.
  - `prompts.ts`: `PromptTemplate`, `PromptVariableMap`, `CompiledPrompt`.
  - `tool.ts`: `ToolDefinition`, `ToolParameters`, `ExecutableTool`, `MCPTool`, `ToolExecutionContext`.
  - `ui.ts`: `GridColumn`, `ThemeMode`, `ColorPalette`, `NavMenuItem`, `ChartSeriesConfig`.
  - `user.ts`, `workflow.ts`, `storage.ts`, `security.ts`, `telemetry.ts`, `pagination.ts`, `health.ts`.

---

### 2.2 `@nexora/config` — HashiCorp Vault Secret Management

- **Zero-Secret Commits**: All sensitive API keys, database credentials, and certificates are fetched dynamically from HashiCorp Vault KV v2 (`http://localhost:8200` in dev).
- **Multi-Environment Auth**: Token auth (dev), AppRole (prod), Kubernetes service account JWT.
- **Typed Secrets Catalog (`SecretsMap`)**:
  - **Database**: PostgreSQL URL, Redis URL, MongoDB URL.
  - **Auth**: JWT Secret, JWT Refresh Secret.
  - **AI (13 Providers)**: OpenAI, Anthropic, Groq, Gemini, DeepSeek, Mistral, Together, Cohere, Perplexity, xAI, Fireworks, Cerebras, SambaNova, Azure OpenAI, AWS Bedrock.
  - **Speech**: Deepgram, Cartesia, ElevenLabs, Sarvam AI, Azure Speech.
  - **Telephony**: Twilio, Telnyx.
  - **Messaging**: Meta WhatsApp Cloud API, Gupshup.
  - **Notifications**: Twilio, Telnyx, Plivo, AWS SNS, Vonage, Infobip, SendGrid, Resend, Postmark, Mailgun, Brevo, SMTP, FCM, APNs, OneSignal, Expo.
  - **Billing**: Stripe, Razorpay, LemonSqueezy, Paddle, PayPal, Paystack, Flutterwave, Coinbase Commerce.
  - **Infrastructure**: LiveKit, S3/MinIO, Sentry, Truecaller.

---

### 2.3 `@nexora/ai` — Multi-LLM Engine (13 Providers)

| Provider | Supported Models | Strengths |
| :--- | :--- | :--- |
| **OpenAI** | GPT-4o, GPT-4o-mini, text-embedding-3 | High reasoning, tool calling, multi-modal |
| **Anthropic** | Claude 3.5 Sonnet, Claude 3.5 Haiku | Complex instructions, structured tool output |
| **Groq** | Llama 3.3 70B, Llama 3.1 8B, Mixtral | Ultra-low-latency real-time voice responses |
| **Google Gemini** | Gemini 1.5 Pro, Gemini 1.5 Flash, 2.0 Flash | Long context windows (1M+ tokens), multimodal |
| **DeepSeek** | DeepSeek-V3, DeepSeek-R1 | Cost-effective advanced mathematical & logical reasoning |
| **Mistral AI** | Mistral Large, Codestral, Pixtral | Multilingual reasoning, European data residency |
| **Together AI** | Llama 3.3 70B, Qwen 2.5 72B | Open-source ecosystem, custom fine-tunes |
| **Cohere** | Command R+, Command R, Embed v3 | Enterprise search, RAG, multilingual embeddings |
| **Perplexity** | Sonar Large, Sonar Online | Live web search citation and verification |
| **xAI** | Grok-2, Grok-Vision | Real-time truth verification, vision reasoning |
| **Fireworks AI** | Fireworks Llama 3.3 70B | High-concurrency low-latency API inference |
| **Cerebras** | Llama 3.3 70B | Industry-leading 2,000+ tokens/sec generation |
| **Ollama** | Local Llama 3.2, Mistral, Qwen | Privacy-first air-gapped on-premise execution |

- **`ModelRouter`**: Automatic fallback chains (e.g., Groq → DeepSeek → OpenAI) on 429 rate limits, timeouts, or transient network errors.

---

### 2.4 `@nexora/shared` — Cryptographic Billing & Platform Utilities

- **Billing Webhook Verifiers**:
  - **Stripe**: `stripe-signature` timestamp + HMAC-SHA256 verification.
  - **Razorpay**: `x-razorpay-signature` HMAC-SHA256 verification.
  - **LemonSqueezy**: `x-signature` raw payload verification.
  - **Paddle Billing**: `Paddle-Signature` timestamped HMAC-SHA256 verifier.
  - **PayPal**: `PAYPAL-TRANSMISSION-SIG` CRC32/RSA signature verification.
  - **Paystack**: `x-paystack-signature` HMAC-SHA512 verification.
  - **Flutterwave**: `verif-hash` secret hash verification.
  - **Coinbase Commerce**: `X-CC-Webhook-Signature` HMAC-SHA256 verifier.
- **Subscription Math**:
  - `calculateOverageCost(used, included, rateCents)`
  - `calculateTax(subtotalCents, ratePercentage)`
  - `calculateProratedAmount(planRateCents, daysRemaining, totalDays)`

---

### 2.5 `@nexora/notifications` — Intelligent Multi-Channel Router

- **Automated Channel Routing**:
  - **SMS**: Twilio → Telnyx → Plivo → Vonage → Infobip → AWS SNS → Bandwidth.
  - **WhatsApp**: Meta Cloud API → Twilio WhatsApp → Gupshup → 360dialog.
  - **Email**: SendGrid → Resend → Postmark → AWS SES → Mailgun → Brevo → SMTP.
  - **Push**: Firebase FCM → Apple APNs → OneSignal → Expo Push → WebPush.
- **Failover Logic**: When a primary provider fails with rate-limit or outage errors, the router automatically re-routes to healthy secondary providers and emits delivery metrics.

---

### 2.6 `@nexora/prompts` & `@nexora/tools`

- **Prompt Compiler**: Fast interpolation with conditional rendering (`{{#if is_vip}}...{{/if}}`).
- **Security Guardrails**:
  - Detects prompt injection (`ignore previous instructions`, `system override`, `DAN mode`).
  - Detects PII (Credit Cards, SSNs, Passwords).
  - Flags harmful or offensive requests.
- **Tool Ecosystem**:
  - Full adherence to OpenAI Function Calling & Anthropic Tool-Use specs.
  - Built-in tools for Call Transfer, End Call, DTMF, Calendar Booking, CRM Lookup, Send Message, Knowledge Base Search, Invoice Status.
  - **MCP Adapter**: Connects external Model Context Protocol servers dynamically.

---

### 2.7 `@nexora/ui` — Zero-Wrapper Component Library

- **Smart Components**: `SmartPaste`, `SmartTextArea` with inline generative autocomplete.
- **DataGrid**: Multi-column sorting, global search, filtering, custom column templates, CSV/Excel export, pagination.
- **Calendars & Schedulers**: Full interactive calendar, date/time pickers, resource scheduling.
- **Data Visualizations**: Line charts, bar charts, area charts, donut/pie charts, sparklines.
- **Controls & Modals**: Complete suite of dark-mode optimized, accessible UI primitives.
