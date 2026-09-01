# 🚀 hiRobin & myEqual Unified AI Telephony Architecture Specification

## 🌟 Executive Overview & Architectural Synthesis

The **Nexora AI Telephony Platform** consolidates and supersedes the state-of-the-art capabilities pioneered by **hiRobin (`hirobin.ai` by TaskGenie Private Limited)** and **Equal AI (`myequal.ai` by Equal Identity Private Limited)**, establishing a unified, enterprise-grade AI Telephony Operating System.

```mermaid
graph TB
    subgraph Inbound["📥 Inbound Call Intelligence (Equal AI + hiRobin)"]
        UnknownShield["🛡️ Unknown Numbers Shield<br/>(Contacts pass, strangers screened)"]
        MMIForward["📲 Dual-SIM Carrier Forwarding<br/>(*004* MMI Trunk Routing)"]
        LiveTrans["📝 Real-Time LiveKit Audio & STT Stream"]
        LiveTakeover["⚡ 1-Tap Seamless Live Take-Over"]
        WhisperGuide["🤫 In-Call Whisper Prompt Injection"]
        SmartMemory["🧠 Contextual Persona Rules<br/>(Deliveries, Gate Passes, Meetings)"]
    end

    subgraph CoreEngine["⚡ Nexora Telephony Core & AI Orchestrator"]
        LiveKitSFU["🎙️ LiveKit WebRTC SFU (:7880)"]
        GoTelephony["📞 Go SIP Trunk Engine (:8085 / :5060)"]
        PyVoice["🐍 Multilingual Voice AI (9+ Indian Dialects + English)"]
        FailureIndex["📊 Carrier Failure Index Monitor (0.0% Drop Target)"]
    end

    subgraph Outbound["📤 Outbound AI Missions (hiRobin Autonomous Execution)"]
        DeliveryCoord["📦 Delivery & Gate Pass Coordination"]
        HoldNav["⏳ Hold Music & IVR DTMF Navigation"]
        ApptBooking["🗓️ Doctor & Restaurant Appointments"]
        SupportRefund["💳 Customer Support & Refund Follow-Ups"]
    end

    Inbound --> CoreEngine
    CoreEngine --> Outbound
```

---

## 🏛️ 1. Comprehensive Feature Comparison Matrix

