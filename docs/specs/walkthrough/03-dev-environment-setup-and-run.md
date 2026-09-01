# 💻 Walkthrough: Local Development Setup & Operations

This walkthrough guides developers through cloning, initializing, running, and testing the entire Nexora development environment.

---

## 1. Prerequisites
- **Node.js**: >= 20.x
- **pnpm**: >= 9.x (`npm install -g pnpm`)
- **Docker & Docker Compose**: For local databases and infrastructure.
- **Python**: >= 3.11 with `venv`
- **Go**: >= 1.22

---

## 2. Universal Auto-Detecting Scripts

All setup and lifecycle scripts auto-detect whether you are on **Windows**, **macOS**, or **Linux**:

```bash
# 1. Initialize environment, install dependencies, build packages, start Docker DBs, seed initial data
pnpm run setup

# 2. Start all microservices, web portals, and agents
pnpm run start:all

# 3. Check live health across all 17 ports
pnpm run healthcheck

# 4. Gracefully stop all services and containers
pnpm run stop:all
```

---

## 3. Local Ports & Dashboards

| Service / Portal | Local URL | Default Credentials |
| :--- | :--- | :--- |
| **User Dashboard** | `http://localhost:3000` | - |
| **Admin Panel** | `http://localhost:3001` | - |
| **Landing Page** | `http://localhost:3002` | - |
| **Documentation** | `http://localhost:3003` | - |
| **API Gateway** | `http://localhost:4000` | - |
| **MinIO Console** | `http://localhost:9001` | `admin` / `adminpassword` |
| **Temporal UI** | `http://localhost:8088` | - |
| **HashiCorp Vault**| `http://localhost:8200` | Token: `root` |
