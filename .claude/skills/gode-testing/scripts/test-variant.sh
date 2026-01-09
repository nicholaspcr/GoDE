#!/bin/bash
# Skill: test-variant
# Description: Run tests for a specific DE variant with coverage reporting
# Usage: /test-variant [variant-name]
#   variant-name: best, rand, pbest, current-to-best, or 'all' for all variants

set -e

VARIANT="${1:-all}"
COVERAGE_DIR=".dev/coverage"

# Create coverage directory
mkdir -p "$COVERAGE_DIR"

echo "=== Testing DE Variants ==="
echo "Variant: $VARIANT"
echo ""

run_variant_tests() {
    local variant_path=$1
    local variant_name=$2

    echo "Testing $variant_name variant..."
    echo "----------------------------------------"

    # Run tests with coverage
    go test -v -race -coverprofile="$COVERAGE_DIR/${variant_name}.out" \
        -covermode=atomic "$variant_path" 2>&1 | tee "$COVERAGE_DIR/${variant_name}.log"

    # Calculate coverage percentage
    if [ -f "$COVERAGE_DIR/${variant_name}.out" ]; then
        coverage=$(go tool cover -func="$COVERAGE_DIR/${variant_name}.out" | \
            grep total | awk '{print $3}')
        echo "Coverage: $coverage"
    fi

    echo ""
}

case "$VARIANT" in
    best)
        run_variant_tests "./pkg/variants/best" "best"
        ;;
    rand)
        run_variant_tests "./pkg/variants/rand" "rand"
        ;;
    pbest)
        run_variant_tests "./pkg/variants/pbest" "pbest"
        ;;
    current-to-best)
        run_variant_tests "./pkg/variants/current-to-best" "current-to-best"
        ;;
    all)
        run_variant_tests "./pkg/variants/best" "best"
        run_variant_tests "./pkg/variants/rand" "rand"
        run_variant_tests "./pkg/variants/pbest" "pbest"
        run_variant_tests "./pkg/variants/current-to-best" "current-to-best"
        run_variant_tests "./pkg/variants" "utils"

        # Generate combined coverage report
        echo "=== Combined Coverage Summary ==="
        for file in "$COVERAGE_DIR"/*.out; do
            if [ -f "$file" ]; then
                name=$(basename "$file" .out)
                coverage=$(go tool cover -func="$file" | grep total | awk '{print $3}')
                printf "%-20s %s\n" "$name:" "$coverage"
            fi
        done
        ;;
    *)
        echo "Error: Unknown variant '$VARIANT'"
        echo "Valid options: best, rand, pbest, current-to-best, all"
        exit 1
        ;;
esac

echo ""
echo "Coverage reports saved to $COVERAGE_DIR/"
echo "View detailed coverage: go tool cover -html=$COVERAGE_DIR/<variant>.out"
