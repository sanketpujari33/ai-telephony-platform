# 🎙️ Inbound AI Call Screening & Outbound Mission Engine

## 🌐 Overview & Core Mission Philosophy

Nexora AI unites the two fundamental vectors of personal and business telephony into a unified autonomous engine:
1. **Inbound Call Screening ("Green Card" System):** Real-time interception, speech transcription, intent classification, whisper guidance, and zero-latency live call takeover (pioneered by **Equal AI**).
2. **Outbound Autonomous Missions:** Multi-step voice goal execution, DTMF IVR navigation, hold music spectral detection, and human negotiation (pioneered by **hiRobin**).

---

## 🛡️ Part 1: Inbound Call Screening Engine

```mermaid
graph TD
    IncomingCall["☎️ Forwarded Call from Carrier (*004*)"] --> Ingress["⚡ LiveKit WebRTC SFU (:7880)"]
    Ingress --> StreamAudio["🎙️ Sub-80ms Inbound RTP Audio Stream"]
    StreamAudio --> DeepgramSTT["📝 Streaming Speech-to-Text (STT)"]
    DeepgramSTT --> IntentClassifier["🧠 Intent & Category Classifier"]
    
    subgraph Categories["Live Intent Categories"]
        C1["📦 Delivery Agent (Swiggy, Zomato, Amazon)"]
        C2["🚪 Gate Security (MyGate, Society Guard)"]
        C3["💼 Recruiter / Business Lead"]
        C4["👨‍👩‍👧 Family / VIP Whitelist"]
        C5["🛡️ Spam / Telemarketing / Loan Bot"]
    end
    
    IntentClassifier --> Categories
    Categories --> PushApp["📱 Mobile App: Real-Time Live Waveform + Live Transcript"]
    PushApp --> UserActions["⚡ User Interaction Interface"]
    
    subgraph Actions["User Live Options"]
        A1["🤫 Whisper Guidance Chip ('Leave at door', 'Call back in 30m')"]
        A2["⚡ 1-Tap Live Take-Over (Barge-In)"]
        A3["🤖 Let AI Handle Completely"]
        A4["🔴 End Call & Block Number"]
    end
    
    UserActions --> Actions
```

### 1. In-Call Whisper Guidance (Prompt Injection)
While the AI assistant is conversing with the caller, the user can guide the conversation in real time without taking over the voice channel:
- **Preset Quick-Chips:**
  - *"Leave with guard at tower 4"*
  - *"Tell them I am driving and will call back in 1 hour"*
  - *"Ask for their company name and callback number"*
  - *"Confirm delivery OTP: 8492"*
- **Custom Whisper Input:** The user can type any custom sentence, which is immediately injected into the Voice AI's prompt context with maximum priority.

### 2. 1-Tap Live Call Take-Over (Barge-In)
- When the user taps **"Take Over Call"**, the mobile client establishes an immediate bidirectional WebRTC audio track to the LiveKit room.
- The Voice AI agent delivers a graceful handover chime (`takeover_trigger.mp3`), smoothly mutes the synthetic speech output, and bridges the user's phone microphone and earpiece to the live PSTN caller in **< 15 milliseconds**.

---

## 🚀 Part 2: Outbound Autonomous Mission Engine

```mermaid
sequenceDiagram
    autonumber
    actor User as 👤 Nexora Mobile User
    participant App as 📱 Mobile App UI
    participant Orch as 🧠 AI Orchestrator (:8086)
    participant Telephony as ☎️ Go Telephony Engine (:8085)
    participant LiveKit as 🎙️ LiveKit SFU (:7880)
    participant VoiceAgent as 🤖 Outbound Voice AI Agent
    actor Target as 🏢 Clinic / Customer Support IVR

    User->>App: Submits Mission: "Book table for 4 at Olive Bistro tonight at 8 PM"
    App->>Orch: POST /api/v1/missions/start (Target DID, Goal, Constraints)
    Orch->>Telephony: Initiate Outbound SIP Call to Target (+91 80 4123 4567)
    Telephony->>Target: Outbound Call Ringing
    Target->>Telephony: Call Answered / IVR Voice Active
    Telephony->>LiveKit: Ingress Audio Stream
    LiveKit->>VoiceAgent: Start Audio Transcription & Intent Parsing

    opt IVR Detected
        VoiceAgent->>VoiceAgent: Parses: "Press 1 for Table Reservations"
        VoiceAgent->>Telephony: Send Out-of-band DTMF Tone "1" (RFC 4733)
        Telephony->>Target: DTMF '1' Transmitted
    end

    opt Placed on Hold
        VoiceAgent->>VoiceAgent: FFT Spectrum detects Hold Music Loop (No human speech)
        VoiceAgent->>App: Push Status: "Waiting on hold music... (Elapsed: 01:45)"
        Target->>VoiceAgent: Human Receptionist joins: "Olive Bistro, how can I help?"
        VoiceAgent->>VoiceAgent: Voice Onset Detected (Latency < 90ms)
    end

    VoiceAgent->>Target: "Hi, I am calling on behalf of Sanket to book a table for 4 tonight at 8 PM."
    Target->>VoiceAgent: "We have an outdoor table at 8:15 PM, does that work?"
    VoiceAgent->>Target: "Yes, 8:15 PM outdoor works perfectly. Please reserve under Sanket."
    Target->>VoiceAgent: "Reserved! Table 14 confirmed."
    VoiceAgent->>Target: "Thank you! Have a great evening."

    VoiceAgent->>Orch: Mission Complete (Confirmation: Table 14, 8:15 PM)
    Orch->>App: Push Notification & In-App Mission Confirmation Card
```

---

## 🧠 Part 3: Shared Memory & Context Rules

Every screening interaction and outbound mission reads and updates the centralized **Memory Service**:
- **Delivery Rules:** *"Always ask delivery agents to leave parcels with security guard Ram Singh at Gate 2."*
- **Gate Pass Rules:** *"Generate MyGate pre-approval token automatically for Amazon delivery."*
- **Meeting Auto-Responses:** *"When in calendar meetings, state I am presenting and will return urgent calls after 5:00 PM."*
- **VIP Contact Whitelist:** Saved contacts (family, team leads, clients) bypass screening or receive personalized greetings.
