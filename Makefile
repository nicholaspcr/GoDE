# vim: set foldmarker={,} foldlevel=0 foldmethod=marker:
#
# GoDE Makefile
# Comprehensive build, test, and development targets
#
# This Makefile is inspired by:
# https://github.com/vincentbernat/hellogopher/blob/master/Makefile

# Configuration
BINARY_NAME_SERVER := deserver
BINARY_NAME_CLI := decli
BUILD_DIR := .dev
COVERAGE_DIR := .dev/coverage
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

.DEFAULT_GOAL := help

##@ General

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\n$(CYAN)Usage:$(RESET)\n  make $(GREEN)<target>$(RESET)\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: version
version: ## Show current version from git
	@git describe --tags --always --dirty

##@ Setup & Dependencies

.PHONY: init
init: ## Initialize development environment
	@echo '$(GREEN)Setting up development directories...$(RESET)'
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)/server
	@mkdir -p $(BUILD_DIR)/cli
	@mkdir -p $(COVERAGE_DIR)
	@mkdir -p .env
	@mkdir -p .env/server
	@mkdir -p .env/cli
	@mkdir -p docs/openapi
	@echo '$(GREEN)✓ Development environment initialized$(RESET)'

.PHONY: deps
deps: ## Download Go dependencies
	@echo '$(GREEN)Downloading Go dependencies...$(RESET)'
	@go mod download
	@go mod verify
	@echo '$(GREEN)✓ Dependencies installed$(RESET)'

.PHONY: deps-update
deps-update: ## Update Go dependencies
	@echo '$(GREEN)Updating Go dependencies...$(RESET)'
	@go get -u ./...
	@go mod tidy
	@echo '$(GREEN)✓ Dependencies updated$(RESET)'

.PHONY: tools
tools: ## Install development tools
	@echo '$(GREEN)Installing development tools...$(RESET)'
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo '$(GREEN)✓ Development tools installed$(RESET)'

##@ Building

.PHONY: build
build: build-server build-cli ## Build all binaries

.PHONY: build-server
build-server: ## Build server binary
	@echo '$(GREEN)Building server...$(RESET)'
	@go build -o $(BUILD_DIR)/$(BINARY_NAME_SERVER) ./cmd/deserver
	@echo '$(GREEN)✓ Server built: $(BUILD_DIR)/$(BINARY_NAME_SERVER)$(RESET)'

.PHONY: build-cli
build-cli: ## Build CLI binary
	@echo '$(GREEN)Building CLI...$(RESET)'
	@go build -o $(BUILD_DIR)/$(BINARY_NAME_CLI) ./cmd/decli
	@echo '$(GREEN)✓ CLI built: $(BUILD_DIR)/$(BINARY_NAME_CLI)$(RESET)'

.PHONY: build-race
build-race: ## Build with race detector enabled
	@echo '$(GREEN)Building with race detector...$(RESET)'
	@go build -race -o $(BUILD_DIR)/$(BINARY_NAME_SERVER)-race ./cmd/deserver
	@go build -race -o $(BUILD_DIR)/$(BINARY_NAME_CLI)-race ./cmd/decli
	@echo '$(GREEN)✓ Race detector binaries built$(RESET)'

.PHONY: install
install: ## Install binaries to $GOPATH/bin
	@echo '$(GREEN)Installing binaries...$(RESET)'
	@go install ./cmd/deserver
	@go install ./cmd/decli
	@echo '$(GREEN)✓ Binaries installed to $(shell go env GOPATH)/bin$(RESET)'

##@ Testing

.PHONY: test
test: ## Run all unit tests (excludes e2e and integration)
	@echo '$(GREEN)Running unit tests...$(RESET)'
	@go test -v -short ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo '$(GREEN)Running unit tests...$(RESET)'
	@go test -v -short -run 'Test[^E2E]' ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo '$(GREEN)Running integration tests...$(RESET)'
	@go test -v -tags=integration ./internal/server/... ./internal/executor/... ./internal/store/...

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests using testcontainers (requires Docker)
	@echo '$(YELLOW)Starting E2E tests (requires Docker)...$(RESET)'
	@go test -v -tags=e2e -timeout=10m ./test/e2e/...

.PHONY: test-all
test-all: test test-e2e ## Run all tests (unit, integration, e2e)

.PHONY: test-race
test-race: ## Run tests with race detector
	@echo '$(GREEN)Running tests with race detector...$(RESET)'
	@go test -race -short ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@go test -v -cover ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo '$(GREEN)Running tests with coverage...$(RESET)'
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo '$(GREEN)✓ Coverage report: $(COVERAGE_DIR)/coverage.html$(RESET)'
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total | awk '{print "Total Coverage: " $$3}'

