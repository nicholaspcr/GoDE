# Development Guide

This guide explains how to run `deserver`, `decli`, and the web frontend locally for development and testing.

## Prerequisites

- Go 1.21+
- Node.js 20+
- Docker and Docker Compose
- Make

## Quick Start (Full Stack)

The easiest way to run everything together:

```bash
# 1. Install frontend dependencies
make web-deps

# 2. Start infrastructure (PostgreSQL + Redis)
docker compose -f docker-compose.test.yml up -d

# 3. Start backend server (in one terminal)
JWT_SECRET="development-secret-key-change-in-production-min-32-chars" \
STORE_TYPE=postgres \
STORE_POSTGRESQL_DNS="postgres://gode:gode123@localhost:5432/gode_test?sslmode=disable" \
REDIS_HOST=localhost \
REDIS_PORT=6379 \
go run ./cmd/deserver start

# 4. Start frontend dev server (in another terminal)
make web-dev
```

The services will be available at:
- **Frontend**: http://localhost:5173
- **Backend HTTP API**: http://localhost:8081
- **Backend gRPC**: localhost:3030

The frontend dev server automatically proxies API requests (`/v1/*`) to the backend.

## Backend Development

### 1. Start the infrastructure

```bash
docker compose -f docker-compose.test.yml up -d
```

This starts:
- PostgreSQL on port 5432
- Redis on port 6379

### 2. Start the server

```bash
JWT_SECRET="development-secret-key-change-in-production-min-32-chars" \
STORE_TYPE=postgres \
STORE_POSTGRESQL_DNS="postgres://gode:gode123@localhost:5432/gode_test?sslmode=disable" \
REDIS_HOST=localhost \
REDIS_PORT=6379 \
go run ./cmd/deserver start
```

The server will be available at:
- gRPC: `localhost:3030`
- HTTP: `localhost:8081`

### 3. Verify the server is running

```bash
curl http://localhost:8081/health
# {"status":"UP"}
```

## Using the CLI

### Register a user

```bash
go run ./cmd/decli auth register --username myuser --email myuser@example.com --password mypassword
```

### Login

```bash
go run ./cmd/decli auth login --username myuser --password mypassword
```

### List available algorithms, variants, and problems

```bash
go run ./cmd/decli de list-algorithms
go run ./cmd/decli de list-variants
go run ./cmd/decli de list-problems
```

### Run a DE execution (synchronous)

```bash
go run ./cmd/decli de run \
  --algorithm gde3 \
  --problem zdt1 \
  --variant rand1 \
  --generations 50 \
  --population-size 100
```

### Run a DE execution (asynchronous)

```bash
# Submit execution
go run ./cmd/decli de run-async \
  --algorithm gde3 \
  --problem zdt1 \
  --variant rand1 \
  --generations 100 \
  --population-size 100

# Check status
go run ./cmd/decli de status --execution-id <execution-id>

# Get results when completed
go run ./cmd/decli de results --execution-id <execution-id>
```

### List your executions

```bash
go run ./cmd/decli de list
```

### Delete an execution

```bash
go run ./cmd/decli de delete --execution-id <execution-id>
```

## HTTP API Examples

### Register via HTTP

```bash
curl -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"user":{"ids":{"username":"testuser"},"email":"test@example.com","password":"password123"}}'
```

### Login via HTTP

```bash
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

Save the `accessToken` from the response for authenticated requests.

### Run DE execution via HTTP

```bash
TOKEN="<your-access-token>"

curl -X POST http://localhost:8081/v1/de/run \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "algorithm": "gde3",
    "problem": "zdt1",
    "variant": "rand1",
    "config": {
      "generations": 50,
      "populationSize": 100,
      "dimensionsSize": 30,
      "objectivesSize": 2,
      "ceilLimiter": 1,
      "gde3": {"cr": 0.5, "f": 0.5, "p": 0.5}
    }
  }'
