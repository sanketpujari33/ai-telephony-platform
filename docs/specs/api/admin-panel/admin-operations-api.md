# 💼 Admin Operations & Platform Modules — Complete API Reference

This document outlines the API specifications for the 9 operational administration modules:
1. **Banner Manager**
2. **Promoter Manager**
3. **Receive Fund Manager**
4. **Profit Loss Manager**
5. **Notification Engine**
6. **Dead Report Manager**
7. **Data Monitoring & Block Keywords**
8. **Feedback & CSAT Review**
9. **General Tabs & Master Settings**

---

## 🖼️ 1. Banner Manager APIs

### `GET /api/v1/admin/banners`
Retrieves configured promotional banners for mobile apps, web portal, or desktop softphone.
- **Required Scope**: `banners:read`

#### Query Parameters:

| Parameter | Type | Required | Default | Description |
| :--- | :--- | :---: | :---: | :--- |
| `platform` | `string` | No | `all` | Filter: `mobile_customer`, `mobile_agent`, `web`, `desktop`. |
| `screen` | `string` | No | - | Filter: `splash`, `home`, `dialer`, `live_call`, `settings`. |
| `status` | `string` | No | `active`| `active`, `inactive`, `scheduled`. |

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": [
    {
      "id": "bnr_summer_promo_01",
      "type": "promotional",
      "bannerType": "slider",
      "platform": "mobile_customer",
      "screenLocation": "home",
      "image": "https://s3.nexora.ai/banners/summer_voice_promo.png",
      "url": "https://nexora.ai/promotions/summer-voice",
      "priority": 1,
      "isActive": true,
      "createdAt": "2026-08-01T00:00:00Z"
    }
  ]
}
```

### `POST /api/v1/admin/banners`
Creates a new promotional banner.
- **Required Scope**: `banners:create`

#### Request Body:
```json
{
  "bannerType": "hero",
  "platform": "web",
  "screenLocation": "home",
  "image": "https://s3.nexora.ai/banners/enterprise_ai_hero.png",
  "url": "https://nexora.ai/enterprise",
  "priority": 2,
  "startDate": "2026-08-15T00:00:00Z",
  "endDate": "2026-09-15T23:59:59Z"
}
```

---

## 🤝 2. Promoter & Affiliate Manager APIs

### `GET /api/v1/admin/promoters`
Lists all affiliate partners, referred customers, and commission earnings.
- **Required Scope**: `promoters:read`

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": [
    {
      "promoterId": "prm_tech_influencer_01",
      "name": "Alex Rivera",
      "email": "alex@techreviews.com",
      "referralCode": "ALEXAI20",
      "commissionRatePercent": 15.0,
      "totalReferredTenants": 84,
      "totalRevenueGeneratedUsd": 38450.00,
      "unpaidEarningsUsd": 1250.00,
      "status": "active"
    }
  ],
  "pagination": { "page": 1, "pageSize": 20, "totalPages": 1, "totalCount": 1 }
}
```

---

## 💰 3. Receive Fund Manager (Withdrawals & Payouts)

### `GET /api/v1/admin/funds/requests`
Lists all affiliate and partner withdrawal payout requests.
- **Required Scope**: `funds:read`

#### Query Parameters:
- `status` (`pending` | `approved` | `rejected` | `paid`)

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": [
    {
      "requestId": "wth_891238472",
      "beneficiaryName": "Alex Rivera",
      "amountUsd": 1250.00,
      "currency": "USD",
      "payoutMethod": "bank_transfer",
      "bankDetails": {
        "bankName": "Chase Bank",
        "accountNumberMasked": "******7891",
        "routingNumber": "021000021"
      },
      "status": "pending",
      "requestedAt": "2026-08-17T14:30:00Z"
    }
  ]
}
```

### `POST /api/v1/admin/funds/requests/:id/approve`
Approves the fund withdrawal and triggers automated bank payout dispatch.
- **Required Scope**: `funds:approve`

---

## 🧮 4. Profit Loss Manager APIs

### `GET /api/v1/admin/profit-loss/summary`
Calculates real-time P&L margins across carrier termination, AI tokens, and retail billing.
- **Required Scope**: `profit_loss:read`

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "grossRevenueUsd": 128450.00,
    "directCosts": {
      "carrierSipTerminationUsd": 23120.00,
      "aiModelTokensUsd": 18450.00,
      "sttAndTtsAudioUsd": 14200.00,
      "livekitServerComputeUsd": 6200.00
    },
    "grossProfitUsd": 66480.00,
    "grossMarginPercent": 51.75,
    "promoterCommissionsUsd": 8450.00,
    "netOperatingProfitUsd": 58030.00,
    "netMarginPercent": 45.18
  }
}
```

