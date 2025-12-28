# Development Guide

This guide explains how to run `deserver` and `decli` locally for development and testing.

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make

## Quick Start

### 1. Build the binaries

```bash
make build
```

This creates `./.dev/deserver` and `./.dev/decli`.

### 2. Start the infrastructure

```bash
docker compose -f docker-compose.test.yml up -d
```

This starts:
- PostgreSQL on port 5432
- Redis on port 6379

### 3. Start the server

```bash
JWT_SECRET="development-secret-key-change-in-production-min-32-chars" \
STORE_TYPE=postgres \
STORE_POSTGRESQL_DNS="postgres://gode:gode123@localhost:5432/gode_test?sslmode=disable" \
REDIS_HOST=localhost \
REDIS_PORT=6379 \
./.dev/deserver start
```

The server will be available at:
- gRPC: `localhost:3030`
- HTTP: `localhost:8081`

### 4. Verify the server is running

```bash
curl http://localhost:8081/health
# {"status":"UP"}
```

## Using the CLI

### Register a user

```bash
./.dev/decli auth register --username myuser --email myuser@example.com --password mypassword
```

### Login

```bash
./.dev/decli auth login --username myuser --password mypassword
```

### List available algorithms, variants, and problems

```bash
./.dev/decli de list-algorithms
./.dev/decli de list-variants
./.dev/decli de list-problems
```

### Run a DE execution (synchronous)

```bash
./.dev/decli de run \
  --algorithm gde3 \
  --problem zdt1 \
  --variant rand1 \
  --generations 50 \
  --population-size 100
```

### Run a DE execution (asynchronous)

```bash
# Submit execution
./.dev/decli de run-async \
  --algorithm gde3 \
  --problem zdt1 \
  --variant rand1 \
  --generations 100 \
  --population-size 100

# Check status
./.dev/decli de status --execution-id <execution-id>

# Get results when completed
./.dev/decli de results --execution-id <execution-id>
```

### List your executions

```bash
./.dev/decli de list
```

### Delete an execution

```bash
./.dev/decli de delete --execution-id <execution-id>
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
