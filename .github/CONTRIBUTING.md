# Contributing to GoDE

Thank you for your interest in contributing to GoDE!

## Development Setup

### Prerequisites

- Go 1.22 or later
- Protocol Buffers compiler (protoc)
- golangci-lint (optional, for local linting)

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/nicholaspcr/GoDE.git
   cd GoDE
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

## Development Workflow

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### Linting

```bash
# Run golangci-lint
make lint

# Or directly
golangci-lint run
```

### Protocol Buffers

```bash
# Regenerate proto files
make proto-generate

# Lint proto files
make proto-lint
```

## Pull Request Process

1. **Fork the repository** and create your branch from `master`

2. **Make your changes** following the coding standards:
   - Write tests for new functionality
   - Ensure all tests pass
   - Follow Go best practices and idioms
   - Add comments for exported functions and types

3. **Test your changes**:
   ```bash
   make test
   make lint
   ```

4. **Commit your changes** with clear, descriptive commit messages:
   ```bash
   git commit -m "feat: add new feature X"
   git commit -m "fix: resolve issue with Y"
   git commit -m "test: add tests for Z"
   ```

   We follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` - New features
   - `fix:` - Bug fixes
   - `test:` - Adding or updating tests
   - `docs:` - Documentation changes
   - `refactor:` - Code refactoring
   - `chore:` - Maintenance tasks
   - `ci:` - CI/CD changes

5. **Push to your fork** and submit a pull request

6. **Wait for CI checks** to pass:
   - All tests must pass
   - Linting must pass
   - Build must succeed

7. **Address review comments** if any

## CI/CD

Our CI/CD pipeline runs automatically on all pull requests and pushes to master:

### CI Pipeline (`.github/workflows/ci.yml`)

1. **Lint**: Runs golangci-lint with comprehensive checks
2. **Test**: Runs tests on Go 1.22 and 1.23 with race detection
3. **Build**: Builds both server and CLI binaries
4. **Proto Check**: Verifies proto-generated files are up to date
5. **Coverage**: Uploads coverage reports to Codecov

### Release Pipeline (`.github/workflows/release.yml`)

Automatically creates releases when you push a version tag:

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

This will:
- Build binaries for Linux and macOS (amd64 and arm64)
- Create SHA256 checksums
- Create a GitHub release with all artifacts

## Code Style

- Follow standard Go formatting (`gofmt`, `goimports`)
- Keep functions focused and concise
- Write table-driven tests where appropriate
- Document exported types and functions
- Handle errors explicitly

## Testing Guidelines

1. **Unit Tests**: Test individual functions and methods in isolation
2. **Handler Tests**: Use mock stores to test handlers
3. **Integration Tests**: Test complete workflows (optional)
4. **Coverage**: Aim for >80% coverage on new code

### Example Test Structure

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "expected", false},
        {"invalid input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Security

- **Never commit secrets** or credentials
- Use environment variables for configuration
- Follow security best practices in the codebase
- Report security vulnerabilities privately

## Questions?

Feel free to open an issue for any questions or clarifications!
