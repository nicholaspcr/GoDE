# Performance Benchmarks

Run and compare performance benchmarks for DE algorithms.

## Usage

```bash
.claude/skills/gode-testing/scripts/benchmark-de.sh [target] [--compare=<file>]
```

## Available Targets

| Target | Packages Benchmarked |
|--------|---------------------|
| `variants` | `./pkg/variants/*` - All mutation variants |
| `algorithms` | `./pkg/de`, `./pkg/de/gde3` - DE algorithms |
| `problems` | `./pkg/problems/multi`, `./pkg/problems/many/*` - Test problems |
| `utils` | `./pkg/de` - Utility functions |
| `all` | Comprehensive benchmarks (default) |

## Examples

```bash
# Benchmark all variants
.claude/skills/gode-testing/scripts/benchmark-de.sh variants

# Full benchmark suite
.claude/skills/gode-testing/scripts/benchmark-de.sh all

# Compare with previous run
.claude/skills/gode-testing/scripts/benchmark-de.sh all --compare=.dev/benchmarks/latest.txt
```

## Output

```
.dev/benchmarks/
├── bench_20240101_120000.txt   # Timestamped results
├── bench_20240102_143000.txt
└── latest.txt                   # Symlink to most recent
```

## Metrics

Benchmarks capture:

| Metric | Description | Unit |
|--------|-------------|------|
| `ns/op` | Time per operation | nanoseconds |
| `B/op` | Memory allocated per operation | bytes |
| `allocs/op` | Heap allocations per operation | count |

Multi-CPU scaling is tested with 1, 2, and 4 cores.

## Benchmark Comparison

### Install benchstat

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

### Compare Results

```bash
# Using the skill
.claude/skills/gode-testing/scripts/benchmark-de.sh all --compare=.dev/benchmarks/bench_20240101.txt

# Manual comparison
benchstat .dev/benchmarks/old.txt .dev/benchmarks/new.txt
```

### Reading Comparison Output

```
name           old time/op    new time/op    delta
Mutate-4       125ns ± 2%     118ns ± 1%    -5.60%  (p=0.008 n=5+5)
```

- Negative delta = performance improvement
- `±` shows variance
- `p` value indicates statistical significance

## Performance Guidelines

### Variant Mutation

Expected performance for single mutation:

| Operation | Target |
|-----------|--------|
| `Mutate` | < 200 ns/op |
| Memory | < 100 B/op |
| Allocs | < 3 allocs/op |

### Problem Evaluation

| Problem | Target (per vector) |
|---------|---------------------|
| ZDT | < 500 ns/op |
| DTLZ | < 1 µs/op |
| WFG | < 2 µs/op |

## Best Practices

1. **Close other applications** during benchmarking for consistent results
2. **Run multiple times** and use `benchstat` for statistical analysis
3. **Save baseline** before optimization work
4. **Use `-benchtime=10s`** for more stable results on fast operations

## Troubleshooting

**High variance**: Run with longer `-benchtime` or close background applications.

**Results differ between runs**: Normal; use `benchstat` for statistical comparison.

**Missing benchmark functions**: Ensure functions are named `Benchmark*` and take `*testing.B`.
