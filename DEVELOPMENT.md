# Development Guide

This guide explains how to develop, test, and run GoDE locally.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Configuration](#configuration)
- [Common Tasks](#common-tasks)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **Make** - Usually pre-installed on Linux/macOS
- **Node.js 20+** (optional, for frontend) - [Install Node.js](https://nodejs.org/)
- **Protocol Buffers** (optional, for proto changes) - [Install Buf](https://buf.build/docs/installation)

## Quick Start

### 1. Clone and Initialize

```bash
git clone <repository-url>
cd GoDE
make init deps
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL and Redis
make db-up
```

### 3. Run Server

```bash
# Run with development defaults
make run-dev
```

The server will be available at:
- **gRPC**: `localhost:3030`
- **HTTP API**: `http://localhost:8081`
- **Health Check**: `http://localhost:8081/health`

### 4. Test the API

```bash
# Register a user
curl -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"user":{"ids":{"username":"testuser"},"email":"test@example.com","password":"password123"}}'

# Login
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Save the token from the response
TOKEN="<your-access-token>"

# List supported algorithms
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/v1/de/algorithms
```

## Development Workflow

### Building

```bash
# Build both server and CLI
make build

# Build only server
make build-server

# Build only CLI
make build-cli

# Build with race detector (for debugging)
make build-race

# Install binaries to $GOPATH/bin
make install
```

Binaries are output to `.dev/` directory.

### Running the Server

```bash
# Quick start with defaults (recommended for development)
make run-dev

# Run from built binary (sources .dev/server/.env if exists)
make run

# Run with race detector
make run-race

# Watch for changes and auto-rebuild (requires entr)
make watch
```

### Using the CLI

Once the server is running:

```bash
# Register and login
./.dev/decli auth register --username myuser --email user@example.com --password mypass
./.dev/decli auth login --username myuser --password mypass

# List available options
./.dev/decli de list-algorithms
./.dev/decli de list-variants
./.dev/decli de list-problems

# Run DE execution (synchronous)
./.dev/decli de run \
  --algorithm gde3 \
  --problem zdt1 \
  --variant rand1 \
  --generations 50 \
  --population-size 100

# Run DE execution (asynchronous)
./.dev/decli de run-async \
  --algorithm gde3 \
  --problem zdt1 \
  --variant rand1 \
  --generations 100 \
  --population-size 100

# Check execution status
./.dev/decli de status --execution-id <id>

# Get results
./.dev/decli de results --execution-id <id>

# List your executions
./.dev/decli de list

# Cancel execution
./.dev/decli de cancel --execution-id <id>

# Delete execution
./.dev/decli de delete --execution-id <id>
```

## Testing

### Unit Tests

```bash
# Run all unit tests
make test

# Run unit tests only (excludes integration/e2e)
make test-unit

# Run with verbose output
make test-verbose

# Run with race detector
make test-race
```

### Integration Tests

```bash
# Run integration tests (requires databases)
make db-up
make test-integration
```

### End-to-End Tests

```bash
# Run E2E tests (uses testcontainers - requires Docker)
make test-e2e

# E2E tests include:
# - Full user workflow tests
# - Authentication failure scenarios
# - Invalid configuration handling
# - Async execution lifecycle
# - Progress streaming
# - Concurrent execution handling
# - Cross-user access control
```

### Test Coverage

```bash
# Run tests with coverage report
make test-coverage

# View coverage summary
make test-coverage-summary

# Coverage report is saved to .dev/coverage/coverage.html
```

### All Tests

```bash
# Run all tests (unit + integration + e2e)
make test-all

# Run full CI pipeline (checks + tests + coverage)
make ci
```

## Benchmarking

```bash
# Run all benchmarks
make bench

# Run variant-specific benchmarks
make bench-variants

# Run DE algorithm benchmarks
make bench-de

# Save benchmark results for comparison
make bench-compare
```

## Code Quality

### Formatting

```bash
# Format all Go code
make fmt

# Check if code is formatted (CI)
make fmt-check
```

### Linting

```bash
# Run go vet
make vet

# Run golangci-lint
make lint

# Auto-fix linting issues
make lint-fix

# Run all quality checks
make check
```

### Pre-commit Checks

```bash
# Run before committing
make pre-commit

# Run before pushing
make pre-push
```

## Configuration

GoDE uses a flexible multi-source configuration system powered by Viper.

### Configuration Priority (highest to lowest)

1. **Environment variables** (with `GODE_` prefix or legacy names)
2. **Config file** (`config.yaml`)
3. **Defaults**

### Using Config Files

```bash
# Copy example config
cp config.yaml.example config.yaml

# Edit config.yaml with your settings
vim config.yaml

# Server automatically searches for config.yaml in:
# - Current directory
# - ./config/
# - /etc/gode/
```

### Environment Variables

**Core Settings:**
```bash
JWT_SECRET="your-secret-min-32-chars"  # Required
STORE_TYPE="postgres"                   # postgres, sqlite, memory, redis
REDIS_HOST="localhost"
REDIS_PORT="6379"
GRPC_PORT="3030"
HTTP_PORT="8081"
```

**Observability:**
```bash
METRICS_ENABLED="true"
METRICS_TYPE="prometheus"              # prometheus, stdout
TRACING_ENABLED="true"
TRACING_EXPORTER="otlp"                # none, stdout, otlp, file
OTLP_ENDPOINT="localhost:4317"
TRACE_SAMPLE_RATIO="1.0"              # 0.0 to 1.0
SLO_ENABLED="true"
PPROF_ENABLED="true"
PPROF_PORT=":6060"
```

**Executor:**
```bash
GODE_EXECUTOR_MAX_WORKERS="10"
GODE_EXECUTOR_QUEUE_SIZE="100"
GODE_EXECUTOR_MAX_VECTORS_IN_PROGRESS="100"
```

**Rate Limiting:**
```bash
GODE_RATE_LIMIT_LOGIN_REQUESTS_PER_MINUTE="5"
GODE_RATE_LIMIT_DE_EXECUTIONS_PER_USER="10"
GODE_RATE_LIMIT_MAX_CONCURRENT_DE_PER_USER="3"
```

See `config.yaml.example` for complete configuration options.

### Database Configuration

**PostgreSQL:**
```bash
STORE_TYPE=postgres
STORE_POSTGRESQL_DNS="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

**SQLite:**
```bash
STORE_TYPE=sqlite
STORE_SQLITE_FILEPATH=".dev/server/sqlite.db"
```

**Redis (for async execution state):**
```bash
REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_PASSWORD=""  # Optional
REDIS_DB="0"
```

## Protocol Buffers

If you modify `.proto` files:

```bash
# Lint proto files
make proto-lint

# Generate Go code from protos
make proto-generate

# Full proto workflow (lint + clean + generate)
make proto

# Generate OpenAPI spec
make openapi
```

## Database Management

```bash
# Start databases
make db-up

# Stop databases
make db-down

# Clean databases (removes volumes)
make db-clean

# View logs
make db-logs

# Connect to PostgreSQL
make db-psql

# Connect to Redis
make db-redis
```

## Frontend Development

```bash
# Install dependencies
make web-deps

# Start dev server (with hot reload)
make web-dev
# Frontend: http://localhost:5173

# Build for production
make web-build

# Run tests
make web-test

# Lint
make web-lint

# Format code
make web-format

# Regenerate API client (after proto changes)
make web-api
```

### Full Stack Development

Terminal 1:
```bash
make db-up
make run-dev
```

Terminal 2:
```bash
make web-dev
```

Access frontend at http://localhost:5173

## Distributed Tracing with Jaeger

### Start Jaeger

```bash
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/jaeger:latest
```

### Configure Server for Jaeger

```bash
TRACING_EXPORTER=otlp \
OTLP_ENDPOINT=localhost:4317 \
make run-dev
```

### View Traces

Open http://localhost:16686 and select "deserver" service.

### Stop Jaeger

```bash
docker stop jaeger && docker rm jaeger
```

## Kubernetes Deployment

### Minikube

```bash
# Build image in minikube
make k8s-build

# Deploy all resources
make k8s-deploy

# Get application URL
make k8s-url

# View logs
make k8s-logs

# Check status
make k8s-status

# Delete resources
make k8s-delete
```

## Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run

# Docker Compose (full stack)
make docker-compose-up
make docker-compose-down
make docker-compose-logs
```

## Common Tasks

### Add a New Problem

1. Create file in `pkg/problems/multi/` or `pkg/problems/many/`
2. Implement `problems.Interface`
3. Register in `init()` function
4. Add tests
5. Update proto if needed

### Add a New Variant

1. Create file in `pkg/variants/<category>/`
2. Implement `variants.Interface`
3. Register in `init()` function
4. Add tests and benchmarks
5. Update validation rules

### Run Specific Tests

```bash
# Run tests in specific package
go test -v ./pkg/de/...

# Run specific test
go test -v -run TestExecutor_SubmitExecution ./internal/executor/...

# Run tests matching pattern
go test -v -run 'TestE2E.*' -tags=e2e ./test/e2e/...
```

### Debugging

```bash
# Run with race detector
make run-race

# Enable pprof
PPROF_ENABLED=true PPROF_PORT=:6060 make run-dev

# View profiles at http://localhost:6060/debug/pprof/

# Enable verbose logging
LOG_LEVEL=debug make run-dev

# Enable file tracing
TRACING_EXPORTER=file TRACE_FILE_PATH=traces.json make run-dev
```

### Clean Everything

```bash
# Clean build artifacts
make clean

# Clean everything including databases
make clean-all
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :3030
lsof -i :8081

# Kill process
pkill -f deserver
```

### Database Connection Failed

```bash
# Verify databases are running
make db-up
docker compose -f docker-compose.test.yml ps

# Check PostgreSQL
docker compose -f docker-compose.test.yml exec postgres psql -U gode -d gode_test -c '\l'

# Check Redis
docker compose -f docker-compose.test.yml exec redis redis-cli ping
```

### Redis Connection Issues

```bash
# Flush Redis cache
docker compose -f docker-compose.test.yml exec redis redis-cli FLUSHALL

# Restart Redis
docker compose -f docker-compose.test.yml restart redis
```

### Build Failures

```bash
# Update dependencies
make deps-update

# Tidy modules
make tidy

# Clean and rebuild
make clean build
```

### Test Failures

```bash
# Run with verbose output
go test -v ./path/to/failing/test

# Run single test
go test -v -run TestName ./path/to/test

# Check for race conditions
make test-race
```

### Protobuf Issues

```bash
# Regenerate protos
make proto

# Check proto lint
make proto-lint

# Verify buf is installed
buf --version
```

### Docker Issues

```bash
# Clean Docker resources
docker system prune -a

# Restart Docker daemon

# Check Docker logs
docker compose -f docker-compose.test.yml logs
```

## Project Structure

```
GoDE/
├── api/                    # Protocol buffer definitions
├── cmd/
│   ├── decli/             # CLI application
│   └── deserver/          # Server application
├── internal/              # Private application code
│   ├── cache/             # Redis cache implementation
│   ├── executor/          # Background DE execution
│   │   ├── executor.go
│   │   ├── worker_pool.go      # NEW: Worker concurrency
│   │   └── progress_tracker.go # NEW: Progress management
│   ├── server/
│   │   ├── handlers/
│   │   │   ├── differential_evolution.go  # Core handler
│   │   │   ├── de_async.go               # NEW: Async execution
│   │   │   ├── de_progress.go            # NEW: Progress streaming
│   │   │   ├── de_results.go             # NEW: Result retrieval
│   │   │   ├── de_management.go          # NEW: Lifecycle management
│   │   │   ├── de_list.go                # NEW: List operations
│   │   │   └── de_conversions.go         # NEW: Data conversions
│   │   └── config.go         # NEW: Viper-based configuration
│   ├── store/
│   │   └── redis/
│   │       ├── execution.go            # Core CRUD
│   │       ├── execution_list.go       # NEW: Query operations
│   │       └── execution_lifecycle.go  # NEW: Progress & state
│   ├── telemetry/         # Observability (metrics, tracing, SLO)
│   └── migrations/        # Database migrations
├── pkg/                   # Public libraries
│   ├── api/              # Generated protobuf code
│   ├── de/               # DE algorithm core (with tracing)
│   ├── problems/         # Optimization problems
│   ├── variants/         # DE mutation variants
│   ├── models/           # Domain models
│   └── validation/       # Input validation
├── test/
│   └── e2e/              # End-to-end tests
│       ├── e2e_test.go           # Happy path tests
│       └── e2e_failures_test.go  # NEW: Failure scenarios
├── web/                  # Frontend (React + TypeScript)
├── k8s/                  # Kubernetes manifests
├── monitoring/           # Grafana dashboards
├── Makefile              # NEW: Comprehensive build system
├── config.yaml.example   # NEW: Configuration template
└── CLAUDE.md             # AI development guide
```

## Recent Improvements

### Configuration Management
- Migrated to Viper with multi-source support (env vars + config files + defaults)
- Added `config.yaml.example` with comprehensive documentation
- Backward compatible with all existing environment variables
- Support for `GODE_` prefixed environment variables

### Code Organization
- Split monolithic files into focused modules:
  - `executor.go`: 543→441 lines (extracted WorkerPool + ProgressTracker)
  - `differential_evolution.go`: 607→49 lines (6 new focused files)
  - Redis store: 519→282 lines (3 new focused files)

### Testing
- Added 26 E2E failure scenarios (763 lines)
- Expanded migration tests (8 functions, 31 sub-tests)
- Added server integration tests (17 sub-tests)
- Comprehensive coverage for edge cases

### Observability
- Added OpenTelemetry tracing to DE algorithm core
- Instrumented 6 critical functions with span attributes
- Support for OTLP, file, and stdout exporters

### Build System
- Comprehensive Makefile with 70+ targets
- Organized sections: build, test, quality, database, k8s, frontend
- Color-coded output for better UX
- Parallel execution support

## Getting Help

```bash
# Show all make targets with descriptions
make help

# Show project information
make info

# Show dependency graph
make deps-graph
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run quality checks: `make pre-commit`
5. Run tests: `make test-all`
6. Submit a pull request

## References

- [CLAUDE.md](CLAUDE.md) - AI development guide with architecture details
- [config.yaml.example](config.yaml.example) - Complete configuration reference
- [Protocol Buffers](https://buf.build/docs/) - API definition language
- [gRPC](https://grpc.io/) - RPC framework
- [OpenTelemetry](https://opentelemetry.io/) - Observability framework
