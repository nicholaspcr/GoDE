---
name: gode-testing
description: Test DE variants, run integration tests, and benchmark performance. Use for testing variants (best, rand, pbest, current-to-best), integration tests across store backends, and DE algorithm benchmarks.
allowed-tools: Read, Bash(go test:*), Bash(.claude/skills/gode-testing/scripts/*:*)
---

# GoDE Testing Skill

Testing and benchmarking skill for the GoDE Differential Evolution framework.

## Quick Start

### Test a specific variant
```bash
.claude/skills/gode-testing/scripts/test-variant.sh best
```

### Run integration tests
```bash
.claude/skills/gode-testing/scripts/test-integration.sh sqlite
```

### Benchmark performance
```bash
.claude/skills/gode-testing/scripts/benchmark-de.sh variants
```

## Available Scripts

| Script | Purpose | Arguments |
|--------|---------|-----------|
| `test-variant.sh` | Test DE variants with coverage | `best`, `rand`, `pbest`, `current-to-best`, `all` |
| `test-integration.sh` | Integration tests for store backends | `sqlite`, `postgres`, `redis`, `memory`, `e2e`, `all` |
| `benchmark-de.sh` | Performance benchmarks | `variants`, `algorithms`, `problems`, `all` |

## Output Locations

All outputs are saved to `.dev/` (gitignored):

```
.dev/
├── coverage/           # Test coverage reports
├── benchmarks/         # Performance benchmark results
└── test/              # Temporary test databases
```

## Detailed Documentation

- [VARIANTS.md](VARIANTS.md) - Variant testing details and coverage analysis
- [INTEGRATION.md](INTEGRATION.md) - Integration testing across store backends
- [BENCHMARKS.md](BENCHMARKS.md) - Performance benchmarking guide

## Requirements

- Go 1.21+
- Docker (for integration tests with postgres/redis)
- `benchstat` (optional, for benchmark comparisons): `go install golang.org/x/perf/cmd/benchstat@latest`