| Feature Capability | hiRobin (`hirobin.ai`) | Equal AI (`myequal.ai`) | Nexora AI Telephony (Unified) | Implementation Module |
| :--- | :---: | :---: | :---: | :--- |
| **Incoming Call Screening** | ✅ Basic Answering | ✅ Full Real-time Screening | 🌟 **Enterprise Shield** | [`lib/features/carrier/`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/carrier) |
| **Unknown Callers Only Mode** | ❌ All calls forwarded | ✅ Address book passthrough | 🌟 **Granular Address Book Filter** | [`lib/features/carrier/presentation/providers/carrier_provider.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/carrier/presentation/providers/carrier_provider.dart) |
| **Dual-SIM Multi-Carrier Forwarding** | ❌ Single SIM only | ❌ Single SIM only | 🌟 **Dual-SIM 1-Click Master Control** | [`lib/features/carrier/presentation/screens/call_forwarding_wizard_screen.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/carrier/presentation/screens/call_forwarding_wizard_screen.dart) |
| **Live Audio & STT Streaming** | ❌ Post-call summary only | ✅ Live Transcript Card | 🌟 **Sub-80ms Live Audio + STT** | [`lib/features/calls/presentation/screens/live_call_screen.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/calls/presentation/screens/live_call_screen.dart) |
| **1-Tap Live Call Take-Over (Barge-In)**| ❌ | ✅ Supported | 🌟 **Zero-Disconnect Audio Swap** | [`lib/features/calls/presentation/widgets/live_waveform_visualizer.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/calls/presentation/widgets/live_waveform_visualizer.dart) |
| **In-Call Whisper Guidance** | ❌ | ⚠️ Limited Chips | 🌟 **Prompt Injection Chips + Custom Text** | [`lib/features/calls/presentation/widgets/in_call_whisper_sheet.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/calls/presentation/widgets/in_call_whisper_sheet.dart) |
| **Autonomous Outbound Calling Missions**| ✅ Full Missions | ❌ Inbound Only | 🌟 **Multi-Step Goal Execution Engine** | [`lib/features/missions/presentation/screens/outbound_task_screen.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/missions/presentation/screens/outbound_task_screen.dart) |
| **Hold Music & IVR DTMF Navigation** | ✅ Supported | ❌ | 🌟 **Autonomous DTMF Dialing & Hold Detection** | [`services/call-service/`](file:///d:/trideocoders/ai-telephony-platform/services/call-service) |
| **Carrier Failure Rate Diagnostic** | ✅ Weekly Web Index | ❌ | 🌟 **Real-Time In-App Handshake Diagnostic** | [`lib/features/carrier/presentation/screens/carrier_activation_journey_screen.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/carrier/presentation/screens/carrier_activation_journey_screen.dart) |
| **Multilingual Real-Time Switching** | ⚠️ English + Hindi | ✅ 9+ Indian Languages | 🌟 **9+ Dialects + 3s Auto-Detect** | [`lib/features/persona/presentation/screens/ai_persona_screen.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/persona/presentation/screens/ai_persona_screen.dart) |
| **Custom Delivery & Gate Pass Memory** | ⚠️ Fixed Prompts | ⚠️ Fixed Prompts | 🌟 **Dynamic Memory Engine & Prompts** | [`lib/features/persona/presentation/providers/persona_provider.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/persona/presentation/providers/persona_provider.dart) |
| **TCPA & TRAI DND Compliance** | ⚠️ Basic | ⚠️ Basic | 🌟 **Enterprise Zero-Retention & Consent Mode** | [`lib/features/help/presentation/screens/privacy_policy_screen.dart`](file:///d:/trideocoders/ai-telephony-platform/apps/mobile-app/lib/features/help/presentation/screens/privacy_policy_screen.dart) |

---

## 🔄 2. Inbound Screening, Whisper Guidance, & Live Take-Over Sequence

```mermaid
sequenceDiagram
    autonumber
    actor Caller as 👤 Unknown Caller (e.g. Swiggy Delivery / Recruiter)
    participant Carrier as 📡 Carrier PSTN (Jio / Airtel / Vi / AT&T)
    participant NexoraTrunk as ☎️ Nexora SIP Trunk (:5060)
    participant LiveKit as 🎙️ LiveKit WebRTC SFU (:7880)
    participant VoiceAI as 🤖 Nexora Voice AI (:8086)
    actor UserApp as 📱 User Mobile App (Flutter)

    Caller->>Carrier: Dials User Mobile Number (+91 98201 88412)
    Note over Carrier: Contact not in Address Book or User Busy / 3 rings timeout
    Carrier->>NexoraTrunk: Conditional Forward (*004*+18008490192#)
    NexoraTrunk->>LiveKit: Establish Low-Latency WebRTC Ingress Room
    NexoraTrunk->>VoiceAI: Trigger Voice Agent with Persona Context
    VoiceAI->>LiveKit: Stream AI greeting audio (<80ms acoustic latency)
    LiveKit->>Caller: Plays natural AI greeting ("Hi, I am Nexora answering for Sanket...")
    VoiceAI->>UserApp: Push Notification & Real-Time Waveform + Transcript Stream
    UserApp->>UserApp: User opens Live Call Screen and sees live speech text

    alt User Injects Whisper Guidance
        UserApp->>VoiceAI: Tap "Leave at Security Gate" chip / Send Whisper
        VoiceAI->>Caller: AI naturally speaks: "Please leave the parcel with the gate security guard."
    else User Decides to Take Over (Barge-In)
        UserApp->>LiveKit: Tap "Take Over Live Call"
        LiveKit->>VoiceAI: Mute AI Synthesizer & Bridge User Audio Track
        LiveKit->>Caller: Seamlessly connects User Mic with 0ms interruption
        UserApp->>Caller: User continues talking with caller directly
    else AI Completes Interaction Autonomously
        VoiceAI->>Caller: Fulfills intent & graceful teardown
        VoiceAI->>UserApp: Generates Post-Call Action Card & Structured Summary
    end
```

---

## 🎯 3. Outbound Mission Execution Lifecycle (hiRobin Parity)

Nexora AI automates complex outbound calling tasks end-to-end:

```mermaid
graph TD
    UserGoal["🎯 User Submits Task: 'Reschedule Dr. Sharma appointment to Friday 5 PM'"] --> TaskPlan["🧠 Orchestrator Deconstructs Mission Steps & Extracts Contact DID"]
    TaskPlan --> DialOut["☎️ SIP Trunk Dials Hospital / Clinic Line"]
    DialOut --> IVRCheck{"Is Multi-Level IVR Detected?"}
    
    IVRCheck -- "Yes" --> ParseMenu["🔊 Deepgram STT Parses Voice Menu Options"]
    ParseMenu --> SendDTMF["🔢 Telephony Service Sends RFC 4733 DTMF ('Press 2 for Appointments')"]
    SendDTMF --> HoldDetect{"Placed on Hold Music?"}
    
    IVRCheck -- "No" --> HoldDetect
    
    HoldDetect -- "Yes (Hold Music Playing)" --> MuteAgent["⏳ Mute AI, Monitor FFT Audio Spectrum for Human Voice Onset"]
    MuteAgent --> HumanAnswers["🗣️ Human Receptionist Answers: 'Hello, Apollo Clinic'"]
    
    HoldDetect -- "No" --> HumanAnswers
    
    HumanAnswers --> ConversationalAI["🤖 Voice AI Negotiates Appointment Time using User Calendar Constraints"]
    ConversationalAI --> ConfirmSlot["✅ Slot Confirmed: Friday 5:00 PM (Token: AP-8821)"]
    ConfirmSlot --> SyncCalendar["📅 Calendar Agent Updates Google/Outlook Calendar"]
    SyncCalendar --> SendAlert["📱 Mobile Push Notification: 'Appointment Confirmed + Audio Snippet'"]
```

---

## 🌐 4. Multilingual Dialect & Audio Intelligence Engine

Nexora AI incorporates deep multilingual synthesis tailored for Indian and global telecom networks:

- **Supported Languages & Dialects:**
  - 🇮🇳 **Indian Regional:** Hindi, Hinglish, Marathi, Gujarati, Tamil, Telugu, Kannada, Bengali, Punjabi.
  - 🌍 **International:** English (US, UK, Indian accents), Spanish, French, German.
- **Dynamic 3-Second Language Auto-Detect:**
  - Fast Fourier Spectral Classifier analyzes the first 3 seconds of caller speech.
  - Switches the LLM system prompt and text-to-speech voice model to match caller dialect on-the-fly without session restart.

---

## 🛡️ 5. Privacy, TCPA, and TRAI Compliance Architecture

- **Zero-Retention Mode:** Purges voice audio recordings immediately upon transcript generation and PII anonymization.
- **Automated PII Entity Redaction:** Credit card numbers, OTPs, Aadhaar numbers, SSNs, and passwords are automatically masked (`[REDACTED_OTP]`, `[REDACTED_CARD]`).
- **Mandatory Disclosure & Consent:** Automatic compliance announcements when forwarding calls (*"This call is assisted by Nexora AI"*), conforming to TCPA (US) and TRAI TCCCPR (India).
