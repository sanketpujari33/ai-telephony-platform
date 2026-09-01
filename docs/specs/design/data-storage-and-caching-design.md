# 💾 Data Storage & Caching Architecture Design Specification

---

## 1. Multi-Model Data Storage Infrastructure

Nexora employs a polyglot persistence strategy optimized for relational integrity, vector similarity search, low-latency caching, and audio binary storage:

```mermaid
graph TD
    Services["⚙️ Microservices & AI Engines"]

    subgraph Relational["🐘 PostgreSQL 16 (:5432)"]
        Users["Users & Tenants"]
        Roles["12 RBAC Roles & Scopes"]
        Billing["Subscriptions & Invoices"]
        pgvector["Vector Embeddings (1536-dim)"]
    end

    subgraph MemoryStore["🍃 MongoDB 6.0 (:27017)"]
        Facts["Customer Memory Facts"]
        Transcripts["Full Raw Transcripts"]
        Workflows["Temporal Workflow States"]
    end

    subgraph CacheLayer["🔴 Redis 7.2 (:6379)"]
        LRU["L2 RBAC Permissions Cache"]
        RateLimit["Sliding Window Rate Limiters"]
        Locks["Redlock Distributed Mutexes"]
    end

    subgraph SearchStore["🔍 OpenSearch 2.13 (:9200)"]
        KnowledgeIdx["Knowledge Base Chunks"]
        CallLogs["Searchable Utterances"]
    end

    subgraph ObjectStore["🪣 MinIO / S3 (:9000)"]
        AudioWAV["WAV/MP3 Call Audio"]
        Waveforms["Waveform JSON Peaks"]
    end

    Services --> Relational
    Services --> MemoryStore
    Services --> CacheLayer
    Services --> SearchStore
    Services --> ObjectStore
```

---

## 2. Multi-Tenant Isolation Pattern

All relational tables enforce strict multi-tenant separation using mandatory `tenant_id` foreign keys indexed with composite B-Trees:

```sql
CREATE TABLE call_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    agent_id UUID NOT NULL REFERENCES agents(id),
    caller_phone VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_call_sessions_tenant ON call_sessions (tenant_id, created_at DESC);
```
