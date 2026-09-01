# ADR 0003: Predictive Dialer Algorithm

## Status
Accepted

## Context
For outbound campaigns, we need to balance the number of concurrent outbound calls with the number of available "AI Agents" (concurrency limits) to minimize the drop rate while maximizing agent utilization.

## Decision
We will use a modified **Erlang C** algorithm combined with real-time Answering Machine Detection (AMD) telemetry to dictate pacing.

## Rationale
1. **Live Connect Rate**: The dialer continuously monitors the real-time "live human pickup" rate (e.g., 20%). 
2. **Abandonment Threshold**: By law (TCPA), abandonment rates must stay below 3%. 
3. **Algorithm**: `Pacing_Lines = (Available_Agents / Connect_Rate) - Error_Margin`.
   - If we have 10 available concurrent LLM channels, and a 50% connect rate, we dial 18 lines (leaving a 10% error margin to avoid abandoning a live human).
4. **AMD Integration**: Calls detected as voicemails or answering machines are immediately disconnected by the SIP proxy and do *not* consume an AI Agent slot, allowing aggressive over-dialing.

## Consequences
- Drastically increases campaign throughput compared to a power dialer (1:1 ratio).
- Requires strict Redis-backed state tracking to ensure we never exceed our AI provider concurrency limits (which would cause dropped calls).
