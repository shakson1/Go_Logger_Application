# Go Logger Application

A simple, modular, and production-ready logger service written in Go. It features log ingestion, a built-in web UI, real-time updates, basic analytics, Prometheus metrics, and is ready for Docker, Kubernetes, and Helm deployment.

## Features
- **Log Ingestion:** Accepts logs via HTTP POST (JSON) on port 9000
- **In-Memory Storage:** Stores logs in a thread-safe in-memory DB (for demo; can be extended)
- **Web UI:** Built-in Go web server (port 8080) with:
  - Search/filter by timestamp, level, and keywords
  - Table view of logs
  - Bar chart (logs per minute, last hour)
  - Pie chart (log level distribution)
  - Real-time updates (polling)
- **Prometheus Metrics:** `/metrics` endpoint for monitoring
- **Docker-ready:** Multi-stage Dockerfile for small images
- **Kubernetes-ready:** Deployment and Service YAMLs
- **Helm Chart:** For easy configuration and deployment

## Quick Start

### 1. Run Locally
```sh
go run main.go
```
- Web UI: [http://localhost:8080](http://localhost:8080)
- Log ingestion: POST to [http://localhost:9000/logs](http://localhost:9000/logs)

### 2. Build & Run with Docker
```sh
docker build -t logger:latest .
docker run -p 8080:8080 -p 9000:9000 logger:latest
```

### 3. Test Log Ingestion
```sh
curl -X POST http://localhost:9000/logs \
  -H "Content-Type: application/json" \
  -d '{"level":"INFO","message":"Hello from curl!","metadata":{"source":"curl"}}'
```

### 4. Kubernetes Deployment
- Build and push your Docker image to a registry accessible by your cluster.
- Apply manifests:
```sh
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```
- Web UI: NodePort 30080
- Log ingestion: NodePort 30900

### 5. Helm Chart
```sh
helm install logger .
# Or to customize:
helm install logger . -f values.yaml
```

## Endpoints
- **Web UI:** `GET /` (port 8080)
- **Log Ingestion:** `POST /logs` (port 9000)
- **API:**
  - `GET /api/logs` (with filters)
  - `GET /api/stats` (for charts)
  - `GET /api/logs/stream` (for polling)
- **Prometheus Metrics:** `GET /metrics` (port 8080)

## Configuration
- `UI_PORT` (default: 8080)
- `LOG_INGEST_PORT` (default: 9000)
- `LOG_RETENTION_HOURS` (default: 24)

## Extending
- Swap in BoltDB, SQLite, or another DB for persistence
- Add authentication, rate limiting, or alerting
- Replace polling with WebSocket for true real-time

## License
MIT 