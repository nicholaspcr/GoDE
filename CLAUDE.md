# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoDE is a Differential Evolution (DE) framework that implements multi-objective optimization algorithms. The project extends [GDE3](https://github.com/nicholaspcr/GDE3) and is architected as both a gRPC/HTTP server and a CLI client, enabling concurrent execution of DE instances across multiple users.

## Build and Development Commands

```bash
# Build both binaries (decli and deserver)
make build

# Run tests
make test

# Run a specific test
go test -v ./pkg/variants/...
go test -run TestGenerateIndices ./pkg/variants/...

# Lint the codebase
make lint

# Protocol buffers (uses remote plugins from Buf Schema Registry)
make proto          # Lint, remove old files, and regenerate
make proto-generate # Generate only
make proto-lint     # Lint only

# Development environment (Docker)
make dev

# Clean build artifacts
make clean
```

**Output locations:** Binaries are built to `./.dev/decli` and `./.dev/deserver`

## Architecture

### Core Components

1. **DE Algorithm Framework** (`pkg/de/`)
   - `Algorithm` interface defines the execution contract for DE algorithms
   - Algorithms execute concurrently across multiple runs (controlled by `Executions` parameter)
   - Results are streamed through channels (`pareto`, `maxObjectives`)
   - The framework filters results using non-dominated sorting and crowding distance

2. **GDE3 Implementation** (`pkg/de/gde3/`)
   - Main multi-objective DE algorithm implementation
   - Configuration: `Constants` struct with CR (crossover rate), F (scaling factor), and P (selection parameter)
   - Implements the `Algorithm` interface

3. **Problems** (`pkg/problems/`)
   - `Interface` defines the contract: `Evaluate(*models.Vector, int) error`
   - Implementations in `multi/` (ZDT, VNT problems) and `many/` (DTLZ, WFG problems)
   - Problems modify the `Objectives` slice of the vector in-place

4. **Variants** (`pkg/variants/`)
   - `Interface` defines mutation strategies: `Mutate(elems, rankZero []Vector, params Parameters) (Vector, error)`
   - Subdirectories contain specific variants: `best/`, `current-to-best/`, `pbest/`, `rand/`
   - `GenerateIndices` utility ensures unique random indices for mutation

5. **Server** (`cmd/deserver/`, `internal/server/`)
   - gRPC server with HTTP gateway (grpc-gateway)
   - Handlers in `internal/server/handlers/` implement the `Handler` interface
   - Session management and authentication middleware
   - OpenTelemetry instrumentation for observability
   - Default ports: gRPC `:3030`, HTTP `:8081`

6. **CLI** (`cmd/decli/`)
   - Client for interacting with the server
   - State management using SQLite (`internal/state/sqlite/`)
   - Configuration in `.dev/cli/` and `.env/cli/`

7. **Storage** (`internal/store/`)
   - Abstract `Store` interface with implementations: memory, SQLite, PostgreSQL, Redis, composite
   - GORM-based implementations in `internal/store/gorm/`
   - Entities: User, Pareto, Vector, Execution
   - Configuration via `store.Config` with type selection
   - Database migrations in `internal/store/migrations/` (2 migrations: 000001_initial_schema, 000002_add_executions_and_indices)
   - Batch operations for performance (CreateInBatches for large Pareto sets)

8. **Async Execution Engine** (`internal/executor/`)
   - Background execution of DE algorithms with worker pool pattern
   - Semaphore-based concurrency control (configurable max workers)
   - Real-time progress tracking via callbacks
   - Execution state: PENDING → RUNNING → COMPLETED/FAILED/CANCELLED
   - Context-based cancellation support
   - TTL configuration: Execution (24h), Results (7d), Progress (1h)

9. **Observability & Monitoring** (`internal/telemetry/`, `internal/slo/`)
   - OpenTelemetry integration: distributed tracing with OTLP exporter
   - Prometheus metrics: request counts, durations, error rates, rate limits
   - Service Level Objective (SLO) tracking: latency (p50/p95/p99), error rate, availability
   - Structured logging with `log/slog`
   - Metrics endpoint: `/metrics`

