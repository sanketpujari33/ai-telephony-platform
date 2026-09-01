# 🔒 Security, Cryptography & Compliance Design Specification

---

## 1. Zero-Secret Architecture (HashiCorp Vault)

Nexora strictly forbids storing API keys, private certificates, or database passwords in source code, Dockerfiles, or environment files.

```text
┌────────────────────────────────────────────────────────┐
│                   Microservice Boot                    │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
           ┌─────────────────────────────────┐
           │ Read VAULT_TOKEN / AppRole JWT │
           └────────────────┬────────────────┘
                            │
                            ▼
           ┌─────────────────────────────────┐
           │ HashiCorp Vault KV v2 (:8200)   │
           │ Path: `secret/data/nexora/dev`  │
           └────────────────┬────────────────┘
                            │ Returns In-Memory Secrets Object
                            ▼
           ┌─────────────────────────────────┐
           │ Freeze Config Object in Memory  │
           │ (`Object.freeze(secrets)`)      │
           └─────────────────────────────────┘
```

---

## 2. Cryptographic Security Standards

- **Symmetric Encryption**: AES-256-GCM authenticated encryption for sensitive database fields (stored with IV and auth tag).
- **Password Hashing**: bcrypt with work factor 12.
- **Webhook HMAC Signatures**: SHA-256 / SHA-512 constant-time comparison (`crypto.timingSafeEqual`) preventing timing attacks.
- **Transport Security**: Mandatory TLS 1.3 for all external ingress and internal service-to-service communication.
- **PII Masking**: Automated tokenization and regex redaction for PAN (Credit Cards), SSNs, and passcodes before logging.
