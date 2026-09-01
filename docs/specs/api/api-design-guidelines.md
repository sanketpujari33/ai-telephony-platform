# 📐 API Design Guidelines & Standards

This document establishes the universal design guidelines, protocol conventions, error models, and pagination standards for all Nexora REST and WebSocket APIs.

---

## 1. URL Path Structure & Naming

All endpoints follow strict RESTful resource-oriented naming:
- **Base Prefix**: `/api/v1/`
- **Resource Nouns**: Lowercase, kebab-case, pluralized (e.g., `/api/v1/call-sessions`, `/api/v1/voice-agents`).
- **Sub-resources**: Natural hierarchical nesting (e.g., `/api/v1/call-sessions/{sessionId}/transcripts`).
- **Actions**: For non-CRUD operations, append a verb action (e.g., `/api/v1/calls/{callId}/transfer`, `/api/v1/calls/{callId}/hangup`).

---

## 2. HTTP Status Code Mapping

| Status Code | Meaning | Usage Scenario |
| :---: | :--- | :--- |
| **`200 OK`** | Request Succeeded | Standard GET, PUT, PATCH responses with body. |
| **`201 Created`** | Resource Created | Successful POST creation with the created entity. |
| **`204 No Content`** | Success Without Body | Successful DELETE or background trigger. |
| **`400 Bad Request`** | Validation Error | Malformed JSON, missing fields, invalid types. |
| **`401 Unauthorized`** | Missing Authentication | Missing, expired, or invalid JWT / API Key. |
| **`403 Forbidden`** | Insufficient Permissions | Valid auth, but missing required RBAC scope. |
| **`404 Not Found`** | Resource Missing | Non-existent entity ID or route. |
| **`409 Conflict`** | State Conflict | Duplicate unique key (e.g., email, phone number). |
| **`422 Unprocessable`** | Semantic Error | Syntactically valid JSON, but violates business rule. |
| **`429 Too Many Requests`** | Rate Limit Exceeded | Tenant token bucket exhausted. |
| **`500 Internal Error`** | Server Fault | Unexpected unhandled exception (logged to Sentry). |
| **`503 Service Unavailable`** | Dependency Down | Downstream carrier, database, or LLM outage. |

---

## 3. Standard Response Envelopes

### A. Success Envelope
```json
{
  "success": true,
  "statusCode": 200,
  "data": { ... },
  "message": "Operation completed successfully",
  "correlationId": "req_89ab-12cd-34ef",
  "timestamp": "2026-08-18T10:00:00.000Z"
}
```

### B. Offset Pagination Envelope (Admin & Tables)
```json
{
  "success": true,
  "statusCode": 200,
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "pageSize": 25,
    "totalPages": 4,
    "totalCount": 92,
    "hasNextPage": true,
    "hasPrevPage": false
  },
  "message": "Items retrieved successfully"
}
```

### C. Cursor Pagination Envelope (Realtime Feeds & Calls)
```json
{
  "success": true,
  "statusCode": 200,
  "data": [ ... ],
  "pagination": {
    "limit": 50,
    "nextCursor": "eyJjcmVhdGVkQXQiOjE3ODY5Nzg5MDB9",
    "prevCursor": null,
    "hasNextPage": true,
    "hasPrevPage": false,
    "totalCount": 1420
  }
}
```

### D. Error Envelope
```json
{
  "success": false,
  "statusCode": 400,
  "error": "VALIDATION_FAILED",
  "message": "The payload contains invalid field values.",
  "details": [
    {
      "field": "callerId",
      "value": "12345",
      "issue": "Phone number must be E.164 compliant."
    }
  ],
  "correlationId": "req_89ab-12cd-34ef",
  "timestamp": "2026-08-18T10:00:00.000Z"
}
```