.PHONY: test-coverage-summary
test-coverage-summary: ## Show coverage summary
	@go test -cover ./... | grep -E '^(ok|FAIL)' | awk '{print $$2 " " $$5}'

##@ Benchmarking

.PHONY: bench
bench: ## Run all benchmarks
	@echo '$(GREEN)Running benchmarks...$(RESET)'
	@go test -bench=. -benchmem -run=^$$ ./...

.PHONY: bench-variants
bench-variants: ## Run variant benchmarks
	@echo '$(GREEN)Running variant benchmarks...$(RESET)'
	@go test -bench=. -benchmem -run=^$$ ./pkg/variants/...

.PHONY: bench-de
bench-de: ## Run DE algorithm benchmarks
	@echo '$(GREEN)Running DE algorithm benchmarks...$(RESET)'
	@go test -bench=. -benchmem -run=^$$ ./pkg/de/...

.PHONY: bench-compare
bench-compare: ## Run benchmarks and save results for comparison
	@echo '$(GREEN)Running benchmarks and saving results...$(RESET)'
	@mkdir -p $(COVERAGE_DIR)
	@go test -bench=. -benchmem -run=^$$ ./... > $(COVERAGE_DIR)/bench-$(shell date +%Y%m%d-%H%M%S).txt
	@echo '$(GREEN)✓ Benchmark results saved to $(COVERAGE_DIR)$(RESET)'

##@ Code Quality

.PHONY: fmt
fmt: ## Format Go code
	@echo '$(GREEN)Formatting code...$(RESET)'
	@gofmt -s -w $(GO_FILES)
	@goimports -w $(GO_FILES)
	@echo '$(GREEN)✓ Code formatted$(RESET)'

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo '$(GREEN)Checking code formatting...$(RESET)'
	@test -z $$(gofmt -l $(GO_FILES)) || (echo "$(RED)Code is not formatted. Run 'make fmt'$(RESET)" && exit 1)
	@echo '$(GREEN)✓ Code is properly formatted$(RESET)'

.PHONY: vet
vet: ## Run go vet
	@echo '$(GREEN)Running go vet...$(RESET)'
	@go vet ./...
	@echo '$(GREEN)✓ Go vet passed$(RESET)'

.PHONY: lint
lint: ## Run golangci-lint
	@echo '$(GREEN)Running golangci-lint...$(RESET)'
	@golangci-lint run
	@echo '$(GREEN)✓ Linting passed$(RESET)'

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo '$(GREEN)Running golangci-lint with auto-fix...$(RESET)'
	@golangci-lint run --fix

.PHONY: check
check: fmt-check vet lint ## Run all code quality checks

.PHONY: tidy
tidy: ## Tidy Go modules
	@echo '$(GREEN)Tidying Go modules...$(RESET)'
	@go mod tidy
	@echo '$(GREEN)✓ Go modules tidied$(RESET)'

##@ Protocol Buffers

.PHONY: proto
proto: proto-lint proto-remove proto-generate ## Lint, clean, and generate proto code

.PHONY: proto-lint
proto-lint: ## Lint proto files
	@echo '$(GREEN)Linting proto files...$(RESET)'
	@buf lint
	@echo '$(GREEN)✓ Proto files linted$(RESET)'

.PHONY: proto-remove
proto-remove: ## Remove generated proto files
	@echo '$(GREEN)Removing generated proto files...$(RESET)'
	@find . -type f -name '*.pb.go' -not -path "./vendor/*" -delete
	@find . -type f -name '*.pb.gw.go' -not -path "./vendor/*" -delete
	@echo '$(GREEN)✓ Generated proto files removed$(RESET)'

.PHONY: proto-generate
proto-generate: ## Generate Go code from proto definitions
	@echo '$(GREEN)Generating code from proto files...$(RESET)'
	@buf generate
	@echo '$(GREEN)✓ Proto code generated$(RESET)'

.PHONY: openapi
openapi: ## Generate OpenAPI specification from proto files
	@echo '$(GREEN)Generating OpenAPI specification...$(RESET)'
	@mkdir -p docs/openapi
	@buf generate --template buf.gen.openapi.yaml
	@echo '$(GREEN)✓ OpenAPI spec generated at docs/openapi/$(RESET)'

##@ Database

.PHONY: db-up
db-up: ## Start PostgreSQL and Redis containers
	@echo '$(GREEN)Starting database containers...$(RESET)'
	@docker compose -f docker-compose.test.yml up -d
	@echo '$(GREEN)✓ Databases started$(RESET)'

.PHONY: db-down
db-down: ## Stop database containers
	@echo '$(GREEN)Stopping database containers...$(RESET)'
	@docker compose -f docker-compose.test.yml down
	@echo '$(GREEN)✓ Databases stopped$(RESET)'

