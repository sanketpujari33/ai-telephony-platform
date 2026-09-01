# 🏃 Sprints 10–12 — Omnichannel Messaging Plan
*Weeks 19–24 · Phase 3*

---

## Sprint 10: WhatsApp Business Cloud API
- **Objective**: Build bidirectional WhatsApp AI messaging.
- **Tasks**:
  - [ ] Implement `agents/whatsapp-agent` with Meta Cloud API.
  - [ ] Support interactive button cards, lists, and media documents.
  - [ ] Manage 24-hour customer service session window.

---

## Sprint 11: Multi-Carrier SMS Communications
- **Objective**: Implement transactional and conversational SMS.
- **Tasks**:
  - [ ] Implement `agents/sms-agent` across Twilio, Telnyx, Sinch, and Infobip.
  - [ ] Automated 10DLC compliance and opt-out keyword handling.
  - [ ] Multi-carrier failover routing in `services/notification-service`.

---

## Sprint 12: Email AI & Unified Customer Timeline
- **Objective**: Implement smart email processing and unified inbox.
- **Tasks**:
  - [ ] Implement `agents/email-agent`: Inbound email triaging, reply auto-drafting.
  - [ ] Unify Voice, WhatsApp, SMS, and Email events into a single customer CRM timeline.
