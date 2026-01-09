#!/bin/bash
# Skill: test-integration
# Description: Run integration tests across different store backends
# Usage: /test-integration [backend]
#   backend: sqlite, postgres, redis, memory, composite, or 'all' (default)

set -e

BACKEND="${1:-all}"
TEST_DIR=".dev/test"
COVERAGE_DIR=".dev/coverage"

# Create test directories
mkdir -p "$TEST_DIR" "$COVERAGE_DIR"

echo "=== Integration Testing ==="
echo "Backend: $BACKEND"
echo ""

# Cleanup function
cleanup() {
    echo "Cleaning up test resources..."
    docker-compose -f docker-compose.test.yml down -v 2>/dev/null || true
    rm -f "$TEST_DIR"/*.db
}

trap cleanup EXIT

# Start required services for integration tests
start_services() {
    local backend=$1

    case "$backend" in
        postgres|all)
            echo "Starting PostgreSQL for testing..."
            docker-compose -f docker-compose.test.yml up -d postgres
            sleep 3
            ;;
        redis|all)
            echo "Starting Redis for testing..."
            docker-compose -f docker-compose.test.yml up -d redis
            sleep 2
            ;;
    esac
}

# Run store tests with specific backend
run_store_tests() {
    local backend=$1
    local store_type=$2

    echo "Testing $backend store..."
    echo "----------------------------------------"

    export GODE_STORE_TYPE="$store_type"

    case "$backend" in
        sqlite)
            export GODE_STORE_PATH="$TEST_DIR/test.db"
            ;;
        postgres)
            export GODE_DB_HOST="localhost"
            export GODE_DB_PORT="5432"
            export GODE_DB_NAME="gode_test"
            export GODE_DB_USER="gode"
            export GODE_DB_PASSWORD="gode123"
            ;;
        redis)
            export GODE_REDIS_ADDR="localhost:6379"
            export GODE_REDIS_PASSWORD=""
            export GODE_REDIS_DB="0"
            ;;
    esac

    # Run GORM store tests
    go test -v -race -coverprofile="$COVERAGE_DIR/store-${backend}.out" \
        -covermode=atomic ./internal/store/gorm/... 2>&1 | tee "$COVERAGE_DIR/store-${backend}.log"

    # Run handler tests (they use stores)
    go test -v -race -coverprofile="$COVERAGE_DIR/handlers-${backend}.out" \
        -covermode=atomic ./internal/server/handlers/... 2>&1 | tee "$COVERAGE_DIR/handlers-${backend}.log"

    # Calculate coverage
    if [ -f "$COVERAGE_DIR/store-${backend}.out" ]; then
        coverage=$(go tool cover -func="$COVERAGE_DIR/store-${backend}.out" | grep total | awk '{print $3}')
        echo "Store coverage: $coverage"
    fi

    echo ""
}

# Run E2E tests
run_e2e_tests() {
    echo "Running E2E tests..."
    echo "----------------------------------------"

    start_services all

    export E2E_SKIP=0
    go test -v -timeout 5m ./test/e2e/... 2>&1 | tee "$COVERAGE_DIR/e2e.log"

    echo ""
}

case "$BACKEND" in
    sqlite)
        run_store_tests "sqlite" "sqlite"
        ;;
    postgres)
        start_services postgres
        run_store_tests "postgres" "postgres"
        ;;
    redis)
        start_services redis
        # Run cache tests
        echo "Testing Redis cache..."
        go test -v -race ./internal/cache/redis/... 2>&1 | tee "$COVERAGE_DIR/redis.log"
        ;;
    memory)
        run_store_tests "memory" "memory"
        ;;
    composite)
        start_services all
        run_store_tests "composite" "composite"
        ;;
    e2e)
        run_e2e_tests
        ;;
    all)
        echo "Running comprehensive integration tests..."
        echo ""

        run_store_tests "memory" "memory"
        run_store_tests "sqlite" "sqlite"

        start_services all
        run_store_tests "postgres" "postgres"

        echo "Testing Redis cache..."
        go test -v -race ./internal/cache/redis/... 2>&1 | tee "$COVERAGE_DIR/redis.log"

        run_e2e_tests

        echo "=== Integration Test Summary ==="
        for file in "$COVERAGE_DIR"/store-*.out; do
            if [ -f "$file" ]; then
                name=$(basename "$file" .out | sed 's/store-//')
                coverage=$(go tool cover -func="$file" | grep total | awk '{print $3}')
                printf "%-15s %s\n" "$name:" "$coverage"
            fi
        done
        ;;
    *)
        echo "Error: Unknown backend '$BACKEND'"
        echo "Valid options: sqlite, postgres, redis, memory, composite, e2e, all"
        exit 1
        ;;
esac

echo ""
echo "Integration test reports saved to $COVERAGE_DIR/"