```

## Stopping the Infrastructure

```bash
docker compose -f docker-compose.test.yml down
```

To also remove volumes:

```bash
docker compose -f docker-compose.test.yml down -v
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | Secret key for JWT signing (min 32 chars) | Required |
| `STORE_TYPE` | Database type: `postgres`, `sqlite`, `memory` | `sqlite` |
| `STORE_POSTGRESQL_DNS` | PostgreSQL connection string | - |
| `STORE_SQLITE_FILEPATH` | SQLite database file path | `.dev/server/sqlite.db` |
| `REDIS_HOST` | Redis server hostname | `localhost` |
| `REDIS_PORT` | Redis server port | `6379` |
| `GRPC_PORT` | gRPC server port | `3030` |
| `HTTP_PORT` | HTTP server port | `8081` |

## Troubleshooting

### Port already in use

If you see "address already in use", kill any existing server:

```bash
pkill -f deserver
```

### Database connection issues

Verify PostgreSQL is running:

```bash
docker compose -f docker-compose.test.yml ps
```

### Redis connection issues

Verify Redis is running and accessible:

```bash
docker compose -f docker-compose.test.yml exec redis redis-cli ping
# PONG
```

### Clear Redis cache

If you experience stale data issues:

```bash
docker compose -f docker-compose.test.yml exec redis redis-cli FLUSHALL
```

## Frontend Development

### Install dependencies

```bash
make web-deps
# or
cd web && npm install
```

### Start development server

```bash
make web-dev
# or
cd web && npm run dev
```

The dev server runs at http://localhost:5173 and proxies `/v1/*` requests to `http://localhost:8081`.

### Build for production

```bash
make web-build
```

### Lint code

```bash
make web-lint
```

### Regenerate API client

If the backend API changes, regenerate the TypeScript client:

```bash
make openapi      # Generate OpenAPI spec from protos
make web-api      # Generate TypeScript client from OpenAPI spec
```

## Distributed Tracing with Jaeger

For visualizing distributed traces during development, you can run Jaeger locally.

### Start Jaeger

```bash
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/jaeger:latest
```

> **Note**: Jaeger v2 uses the `jaegertracing/jaeger` image. The previous `jaegertracing/all-in-one` image is v1 (end-of-life).

Ports:
- **16686**: Jaeger UI
- **4317**: OTLP gRPC receiver
- **4318**: OTLP HTTP receiver

### Configure the Server for Jaeger

```bash
JWT_SECRET="development-secret-key-change-in-production-min-32-chars" \
STORE_TYPE=postgres \
STORE_POSTGRESQL_DNS="postgres://gode:gode123@localhost:5432/gode_test?sslmode=disable" \
TRACING_EXPORTER=otlp \
OTLP_ENDPOINT=localhost:4317 \
go run ./cmd/deserver start
```

### View Traces

Open http://localhost:16686 to access the Jaeger UI. Select "deserver" from the Service dropdown to view traces.

### Tracing Configuration

| Variable | Values | Description |
|----------|--------|-------------|
| `TRACING_ENABLED` | `true`/`false` | Enable/disable tracing (default: `true`) |
| `TRACING_EXPORTER` | `none`, `file`, `stdout`, `otlp` | Exporter type (default: `none`) |
| `TRACE_FILE_PATH` | path | Output file for `file` exporter (default: `traces.json`) |
| `OTLP_ENDPOINT` | host:port | OTLP endpoint for `otlp` exporter (default: `localhost:4317`) |
| `TRACE_SAMPLE_RATIO` | 0.0-1.0 | Sampling ratio (default: `1.0` = all traces) |

### Stop Jaeger

```bash
docker stop jaeger && docker rm jaeger
```

## Docker Development

Run the full stack with Docker Compose:

```bash
# Build and start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down
```

Services:
- **Frontend**: http://localhost:3001
- **Backend API**: http://localhost:8081
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
