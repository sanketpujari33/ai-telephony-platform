# ADR 0002: STT, LLM, and TTS Provider Selection

## Status
Accepted

## Context
Conversational AI over telephony requires absolute minimum latency. Any delay over 800ms breaks conversational flow. We need to select the primary and fallback providers for Speech-to-Text (STT), Large Language Models (LLM), and Text-to-Speech (TTS).

## Decision
We implemented a **Circuit-Breaker Router Pattern** with the following provider priority:

### STT (Speech-to-Text)
1. **Deepgram** (Primary): Selected for its sub-300ms streaming latency and Nova-2 model accuracy on 8kHz telephony audio.
2. **AssemblyAI** (Fallback): High accuracy but slightly higher latency, used if Deepgram API experiences 5xx errors.

### LLM (Large Language Model)
1. **OpenAI GPT-4o-mini** (Primary): Unbeatable time-to-first-token (TTFT) and instruction following for dynamic conversations.
2. **Anthropic Claude 3 Haiku** (Fallback): Extremely fast fallback option if OpenAI rate limits are hit.

### TTS (Text-to-Speech)
1. **ElevenLabs** (Primary): Selected for the highest MOS (Mean Opinion Score) for human-like emotional intonation and low-latency websocket streaming (turbo v2.5 model).
2. **Azure Cognitive Services TTS** (Fallback): Highly reliable enterprise SLA, used if ElevenLabs streaming drops.

## Rationale
Relying on a single vendor for AI primitives poses a critical uptime risk. By writing interface abstractions (`STTRouter`, `LLMRouter`, `TTSRouter`), we can dynamically failover mid-call if a websocket drops or a 429 Rate Limit is encountered.

## Consequences
- Requires maintaining API integrations and SDK wrappers for 6 different external vendors.
- Ensures ~99.99% uptime for the AI conversational loop.
