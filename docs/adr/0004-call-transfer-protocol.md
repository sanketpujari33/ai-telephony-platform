# ADR 0004: Call Transfer Protocol

## Status
Accepted

## Context
When the AI assistant encounters a situation it cannot handle (or the user demands a human), the call must be transferred to a live agent's SIP phone or a contact center queue (e.g., Genesys, Twilio Flex).

## Decision
We use the **SIP REFER (RFC 3515)** method for Blind and Warm transfers, injecting transcript context into custom SIP Headers (`X-AI-Transcript-Context`).

## Rationale
1. **Standard Compliance**: SIP REFER is universally supported by SBCs (Session Border Controllers) and PBXs.
2. **Context Passing**: By injecting `X-AI-Sentiment` and `X-AI-Context` into the REFER message, the downstream contact center's CTI (Computer Telephony Integration) screen can pop up the AI conversation history for the human agent.
3. **B2BUA Fallback**: If the upstream trunk does not support SIP REFER, the SIP Ingress node acts as a B2BUA (Back-to-Back User Agent), dialing out to the human agent and bridging the RTP streams internally.

## Consequences
- Offloads the media processing overhead from our platform once the transfer completes (in the case of a successful REFER).
- Requires complex SIP state machine handling to manage `REFER`, `NOTIFY` (100 Trying, 200 OK), and tear down the AI media session gracefully.
