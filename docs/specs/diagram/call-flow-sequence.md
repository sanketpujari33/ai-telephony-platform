# 📞 Real-Time Call Flow & Telephony Sequences

---

## 1. Outbound AI Call Dispatch Sequence

```mermaid
sequenceDiagram
    autonumber
    actor User as 👤 User / Dashboard
    participant GW as ⚡ API Gateway (:4000)
    participant Auth as 🔐 Auth Guard
    participant CallSvc as 📞 Call Service (:4002)
    participant LK as 🎙️ LiveKit SFU (:7880)
    participant SIP as 📡 SIP Trunk / Telnyx
    participant PythonAI as 🤖 Python AI Agent
    actor Recipient as 📱 Recipient Phone

    User->>GW: 1. POST /api/v1/calls/dispatch<br/>{agentId, to: "+15551234567"}
    GW->>Auth: 2. Validate Tenant API Key & Quota
    Auth-->>GW: 3. Token Authorized
    GW->>CallSvc: 4. Create Call Record (Status: 'initiated')
    CallSvc->>LK: 5. Create LiveKit Audio Room (room_call_xyz)
    LK-->>CallSvc: 6. Room Created + Agent Token
    CallSvc->>SIP: 7. Originate SIP Trunk Call to Recipient
    SIP->>Recipient: 8. PSTN Ringing Tone
    CallSvc->>PythonAI: 9. Launch Agent Worker in Room
    PythonAI->>LK: 10. Agent Joins Room as Audio Participant
    
    Recipient->>SIP: 11. Recipient Answers (Status: 'in-progress')
    SIP->>LK: 12. Bridge PSTN Audio to LiveKit Room
    PythonAI->>LK: 13. Speak Opening Greeting ("Hello, this is Nexora AI...")
    LK->>Recipient: 14. Real-time Audio Playback
```

---

## 2. Inbound SIP Call Routing & Intent Handoff Sequence

```mermaid
sequenceDiagram
    autonumber
    actor Caller as ☎️ Inbound Caller
    participant Carrier as 📡 Carrier / Twilio SIP
    participant SIPBridge as 📞 SIP Gateway (Go)
    participant LK as 🎙️ LiveKit SFU
    participant Orch as 🧠 Orchestrator Agent
    participant SupportAgent as 🎧 Support Agent
    participant Tools as 🔌 Tool Registry

    Caller->>Carrier: 1. Dials Nexora DID (+18005559000)
    Carrier->>SIPBridge: 2. SIP INVITE (TLS / SRTP)
    SIPBridge->>LK: 3. Create WebRTC Room & Connect RTP
    SIPBridge->>Orch: 4. Route Inbound Call Event
    Orch->>LK: 5. Connect as Initial Greeting Agent
    Orch->>Caller: 6. "Thank you for calling Nexora. How may I help you today?"
    Caller->>Orch: 7. "I have a question about my invoice billing."

    note over Orch: Intent Classification (Intent: 'billing_support')
    Orch->>SupportAgent: 8. Warm Handoff with Call Context & History
    SupportAgent->>Tools: 9. Execute 'invoice_check_status'
    Tools-->>SupportAgent: 10. Invoice #1049 is Pending ($120.00)
    SupportAgent->>Caller: 11. "I see your invoice #1049 for $120 is currently pending..."
```

---

## 3. Post-Call Async Lifecycle Sequence

```mermaid
sequenceDiagram
    autonumber
    actor Caller as ☎️ Caller
    participant LK as 🎙️ LiveKit SFU
    participant CallSvc as 📞 Call Service
    participant NATS as ⚡ NATS JetStream
    participant TransSvc as 📝 Transcript Service
    participant MemAgent as 🧠 Memory Agent
    participant BillSvc as 💳 Billing Service
    participant S3 as 🪣 MinIO / S3

    Caller->>LK: 1. Caller Hangs Up (Disconnect)
    LK->>CallSvc: 2. Room Closed Event (Duration: 184s)
    CallSvc->>S3: 3. Save Call Audio Recording (recording.wav)
    CallSvc->>CallSvc: 4. Mark Status: 'completed'
    CallSvc->>NATS: 5. Publish `call.completed` event
    
    par Async Processing
        NATS->>TransSvc: 6. Generate Full Formatted Transcript
        TransSvc->>TransSvc: 7. Index in OpenSearch
    and Memory Extraction
        NATS->>MemAgent: 8. Extract Conversation Facts & Sentiment
        MemAgent->>MemAgent: 9. Persist Facts in MongoDB
    and Usage Invoicing
        NATS->>BillSvc: 10. Calculate Minute Usage & Deduct Wallet
    end
```
