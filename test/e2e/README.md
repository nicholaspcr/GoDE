# End-to-End Integration Tests

This directory contains end-to-end integration tests for the GoDE server.

## Prerequisites

- Docker and docker-compose installed
- Go 1.25 or later

## Running E2E Tests

### Option 1: Against Running Server

If you already have a server running:

```bash
# Run tests against default server (localhost:3030)
go test -v ./test/e2e/...

# Or specify custom server address
E2E_SERVER_ADDR=localhost:9090 go test -v ./test/e2e/...
```

### Option 2: Using Docker Compose

Start the test environment:

```bash
# From repository root
docker-compose -f test/e2e/docker-compose.e2e.yaml up -d

# Wait for services to be ready (about 5 seconds)
sleep 5

# Run tests
go test -v ./test/e2e/...

# Cleanup
docker-compose -f test/e2e/docker-compose.e2e.yaml down
```

### Option 3: Using Make (Recommended)

```bash
# From repository root
make test-e2e
```

## Skipping E2E Tests

If you want to skip e2e tests in CI or local development:

```bash
E2E_SKIP=1 go test -v ./test/e2e/...
```

## Test Coverage

The e2e tests cover:

1. **Full User Workflow**
   - User registration
   - User login
   - Get user details
   - Update user information
   - Delete user

2. **DE Execution**
   - List available algorithms
   - List available problems
   - List available variants
   - Run differential evolution
   - Retrieve results

3. **Rate Limiting**
   - Login rate limits
   - Registration rate limits

4. **Authorization**
   - Unauthorized access attempts
   - Protected endpoint validation

## Environment Variables

- `E2E_SERVER_ADDR`: Server address (default: `localhost:3030`)
- `E2E_SKIP`: Skip e2e tests if set to any value

## Notes

- Tests use unique usernames based on timestamps to avoid conflicts
- Each test suite creates and cleans up its own test data
- Tests have a 30-second timeout per test case
- Rate limiting tests are best-effort and may not always trigger limits
