# Runbook: Deployment & Scaling

## Deployment Procedure
The AI Telephony Platform uses a stateless microservices architecture, deploying via Helm charts onto Kubernetes.

1. **Pre-flight Checks**:
   - Ensure the PostgreSQL schema migrations have been applied (`migrate up`).
   - Validate API Keys for OpenAI, Deepgram, and ElevenLabs are present in the Kubernetes Secrets payload.
2. **Rolling Update**:
   - Execute `helm upgrade ai-telephony ./deploy/helm/ai-telephony --set image.tag=<new_tag>`
   - Kubernetes will execute a rolling update. The `SIP Ingress` pods will stop accepting *new* calls but will allow existing calls to finish before terminating (graceful shutdown up to 5 minutes).

## Scaling Guidelines

### Horizontal Pod Autoscaler (HPA)
- **API Gateway & Orchestrator**: Scaled based on standard CPU utilization (Target 70%).
- **Media Proxy**: 
  - **CRITICAL**: The Media Proxy cannot scale purely on CPU. WebRTC and RTP processing require dedicated port bindings.
  - HPA is configured on **Custom Metrics** via Prometheus: `active_rtp_sessions`.
  - Triggers a scale-up when `active_rtp_sessions > 400` per pod.

### Database Pooling
- Ensure PgBouncer is running in front of PostgreSQL.
- The `Campaign Manager` will exhaust DB connections rapidly during list hydration. Ensure `max_open_conns` in `config.yaml` does not exceed PgBouncer's pool limit.
