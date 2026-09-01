# Runbook: Backup & Disaster Recovery (DR)

## Database Backups (PostgreSQL)
- **RPO (Recovery Point Objective)**: 5 Minutes
- **Mechanism**: Continuous WAL (Write-Ahead Log) archiving to S3 via `pgBackRest`.
- **Restoration**:
  To restore a database to a point in time:
  ```bash
  pgbackrest --stanza=aitelephony --type=time --target="2026-08-13 14:00:00" restore
  ```

## Cache Rebuilding (Redis)
- Redis is treated as ephemeral. It holds session state, tenant caching, and port allocation maps.
- **In event of Redis wipe**:
  1. All active calls will drop (state lost).
  2. The system will automatically self-heal on the next inbound call by fetching tenant configs from PostgreSQL and re-seeding the cache.
  3. No manual intervention required.

## Multi-Region Failover
- **Architecture**: Active-Passive setup (e.g., US-East-1 Active, US-West-2 Passive).
- **Trigger**: Complete region failure.
- **Steps**:
  1. Promote the PostgreSQL Read-Replica in US-West-2 to Primary.
  2. Update Global DNS (Route53) to point `sip.telephony.example.com` to the US-West-2 Ingress load balancers.
  3. Notify SIP trunk providers to failover traffic to secondary IPs.
  4. Scale up the US-West-2 HPA limits (which run at 10% capacity normally).
