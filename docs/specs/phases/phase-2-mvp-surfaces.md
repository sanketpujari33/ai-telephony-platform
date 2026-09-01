# 📱 Phase 2 — MVP Client Surfaces Specification
*Sprints 7–9 · Weeks 13–18*

---

## 1. Phase Goals & Scope

Phase 2 builds the primary client application interfaces that allow end-users, agents, and administrators to interact with the Nexora voice engine:
1. **`apps/mobile-app`**: Flutter cross-platform mobile application (Android & iOS).
2. **`apps/user-dashboard`**: Next.js web application for tenant operations and agent management.
3. **`apps/admin-panel`**: Next.js internal administration control center.
4. **`apps/landing-page`**: Public marketing website with interactive 3-minute voice AI sandbox.

---

## 2. Sprint Deliverables

### Sprint 7: Flutter Mobile App (Android / iOS)
- **LiveKit WebRTC Audio**: Low-latency voice calls directly from mobile softphone.
- **Dialpad & Contacts**: Keypad dialer, native contact integration, recent call logs.
- **Push Notifications**: Firebase Cloud Messaging (FCM) & Apple APNs incoming call alerts.
- **Biometric Auth & Phone OTP**: Fast login with fingerprint/FaceID.

### Sprint 8: User Dashboard & Agent Control
- **Next.js 14 App Router**: Responsive dark-mode dashboard built with `@nexora/ui`.
- **Agent Studio**: Configure agent prompt instructions, temperature, voice selection (Cartesia/ElevenLabs), and tool bindings.
- **Call Analytics & History**: Audio player with synchronized waveform, transcript viewer, and sentiment badges.
- **Workspace Settings**: Team member invitations, API key generation, webhook settings.

### Sprint 9: Admin Panel v1 & Public Landing Page
- **Admin Panel**: Tenant provisioning, platform resource monitoring, global kill-switches.
- **Landing Page**: Interactive 3-minute WebRTC voice sandbox allowing prospective users to talk to the AI live.
- **Documentation Portal (`apps/documentation`)**: Nextra/Next.js interactive API and SDK documentation.

---

## 3. Exit Criteria

- Mobile softphone places and receives high-fidelity WebRTC voice calls.
- Web dashboard allows prompt updates and displays live call logs in real-time.
- Landing page voice demo operates seamlessly in browser without plugins.
