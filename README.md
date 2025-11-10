# GoDE - Differential Evolution Framework

GoDE is a production-ready Differential Evolution (DE) framework built as a gRPC/HTTP server, enabling concurrent execution of multi-objective optimization algorithms across multiple users.

This project extends [GDE3](https://github.com/nicholaspcr/GDE3), which originated from scientific research at [CEFET-MG](https://www.cefetmg.br/).

## Features

### Multi-Objective Optimization
- **GDE3 Algorithm**: Generalized Differential Evolution
- **6 Mutation Variants**: rand/1, rand/2, best/1, best/2, pbest, current-to-best/1
- **22 Benchmark Problems**: ZDT, DTLZ, WFG families

### Async Execution Architecture
- **Background Job Processing**: Long-running optimizations don't block API requests
- **Redis-Backed State**: Fast access to execution status and progress
- **Real-time Progress Streaming**: Server-sent events for live progress updates
- **Cancellation Support**: Stop running executions on demand
- **Composite Storage**: Hybrid Redis/database architecture for performance and persistence
- **Worker Pool**: Configurable concurrency limits for resource management

### Production-Ready Architecture
- **gRPC + HTTP Gateway**: Dual protocol support
- **JWT Authentication**: Secure user authentication
- **Database Support**: PostgreSQL, SQLite, in-memory
- **Redis Integration**: Required for async execution support
- **Database Migrations**: Version-controlled schema evolution
- **Rate Limiting**: Per-IP auth limiting, per-user DE execution limiting
- **TLS/HTTPS Support**: Secure communication
- **Health Checks**: Liveness and readiness probes (includes Redis health)
- **Graceful Shutdown**: Zero-downtime deployments

### Observability
- **OpenTelemetry Tracing**: Distributed tracing support
- **Prometheus Metrics**: Comprehensive metrics collection
- **Structured Logging**: JSON logging with slog
- **Panic Recovery**: Automatic recovery with stack traces

## Quick Start

### Prerequisites
- Go 1.25 or later
- Redis 6.0 or later (required for async execution)
- Make (optional, for convenience commands)
- PostgreSQL 12+ (optional, recommended for production)

### Installation

```bash
# Clone the repository
git clone https://github.com/nicholaspcr/GoDE.git
cd GoDE

# Build binaries
make build

# Binaries will be in .dev/decli and .dev/deserver
```

### Configuration

1. Start Redis (required for async execution):
```bash
# Using Docker
docker run -d -p 6379:6379 redis:latest

# Or using package manager
# Ubuntu/Debian
sudo apt-get install redis-server
sudo systemctl start redis

# macOS
brew install redis
brew services start redis
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Generate a secure JWT secret:
```bash
# Generate a random 32+ character secret
openssl rand -base64 32
```

4. Update `.env` with your configuration:
```bash
JWT_SECRET=your-generated-secret-here
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Running the Server

```bash
# Run with default configuration (SQLite)
./dev/deserver start

# Or using make
make dev
```

The server will start on:
- gRPC: `localhost:3030`
- HTTP: `localhost:8081`

### Health Checks

```bash
# Liveness probe
curl http://localhost:8081/health

# Readiness probe (includes database and Redis health checks)
curl http://localhost:8081/readiness
```

Health checks verify:
- Server is running (`/health`)
- Database connectivity (`/readiness`)
- Redis connectivity (`/readiness`)

### Metrics

Prometheus metrics are available at: `http://localhost:8081/metrics` (when using Prometheus exporter)

## Database Migrations

GoDE uses golang-migrate for database schema management.

### Running Migrations

```bash
# Apply all pending migrations
./dev/deserver migrate up

# Check current migration version
./dev/deserver migrate version

# Rollback last migration
./dev/deserver migrate down -n 1
```

Migrations run automatically on server startup (except for memory stores).

### Migration Files

Located in `db/migrations/`:
- `000001_initial_schema.up.sql` - Create tables
- `000001_initial_schema.down.sql` - Drop tables

## Configuration

### Environment Variables

See `.env.example` for all available configuration options.

#### Required
- `JWT_SECRET` - JWT signing secret (min 32 characters)

#### Security
- `TLS_ENABLED` - Enable TLS/HTTPS (default: false)
- `TLS_CERT_FILE` - Path to TLS certificate
- `TLS_KEY_FILE` - Path to TLS private key

#### Rate Limiting
- `RATE_LIMIT_AUTH_PER_MINUTE` - Auth requests per IP (default: 5)
- `RATE_LIMIT_DE_PER_USER` - DE executions per user (default: 10)
- `RATE_LIMIT_MAX_CONCURRENT_DE` - Concurrent DEs per user (default: 3)

#### Database
- `STORE_TYPE` - Database type: sqlite, postgresql, memory (default: sqlite)
- `STORE_SQLITE_FILEPATH` - SQLite file path
- `STORE_POSTGRESQL_DNS` - PostgreSQL connection string

#### Redis (required for async execution)
- `REDIS_HOST` - Redis server host (default: localhost)
- `REDIS_PORT` - Redis server port (default: 6379)
- `REDIS_PASSWORD` - Redis password (default: empty)
- `REDIS_DB` - Redis database number (default: 0)

#### Async Executor
- `EXECUTOR_MAX_WORKERS` - Maximum concurrent workers (default: 10)
- `EXECUTOR_QUEUE_SIZE` - Execution queue size (default: 100)
- `EXECUTOR_EXECUTION_TTL` - Execution metadata TTL (default: 24h)
- `EXECUTOR_RESULT_TTL` - Result data TTL (default: 168h / 7 days)
- `EXECUTOR_PROGRESS_TTL` - Progress update TTL (default: 1h)

#### Observability
- `METRICS_ENABLED` - Enable metrics collection (default: true)
- `METRICS_TYPE` - Exporter type: prometheus, stdout (default: prometheus)
- `LOG_LEVEL` - Logging level: debug, info, warn, error
- `LOG_TYPE` - Log format: json, text

### Configuration File

You can also use a YAML configuration file:

```bash
./dev/deserver start --config=/path/to/config.yaml
```

## API

### Authentication

Register a new user:
```bash
curl -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user",
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

Login:
```bash
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user",
    "password": "securepassword"
  }'
