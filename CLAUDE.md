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
   - Abstract `Store` interface with implementations: memory, SQLite, PostgreSQL
   - GORM-based implementations in `internal/store/gorm/`
   - Entities: User, Pareto, Vector
   - Configuration via `store.Config` with type selection

### Key Design Patterns

- **Interface-driven design**: Problems, Variants, Algorithms, Handlers, and Store all use interfaces for extensibility
- **Channel-based concurrency**: DE executions communicate via channels for streaming results
- **Middleware architecture**: Server uses composable middleware for auth, logging, telemetry
- **Multi-tenancy support**: Tenant context propagation through `internal/tenant/`

## Protocol Buffers and API

API definitions are in `api/v1/*.proto`. Generated code goes to `pkg/api/v1/`.

Key services:
- `DifferentialEvolutionService`: List algorithms/variants/problems, run DE
- `AuthService`: User authentication
- `UserService`: User management
- `ParetoSetService`: Access Pareto results

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
- CLI state stored in SQLite at `.dev/cli/`

## Important Context

- **Vector modification**: Problem evaluation modifies vectors in-place for performance
- **Pareto filtering**: Uses non-dominated sorting with crowding distance reduction
- **Concurrent executions**: Multiple DE runs execute in goroutines, results aggregated through channels
- **OpenTelemetry**: Instrumentation is enabled for gRPC and GORM operations

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
