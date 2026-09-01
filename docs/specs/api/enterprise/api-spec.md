# 🏢 Next.js Enterprise Portal — Complete API Reference

---

## 📌 1. Enterprise SSO & SAML 2.0 Integration

### `POST /api/v1/enterprise/sso/config`
Configures or updates Enterprise Single Sign-On (Okta, Microsoft Entra ID, PingIdentity).

#### Request Headers:
```http
Authorization: Bearer <ENTERPRISE_ADMIN_TOKEN>
x-tenant-id: tenant_enterprise_01
Content-Type: application/json
```

#### Request Body:
```json
{
  "provider": "saml",
  "idpIssuer": "http://www.okta.com/exk192837465",
  "idpSsoUrl": "https://company.okta.com/app/nexora/exk192837465/sso/saml",
  "idpCertificate": "-----BEGIN CERTIFICATE-----\nMIIDpDCCAoygAwIBAgIGAX...\n-----END CERTIFICATE-----",
  "enforceSso": true,
  "allowedDomains": ["company.com", "subsidiary.com"],
  "defaultRole": "agent_operator"
}
```

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": {
    "ssoConfigId": "sso_okta_enterprise_01",
    "status": "active",
    "acsUrl": "https://api.nexora.ai/api/v1/auth/sso/saml/acs",
    "entityId": "https://api.nexora.ai/saml/metadata/tenant_enterprise_01",
    "updatedAt": "2026-08-17T15:30:00Z"
  },
  "message": "SSO configuration saved and activated"
}
```

---

## 📞 2. Dedicated Carrier SIP Trunks & DIDs

### `GET /api/v1/enterprise/sip/trunks`
Lists all dedicated private SIP trunks with capacity, MOS quality score, and jitter telemetry.

#### Query Parameters:

| Parameter | Type | Required | Default | Description |
| :--- | :--- | :---: | :---: | :--- |
| `page` | `integer` | No | `1` | Page number. |
| `pageSize` | `integer` | No | `20` | Page size. |
| `status` | `string` | No | - | Filter: `connected`, `degraded`, `disconnected`. |

#### Response `200 OK`:
```json
{
  "success": true,
  "statusCode": 200,
  "data": [
    {
      "id": "trunk_ent_01",
      "name": "Primary Telnyx US-East",
      "sipHost": "sip.telnyx.com:5061",
      "transport": "TLS",
      "srtpEnabled": true,
      "status": "connected",
      "activeChannels": 42,
      "maxChannels": 250,
      "telemetry": {
        "mosScore": 4.38,
        "roundTripLatencyMs": 24.5,
        "packetLossPercent": 0.01,
        "jitterMs": 2.1
      },
      "createdAt": "2026-02-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "totalPages": 1,
    "totalCount": 1,
    "hasNextPage": false,
    "hasPrevPage": false
  }
}
```

---

## 🧠 3. Dedicated / Self-Hosted LLM Endpoints

### `POST /api/v1/enterprise/models/custom`
Registers an internal enterprise on-premise model cluster (vLLM / TensorRT).

#### Request Body:
```json
{
  "name": "Internal Llama-3.3-70B Cluster",
  "baseUrl": "https://vllm-internal.company.com/v1",
  "apiKey": "sec_vllm_live_private_token",
  "modelIdentifier": "meta-llama/Llama-3.3-70B-Instruct",
  "maxTokens": 2048,
  "temperature": 0.6,
  "concurrencyLimit": 80,
  "timeoutMs": 12000
}
```

#### Response `201 Created`:
```json
{
  "success": true,
  "statusCode": 201,
  "data": {
    "customModelId": "mod_llama70b_internal_01",
    "name": "Internal Llama-3.3-70B Cluster",
    "status": "active",
    "benchmarkLatencyMs": 118,
    "registeredAt": "2026-08-17T15:35:00Z"
  }
}
```

---

## 🔒 4. Compliance & Audit Exports

### `POST /api/v1/enterprise/compliance/export`
Initiates an asynchronous encrypted audit package export for SOC 2 or HIPAA compliance.

#### Request Body:
```json
{
  "startDate": "2026-01-01T00:00:00Z",
  "endDate": "2026-06-30T23:59:59Z",
  "includeAudioRecordings": false,
  "includeTranscripts": true,
  "includeAuditLogs": true,
  "encryptionAlgorithm": "PGP-RSA-4096",
  "pgpPublicKey": "-----BEGIN PGP PUBLIC KEY BLOCK-----\nmQINBF..."
}
```

#### Response `202 Accepted`:
```json
{
  "success": true,
  "statusCode": 202,
  "data": {
    "exportJobId": "exp_soc2_99182374921",
    "status": "processing",
    "estimatedCompletionMinutes": 8,
    "downloadUrlExpiryHours": 24
  },
  "message": "Compliance export job queued. Notification will be dispatched upon completion."
}
```