.PHONY: db-clean
db-clean: ## Stop and remove database containers and volumes
	@echo '$(YELLOW)Cleaning database containers and volumes...$(RESET)'
	@docker compose -f docker-compose.test.yml down -v
	@echo '$(GREEN)✓ Databases cleaned$(RESET)'

.PHONY: db-logs
db-logs: ## Show database container logs
	@docker compose -f docker-compose.test.yml logs -f

.PHONY: db-psql
db-psql: ## Connect to PostgreSQL via psql
	@docker compose -f docker-compose.test.yml exec postgres psql -U gode -d gode_test

.PHONY: db-redis
db-redis: ## Connect to Redis via redis-cli
	@docker compose -f docker-compose.test.yml exec redis redis-cli

##@ Development

.PHONY: run
run: build-server ## Build and run server (sources .dev/server/.env)
	@echo '$(GREEN)Starting server...$(RESET)'
	@set -a && [ -f ./.dev/server/.env ] && . ./.dev/server/.env && set +a && $(BUILD_DIR)/$(BINARY_NAME_SERVER) start

.PHONY: run-dev
run-dev: db-up ## Start database and run server with development settings
	@echo '$(GREEN)Starting development server...$(RESET)'
	@sleep 2
	@JWT_SECRET="development-secret-key-change-in-production-min-32-chars" \
		STORE_TYPE=postgres \
		STORE_POSTGRESQL_DNS="postgres://gode:gode123@localhost:5432/gode_test?sslmode=disable" \
		REDIS_HOST=localhost \
		REDIS_PORT=6379 \
		TRACING_ENABLED=true \
		METRICS_ENABLED=true \
		go run ./cmd/deserver start

.PHONY: run-race
run-race: build-race ## Run server with race detector
	@echo '$(GREEN)Starting server with race detector...$(RESET)'
	@set -a && [ -f ./.dev/server/.env ] && . ./.dev/server/.env && set +a && $(BUILD_DIR)/$(BINARY_NAME_SERVER)-race start

.PHONY: dev
dev: ## Run full development environment (Docker Compose)
	@echo '$(GREEN)Starting full development environment...$(RESET)'
	@docker compose -f docker-compose.yml up

.PHONY: dev-full
dev-full: db-up ## Start databases for full stack development
	@echo '$(GREEN)Databases started. Run "make run-dev" in one terminal and "make web-dev" in another$(RESET)'

.PHONY: watch
watch: ## Watch for changes and rebuild (requires entr or similar)
	@which entr > /dev/null || (echo "$(RED)entr not found. Install with: brew install entr (macOS) or apt-get install entr (Linux)$(RESET)" && exit 1)
	@echo '$(GREEN)Watching for changes...$(RESET)'
	@find . -name '*.go' | entr -r make run-dev

##@ Cleaning

.PHONY: clean
clean: ## Clean build artifacts and caches
	@echo '$(GREEN)Cleaning build artifacts...$(RESET)'
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME_SERVER) $(BUILD_DIR)/$(BINARY_NAME_CLI)
	@rm -rf $(BUILD_DIR)/*-race
	@rm -rf $(COVERAGE_DIR)
	@rm -f profile.cov profile.cov.tmp
	@rm -rf ./bin
	@echo '$(GREEN)✓ Build artifacts cleaned$(RESET)'

.PHONY: clean-all
clean-all: clean db-clean ## Clean everything including databases
	@echo '$(GREEN)✓ Complete cleanup done$(RESET)'

##@ Frontend

.PHONY: web-deps
web-deps: ## Install frontend dependencies
	@echo '$(GREEN)Installing frontend dependencies...$(RESET)'
	@cd web && npm install
	@echo '$(GREEN)✓ Frontend dependencies installed$(RESET)'

.PHONY: web-dev
web-dev: ## Run frontend development server
	@echo '$(GREEN)Starting frontend development server...$(RESET)'
	@cd web && npm run dev

.PHONY: web-build
web-build: ## Build frontend for production
	@echo '$(GREEN)Building frontend...$(RESET)'
	@cd web && npm run build
	@echo '$(GREEN)✓ Frontend built$(RESET)'

.PHONY: web-test
web-test: ## Run frontend tests
	@cd web && npm run test

.PHONY: web-lint
web-lint: ## Lint frontend code
	@cd web && npm run lint

.PHONY: web-format
web-format: ## Format frontend code with Prettier
	@cd web && npx prettier --write "src/**/*.{ts,tsx}"

