# Compliance Configuration Guide

The AI Telephony Platform deals with highly regulated telecom traffic. Strict adherence to compliance configurations is mandatory.

## Do Not Call (DNC) List Management
- The `compliance.DNCCheckEnabled` flag must ALWAYS be `true` in production `config.yaml`.
- The system checks both internal DNC tables and the federal registry API before initiating any outbound dial.
- **Opt-Outs**: If a user says "Put me on your do not call list" or "Stop calling me", the LLM detects the intent and fires a system function. The orchestrator intercepts this function and writes the number to the DNC database in real-time, instantly severing the call.

## Calling Windows (TCPA)
- Telemarketing calls are restricted to 08:00 - 21:00 local time of the recipient.
- The `campaign_manager` relies on the `timezone` field provided in the `Contact` payload. 
- If no timezone is provided, the system does a Numbering Plan Area (NPA) area code lookup to estimate the timezone. If ambiguous, the system defaults to the most restrictive overlap.

## PCI DSS Scope Reduction
- The platform is **NOT** designed to record or transcribe credit card numbers.
- **Implementation**: When the LLM reaches a payment collection state, the `STT` stream is paused. The user is prompted to enter their credit card via DTMF (keypad). 
- DTMF tones are parsed by the Media Proxy (`RFC4733`) and sent securely to the payment gateway tokenization endpoint. The raw audio of the DTMF tones is discarded, completely removing the system from full PCI DSS compliance scope.

## Call Recording Consent
- Some states require two-party consent.
- The platform supports a `play_disclaimer_first` tenant flag. If true, the SIP Ingress node plays a pre-recorded standard audio file ("This call may be recorded...") *before* bridging the media to the AI Orchestrator.

## Audit Logs & Data Retention
- All SIP signaling metadata and transcript summaries are kept for 90 days in PostgreSQL, then moved to cold storage (S3 Glacier).
- Access to the `audit_logs` table requires `Admin` RBAC roles. `Operator` roles can only view aggregated metrics.
