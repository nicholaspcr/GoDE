#!/bin/bash
# Skill: benchmark-de
# Description: Run and compare performance benchmarks for DE algorithms
# Usage: /benchmark-de [target] [options]
#   target: variants, algorithms, problems, utils, or 'all' (default)
#   options: --compare=<file> to compare with previous results

set -e

TARGET="${1:-all}"
BENCH_DIR=".dev/benchmarks"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BENCH_FILE="$BENCH_DIR/bench_${TIMESTAMP}.txt"
COMPARE_FILE=""

# Parse options
shift || true
for arg in "$@"; do
    case $arg in
        --compare=*)
            COMPARE_FILE="${arg#*=}"
            ;;
    esac
done

# Create benchmark directory
mkdir -p "$BENCH_DIR"

echo "=== DE Performance Benchmarks ==="
echo "Target: $TARGET"
echo "Timestamp: $TIMESTAMP"
echo ""

run_benchmark() {
    local package=$1
    local name=$2

    echo "Benchmarking $name..."
    echo "----------------------------------------"

    # Run benchmarks with memory stats
    go test -bench=. -benchmem -benchtime=5s -cpu=1,2,4 \
        -run=^$ "$package" 2>&1 | tee -a "$BENCH_FILE"

    echo ""
}

case "$TARGET" in
    variants)
        echo "Running variant benchmarks..."
        run_benchmark "./pkg/variants" "variants"
        run_benchmark "./pkg/variants/best" "best"
        run_benchmark "./pkg/variants/rand" "rand"
        run_benchmark "./pkg/variants/pbest" "pbest"
        run_benchmark "./pkg/variants/current-to-best" "current-to-best"
        ;;
    algorithms)
        echo "Running algorithm benchmarks..."
        run_benchmark "./pkg/de" "de-utils"
        run_benchmark "./pkg/de/gde3" "gde3"
        ;;
    problems)
        echo "Running problem benchmarks..."
        run_benchmark "./pkg/problems/multi" "multi-objective"
        run_benchmark "./pkg/problems/many/dtlz" "dtlz"
        run_benchmark "./pkg/problems/many/wfg" "wfg"
        ;;
    utils)
        echo "Running utility benchmarks..."
        run_benchmark "./pkg/de" "de-utils"
        ;;
    all)
        echo "Running comprehensive benchmarks..."
        echo ""

        # Variants
        echo "=== Variants ===" | tee -a "$BENCH_FILE"
        run_benchmark "./pkg/variants" "variants"
        run_benchmark "./pkg/variants/best" "best"
        run_benchmark "./pkg/variants/rand" "rand"
        run_benchmark "./pkg/variants/pbest" "pbest"
        run_benchmark "./pkg/variants/current-to-best" "current-to-best"

        # Algorithms
        echo "=== Algorithms ===" | tee -a "$BENCH_FILE"
        run_benchmark "./pkg/de" "de-utils"
        run_benchmark "./pkg/de/gde3" "gde3"

        # Problems
        echo "=== Problems ===" | tee -a "$BENCH_FILE"
        run_benchmark "./pkg/problems/multi" "multi-objective"
        run_benchmark "./pkg/problems/many/dtlz" "dtlz"
        run_benchmark "./pkg/problems/many/wfg" "wfg"
        ;;
    *)
        echo "Error: Unknown target '$TARGET'"
        echo "Valid options: variants, algorithms, problems, utils, all"
        exit 1
        ;;
esac

echo ""
echo "=== Benchmark Summary ==="
echo "Results saved to: $BENCH_FILE"

# Extract key metrics
echo ""
echo "Key Performance Indicators:"
grep -E "Benchmark.*-[0-9]+" "$BENCH_FILE" | \
    awk '{print $1, $3, $4, $5, $6}' | \
    column -t || true

# Compare with previous results if requested
if [ -n "$COMPARE_FILE" ] && [ -f "$COMPARE_FILE" ]; then
    echo ""
    echo "=== Performance Comparison ==="
    echo "Comparing with: $COMPARE_FILE"

    if command -v benchstat &> /dev/null; then
        benchstat "$COMPARE_FILE" "$BENCH_FILE"
    else
        echo "Install benchstat for detailed comparison: go install golang.org/x/perf/cmd/benchstat@latest"
        echo ""
        echo "Basic comparison:"
        echo "Old results:"
        grep -E "Benchmark.*-[0-9]+" "$COMPARE_FILE" | head -5
        echo ""
        echo "New results:"
        grep -E "Benchmark.*-[0-9]+" "$BENCH_FILE" | head -5
    fi
fi

# Save as latest
cp "$BENCH_FILE" "$BENCH_DIR/latest.txt"

echo ""
echo "Saved as latest: $BENCH_DIR/latest.txt"
echo ""
echo "To compare future benchmarks, use:"
echo "  /benchmark-de $TARGET --compare=$BENCH_FILE"
echo ""
echo "For detailed comparison, install benchstat:"
echo "  go install golang.org/x/perf/cmd/benchstat@latest"
