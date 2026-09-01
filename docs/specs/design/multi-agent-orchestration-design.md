# 🤖 Multi-Agent Orchestration & Tool Registry Design Specification

---

## 1. Multi-Agent Routing Topology

Nexora coordinates 20 specialized Python AI agents through a master intent router:

```mermaid
graph TD
    User["🗣️ User Utterance (Voice / Text)"]
    Orch["🧠 Master Orchestrator Agent<br/>(Intent Classifier)"]

    User --> Orch

    subgraph VoiceCore["📞 Realtime Voice & Support"]
        Phone["📞 phone-agent"]
        Support["🎧 support-agent"]
        Calendar["📅 calendar-agent"]
        Appoint["⏰ appointment-agent"]
    end

    subgraph Messaging["💬 Omnichannel Messaging"]
        WhatsApp["💬 whatsapp-agent"]
        SMS["📱 sms-agent"]
        Email["✉️ email-agent"]
        Reminder["🔔 reminder-agent"]
    end

    subgraph Business["💼 Enterprise & Operations"]
        Finance["💰 finance-agent"]
        Shopping["🛒 shopping-agent"]
        Delivery["📦 delivery-agent"]
        Travel["✈️ travel-agent"]
        Enterprise["🏢 enterprise-agent"]
    end

    subgraph Intelligence["🔍 Research & Automation"]
        Research["🌐 research-agent"]
        Knowledge["📚 knowledge-agent"]
        Memory["🧠 memory-agent"]
        Browser["🌐 browser-agent"]
        Coding["💻 coding-agent"]
        Automation["⚡ automation-agent"]
    end

    Orch --> VoiceCore
    Orch --> Messaging
    Orch --> Business
    Orch --> Intelligence

    subgraph MCPBridge["🔌 Dynamic Tool Registry & MCP Bridge"]
        Tools["🛠️ Built-in Tools (Transfer, Hangup, DTMF, Booking)"]
        RemoteMCP["🔌 Remote MCP Servers (JSON-RPC)"]
    end

    VoiceCore --> MCPBridge
    Messaging --> MCPBridge
    Business --> MCPBridge
    Intelligence --> MCPBridge
```

---

## 2. Dynamic Tool Registry & MCP Integration

- **Standard Schemas**: Tools define parameters using Zod schemas compiled to JSON Schema Draft 7 via `@nexora/tools`.
- **Model Context Protocol (MCP)**: Agents dynamically connect to external MCP servers over stdio or SSE/JSON-RPC without local code modifications.
