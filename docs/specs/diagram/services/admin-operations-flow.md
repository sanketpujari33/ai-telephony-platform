# 📊 Admin Operations & Management Diagrams Flow

This document details the architectural flows and data structures for the 9 core management modules:
1. **Banner Manager**
2. **Promoter Manager**
3. **Receive Fund Manager**
4. **Profit Loss Manager**
5. **Notifications Engine**
6. **Dead Report Manager**
7. **Data Monitoring & Block Keywords**
8. **Feedback & CSAT Review**
9. **General Tabs & Master Configuration**

---

## 1. Banner Management & Dynamic App Delivery Flow

```mermaid
graph TD
    Admin["👨‍💼 Admin / Marketing Lead"] --> BannerAPI["⚡ Banner Manager API<br/><code>/api/v1/admin/banners</code>"]
    BannerAPI --> S3["🪣 MinIO / S3 Storage (Banner Assets)"]
    BannerAPI --> DB[("🐘 PostgreSQL / Mongo (Banner Metadata)")]
    BannerAPI --> RedisCache[("🔴 Redis Cache <code>nexora:banners:{platform}:{screen}</code>")]

    subgraph ClientDelivery["📱 Dynamic Multi-Platform Ingress"]
        MobileApp["📱 Flutter Mobile (Splash, Home, Dialer, Live)"]
        WebApp["🌐 Next.js Portal (Header, Sidebar, Hero)"]
        DesktopApp["💻 Flutter Desktop (Dialer Banner)"]
    end

    RedisCache -.->|Cached Query <1ms| MobileApp
    RedisCache -.->|Cached Query <1ms| WebApp
    RedisCache -.->|Cached Query <1ms| DesktopApp
```

---

## 2. Promoter, Affiliate & Payout Flow (`Promoter Manager` & `Receive Fund Manager`)

```mermaid
sequenceDiagram
    autonumber
    actor Promoter as 🤝 Promoter / Affiliate
    participant Platform as 🌐 Nexora Platform
    participant PromSvc as 📈 Promoter Manager
    participant FundSvc as 💰 Receive Fund Manager
    participant BankGateway as 🏦 Stripe / Razorpay / Wise
    actor Admin as 👨‍💼 Billing Admin

    Promoter->>Platform: 1. Shares Referral Link / Promo Code
    Platform->>PromSvc: 2. Tracks New Customer Sign-Up & Voice Minutes
    PromSvc->>PromSvc: 3. Calculates Earning Commission (e.g., 10% of Spend)
    PromSvc-->>Promoter: 4. Updates Balance in Promoter Wallet
    
    Promoter->>FundSvc: 5. Submits Withdrawal Request ($500.00)
    FundSvc->>Admin: 6. Creates Pending Fund Request with Beneficiary Bank Details
    Admin->>FundSvc: 7. Approves Withdrawal Request
    FundSvc->>BankGateway: 8. Dispatches Automated Payout Transfer
    BankGateway-->>Promoter: 9. Funds Deposited to Bank Account
    FundSvc-->>FundSvc: 10. Marks Request as 'Paid' & Emits Event
```

---

## 3. Real-Time Profit & Loss Analysis Engine (`Profit Loss Manager`)

```mermaid
graph LR
    subgraph RevenueIngress["💵 Revenue Streams"]
        R1["Call Minute Billing ($0.05/min)"]
        R2["Monthly Subscriptions ($49/mo)"]
        R3["DID Phone Number Leasing ($2/no)"]
    end

    subgraph CostIngress["📉 Direct Operating Costs"]
        C1["Carrier SIP Termination ($0.009/min)"]
        C2["Deepgram STT & Cartesia TTS ($0.012/min)"]
        C3["Groq / DeepSeek LLM Tokens ($0.003/min)"]
        C4["AWS / LiveKit SFU Compute ($0.004/min)"]
    end

    subgraph Engine["🧮 Profit Loss Manager"]
        Calc["Real-Time Margin Aggregator<br/><code>Gross Margin = (Revenue - Direct Costs)</code>"]
        NetCalc["Net Margin = Gross Margin - Affiliate Commissions"]
    end

    subgraph Reports["📊 Executive Financial Reports"]
        Rep1["📈 Daily Profit / Loss Trend"]
        Rep2["📉 Carrier Cost Variance Report"]
        Rep3["🧾 Partner Earning & Margin Report"]
    end

    RevenueIngress --> Calc
    CostIngress --> Calc
    Calc --> NetCalc
    NetCalc --> Reports
```

