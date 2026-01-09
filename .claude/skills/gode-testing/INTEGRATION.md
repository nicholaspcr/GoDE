# Integration Testing

Run integration tests across different storage backends.

## Usage

```bash
.claude/skills/gode-testing/scripts/test-integration.sh [backend]
```

## Available Backends

| Backend | Description | Requires |
|---------|-------------|----------|
| `sqlite` | SQLite file-based store | None |
| `postgres` | PostgreSQL store | Docker |
| `redis` | Redis cache layer | Docker |
| `memory` | In-memory store | None |
| `composite` | Multi-backend composite store | Docker |
| `e2e` | End-to-end server tests | Docker |
| `all` | All backends (default) | Docker |

## Examples

```bash
# Quick local test with SQLite
.claude/skills/gode-testing/scripts/test-integration.sh sqlite

# Test PostgreSQL (starts Docker container)
.claude/skills/gode-testing/scripts/test-integration.sh postgres

# Full integration suite
.claude/skills/gode-testing/scripts/test-integration.sh all
```

## Prerequisites

### Docker Setup

For `postgres`, `redis`, `composite`, and `e2e` tests:

```bash
# Start Docker daemon
sudo systemctl start docker  # Linux
# or launch Docker Desktop on macOS/Windows

# Verify Docker is running
docker ps
```

### Test Configuration

The script automatically configures environment variables:

| Backend | Environment Variables |
|---------|----------------------|
| SQLite | `GODE_STORE_TYPE=sqlite`, `GODE_STORE_PATH=.dev/test/test.db` |
| PostgreSQL | `GODE_DB_HOST=localhost`, `GODE_DB_PORT=5432`, etc. |
| Redis | `GODE_REDIS_ADDR=localhost:6379` |

## Output

```
.dev/
├── coverage/
│   ├── store-sqlite.out    # Store coverage per backend
│   ├── store-sqlite.log
│   ├── store-postgres.out
│   ├── handlers-sqlite.out # Handler coverage per backend
│   ├── redis.log
│   └── e2e.log
└── test/
    └── test.db             # SQLite test database
```

## Docker Compose

Tests use `docker-compose.test.yml` for services:

```bash
# View running containers
docker-compose -f docker-compose.test.yml ps

# Manual cleanup
docker-compose -f docker-compose.test.yml down -v

# View logs
docker-compose -f docker-compose.test.yml logs postgres
```

## Troubleshooting

**Docker connection refused**: Ensure Docker daemon is running.

**Port already in use**: Run `docker-compose -f docker-compose.test.yml down -v` to clean up.

**Tests timeout**: Increase timeout in the script or check system resources.

**PostgreSQL syntax errors**: Ensure migrations are compatible with both PostgreSQL and SQLite (avoid `DO $$` blocks).

## E2E Test Structure

End-to-end tests are located in `test/e2e/` and test:

- Server startup and shutdown
- gRPC and HTTP endpoint functionality
- Authentication flow
- DE execution lifecycle
- Multi-tenancy isolation
