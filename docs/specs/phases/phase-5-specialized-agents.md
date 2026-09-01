# 🤖 Phase 5 — Specialized AI Agents Specification
*Sprints 16–20 · Weeks 31–40*

---

## 1. Phase Goals & Scope

Phase 5 scales the AI Platform from the initial 5-agent MVP to the full **20-agent ecosystem**, empowering verticalized workflows across finance, research, travel, e-commerce, delivery, and web automation.

---

## 2. Sprint Deliverables

### Sprint 16: Finance, Delivery & Shopping Agents
- **`agents/finance-agent`**: Invoices, balance inquiries, payment links, refund requests.
- **`agents/delivery-agent`**: Real-time package tracking, delivery ETA recalculations, address corrections.
- **`agents/shopping-agent`**: E-commerce catalog search, cart operations, order status lookups.

### Sprint 17: Travel, Research & Reminder Agents
- **`agents/travel-agent`**: Flight status lookups, hotel booking confirmations, itinerary updates.
- **`agents/research-agent`**: Perplexity / Grok live web search, multi-source factual synthesis, source citations.
- **`agents/reminder-agent`**: Temporal durable scheduled alarms, follow-up callbacks, customer re-engagement.

### Sprint 18: Browser, Coding & Automation Agents
- **`agents/browser-agent`**: Headless Playwright automation, web form submission, portal data extraction.
- **`agents/coding-agent`**: Code generation, API snippet synthesis, webhook payload debugging.
- **`agents/automation-agent`**: Multi-step business workflow execution triggered by NATS events or incoming webhooks.

### Sprint 19: Enterprise & Knowledge Agents
- **`agents/enterprise-agent`**: SLA compliance checks, policy enforcement, enterprise routing rules.
- **`agents/knowledge-agent`**: OpenSearch 2.13 hybrid search + pgvector semantic chunk retrieval with source attribution.

### Sprint 20: 20-Agent Graph Routing & MCP Manager
- **Model Context Protocol (MCP)**: Universal MCP client adapter dynamically discovering local and remote MCP tool servers.
- **Master Routing Matrix**: Orchestrator sub-millisecond intent classification dynamically delegating between all 20 agents.

---

## 3. Exit Criteria

- All 20 agents execute independently in Python with fallback configurations.
- Master orchestrator dynamically routes complex multi-turn user requests across specialist agents.
- External MCP servers can be plugged in without recompiling agent code.
