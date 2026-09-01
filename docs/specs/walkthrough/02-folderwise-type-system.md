# 📁 Walkthrough: Folderwise Type System in `@nexora/types`

This walkthrough explains the three-tier domain categorization implemented in `@nexora/types`.

---

## 1. Directory Structure Overview

The [`@nexora/types`](file:///d:/trideocoders/ai-telephony-platform/packages/types) package is organized into three clean folder-based categories:

```text
packages/types/src/
├── index.ts                      ── Root barrel exporting packages/, services/, and apps/
├── packages/                     ── Core shared @nexora/* package domain types
│   ├── agent.ts
│   ├── ai.ts
│   ├── billing.ts
│   ├── call.ts
│   ├── config.ts
│   ├── database.ts
│   ├── events.ts
│   ├── health.ts
│   ├── logger.ts
│   ├── memory.ts
│   ├── notification.ts
│   ├── pagination.ts
│   ├── prompts.ts
│   ├── sdk.ts
│   ├── security.ts
│   ├── shared.ts
│   ├── storage.ts
│   ├── telemetry.ts
│   ├── tool.ts
│   ├── ui.ts
│   ├── user.ts
│   ├── workflow.ts
│   └── index.ts
├── services/                     ── Backend microservices contracts (API Gateway, Telephony, Voice, Orchestrator...)
│   ├── admin.ts
│   ├── analytics.ts
│   ├── auth.ts
│   ├── gateway.ts
│   ├── orchestrator.ts
│   ├── telephony.ts
│   ├── voice.ts
│   └── index.ts
└── apps/                         ── Frontend client application contracts (Dashboard, Admin Panel, Mobile, Docs)
    ├── admin.ts
    ├── client-app.ts
    ├── dashboard.ts
    ├── docs.ts
    └── index.ts
```

---

## 2. Verification Status

- `pnpm run type-check`: **28 / 28 tasks passed (0 errors)**.
- `pnpm --filter @nexora/tests test`: **15 / 15 test suites passed (56 tests passed)**.
- `pnpm test`: **27 / 27 Turbo tasks passed (FULL TURBO / 0 warnings)**.
