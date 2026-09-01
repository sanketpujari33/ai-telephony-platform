# 📞 Truecaller Intelligence & Caller Identity Integration Specification

This document details the external caller identity verification and reputation scoring architecture integrating with Truecaller Developer and Truecaller for Business APIs.

---

## 1. Integration Scope & Security

Nexora integrates with Truecaller services to provide:
1. **Caller Identity**: Name, profile verification, and business badge lookup during inbound calls.
2. **Phone Number Verification**: Instant user profile verification with user consent.
3. **Reputation Signals**: Verified business caller-ID and spam reputation scoring where enabled.

> [!IMPORTANT]
> Truecaller API credentials must NEVER be stored in client applications (Flutter, Next.js). All requests are proxied securely through the backend `spam-service` using credentials retrieved dynamically from HashiCorp Vault.

---

## 2. Architecture & Call Ingress Flow

```mermaid
sequenceDiagram
    autonumber
    actor Caller as 👤 Inbound PSTN Caller
    participant SIP as 📞 Go SIP Gateway (:5060)
    participant SpamSvc as 🛡️ spam-service (:4012)
    participant Vault as 🔒 HashiCorp Vault (:8200)
    participant Truecaller as 🌐 Truecaller API
    participant LiveKit as 🎙️ LiveKit Media SFU

    Caller->>SIP: Inbound Call (+15551234567)
    SIP->>SpamSvc: Validate Caller Reputation (+15551234567)
    SpamSvc->>Vault: Read TRUECALLER_API_KEY
    SpamSvc->>Truecaller: Query Caller Identity & Spam Score
    Truecaller-->>SpamSvc: { name: "Acme Corp", score: 0.05, verified: true }
    SpamSvc-->>SIP: Caller Allowed (Reputation Clean)
    SIP->>LiveKit: Route to AI Phone Agent with Caller Context
```

---

## 3. Configuration & Vault Keys

```env
TRUECALLER_CLIENT_ID=
TRUECALLER_APP_KEY=
TRUECALLER_API_KEY=
TRUECALLER_API_SECRET=
TRUECALLER_BASE_URL=https://api4.truecaller.com/v1
```