.PHONY: web-api
web-api: openapi ## Generate TypeScript API client from OpenAPI spec
	@echo '$(GREEN)Generating TypeScript API client...$(RESET)'
	@cd web && npx @openapitools/openapi-generator-cli generate \
		-i ../docs/openapi/api.swagger.json \
		-g typescript-fetch \
		-o src/api/generated \
		--additional-properties=typescriptThreePlus=true,supportsES6=true
	@echo '$(GREEN)✓ TypeScript API client generated$(RESET)'

##@ Kubernetes

.PHONY: k8s-build
k8s-build: ## Build Docker image for Kubernetes
	@echo '$(GREEN)Building Docker image for Kubernetes...$(RESET)'
	@eval $$(minikube docker-env) && docker build -t gode-server:latest .
	@echo '$(GREEN)✓ Docker image built$(RESET)'

.PHONY: k8s-deploy
k8s-deploy: ## Deploy to Kubernetes (minikube)
	@echo '$(GREEN)Deploying to Kubernetes...$(RESET)'
	@kubectl apply -f k8s/configmap.yaml
	@kubectl apply -f k8s/secret.yaml
	@kubectl apply -f k8s/postgres.yaml
	@kubectl apply -f k8s/redis.yaml
	@echo 'Waiting for databases...'
	@kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s || true
	@kubectl wait --for=condition=ready pod -l app=redis --timeout=60s || true
	@kubectl apply -f k8s/deserver.yaml
	@kubectl wait --for=condition=ready pod -l app=deserver --timeout=120s || true
	@echo '$(GREEN)✓ Deployment complete$(RESET)'

.PHONY: k8s-delete
k8s-delete: ## Delete all Kubernetes resources
	@echo '$(GREEN)Deleting Kubernetes resources...$(RESET)'
	@kubectl delete -f k8s/deserver.yaml --ignore-not-found=true
	@kubectl delete -f k8s/redis.yaml --ignore-not-found=true
	@kubectl delete -f k8s/postgres.yaml --ignore-not-found=true
	@kubectl delete -f k8s/secret.yaml --ignore-not-found=true
	@kubectl delete -f k8s/configmap.yaml --ignore-not-found=true
	@kubectl delete pvc postgres-pvc --ignore-not-found=true
	@echo '$(GREEN)✓ Kubernetes resources deleted$(RESET)'

.PHONY: k8s-logs
k8s-logs: ## Show logs from deserver pods
	@kubectl logs -l app=deserver -f --max-log-requests=10

.PHONY: k8s-status
k8s-status: ## Show status of Kubernetes resources
	@echo '$(CYAN)Deployments:$(RESET)'
	@kubectl get deployments
	@echo '\n$(CYAN)Pods:$(RESET)'
	@kubectl get pods
	@echo '\n$(CYAN)Services:$(RESET)'
	@kubectl get services

.PHONY: k8s-url
k8s-url: ## Get application URL (minikube)
	@echo '$(CYAN)HTTP Gateway URL:$(RESET)'
	@minikube service deserver-http --url

##@ CI/CD

.PHONY: ci
ci: deps check test test-coverage ## Run CI pipeline (checks + tests + coverage)

.PHONY: ci-quick
ci-quick: check test ## Run quick CI checks (no coverage)

.PHONY: pre-commit
pre-commit: fmt vet lint test-unit ## Run pre-commit checks

.PHONY: pre-push
pre-push: check test ## Run pre-push checks

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo '$(GREEN)Building Docker image...$(RESET)'
	@docker build -t gode-server:latest .
	@echo '$(GREEN)✓ Docker image built$(RESET)'

.PHONY: docker-run
docker-run: docker-build ## Build and run Docker container
	@echo '$(GREEN)Running Docker container...$(RESET)'
	@docker run --rm -p 3030:3030 -p 8081:8081 \
		-e JWT_SECRET="development-secret-key-change-in-production-min-32-chars" \
		gode-server:latest

.PHONY: docker-compose-up
docker-compose-up: ## Start all services with Docker Compose
	@docker compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop all Docker Compose services
	@docker compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## Show Docker Compose logs
	@docker compose logs -f

##@ Information

.PHONY: info
info: ## Show project information
	@echo '$(CYAN)GoDE Project Information$(RESET)'
	@echo '$(YELLOW)Go Version:$(RESET)'
	@go version
	@echo '\n$(YELLOW)Project Structure:$(RESET)'
	@echo 'Binaries: $(BUILD_DIR)/'
	@echo 'Coverage: $(COVERAGE_DIR)/'
	@echo '\n$(YELLOW)Available Commands:$(RESET)'
	@echo 'Run "make help" to see all available targets'

.PHONY: deps-graph
deps-graph: ## Show dependency graph (requires graphviz)
	@go mod graph | grep -v '@' | awk '{print $$1}' | sort | uniq
