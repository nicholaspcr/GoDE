# Variant Testing

Test DE mutation variants with coverage reporting.

## Usage

```bash
.claude/skills/gode-testing/scripts/test-variant.sh [variant-name]
```

## Available Variants

| Variant | Path | Description |
|---------|------|-------------|
| `best` | `./pkg/variants/best` | DE/best/1 and DE/best/2 mutation strategies |
| `rand` | `./pkg/variants/rand` | DE/rand/1 and DE/rand/2 mutation strategies |
| `pbest` | `./pkg/variants/pbest` | DE/p-best mutation strategy |
| `current-to-best` | `./pkg/variants/current-to-best` | DE/current-to-best/1 mutation strategy |
| `all` | All variants + utils | Run comprehensive tests (default) |

## Examples

```bash
# Test best variants only
.claude/skills/gode-testing/scripts/test-variant.sh best

# Test all variants with full coverage
.claude/skills/gode-testing/scripts/test-variant.sh all
```

## Output

Coverage files are saved to `.dev/coverage/`:

```
.dev/coverage/
├── best.out           # Coverage data
├── best.log           # Test output
├── rand.out
├── rand.log
├── pbest.out
├── pbest.log
├── current-to-best.out
├── current-to-best.log
└── utils.out
```

## Viewing Coverage

```bash
# View coverage in browser
go tool cover -html=.dev/coverage/best.out

# Show coverage summary
go tool cover -func=.dev/coverage/best.out
```

## Coverage Targets

Recommended coverage thresholds:

| Package | Target |
|---------|--------|
| `best` | > 80% |
| `rand` | > 80% |
| `pbest` | > 80% |
| `current-to-best` | > 80% |
| `utils` | > 90% |

## Test Patterns

Tests use table-driven patterns with deterministic random seeds:

```go
func TestMutate(t *testing.T) {
    rng := rand.New(rand.NewSource(1)) // Deterministic
    // ...
}
```

## Troubleshooting

**Tests skipped**: Check that the variant package has `_test.go` files.

**Coverage file empty**: Ensure tests are actually running and not failing early.

**Race condition detected**: The `-race` flag is enabled by default; fix data races before proceeding.
