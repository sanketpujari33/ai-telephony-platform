# 📋 Platform Tasks Catalog: Letters N – S

Comprehensive specification and tracking for platform capabilities and components categorized alphabetically from N through S.

---

## 🅽 Letter N
- **`notification-service`** (`services/notification-service`): Multi-channel dispatcher for SMS, WhatsApp, Email, and Push notifications.
- **`@nexora/notifications`**: Channel router supporting 22+ communication providers with automatic failover.
- **NATS JetStream**: High-throughput event streaming bus for domain events and real-time pub/sub.

---

## 🅾️ Letter O
- **`orchestrator-agent`** (`agents/orchestrator-agent`): Master coordinator classifying user intent and delegating turns to specialist agents.
- **OpenSearch 2.13**: Full-text and BM25 vector hybrid search over transcripts and knowledge bases.

---

## 🅿️ Letter P
- **`payment-service`** (`services/payment-service`): Stripe, Razorpay, Paddle, and LemonSqueezy webhook handlers and checkout sessions.
- **`phone-agent`** (`agents/phone-agent`): WebRTC live voice conversational agent with Silero VAD, Deepgram STT, and Cartesia TTS.
- **`@nexora/prompts`**: Versioned template compiler and prompt injection security guardrails.
- **Profit & Loss Manager** (`apps/admin-panel`): Real-time per-minute carrier telephony margin calculation.
- **Promoters Manager** (`apps/admin-panel`): Affiliate tracking, referral links, and commission payout manager.

---

## 🆀 Letter Q
- **Queue Manager** (`services/call-service`): Inbound call hold queues, ACD (Automatic Call Distribution), and agent routing.

---

## 🆁 Letter R
- **`recording-service`** (`services/recording-service`): Call audio recording ingestion, S3/MinIO upload, and waveform extraction.
- **`reminder-agent`** (`agents/reminder-agent`): Temporal durable scheduled follow-up alerts and appointment reminders.
- **`research-agent`** (`agents/research-agent`): Live web search and fact verification via Perplexity and Grok APIs.
- **Role-Based Access Control (RBAC)**: 12-Role permission matrix with two-tier Redis permission caching.

---

## 🆂 Letter S
- **`scheduler-service`** (`services/scheduler-service`): Scheduled outbound dialer campaigns and recurring job execution.
- **`search-service`** (`services/search-service`): Search index abstraction over OpenSearch and PostgreSQL.
- **`shopping-agent`** (`agents/shopping-agent`): E-commerce product catalog search, cart operations, and order inquiries.
- **`sms-agent`** (`agents/sms-agent`): Concise transactional and conversational SMS messaging.
- **`spam-service`** (`services/spam-service`): Caller reputation scoring, spam detection, and Truecaller verification.
- **`storage-service`** (`services/storage-service`): MinIO / S3 asset management and presigned URL generation.
- **`support-agent`** (`agents/support-agent`): Customer service CRM resolution, dispute management, and policy FAQ.
