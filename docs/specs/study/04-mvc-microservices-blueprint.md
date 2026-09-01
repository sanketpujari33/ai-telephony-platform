# ⚙️ MVC Microservices Blueprint & Tier Mapping Deep Study

This document details the standardized MVC + Layered Service Architecture across all 24 backend microservices under `services/`.

---

## 1. Architectural Philosophy & Layering

All microservices in [`services/`](file:///d:/trideocoders/ai-telephony-platform/services) adopt a standardized **MVC + Layered Service Architecture** to ensure strict separation of concerns, high maintainability, testability, and horizontal scaling.

```text
                              ┌─────────────────────────────┐
                              │   Client Request (HTTP/WS)  │
                              └──────────────┬──────────────┘
                                             │
                              ┌──────────────▼──────────────┐
                              │       Middlewares Layer     │
                              │ Auth, TenantCtx, RateLimit  │
                              └──────────────┬──────────────┘
                                             │
                              ┌──────────────▼──────────────┐
                              │         Routes Layer        │
                              │ URL Endpoints & HTTP Verbs  │
                              └──────────────┬──────────────┘
                                             │
                              ┌──────────────▼──────────────┐
                              │      Controllers Layer      │
                              │ Req Validation & Res Format │
                              └──────────────┬──────────────┘
                                             │
                              ┌──────────────▼──────────────┐
                              │       Services Layer        │
                              │ Business Logic, AI, Vault   │
                              └──────────────┬──────────────┘
                                             │
                              ┌──────────────▼──────────────┐
                              │        Models Layer         │
                              │  Repositories, DB, Schemas  │
                              └──────────────┬──────────────┘
                                             │
                              ┌──────────────┴──────────────┐
                              │                             │
                        ┌─────▼─────┐                 ┌─────▼─────┐
                        │ PostgreSQL│                 │   Redis   │
                        │ +pgvector │                 │ + NATS    │
                        └───────────┘                 └───────────┘
```

---

## 2. Standard Directory Layout

```text
services/<service-name>/src/
├── routes/             ── Express Router definitions; maps URI paths to Controller methods
├── controllers/        ── Handles HTTP/WS requests, input validation, status codes, and JSON responses
├── services/           ── Core business logic, external API integrations, Vault secrets, EventBus publishing
├── models/             ── Database schemas, TypeORM/Prisma/Postgres repository layer, data entity structures
├── utils/              ── Custom error classes, response formatters, crypto helpers, validation schemas
├── middlewares/        ── JWT verification, tenant context extraction, error handling, rate limiting
└── index.ts            ── Microservice bootstrapper, Express app setup, graceful shutdown hooks
```

---

## 3. Structural Mapping of All 24 Microservices

### 3.1 Tier 1: Gateway, Authentication & Core Telephony (MVP Loop)

| Service | Responsibility | Routes (`routes/`) | Controllers (`controllers/`) | Services (`services/`) | Models (`models/`) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **`api-gateway`** | Central reverse proxy, route aggregation, auth enforcement, rate limiting | `call.routes.ts`, `agent.routes.ts`, `auth.routes.ts`, `health.routes.ts` | `CallGatewayController`, `AgentGatewayController`, `HealthController` | `ProxyService`, `RateLimitService`, `DiscoveryService` | `RouteConfigModel`, `ApiKeyModel` |
| **`auth-service`** | User authentication, JWT tokens, OTP phone verification, RBAC permissions | `auth.routes.ts`, `otp.routes.ts`, `token.routes.ts` | `AuthController`, `OTPController`, `TokenController` | `AuthService`, `OTPService`, `PasswordService`, `VaultTokenService` | `UserModel`, `SessionModel`, `OTPRecordModel` |
| **`call-service`** | Call session management, call state machine (`initiated` → `ringing` → `answered` → `completed`), hold queue | `call.routes.ts`, `session.routes.ts`, `queue.routes.ts` | `CallController`, `SessionController`, `QueueController` | `CallSessionService`, `QueueManagerService`, `RoutingService` | `CallSessionModel`, `CallQueueModel`, `CallHistoryModel` |
| **`telephony-service`** | Real-time SIP trunking, LiveKit WebRTC audio room management, Twilio/Telnyx bridge (Go) | `sip_routes.go`, `webrtc_routes.go`, `health_routes.go` | `SIPController`, `LiveKitController` | `LiveKitRoomService`, `SIPTrunkService`, `PacketRouterService` | `LiveKitRoomModel`, `SIPSessionModel` |
| **`voice-service`** | Streaming STT (Deepgram), TTS (ElevenLabs, Cartesia), VAD audio interruption detection (Python) | `audio_routes.py`, `tts_routes.py`, `stt_routes.py` | `AudioStreamController`, `TTSController`, `STTController` | `PipecatPipelineService`, `DeepgramSTTService`, `CartesiaTTSService` | `AudioFrameModel`, `TranscriptModel` |
| **`ai-orchestrator`** | Multi-agent reasoning loop, tool execution with `@nexora/tools`, prompt compiling with `@nexora/prompts` | `turn.routes.ts`, `agent.routes.ts`, `tool.routes.ts` | `TurnController`, `AgentOrchestratorController` | `AgentTurnService`, `ToolExecutionService`, `MemoryExtractionService` | `AgentSessionModel`, `ConversationStateModel` |

---

### 3.2 Tier 2: Storage, Memory, Transcripts & Recordings

| Service | Responsibility | Routes (`routes/`) | Controllers (`controllers/`) | Services (`services/`) | Models (`models/`) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **`memory-service`** | Semantic conversation memory, vector embeddings via `pgvector` / Qdrant | `memory.routes.ts`, `embedding.routes.ts` | `MemoryController`, `EmbeddingController` | `MemoryStoreService`, `VectorSearchService`, `MemoryExtractionService` | `MemoryRecordModel`, `EmbeddingVectorModel` |
| **`recording-service`** | Call audio recording ingestion, S3/MinIO upload, waveform generation | `recording.routes.ts` | `RecordingController` | `AudioUploadService`, `MinioStorageService`, `WaveformService` | `RecordingModel`, `StorageMetadataModel` |
| **`transcript-service`** | Speech-to-text post-processing, diarization, speaker labeling, summary extraction | `transcript.routes.ts` | `TranscriptController` | `DiarizationService`, `SummarizationService`, `SentimentService` | `TranscriptUtteranceModel`, `CallSummaryModel` |
| **`storage-service`** | MinIO / S3 multi-tenant file asset storage, presigned URLs, retention policies | `storage.routes.ts` | `StorageController` | `S3StorageService`, `PresignedUrlService`, `RetentionService` | `FileAssetModel`, `BucketConfigModel` |

---

### 3.3 Tier 3: Business Logic, Billing, Notifications & Contacts

| Service | Responsibility | Routes (`routes/`) | Controllers (`controllers/`) | Services (`services/`) | Models (`models/`) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **`billing-service`** | Tiered subscription management, minute usage metering, invoice generation | `subscription.routes.ts`, `usage.routes.ts`, `invoice.routes.ts` | `SubscriptionController`, `UsageController`, `InvoiceController` | `SubscriptionService`, `UsageMeterService`, `InvoiceService` | `SubscriptionModel`, `UsageLedgerModel`, `InvoiceModel` |
| **`payment-service`** | Stripe, Razorpay, Paddle, PayPal webhook processing and checkout sessions | `webhook.routes.ts`, `checkout.routes.ts` | `WebhookController`, `CheckoutController` | `StripePaymentService`, `RazorpayService`, `PaddleService` | `PaymentTransactionModel`, `PaymentMethodModel` |
| **`notification-service`** | SMS, WhatsApp, Email, and Push delivery via `@nexora/notifications` | `notification.routes.ts`, `template.routes.ts` | `NotificationController`, `TemplateController` | `DispatchService`, `ChannelFailoverService`, `DeliveryReceiptService` | `NotificationLogModel`, `TemplateModel` |
| **`contact-service`** | Multi-tenant CRM contact directory, address book, caller ID lookup | `contact.routes.ts`, `tag.routes.ts` | `ContactController`, `TagController` | `ContactService`, `CallerIDService`, `ImportExportService` | `ContactModel`, `ContactTagModel` |
| **`user-service`** | User account lifecycle, tenant provisioning, team member invitations | `user.routes.ts`, `tenant.routes.ts`, `member.routes.ts` | `UserController`, `TenantController`, `MemberController` | `UserService`, `TenantProvisioningService`, `InviteService` | `UserModel`, `TenantModel`, `MembershipModel` |

---

### 3.4 Tier 4: Enterprise, Intelligence & Operational Governance

| Service | Responsibility | Routes (`routes/`) | Controllers (`controllers/`) | Services (`services/`) | Models (`models/`) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **`admin-service`** | Platform administrator controls, global tenant oversight, feature flags | `admin.routes.ts`, `tenant.routes.ts`, `feature.routes.ts` | `AdminController`, `TenantAdminController` | `TenantManagementService`, `FeatureFlagService`, `SystemMetricsService` | `AdminAuditModel`, `FeatureFlagModel` |
| **`analytics-service`** | Real-time call dashboards, sentiment analytics, latency metrics | `analytics.routes.ts`, `dashboard.routes.ts` | `AnalyticsController`, `DashboardController` | `CallMetricsService`, `SentimentAnalyticsService`, `AggregationService` | `AggregatedMetricModel`, `DashboardSummaryModel` |
| **`audit-service`** | Immutable compliance audit logs, SOC2/HIPAA access records | `audit.routes.ts`, `compliance.routes.ts` | `AuditController`, `ComplianceController` | `AuditLogService`, `TamperProofService`, `ExportService` | `AuditEntryModel`, `ComplianceReportModel` |
| **`enterprise-service`** | Enterprise SSO (SAML, OIDC), SCIM user provisioning, custom domain mapping | `sso.routes.ts`, `scim.routes.ts`, `domain.routes.ts` | `SSOController`, `SCIMController`, `DomainController` | `SAMLService`, `OIDCService`, `SCIMProvisioningService` | `EnterpriseSSOConfigModel`, `DomainMappingModel` |
| **`scheduler-service`** | Scheduled outbound dialer campaigns, automated follow-up reminders | `schedule.routes.ts`, `campaign.routes.ts` | `ScheduleController`, `CampaignController` | `CronSchedulerService`, `OutboundCampaignService` | `ScheduledJobModel`, `CampaignModel` |
| **`search-service`** | Full-text and vector search across call transcripts, contacts, and logs | `search.routes.ts` | `SearchController` | `OpenSearchService`, `HybridSearchService` | `SearchIndexModel`, `SearchFilterModel` |
| **`spam-service`** | Robocall, fraud detection, Truecaller verification, STIR/SHAKEN validation | `spam.routes.ts`, `verify.routes.ts` | `SpamController`, `VerifyController` | `SpamDetectionService`, `TruecallerVerificationService` | `SpamRuleModel`, `CallerReputationModel` |
| **`websocket-service`** | Real-time WebSocket connection hub for live call monitoring & web dialer | `ws.routes.ts` | `WebSocketController` | `ConnectionManagerService`, `LiveStreamBroadcaster` | `WebSocketSessionModel` |
| **`workflow-service`** | Durable business workflows, LangGraph agent steps, Temporal orchestrations | `workflow.routes.ts`, `execution.routes.ts` | `WorkflowController`, `ExecutionController` | `TemporalWorkflowService`, `StepExecutionService` | `WorkflowDefinitionModel`, `ExecutionStateModel` |
