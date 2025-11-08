# End-to-End Integration Tests

This directory contains end-to-end integration tests for the GoDE server using [testcontainers-go](https://github.com/testcontainers/testcontainers-go).

## Overview

These tests automatically spin up Docker containers for:
- PostgreSQL database
- GoDE server (deserver)

Containers are automatically created before tests run and cleaned up after tests complete.

## Prerequisites

- Docker daemon running locally
- Go 1.25 or later
- Docker socket accessible (default: `/var/run/docker.sock`)

## Running E2E Tests

### Using Make (Recommended)

```bash
# From repository root
make test-e2e
```

### Using Go Directly

```bash
# Run all e2e tests
go test -v -tags=e2e -timeout=10m ./test/e2e/...

# Run specific test
go test -v -tags=e2e -timeout=10m -run TestE2E_FullUserWorkflow ./test/e2e/...
```

## How It Works

1. **TestMain setup**: The `TestMain` function in `setup_test.go` runs before any tests
2. **Container creation**: PostgreSQL and deserver containers are started using testcontainers
3. **Automatic cleanup**: Containers are automatically terminated after tests complete
4. **Test execution**: Individual tests connect to the containerized server

## Configuration

The tests use the following default configuration:

- **PostgreSQL**:
  - Database: `gode_test`
  - User: `gode`
  - Password: `gode_password`

- **Server**:
  - gRPC port: dynamically assigned by testcontainers
  - HTTP port: dynamically assigned by testcontainers
  - JWT Secret: `e2e-test-secret-key-with-sufficient-length-for-security`

## Alternative: Manual Server Setup

If you prefer to run against a manually started server (e.g., for debugging):

```bash
# Start your server manually, then:
E2E_SERVER_ADDR=localhost:3030 go test -v -tags=e2e ./test/e2e/...
```

When `E2E_SERVER_ADDR` is set, the tests will use that address instead of starting containers.

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

- `E2E_SERVER_ADDR`: Server address to connect to instead of starting containers (optional)

## Troubleshooting

### Tests hang or timeout

- Check that Docker is running: `docker ps`
- Increase timeout: `go test -v -tags=e2e -timeout=15m ./test/e2e/...`
- Check container logs during test execution

### Container cleanup issues

If containers aren't cleaned up properly:

```bash
docker ps -a | grep testcontainers
docker rm -f <container-id>
```

### Port conflicts

Testcontainers automatically assigns random available ports, so port conflicts should not occur.

## Notes

- Tests use unique usernames based on timestamps to avoid conflicts
- Each test suite creates and cleans up its own test data
- Tests have a 30-second timeout per test case
- Rate limiting tests are best-effort and may not always trigger limits
- E2E tests are excluded from `go test ./...` by default (requires `-tags=e2e` flag)
