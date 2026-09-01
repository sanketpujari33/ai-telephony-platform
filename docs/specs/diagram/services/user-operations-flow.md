# 👤 User & Agent Operator Lifecycle Diagrams Flow

This document details the visual architecture, call dispatching loops, and data flows for **Tenant Admins**, **Agent Operators**, **Team Members**, and **End Users**.

---

## 1. User & Team Member Role Hierarchy in Workspace

```mermaid
graph TD
    subgraph WorkspaceBoundary["🏢 Tenant Workspace Isolation"]
        TenantAdmin["👔 TENANT_ADMIN<br/>(Full Workspace Control & Billing)"]
        
        TenantAdmin --> TeamLead["👥 TEAM_LEAD<br/>(Assigns Queues & Agent Seats)"]
        TeamLead --> AgentOperator["📞 AGENT_OPERATOR<br/>(Softphone Calling & Contact Notes)"]
        TeamLead --> Marketer["📈 MARKETING_MEMBER<br/>(Campaign Broadcasts & Lists)"]
        
        TenantAdmin --> Viewer["👁️ VIEWER<br/>(Analytics & Read-Only Metrics)"]
    end

    subgraph ExternalEndUser["📱 External Customer / End-User"]
        Customer["👤 End-User / Caller<br/>(PSTN / Mobile App / Web Widget)"]
    end

    Customer <-->|Inbound / Outbound Calls| AgentOperator
    Marketer -.->|Dispatches Campaigns| Customer
```

---

## 2. Agent Operator Softphone Inbound & Outbound Calling Flow

```mermaid
sequenceDiagram
    autonumber
    actor Operator as 📞 Agent Operator (Softphone)
    participant WebRTC as 🎙️ LiveKit SFU (:7880)
    participant Gateway as ⚡ API Gateway (:4000)
    participant CallSvc as 📞 Call Service (:4002)
    participant SIPBridge as 📡 Carrier SIP Trunk
    actor Caller as 👤 Customer Phone

    Note over Operator,Caller: Outbound Softphone Call Flow
    Operator->>Gateway: 1. POST /api/v1/user/dialer/call<br/>{to: "+15552345678", agentId: "agent_phone_default"}
    Gateway->>CallSvc: 2. Create Call Session & Allocate Room
    CallSvc->>WebRTC: 3. Join Operator to LiveKit Audio Room
    CallSvc->>SIPBridge: 4. Dial Customer via PSTN Trunk
    SIPBridge->>Caller: 5. Ringing Customer Handset
    Caller->>SIPBridge: 6. Customer Answers Call
    SIPBridge->>WebRTC: 7. Bridge Audio Stream (Full-Duplex Opus 48kHz)
    
    rect rgb(240, 248, 255)
        Note over Operator,Caller: Live In-Call Assistance & Assistance AI
        WebRTC-->>Operator: 8. Live Real-Time Transcript Stream
        Operator->>Gateway: 9. Add In-Call Note: "Customer interested in upgraded plan"
        Operator->>Gateway: 10. Warm Transfer to Support Supervisor
    end

    Caller->>SIPBridge: 11. Call Ended (Hangup)
    CallSvc-->>Gateway: 12. Persist Call Record & Audio Recording
```

---

## 3. Team Invitation & Dynamic RBAC Provisioning Flow

```mermaid
graph TD
    Admin["👔 Tenant Admin"] --> Invite["POST /api/v1/user/team/invite<br/>{email, role: 'agent_operator'}"]
    Invite --> EmailSvc["📧 Notification Service (SES Email)"]
    EmailSvc --> Member["👤 Invited Team Member"]
    
    Member --> Accept["Accept Invitation Link"]
    Accept --> SetPass["Set Password & Enable 2FA"]
    SetPass --> AssignSeat["Assign Dedicated Softphone Seat & DID"]
    
    subgraph ActiveWorkspace["⚡ Active Workspace Permissions"]
        AssignSeat --> Perms["Grant Granular Permissions<br/><code>[calls:dial, contacts:edit, recordings:listen]</code>"]
    end
```

---

## 4. End-User Communication Preferences & Self-Service Flow

```mermaid
graph LR
    subgraph EndUserPortal["🌐 Customer Self-Service Portal"]
        Prefs["Communication Preferences<br/>(Call vs. WhatsApp vs. SMS)"]
        TimeWindow["Preferred Time Window<br/>(09:00 AM - 12:00 PM EST)"]
        Consent["DND / Marketing Consent Toggle"]
    end

    subgraph BackendSync["💾 Backend Memory Sync"]
        MongoMemory[("🍃 MongoDB Memory Facts")]
        RedisDND[("🔴 Redis DND State")]
    end

    Prefs --> MongoMemory
    TimeWindow --> MongoMemory
    Consent --> RedisDND
```
