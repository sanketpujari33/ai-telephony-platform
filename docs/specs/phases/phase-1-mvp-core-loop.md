# 🎯 Phase 1 — MVP Core Loop Specification
*Sprints 3–6 · Weeks 5–12*

---

## 1. Phase Goals & Scope

Phase 1 delivers the core real-time conversational voice loop — the primary value proposition of Nexora. It focuses exclusively on the **5-Agent MVP**:
1. **`phone-agent`**: Real-time voice call handling over WebRTC and SIP.
2. **`orchestrator-agent`**: Intent routing and specialist agent delegation.
3. **`calendar-agent`**: Scheduling, availability checks, and meeting bookings.
4. **`support-agent`**: CRM customer support, ticket resolution, and policy QA.
5. **`memory-agent`**: Semantic conversation facts extraction and storage in MongoDB & pgvector.

---

## 2. Sprint Deliverables

### Sprint 3: Telephony, Voice & Realtime Pipeline
- **`services/telephony-service`**: Go 1.22 SIP protocol bridge, LiveKit WebRTC audio room allocation, carrier trunking.
- **`services/voice-service`**: Python streaming audio pipeline:
  - Silero VAD (Voice Activity Detection) with barge-in interruption (<20ms).
  - Deepgram Nova-2 streaming STT with interim results.
  - Cartesia Sonic / ElevenLabs Flash streaming TTS with chunked audio playback.
- **`services/call-service`**: State machine transitions (`initiated` → `ringing` → `in-progress` → `completed`).

### Sprint 4: Multi-Agent Orchestrator & Tool Calling
- **`services/ai-orchestrator`**: High-performance intent router with LangGraph/Temporal coordination.
- **`agents/orchestrator-agent`**: Classifies user utterance and delegates turns.
- **`packages/tools`**: Dynamic Tool Registry conforming to OpenAI / Anthropic function calling and MCP (Model Context Protocol).
- **Core Tools**: `hangup_call`, `transfer_call`, `send_dtmf`, `book_appointment`, `lookup_knowledge`.

### Sprint 5: Memory, RAG & Specialized Core Agents
- **`services/memory-service`**: Conversation fact extraction into MongoDB and pgvector similarity embeddings.
- **`agents/support-agent`**: Knowledge base retrieval and CRM ticket resolution.
- **`agents/calendar-agent`**: Google Calendar / Outlook integration via CalDAV/OAuth.
- **`services/transcript-service`**: Live interim transcripts broadcast over WebSocket and post-call speaker diarization.

### Sprint 6: End-to-End Voice Integration & Latency Tuning
- End-to-end phone call test: Inbound PSTN call → SIP Ingress → LiveKit SFU → Python Phone Agent → Groq/Claude LLM → Cartesia TTS → Audio Egress.
- Optimization goal: **Sub-800ms Time-to-First-Audio-Byte (TTFAB)**.

---

## 3. Exit Criteria & Definition of Done

- Inbound and outbound phone calls successfully connect with natural bi-directional voice.
- Barge-in voice interruption immediately halts TTS playback and cancels active LLM generation.
- Post-call transcripts, summary, sentiment, and facts are persisted within 5 seconds of call completion.
