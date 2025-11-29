# procspy

> The secret agent for your system: real-time CPU, RAM, Disk, Network, and GPU metrics backend in Go.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Docker Setup](#docker-setup)
- [Examples](#examples)
- [Next Steps](#next-steps)

---

## Overview

`procspy` is a self-hosted, cross-platform system monitoring backend written in Go. It collects system metrics like CPU, RAM, Disk, Network, and GPU usage from multiple machines in real-time. The metrics are stored in a database and can be queried via REST API or streamed via WebSocket.

Perfect for devs, sysadmins, and hobbyists who want a lightweight, hacker-cool monitoring solution.

---

## Features

- Cross-platform agent (Linux, Windows, Mac) collects system metrics
- REST API for registering machines and pushing metrics
- WebSocket endpoint for real-time metric streams
- Metrics storage in Postgres/SQLite
- Prometheus endpoint for integration
- Configurable collection intervals
- Optional alerts (CPU>90%, RAM>80%, etc.)

---

## Architecture

**Components:**

1. **Agent (`pkg/agent`)**

   - Runs on each machine
   - Collects CPU, RAM, Disk, Network, GPU metrics
   - Sends metrics to backend via REST/WebSocket

2. **Collector (`pkg/collector`)**

   - Abstracts OS-specific metric collection

3. **Backend Server (`cmd/procspy-server`)**

   - REST API to register agents and push metrics
   - WebSocket endpoint for real-time streams
   - Prometheus metrics endpoint

4. **Store (`pkg/store`)**

   - Persists machines and metrics
   - Tables: `machines`, `metrics`, `alerts`

5. **Scheduler (`internal/scheduler`)**

   - Cleans old metrics
   - Optional alert evaluation

6. **Auth (`internal/auth`)**
   - API key or JWT for agents

---

## Getting Started

### Clone the repo

```bash
git clone git@github.com:upayanmazumder/procspy.git
cd procspy
```

### Start Backend & DB

```bash
docker-compose up -d
# then run server in dev mode
go run ./cmd/procspy-server
```

### Agent Example

Run the agent example to start sending metrics:

```bash
go run ./examples/agent-example/main.go
```

---

## API Endpoints

### POST `/api/v1/machines/register`

Register a new agent machine.

Request body:

```json
{
	"machine_name": "laptop-01",
	"os": "linux",
	"agent_version": "0.1.0"
}
```

Response:

```json
{ "status": "registered", "machine_id": "uuid" }
```

### POST `/api/v1/metrics`

Push metrics from agent.

Request body:

```json
{
	"machine_id": "uuid",
	"cpu": 25.3,
	"ram": 40.5,
	"disk": 70.2,
	"network_in": 1200,
	"network_out": 800
}
```

Response:

```json
{ "status": "metrics_received" }
```

### GET `/health`

Simple health check.

Response:

```
procspy is live!
```

---

## Docker Setup

```yaml
version: "3.8"
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: procspy
      POSTGRES_PASSWORD: procspy
      POSTGRES_DB: procspy
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  procspy-server:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres

volumes:
  pgdata:
```

---

## Examples

Check out `examples/agent-example/` to see how to run an agent locally and start sending metrics to the backend.

---

## Next Steps

- Implement `pkg/agent` using `gopsutil` for cross-platform metrics collection
- Add DB integration for `machines` and `metrics` tables
- Implement WebSocket streaming for live metric updates
- Prometheus `/metrics` endpoint for observability
- Alerts system (CPU, RAM, Disk thresholds)
- Optional authentication with API keys or JWTs