```

### Async Execution API

The server provides async execution APIs that allow long-running optimizations to run in the background.

#### Submit Async Execution

```bash
curl -X POST http://localhost:8081/v1/de/async/run \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "algorithm": "gde3",
    "problem": "zdt1",
    "variant": "rand1",
    "de_config": {
      "population_size": 100,
      "dimensions_size": 30,
      "objectives_size": 2,
      "executions": 10,
      "generations": 100,
      "floor_limiter": 0.0,
      "ceil_limiter": 1.0,
      "gde3": {
        "cr": 0.5,
        "f": 0.5,
        "p": 0.1
      }
    }
  }'
```

Returns: `{"execution_id": "uuid-here"}`

#### Check Execution Status

```bash
curl http://localhost:8081/v1/de/executions/EXECUTION_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Execution Results

```bash
curl http://localhost:8081/v1/de/executions/EXECUTION_ID/results \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### List User's Executions

```bash
# All executions
curl http://localhost:8081/v1/de/executions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Filter by status
curl http://localhost:8081/v1/de/executions?status=EXECUTION_STATUS_RUNNING \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Cancel Running Execution

```bash
curl -X POST http://localhost:8081/v1/de/executions/EXECUTION_ID/cancel \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Delete Execution

```bash
curl -X DELETE http://localhost:8081/v1/de/executions/EXECUTION_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### CLI Async Commands

The CLI provides convenient commands for async execution:

```bash
# Submit and wait for completion (polls status)
./dev/decli de run --algorithm gde3 --variant rand1 --problem zdt1 \
  --generations 100 --population-size 100

# Submit and return immediately
./dev/decli de run-async --algorithm gde3 --variant rand1 --problem zdt1 \
  --generations 100 --population-size 100

# Check status
./dev/decli de status --execution-id EXECUTION_ID

# Stream real-time progress
./dev/decli de stream --execution-id EXECUTION_ID

# List executions
./dev/decli de list
./dev/decli de list --status running

# Get results
./dev/decli de results --execution-id EXECUTION_ID
./dev/decli de results --execution-id EXECUTION_ID --format json --output results.json

# Cancel execution
./dev/decli de cancel --execution-id EXECUTION_ID

# Delete execution
./dev/decli de delete --execution-id EXECUTION_ID
./dev/decli de delete --execution-id EXECUTION_ID --force  # Cancel first if running
```

## Development

### Project Structure