10. **Caching & Performance** (`internal/cache/redis/`)
   - Redis client with circuit breaker (uses `github.com/sony/gobreaker`)
   - Execution metadata caching
   - Progress tracking with TTL
   - Graceful degradation when Redis unavailable

### Key Design Patterns

- **Interface-driven design**: Problems, Variants, Algorithms, Handlers, and Store all use interfaces for extensibility
- **Channel-based concurrency**: DE executions communicate via channels for streaming results
- **Worker pool pattern**: Async executor uses semaphore-based concurrency limiting
- **Registry pattern**: Problems and variants auto-register via blank imports
- **Middleware architecture**: Server uses composable middleware stack (panic recovery, tracing, metrics, logging, CORS, auth, session)
- **Circuit breaker**: Redis operations protected against cascade failures
- **Multi-tenancy support**: Tenant context propagation through `internal/tenant/`
- **Factory pattern**: Store creation via `storefactory` based on configuration

## Protocol Buffers and API

API definitions are in `api/v1/*.proto` (8 proto files). Generated code goes to `pkg/api/v1/`.

### Key Services

**DifferentialEvolutionService** (`differential_evolution.proto`):
- List supported: `ListSupportedAlgorithms`, `ListSupportedVariants`, `ListSupportedProblems`
- Async execution: `RunAsync`, `StreamProgress`, `GetExecutionStatus`, `GetExecutionResults`
- Management: `ListExecutions`, `CancelExecution`, `DeleteExecution`
- HTTP endpoints: `/v1/de/run` (POST), `/v1/de/executions` (GET), `/v1/de/executions/{id}` (GET/DELETE)

**AuthService** (`auth.proto`):
- Authentication: `Register`, `Login`, `Logout`, `RefreshToken`
- JWT-based with separate access and refresh tokens
- HTTP endpoints: `/v1/auth/login`, `/v1/auth/register`, `/v1/auth/refresh`

**UserService** (`user.proto`):
- User management: `Create`, `Get`, `Update`, `Delete`
- Uses `UserResponse` message in responses (excludes password for security)

**ParetoService** (`pareto_set.proto`):
- Access Pareto results: `Get`, `Delete`, `ListByUser`

### API Data Types
- `definitions.proto` - Core types: Vector, Pareto, PopulationParameters
- `differential_evolution_config.proto` - DEConfig, GDE3Config
- `errors.proto` - ErrorDetail, FieldViolation for structured error responses
- `tenant.proto` - Multi-tenancy support

## Testing Conventions

- Test files follow Go conventions: `*_test.go`
- Use `github.com/stretchr/testify/assert` for assertions
- Table-driven tests are preferred (see `pkg/variants/utils_test.go`)
- Use `rand.New(rand.NewSource(1))` for deterministic random testing

## Error Handling

- Use sentinel errors from `pkg/variants/errors.go` and `internal/store/errors/`
- Return wrapped errors with context using `fmt.Errorf` or `errors.Join`
- Log errors using `log/slog` with structured context

## Configuration

- Server config: `internal/server/config.go` - default ports and DE limiters
- Store config: `internal/store/config.go` - database type and connection details
- DE config: `pkg/de/config.go` - channel limiters and result limits
- Executor config: `internal/server/config.go` - ExecutorConfig with MaxWorkers (10), MaxVectorsInProgress (100), TTLs (execution: 24h, results: 7d, progress: 1h)
- CLI state stored in SQLite at `.dev/cli/`

## Important Context

- **Vector modification**: Problem evaluation modifies vectors in-place for performance
- **Pareto filtering**: Uses non-dominated sorting with crowding distance reduction
- **Concurrent executions**: Multiple DE runs execute in goroutines, results aggregated through channels
- **OpenTelemetry**: Instrumentation is enabled for gRPC and GORM operations
- **Field naming**: Use `ObjectivesSize` (not `ObjetivesSize`) - typo was fixed in API and codebase
- **Password security**: User responses exclude password field via `UserResponse` message
- **Database compatibility**: Migrations work with both PostgreSQL and SQLite (avoid PostgreSQL-specific syntax)
- **Async execution**: All long-running DE operations use async pattern (RunAsync → poll status → get results)

