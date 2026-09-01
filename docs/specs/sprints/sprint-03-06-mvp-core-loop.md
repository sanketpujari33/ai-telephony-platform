# 🏃 Sprints 3–6 — MVP Core Loop Plan
*Weeks 5–12 · Phase 1*

---

## Sprint 3: Telephony, Voice & Realtime
- **Objective**: Establish low-latency WebRTC and SIP voice infrastructure.
- **Tasks**:
  - [ ] Implement `services/telephony-service` (Go 1.22): SIP trunking, LiveKit room orchestration.
  - [ ] Implement `services/voice-service` (Python 3.11): Silero VAD, Deepgram STT, Cartesia TTS.
  - [ ] Implement `services/call-service`: Call session state machine transitions.

---

## Sprint 4: Multi-Agent Orchestrator & Tool Calling
- **Objective**: Build intent routing engine and tool execution registry.
- **Tasks**:
  - [ ] Implement `services/ai-orchestrator`: Intent classification and agent handover.
  - [ ] Build `agents/orchestrator-agent`: Master coordinator agent.
  - [ ] Build `@nexora/tools`: Zod to JSON Schema converter, MCP adapter.
  - [ ] Build core tools: `transfer_call`, `hangup_call`, `send_dtmf`, `book_appointment`.

---

## Sprint 5: Memory, RAG & Specialized Core Agents
- **Objective**: Build long-term memory extraction, knowledge retrieval, and calendar booking.
- **Tasks**:
  - [ ] Implement `services/memory-service`: pgvector embeddings & MongoDB fact store.
  - [ ] Build `agents/support-agent`: Knowledge base retrieval and CRM ticket resolution.
  - [ ] Build `agents/calendar-agent`: Google / Outlook calendar integration.
  - [ ] Build `services/transcript-service`: Live WebSocket interim transcripts.

---

## Sprint 6: End-to-End Voice Integration & Latency Tuning
- **Objective**: Full end-to-end integration and latency optimization.
- **Tasks**:
  - [ ] Connect PSTN SIP trunk to LiveKit audio room and Python voice agent.
  - [ ] Validate sub-20ms barge-in voice interruption.
  - [ ] Tune Time-to-First-Audio-Byte (TTFAB) to under 800ms.