---

## 🔔 5. Notification Engine APIs

### `POST /api/v1/admin/notifications/send`
Dispatches multi-channel announcements respecting recipient Do-Not-Disturb (DND) windows.
- **Required Scope**: `notifications:create`

#### Request Body:
```json
{
  "title": "Platform Scheduled Maintenance",
  "body": "Nexora SIP routing engine will undergo a 5-minute zero-downtime maintenance on Sunday at 02:00 AM UTC.",
  "channels": ["push", "in_app", "email"],
  "targetAudience": {
    "tier": ["pro", "enterprise"]
  },
  "respectDND": true
}
```

---

## 💀 6. Dead Report Manager APIs

### `GET /api/v1/admin/dead-reports`
Retrieves call records that failed due to SIP carrier errors, dead lines, or network drops.
- **Required Scope**: `dead_reports:read`

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": [
    {
      "deadReportId": "dr_89123847",
      "callId": "call_err_991823749",
      "destinationNumber": "+15559870000",
      "sipErrorCode": "503_SERVICE_UNAVAILABLE",
      "failureReason": "carrier_route_congestion",
      "carrier": "Telnyx US-East",
      "autoReEngaged": true,
      "reEngagementAction": "auto_sms_dispatched",
      "recordedAt": "2026-08-17T15:10:00Z"
    }
  ]
}
```

---

## 🛡️ 7. Data Monitoring & Blocked Keywords

### `GET /api/v1/admin/data-monitoring/blocked-keywords`
Lists active prohibited words and regex patterns for real-time redaction.
- **Required Scope**: `data_monitoring:read`

### `POST /api/v1/admin/data-monitoring/blocked-keywords`
Adds a new blocked keyword or pattern.
- **Required Scope**: `data_monitoring:create`

#### Request Body:
```json
{
  "pattern": "cvv|card_number|ssn",
  "category": "compliance_pii",
  "action": "redact_and_alert"
}
```

---

## ⭐ 8. Feedback & CSAT Review APIs

### `GET /api/v1/admin/feedback`
Lists customer ratings and post-call feedback comments.
- **Required Scope**: `feedback:read`

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": [
    {
      "feedbackId": "fb_99182374",
      "callId": "call_dVMlFQkC2YLHGRWstTNv_",
      "rating": 5,
      "npsScore": 10,
      "comment": "The AI voice sounded completely natural and resolved my billing issue in seconds!",
      "customerPhone": "+15552345678",
      "createdAt": "2026-08-17T15:20:00Z"
    }
  ]
}
```

---

## ⚙️ 9. General Tabs & Master Settings APIs

### `GET /api/v1/admin/general-tabs`
Fetches all global master application configurations.
- **Required Scope**: `general_tabs:read`

### `PUT /api/v1/admin/general-tabs/payment-gateways`
Updates platform payment gateway credentials.
- **Required Scope**: `general_tabs:edit`

#### Request Body:
```json
{
  "stripe": {
    "publishableKey": "pk_live_51Nx...",
    "secretKey": "sk_live_51Nx...",
    "webhookSecret": "whsec_..."
  },
  "razorpay": {
    "keyId": "rzp_live_...",
    "keySecret": "sec_..."
  }
}
```
