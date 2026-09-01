# 🛠️ Platform Tech Stack Decisions & Runtimes

Definitive summary of runtime versions, frameworks, databases, and third-party AI provider ecosystems powering Nexora.

---

## 1. Runtime Versions & Frameworks

- **Node.js**: `v20.x LTS` (TypeScript `5.x`, pnpm `9.x`)
- **Python**: `3.11.x` (LiveKit Agents, PyTorch, Pydantic, NumPy)
- **Go**: `1.22.x` (SIP stack, gRPC Protobuf 3)
- **Flutter**: `3.x` (Dart `3.x`)
- **Next.js**: `14.x` (React 19, TailwindCSS, `@nexora/ui`)

---

## 2. Infrastructure & Database Engines

- **PostgreSQL**: `16.x` with `pgvector` extension
- **Redis**: `7.2.x` (Cache, token bucket rate limits, Redlock)
- **MongoDB**: `6.0.x` (Semi-structured memory facts and conversation logs)
- **NATS JetStream**: `2.10.x` (Event streaming and internal pub/sub)
- **Apache Kafka**: `3.6.x` (Transcript stream ingestion)
- **MinIO**: `RELEASE.2024` (S3 API compatible object storage)
- **OpenSearch**: `2.13.x` (BM25 and vector hybrid search)
- **Temporal**: `1.24.x` (Durable multi-step workflow engine)
- **HashiCorp Vault**: `1.15.x` (KV v2 secrets management)

---

## 3. External AI & Telephony Providers

- **Speech-to-Text (STT)**: Deepgram Nova-2 (Streaming & Prerecorded)
- **Text-to-Speech (TTS)**: Cartesia Sonic (<100ms TTFB) & ElevenLabs Flash v2.5
- **Large Language Models (LLM)**: OpenAI (GPT-4o), Anthropic (Claude 3.5 Sonnet), Groq (Llama 3.3 70B), Google Gemini (2.0 Flash), DeepSeek (V3/R1), Mistral, Together AI, Cohere, Perplexity, xAI Grok.
- **Telephony Carriers**: Twilio, Telnyx, Sinch, Plivo.
- **Messaging**: Meta WhatsApp Cloud API, Twilio WhatsApp, Gupshup.
- **Payment Gateways**: Stripe, Razorpay, LemonSqueezy, Paddle.
