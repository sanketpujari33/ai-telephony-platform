# 📊 Diagram Inventory & Source Catalog

This document provides the complete catalog and inventory of all visual architecture blueprints, Mermaid models, sequence diagrams, and Entity-Relationship Diagrams (ERDs) across the platform.

---

## 1. Topologies & Sequences Master Catalog

| Diagram Title | Document Path | Type | Key Subsystems Represented |
| :--- | :--- | :---: | :--- |
| **Master System Architecture** | [`./system-architecture.md`](./system-architecture.md) | Flowchart / Topology | Client Apps, Ingress Gateways, 24 Microservices, LiveKit Media SFU, Storage Engines. |
| **Realtime Call Flows** | [`./call-flow-sequence.md`](./call-flow-sequence.md) | Sequence Flow | Outbound call dispatch, PSTN carrier webhook, LiveKit audio stream, VAD barge-in. |
| **Multi-Agent Coordination** | [`./multi-agent-coordination.md`](./multi-agent-coordination.md) | Graph / Architecture | 20 specialized Python agents, Intent routing, Tool Registry, MCP bridge. |
| **Database Schemas ERD** | [`./db/database-schemas-erd.md`](./db/database-schemas-erd.md) | Entity-Relationship | PostgreSQL 16 relational ERD, MongoDB memory schemas, Redis key-spaces, OpenSearch mappings. |
| **Microservices Architecture** | [`./services/microservices-architecture.md`](./services/microservices-architecture.md) | Architecture / Flow | Standard MVC directory mapping, NATS JetStream event bus taxonomy, route aggregation. |
| **Admin Operations Workflows**| [`./services/admin-operations-flow.md`](./services/admin-operations-flow.md) | Workflow Flowchart | 9 Operations: Banners, Promoters, Receive Funds, P&L, Dead Reports, Keywords, Feedback, General Tabs. |
| **User & Operator Softphone** | [`./services/user-operations-flow.md`](./services/user-operations-flow.md) | Sequence Flow | Softphone dialer, live in-call agent co-pilot, transcript streaming, and wrap-up notes. |