---

## 4. Multi-Channel Notification Pipeline with DND Protection

```mermaid
graph TD
    Trigger["📢 Trigger Event<br/>(Campaign / Low Balance / Call Missed / Promo)"]
    
    subgraph NotifEngine["🔔 Notification Engine"]
        CheckDND{"Is Recipient in DND (Do-Not-Disturb) Window?"}
        Router["Channel Router"]
        Templates["Template Engine (Variables Injection)"]
    end

    subgraph Channels["📡 Multi-Channel Dispatchers"]
        PushCh["📱 Push Notification (FCM / APNs)"]
        SMSCh["💬 SMS Gateway (Twilio / Sinch)"]
        EmailCh["📧 Email Gateway (SES / Resend)"]
        InAppCh["🔔 In-App WebSocket Alert"]
    end

    Trigger --> CheckDND
    CheckDND -- Yes --> QueueDND["⏳ Queue for Morning Dispatch (09:00 AM)"]
    CheckDND -- No --> Templates
    Templates --> Router
    Router --> PushCh
    Router --> SMSCh
    Router --> EmailCh
    Router --> InAppCh
```

---

## 5. Dead Report & Failed Call Re-Engagement Engine

```mermaid
graph TD
    FailedCall["📞 Failed Call Event (SIP 404 / 503 / Busy / No Answer)"]
    
    subgraph DeadReportManager["💀 Dead Report Manager"]
        Classifier["Error Classifier<br/>(Network Timeout vs. Number Unreachable)"]
        DeadStore[("💾 Dead Reports Database")]
        ReEngageRule{"Automatic Re-engagement Eligible?"}
    end

    subgraph Actions["⚡ Re-engagement Actions"]
        AutoSMS["📱 Auto-dispatch SMS: 'We tried calling you regarding...'"]
        WhatsAppMsg["💬 WhatsApp Reconnect Link"]
        AutoRedial["🔁 Queue Delayed Redial (in 15 mins)"]
    end

    FailedCall --> Classifier
    Classifier --> DeadStore
    Classifier --> ReEngageRule
    ReEngageRule -- SMS Fallback --> AutoSMS
    ReEngageRule -- WhatsApp Fallback --> WhatsAppMsg
    ReEngageRule -- Redial Fallback --> AutoRedial
```

---

## 6. Real-Time Data Monitoring & Keyword Blockage Pipeline

```mermaid
graph TD
    LiveStream["🎙️ Live Call / Chat Utterance Stream"]
    
    subgraph DataMonitoring["🛡️ Data Monitoring Engine"]
        Matcher["Regex & Semantic Keyword Matcher"]
        BlockList[("🚫 Blocked Keywords Database<br/>(Profanity, Fraud, Card Numbers, PII)")]
        SentimentTracker["Real-Time Sentiment Analyzer"]
    end

    subgraph Actions["⚡ Intervention Actions"]
        RedactText["📝 Redact Sensitive Text (***-**-****)"]
        AlertAdmin["🚨 Trigger Supervisor Softphone Alert"]
        DropCall["⛔ Force-Terminate Suspicious Call"]
        FlagReview["🚩 Mark for QA Compliance Review"]
    end

    LiveStream --> Matcher
    BlockList --> Matcher
    Matcher --> RedactText
    Matcher --> AlertAdmin
    LiveStream --> SentimentTracker
    SentimentTracker -- Sentiment < 0.20 --> FlagReview
    SentimentTracker -- Fraud Detected --> DropCall
```