```
GoDE/
├── api/v1/              # Protocol buffer definitions
├── cmd/
│   ├── decli/          # CLI client
│   └── deserver/       # Server
├── db/migrations/      # Database migrations
├── internal/
│   ├── migrations/     # Migration management
│   ├── server/         # Server implementation
│   │   ├── auth/       # Authentication (JWT)
│   │   ├── handlers/   # gRPC handlers
│   │   └── middleware/ # Middleware (auth, rate limit, metrics, recovery)
│   ├── store/          # Database abstraction
│   ├── telemetry/      # Metrics and tracing
│   └── tenant/         # Multi-tenancy support
├── pkg/
│   ├── de/             # DE algorithm framework
│   ├── models/         # Data models
│   ├── problems/       # Optimization problems (ZDT, DTLZ, WFG)
│   ├── validation/     # Input validation
│   └── variants/       # DE mutation variants
└── CLAUDE.md           # Project documentation for AI

```

### Testing

```bash
# Run all tests
make test

# Run tests for a specific package
go test ./pkg/variants/...

# Run with coverage
go test -cover ./...
```

### Linting

```bash
make lint
```

### Protocol Buffers

```bash
# Regenerate proto files
make proto
```

## Production Deployment

### Prerequisites
1. PostgreSQL database (recommended over SQLite)
2. TLS certificates
3. Reverse proxy (nginx, Traefik) for HTTPS termination (optional)
4. Prometheus for metrics collection
5. Container orchestrator (Kubernetes, Docker Swarm) for high availability

### Deployment Checklist

- [ ] Generate strong JWT secret (32+ characters)
- [ ] Enable TLS (`TLS_ENABLED=true`)
- [ ] Configure PostgreSQL connection
- [ ] Set appropriate rate limits
- [ ] Configure structured logging (`LOG_TYPE=json`)
- [ ] Run database migrations (`deserver migrate up`)
- [ ] Set up Prometheus scraping
- [ ] Configure liveness probe: `GET /health`
- [ ] Configure readiness probe: `GET /readiness`
- [ ] Test graceful shutdown (SIGTERM handling)
- [ ] Set up log aggregation
- [ ] Configure alerting on metrics

### Docker

```dockerfile
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN make build

FROM debian:bookworm-slim
COPY --from=builder /app/.dev/deserver /usr/local/bin/
COPY db/migrations /app/db/migrations
WORKDIR /app
ENTRYPOINT ["deserver"]
CMD ["start"]
```

### Kubernetes Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gode-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gode-server
  template:
    metadata:
      labels:
        app: gode-server
    spec:
      containers:
      - name: gode-server
        image: gode-server:latest
        ports:
        - containerPort: 3030
          name: grpc
        - containerPort: 8081
          name: http
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: gode-secrets
              key: jwt-secret
        - name: STORE_TYPE
          value: "postgresql"
        - name: STORE_POSTGRESQL_DNS
          valueFrom:
            secretKeyRef:
              name: gode-secrets
              key: database-url
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

## Monitoring

### Prometheus Metrics

Key metrics to monitor:

- `api_requests_total` - Total API requests by method and status
- `api_request_duration_seconds` - Request duration histogram
- `api_requests_in_flight` - Current active requests
- `de_executions_total` - Total DE executions by algorithm/variant/problem
- `de_execution_duration_seconds` - DE execution duration histogram
- `de_executions_in_flight` - Currently running DE executions
- `auth_attempts_total` - Authentication attempts
- `auth_success_total` - Successful authentications
- `rate_limit_exceeded_total` - Rate limit violations
- `panics_total` - Recovered panics by location

### Example Prometheus Alerts

```yaml
groups:
- name: gode_alerts
  rules:
  - alert: HighErrorRate
    expr: rate(api_errors_total[5m]) > 0.05
    annotations:
      summary: "High API error rate"

  - alert: DEExecutionSlow
    expr: histogram_quantile(0.95, rate(de_execution_duration_seconds_bucket[5m])) > 300
    annotations:
      summary: "95th percentile DE execution time > 5 minutes"
```

## Troubleshooting

### Server won't start

**Error**: `JWT_SECRET environment variable is required`
- **Solution**: Set `JWT_SECRET` in `.env` or environment

**Error**: `database migration failed`
- **Solution**: Check database connectivity and run `deserver migrate version`

### Rate limiting issues

**Error**: `too many authentication attempts`
- **Solution**: Increase `RATE_LIMIT_AUTH_PER_MINUTE` or wait before retrying

**Error**: `maximum concurrent DE executions reached`
- **Solution**: Wait for existing executions to complete or increase `RATE_LIMIT_MAX_CONCURRENT_DE`

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linter (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

See [LICENSE](LICENSE) file for details.

## Acknowledgments

- Based on [GDE3](https://github.com/nicholaspcr/GDE3)
- Research from [CEFET-MG](https://www.cefetmg.br/)
- Built with Go, gRPC, and OpenTelemetry
