# 💳 Phase 4 — Admin, Analytics & Billing Maturity Specification
*Sprints 13–15 · Weeks 25–30*

---

## 1. Phase Goals & Scope

Phase 4 introduces complete financial automation, usage metering, subscription tiering, carrier P&L margins, and the 9 specialized Admin Operations workflows.

---

## 2. Sprint Deliverables

### Sprint 13: Subscriptions, Wallets & Metered Billing
- **`services/billing-service`**: Minute-level voice metering, LLM token usage tracking, and multi-currency ledger.
- **`services/payment-service`**: Stripe Billing, Razorpay Subscriptions, LemonSqueezy, Paddle webhook verifications.
- **Auto-Recharge**: Automated balance replenishment when prepaid minutes fall below threshold.

### Sprint 14: 9 Admin Operations Modules
- **`services/admin-service`** & **`apps/admin-panel`**:
  1. **Banners Manager**: Targeted marketing & scheduled announcements.
  2. **Promoters & Affiliates**: Referral codes, tier commission calculation, payout ledger.
  3. **Receive Funds**: Manual wire top-ups, dispute refunds, offline transaction approval.
  4. **Profit & Loss**: Per-minute carrier cost vs retail margin calculations.
  5. **Dead Reports**: Inactive campaigns, broken trunks, dead-letter jobs.
  6. **Data Monitoring**: Real-time spam keyword filters, regulatory blacklists.
  7. **Feedback & CSAT**: Post-call survey scoring and sentiment flags.
  8. **General Tabs**: Dynamic portal settings, feature toggles.
  9. **RBAC CRUD**: Custom role and permission scope creation.

### Sprint 15: Deep Analytics & Observability Dashboards
- **`services/analytics-service`**: Real-time call concurrency, Mean Opinion Score (MOS) metrics, carrier latency graphs.
- **Grafana & Prometheus**: Pre-configured dashboards for CPU, memory, voice latency (TTFAB), and billing revenue.

---

## 3. Exit Criteria

- Subscriptions, overages, and top-ups process automatically via Stripe/Razorpay.
- Admin Panel operates all 9 modules with strict RBAC permission enforcement.
- Carrier telephony costs and profit margins are calculated in real time.
