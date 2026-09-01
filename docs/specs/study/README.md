# 📚 Architecture & System Deep Studies

Comprehensive deep architectural studies, technical breakdowns, and implementation blueprints for the Nexora platform.

---

## 📑 Studies Index

| Document | Focus Area | Key Architectural Content |
| :--- | :--- | :--- |
| [**01. Core Packages & High-Level Architecture**](./01-core-packages-and-architecture.md) | `@nexora/*` Packages & Polyglot Runtime | 12 foundation packages, HashiCorp Vault KV v2 secret maps, 13 LLM providers, multi-channel router, UI tokens. |
| [**02. Packages & Services Test Architecture**](./02-packages-test-architecture.md) | Test Strategy & Verification Matrix | Centralized test suites in `tests/packages/` and `tests/services/`, 100% pass verification. |
| [**03. Infrastructure & Agent Ecosystem**](./03-infrastructure-and-agent-ecosystem.md) | Docker Topology & 20 AI Agents | 10 Docker infrastructure containers, ports, database adapters, and 20 specialized Python agents. |
| [**04. MVC Microservices Blueprint**](./04-mvc-microservices-blueprint.md) | Microservices Tier Mapping & MVC Structure | Standard MVC layers (`routes`, `controllers`, `services`, `models`), 4 service tiers, and 24 microservices. |
| [**05. Admin RBAC Model Study**](./05-admin-rbac-model-study.md) | 12-Role RBAC & Two-Tier Cache | 12 system roles, granular permission scopes, LRU/Redis caching, and real-time NATS invalidation. |
| [**06. Monorepo Architecture, Infrastructure & Security Alignment**](./06-monorepo-architecture-infrastructure-security-master-alignment.md) | Master System & Security Blueprint | Complete localhost loopback isolation, Redis requirepass auth, 3-broker KRaft Kafka cluster, 58 Vault paths, 26 SQL migrations, and 21 autonomous AI agents. |
| [**07. Polyglot Packages & Infrastructure Deep Dive**](./07-polyglot-packages-and-infrastructure-deep-dive.md) | Polyglot Packages & Infrastructure | Comprehensive deep dive into all 14 foundation packages across TypeScript, Go, and Python. |
| [**08. 24-Microservice Master Architecture & Roadmap**](./08-microservices-master-architecture-and-roadmap.md) | 24 Microservices Mesh & Implementation Roadmap | Complete 8-phase execution roadmap, port matrix, clean architecture blueprint, and package bindings. |

