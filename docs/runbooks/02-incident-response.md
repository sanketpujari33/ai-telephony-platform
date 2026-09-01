# Runbook: Incident Response

## Scenario 1: STT Provider Degradation
**Symptoms**: Increased call latency, caller saying "hello?" multiple times, dashboard shows `stt_latency_ms` > 800ms.
**Action**:
1. Check Prometheus metrics for `circuit_breaker_state{provider="deepgram"}`.
2. If the circuit breaker hasn't tripped automatically, forcefully trip it via the Admin API:
   ```bash
   curl -X POST https://admin.telephony.example.com/api/v1/circuit-breakers/stt/deepgram/open
   ```
3. This routes all traffic to the fallback provider (AssemblyAI).
4. Monitor fallback provider quotas.

## Scenario 2: SIP Trunk 503 / Timeout Responses
**Symptoms**: Inbound calls are dropping with fast-busies. Outbound campaigns show 0% connect rate.
**Action**:
1. Verify the `Keepalive (OPTIONS)` ping status in the SIP Ingress logs.
2. If the SIP Provider (e.g. Twilio, Bandwidth) is unresponsive, update the DNS SRV records to route inbound DIDs to the secondary region proxy.
3. For outbound, adjust the `weight` parameter in `config.yaml` to push all SIP traffic to the secondary trunk, then trigger a config reload:
   ```bash
   kill -SIGHUP <pid_of_sip_ingress>
   ```

## Scenario 3: DNC Violation Alert
**Symptoms**: System alerts that a dialed number matched the National Do Not Call registry.
**Action**:
1. Immediately halt the associated campaign via API: `POST /campaign/stop`.
2. Pull the Audit Logs for the specific `call_id` to determine why the compliance check failed (e.g., Redis cache out of sync, DB connection dropped causing a fail-open).
3. Do not resume campaigns until the compliance DB synchronization is verified.
