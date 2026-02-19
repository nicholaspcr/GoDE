# CLAUDE.md

This file provides architectural guidance and technical context for Claude Code when working with this repository.

> **For development workflows, build commands, and testing:** See [DEVELOPMENT.md](DEVELOPMENT.md)
> **For project structure:** See [DEVELOPMENT.md#project-structure](DEVELOPMENT.md#project-structure)

## Project Overview

GoDE is a Differential Evolution (DE) framework implementing multi-objective optimization algorithms. The project extends [GDE3](https://github.com/nicholaspcr/GDE3) and is architected as a multi-tenant gRPC/HTTP server with CLI client, enabling concurrent execution of DE instances across multiple users.

**Key Technologies:**
- Go 1.21+ with generics and context propagation
- gRPC + grpc-gateway for dual API exposure
- Protocol Buffers for API definitions
- PostgreSQL/SQLite/Redis for storage
- OpenTelemetry for distributed tracing
- Prometheus for metrics
- Viper for configuration management

## Quick Reference

```bash
# Build
make build              # Build server and CLI
make build-race         # Build with race detector

# Test
make test               # Unit tests
make test-e2e           # E2E tests (requires Docker)
make test-coverage      # Coverage report
make bench              # Run benchmarks

# Quality
make check              # All quality checks
make lint               # Lint code
make fmt                # Format code

# Development
make run-dev            # Start server with dev config
make db-up              # Start PostgreSQL + Redis
make watch              # Watch and rebuild

# Full reference: make help
```

## Architecture

### System Design

```
┌─────────────┐         ┌──────────────┐
│   decli     │◄───────►│  deserver    │
│  (Client)   │  gRPC   │   (Server)   │
└─────────────┘         └──────┬───────┘
                               │
                 ┌─────────────┼────────────┐
                 │             │            │
           ┌─────▼────┐   ┌────▼────┐  ┌────▼─────┐
           │PostgreSQL│   │  Redis  │  │ Executor │
           │  Store   │   │  Cache  │  │  Pool    │
           └──────────┘   └─────────┘  └──────────┘
```

### Core Components

#### 1. DE Algorithm Framework (`pkg/de/`)

The heart of the optimization engine:

- **`Algorithm` interface**: Defines execution contract for all DE algorithms
- **Concurrent execution**: Multiple runs execute in parallel (controlled by `Executions`)
- **Channel-based streaming**: Results flow through `pareto` and `maxObjectives` channels
- **Non-dominated sorting**: Filters results using NSGA-II-style sorting + crowding distance
- **OpenTelemetry tracing**: Spans instrument critical algorithm paths

**Key files:**
- `de.go`: Core execution engine
- `utils.go`: NSGA-II utilities (ReduceByCrowdDistance, FastNonDominatedRanking)
- `config.go`: Channel buffer configuration

#### 2. GDE3 Implementation (`pkg/de/gde3/`)

Multi-objective DE algorithm with Pareto-based selection:

- **Constants**: CR (crossover rate), F (scaling factor), P (selection parameter)
- **Generational evolution**: Population evolves over N generations
- **Variant-based mutation**: Pluggable mutation strategies
- **Progress callbacks**: Real-time generation progress reporting

#### 3. Problems (`pkg/problems/`)

Optimization problem definitions with benchmark functions:

- **Interface**: `Evaluate(*models.Vector, int) error`
- **In-place evaluation**: Modifies `Objectives` slice directly for performance
- **Multi-objective**: `multi/` (ZDT, VNT - 2 objectives)
- **Many-objective**: `many/` (DTLZ, WFG - 3+ objectives)
- **Auto-registration**: Problems register via `init()` in blank imports

#### 4. Variants (`pkg/variants/`)

DE mutation strategies:

- **Interface**: `Mutate(elems, rankZero []Vector, params Parameters) (Vector, error)`
- **Categories**: `best/`, `current-to-best/`, `pbest/`, `rand/`
- **Index generation**: `GenerateIndices` ensures unique random indices
- **Metadata**: Each variant provides name, description, min population requirements
- **Validation**: Population size validation prevents index out-of-bounds

**Common variants:**
- `rand/1`: `v = r1 + F * (r2 - r3)`
- `best/1`: `v = best + F * (r1 - r2)`
- `current-to-best/1`: `v = current + F * (best - current) + F * (r1 - r2)`

#### 5. Server (`internal/server/`)

Multi-tenant gRPC server with HTTP gateway:

**Architecture:**
- **Middleware stack**: Panic recovery → Tracing → Metrics → Logging → CORS → Auth → Session
- **Handler pattern**: Each service has focused handler files
- **OpenTelemetry**: Automatic span creation for all gRPC calls
- **Multi-tenancy**: Tenant context propagates through request chain

**Handler organization** (`internal/server/handlers/`):
- `differential_evolution.go`: Core handler infrastructure
- `de_async.go`: Async execution submission
- `de_progress.go`: Real-time progress streaming
- `de_results.go`: Status and result retrieval
- `de_management.go`: Execution lifecycle (list, cancel, delete)
- `de_list.go`: Registry queries (algorithms, variants, problems)
- `de_conversions.go`: Proto/store data transformations
- `auth.go`: Authentication (register, login, JWT)
- `user.go`: User management
- `pareto_set.go`: Pareto result access

#### 6. Async Execution Engine (`internal/executor/`)

Background DE execution with worker pool pattern:

**Structure:**
- `executor.go`: Main executor with lifecycle management
- `worker_pool.go`: Semaphore-based concurrency control
- `progress_tracker.go`: Progress callback management

**Features:**
- **Worker pool**: Configurable max workers (default: 10)
- **Queue management**: Blocks when pool exhausted
- **Progress tracking**: Atomic counters for completion status
- **State management**: PENDING → RUNNING → COMPLETED/FAILED/CANCELLED
- **Context cancellation**: Graceful shutdown and execution cancellation
- **TTL configuration**: Execution (24h), Results (7d), Progress (1h)
- **Metrics**: Queue wait time, active workers, utilization percentage

#### 7. Storage (`internal/store/`)

Abstract storage with multiple backends:

**Implementations:**
- `gorm/`: PostgreSQL/SQLite via GORM
- `redis/`: Redis for execution state and progress
- `memory/`: In-memory for testing
- `composite/`: Combines multiple stores

**Redis store organization** (`internal/store/redis/`):
- `execution.go`: Core CRUD operations
- `execution_list.go`: Query and pagination (HSCAN-based)
- `execution_lifecycle.go`: Progress, cancellation, pub/sub

**Entities:**
- User: Authentication and authorization
- Execution: DE run metadata and state
- ParetoSet: Optimization results
- Vector: Individual solutions

**Migrations:**
- Location: `internal/store/migrations/`
- 7 migrations total (initial schema through execution metadata)
- Compatible with both PostgreSQL and SQLite

#### 8. Configuration (`internal/server/config.go`)

Viper-based multi-source configuration:

**Priority (highest to lowest):**
1. Environment variables (GODE_ prefix or legacy names)
2. Config file (`config.yaml`)
3. Defaults

**Config categories:**
- Server (ports, JWT settings)
- Observability (metrics, tracing, SLO)
- Rate limiting (per-IP, per-user)
- Redis (host, port, DB)
- Executor (workers, TTLs)
- Database (type, connection)

**See:** `config.yaml.example` for complete reference

#### 9. Observability (`internal/telemetry/`, `internal/slo/`)

Comprehensive observability stack:

**OpenTelemetry:**
- Distributed tracing with OTLP exporter
- Automatic span creation for gRPC, GORM, algorithm operations
- Span attributes for request context
- Support for Jaeger, file, stdout exporters

**Prometheus Metrics:**
- Request counts, durations, error rates
- Rate limit violations
- Worker pool utilization
- Execution queue wait times
- Active workers and executions

**SLO Tracking:**
- Latency percentiles (p50, p95, p99)
- Error rate monitoring
- Availability tracking

**Endpoints:**
- `/metrics`: Prometheus metrics
- `/health`: Liveness probe
- `/readiness`: Readiness probe (checks DB + Redis)

#### 10. Caching (`internal/cache/redis/`)

Redis client with circuit breaker pattern:

- **Circuit breaker**: `github.com/sony/gobreaker` prevents cascade failures
- **Graceful degradation**: Continues operation when Redis unavailable
- **Use cases**: Execution metadata, progress updates, cancellation flags
- **TTL management**: Automatic expiration for transient data

### Design Patterns

| Pattern | Usage | Location |
|---------|-------|----------|
| **Interface-driven** | Extensibility for problems, variants, algorithms, stores | `pkg/problems/`, `pkg/variants/`, `pkg/de/`, `internal/store/` |
| **Registry** | Auto-registration via blank imports | `pkg/problems/`, `pkg/variants/` |
| **Worker pool** | Semaphore-based concurrency limiting | `internal/executor/worker_pool.go` |
| **Channel-based** | Streaming results from concurrent executions | `pkg/de/de.go` |
| **Middleware chain** | Composable request processing | `internal/server/middleware/` |
| **Circuit breaker** | Redis fault tolerance | `internal/cache/redis/` |
| **Factory** | Store creation based on config | `internal/storefactory/` |
| **Repository** | Data access abstraction | `internal/store/` |

## Protocol Buffers and API

### Services

**DifferentialEvolutionService** (`differential_evolution.proto`):
```protobuf
// List registries
rpc ListSupportedAlgorithms()
rpc ListSupportedVariants()
rpc ListSupportedProblems()

// Async execution
rpc RunAsync(RunAsyncRequest) returns (RunAsyncResponse)
rpc StreamProgress(StreamProgressRequest) returns (stream StreamProgressResponse)
rpc GetExecutionStatus(GetExecutionStatusRequest) returns (GetExecutionStatusResponse)
rpc GetExecutionResults(GetExecutionResultsRequest) returns (GetExecutionResultsResponse)

// Management
rpc ListExecutions(ListExecutionsRequest) returns (ListExecutionsResponse)
rpc CancelExecution(CancelExecutionRequest) returns (Empty)
rpc DeleteExecution(DeleteExecutionRequest) returns (Empty)
```

**HTTP Mapping:**
- `POST /v1/de/run` → RunAsync
- `GET /v1/de/executions/{execution_id}/progress` → StreamProgress (SSE)
- `GET /v1/de/executions/{execution_id}` → GetExecutionStatus
- `GET /v1/de/executions/{execution_id}/results` → GetExecutionResults
- `GET /v1/de/executions` → ListExecutions
- `POST /v1/de/executions/{execution_id}/cancel` → CancelExecution
- `DELETE /v1/de/executions/{execution_id}` → DeleteExecution

**AuthService** (`auth.proto`):
- JWT-based authentication
- Access + refresh token pattern
- Bcrypt password hashing

**UserService** (`user.proto`):
- User CRUD operations
- `UserResponse` excludes password field

**ParetoService** (`pareto_set.proto`):
- Pareto result retrieval
- Per-user isolation

### Proto File Organization

```
api/v1/
├── definitions.proto          # Core types (Vector, Pareto)
├── differential_evolution.proto         # Main DE service
├── differential_evolution_config.proto  # DE configuration
├── auth.proto                 # Authentication
├── user.proto                 # User management
├── pareto_set.proto          # Results access
├── errors.proto              # Error details
└── tenant.proto              # Multi-tenancy
```

## Testing Strategy

### Test Types

**Unit Tests** (`*_test.go`):
- Test individual functions and methods
- Use table-driven tests
- Mock external dependencies
- Deterministic randomness: `rand.New(rand.NewSource(1))`

**Integration Tests** (`-tags=integration`):
- Test component interactions
- Real database connections
- Full server startup

**E2E Tests** (`test/e2e/`, `-tags=e2e`):
- Full user workflows via gRPC/HTTP
- Uses testcontainers (PostgreSQL + Redis + deserver)
- Happy paths + failure scenarios (26 scenarios)
- Authentication, authorization, concurrent executions

**Benchmarks** (`*_test.go` with `Benchmark*`):
- Variant performance comparison
- Algorithm execution time
- Channel throughput

### Running Tests

```bash
make test              # Unit tests only
make test-integration  # Integration tests
make test-e2e          # E2E tests (requires Docker)
make test-all          # All tests
make test-coverage     # Coverage report
make bench             # All benchmarks
```

## Error Handling

**Patterns:**
- Sentinel errors in `pkg/variants/errors.go`, `internal/store/errors/`
- Wrapped errors with context: `fmt.Errorf("context: %w", err)`
- Structured logging: `log/slog` with attributes
- gRPC status codes: Use appropriate codes (NotFound, InvalidArgument, etc.)

**Example:**
```go
if execution.Status != store.ExecutionStatusCompleted {
    return nil, status.Error(codes.FailedPrecondition, "execution is not completed")
}
```

## Important Implementation Details

### Performance Considerations

1. **In-place vector evaluation**: Problems modify `Objectives` slice directly
2. **Batch operations**: Use `CreateInBatches` for large Pareto sets
3. **Channel buffering**: Configure channel limiters to prevent goroutine leaks
4. **Worker pool sizing**: Balance concurrency vs memory pressure
5. **Redis circuit breaker**: Prevents cascade failures under Redis outage

### Data Flow

**Async execution:**
```
Client → RunAsync → Executor.SubmitExecution → Worker Pool
                                                     ↓
                    Progress Tracker ← Algorithm ← Worker
                           ↓
                    Redis pub/sub → StreamProgress → Client
```

**Result retrieval:**
```
Client → GetExecutionStatus → Redis (metadata)
                            → PostgreSQL (state)

Client → GetExecutionResults → PostgreSQL (ParetoSet + Vectors)
```

### Multi-Tenancy

- User ID propagates through context
- Store methods accept userID parameter
- Cross-user access blocked at store layer
- Execution ownership verified on all operations

### Database Compatibility

**Migrations:**
- Must work with both PostgreSQL and SQLite
- Avoid PostgreSQL-specific syntax (`DO $$` blocks)
- Use inline `REFERENCES` for foreign keys
- Test migrations on both databases

### Security

- JWT secrets: Min 32 characters (enforced in validation)
- Password hashing: Bcrypt
- Rate limiting: Per-IP (login/register) and per-user (DE executions)
- Authorization: Scope-based with middleware enforcement
- Input validation: All requests validated before processing

## Troubleshooting

### Common Issues

**Migration errors:**
- Ensure migrations avoid PostgreSQL-specific syntax
- Use inline REFERENCES, not separate ALTER TABLE

**Field naming:**
- Use `ObjectivesSize` (not `ObjetivesSize`) - typo was fixed

**Redis failures:**
- Circuit breaker provides graceful degradation
- Check Redis connectivity
- Flush cache: `docker compose exec redis redis-cli FLUSHALL`

**Password in responses:**
- Use `UserResponse` (excludes password), not `User`

### Performance Tuning

- Monitor worker pool utilization via metrics
- Adjust `maxWorkers` if executions queue up
- Review TTL settings for execution metadata
- Check SLO metrics for latency violations
- Use batch operations for large result sets

## Development Workflow

1. **Setup**: `make init deps db-up`
2. **Development**: `make run-dev` (or `make watch`)
3. **Testing**: `make test`, `make test-e2e`
4. **Quality**: `make check` (fmt + vet + lint)
5. **Benchmarking**: `make bench`
6. **Proto changes**: `make proto`

**Pre-commit:** `make pre-commit` (fmt + vet + lint + test)

## References

- [DEVELOPMENT.md](DEVELOPMENT.md) - Complete development guide
- [config.yaml.example](config.yaml.example) - Configuration reference
- [Protocol Buffers](https://buf.build/docs/) - API definition
- [gRPC](https://grpc.io/) - RPC framework
- [OpenTelemetry](https://opentelemetry.io/) - Observability
- [NSGA-II](https://ieeexplore.ieee.org/document/996017) - Multi-objective optimization algorithm
