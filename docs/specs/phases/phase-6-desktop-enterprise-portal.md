# 💻 Phase 6 — Desktop App & Enterprise Portal Specification
*Sprints 21–23 · Weeks 41–46*

---

## 1. Phase Goals & Scope

Phase 6 delivers dedicated desktop softphone capabilities for contact center agents and an Enterprise Portal for large-scale enterprise organizations requiring SAML SSO, SCIM provisioning, and custom SIP trunking.

---

## 2. Sprint Deliverables

### Sprint 21: Flutter Desktop Application (macOS & Windows)
- **Multi-Line Softphone**: Support for concurrent call hold, transfer, and conference bridging.
- **Global Shortcuts**: System-wide push-to-talk, mute, and call answer hotkeys.
- **Local Audio Loopback**: Real-time microphone/headset diagnostic tests and device selection.
- **Live Interim Transcripts**: Low-latency WebSocket streaming of agent and customer speech.

### Sprint 22: Enterprise SSO & SCIM Provisioning
- **`services/enterprise-service`**: SAML 2.0 (Okta, Azure AD, Ping) & OIDC authentication flows.
- **SCIM 2.0**: Automated user provisioning, deprovisioning, and role syncing from enterprise IdPs.
- **Custom Domains**: Automated SSL certificate provisioning via Let's Encrypt for white-label tenants.

### Sprint 23: Dedicated SIP Trunks & Private vLLM Endpoints
- **Dedicated SIP Trunks**: Direct IP-to-IP peering with enterprise SBCs (Session Border Controllers).
- **Private LLM Gateways**: Enterprise tenant routing to self-hosted vLLM / TGI GPU clusters for zero data retention.

---

## 3. Exit Criteria

- Desktop softphone builds natively for macOS (.dmg) and Windows (.msi/.exe).
- Enterprise tenants successfully authenticate via SAML SSO and manage users via SCIM.