## Database Migrations

Location: `internal/store/migrations/`

**Migration files:**
1. `000001_initial_schema.{up,down}.sql` - Users, Pareto sets, Vectors tables
2. `000002_add_executions_and_indices.{up,down}.sql` - Executions table, indices for query performance

**Running migrations:**
```bash
# Automatic on server startup
./deserver

# Manual migration commands
deserver migrate up      # Apply all pending migrations
deserver migrate down 1  # Rollback last migration
deserver migrate version # Check current version
```

**Database compatibility notes:**
- Migrations must work with both PostgreSQL and SQLite
- Avoid PostgreSQL-specific syntax: `DO $$` blocks, `ALTER TABLE ADD CONSTRAINT` (use inline `REFERENCES` instead)
- Use `REFERENCES` in CREATE TABLE for foreign keys, not separate ALTER TABLE statements

## Deployment

**Kubernetes:**
- Manifests in `k8s/` directory
- Includes: Deployment, Service, HPA, PostgreSQL StatefulSet, Redis, ConfigMaps, Secrets
- Minikube commands: `make k8s-build`, `make k8s-deploy`, `make k8s-logs`, `make k8s-status`

**Health checks:**
- Liveness probe: `/health` (simple ping)
- Readiness probe: `/readiness` (DB + Redis connectivity)

**Monitoring:**
- Grafana dashboards in `monitoring/grafana/`
- Prometheus metrics at `/metrics`

## Troubleshooting

### Common Issues

**1. Migration errors with SQLite**
- Symptom: `syntax error` during migration
- Fix: Ensure migrations avoid PostgreSQL-specific syntax (no `DO $$` blocks, use inline `REFERENCES`)

**3. Compilation errors with `ObjetivesSize`**
- Symptom: `undefined: config.ObjetivesSize`
- Fix: Field was renamed to `ObjectivesSize` (correct spelling) - regenerate protobuf with `make proto-generate`

**4. Password exposed in API responses**
- Symptom: User GET endpoint returns password field
- Fix: Use `UserResponse` message (excludes password), not `User` message in responses

**5. Redis connection failures**
- Symptom: Execution progress not updating
- Behavior: Circuit breaker provides graceful degradation
- Fix: Check Redis connectivity, circuit breaker prevents cascade failures

### Performance Tips

- Use batch operations for large Pareto sets (automatic via CreateInBatches)
- Monitor SLO metrics for latency violations
- Check worker pool configuration (maxWorkers) if executions queue up
- Review TTL settings for execution metadata (default 24h)

## Commit Message Guidelines

When creating commits for this repository:

1. **Keep commits compact and focused**: Each commit should address a single concern or feature
2. **Use conventional commit format**:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `test:` for test additions/changes
   - `refactor:` for code refactoring
   - `chore:` for maintenance tasks
   - `perf:` for performance improvements

3. **Do NOT include AI attribution**:
   - Do not add "Generated with Claude Code" or similar AI references
   - Do not add "Co-Authored-By: Claude" lines
   - The presence of CLAUDE.md in the repository already indicates AI assistance is used
   - Commit messages should focus on the technical changes, not the tools used to create them

4. **Write clear, descriptive messages**:
   - First line: concise summary (50-72 chars)
   - Optional body: detailed explanation of what and why (not how)
   - Use bullet points for multiple changes
   - Reference issue numbers when applicable

**Example of a good commit message:**
```
fix: correct best/1 variant mutation index calculation

- Changed index[i] to index[1] in best_1.go line 36
- Bug caused index out of bounds when DIM >= 3
- Add comprehensive unit tests for best/1 and best/2 variants
- All 77 tests now passing
```
