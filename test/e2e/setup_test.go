//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	testServerAddr string
	serverConn     *grpc.ClientConn
)

// TestMain sets up the test environment using testcontainers
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, pgConnStr, err := setupPostgres(ctx)
	if err != nil {
		fmt.Printf("Failed to start PostgreSQL container: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			fmt.Printf("Failed to terminate PostgreSQL container: %v\n", err)
		}
	}()

	// Start deserver container
	serverContainer, serverAddr, err := setupServer(ctx, pgConnStr)
	if err != nil {
		fmt.Printf("Failed to start server container: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := serverContainer.Terminate(ctx); err != nil {
			fmt.Printf("Failed to terminate server container: %v\n", err)
		}
	}()

	testServerAddr = serverAddr

	// Wait for server to be ready
	if err := waitForServer(ctx, serverAddr); err != nil {
		fmt.Printf("Server failed to become ready: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// setupPostgres creates and starts a PostgreSQL container
func setupPostgres(ctx context.Context) (testcontainers.Container, string, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("gode_test"),
		postgres.WithUsername("gode"),
		postgres.WithPassword("gode_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to start postgres container: %w", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get postgres host: %w", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get postgres port: %w", err)
	}

	connStr := fmt.Sprintf("postgres://gode:gode_password@%s:%s/gode_test?sslmode=disable", host, port.Port())

	return pgContainer, connStr, nil
}

// setupServer builds and starts the deserver container
func setupServer(ctx context.Context, pgConnStr string) (testcontainers.Container, string, error) {
	// Get the project root directory
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get project root: %w", err)
	}

	// Build request with Dockerfile
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    projectRoot,
			Dockerfile: "Dockerfile.server",
		},
		ExposedPorts: []string{"3030/tcp", "8081/tcp"},
		Env: map[string]string{
			"STORE_TYPE":           "postgres",
			"STORE_POSTGRESQL_DNS": pgConnStr,
			"JWT_SECRET":           "e2e-test-secret-key-with-sufficient-length-for-security",
			"GRPC_PORT":            "3030",
			"HTTP_PORT":            "8081",
			"METRICS_ENABLED":      "false",
		},
		WaitingFor: wait.ForHTTP("/health").
			WithPort("8081/tcp").
			WithStartupTimeout(60 * time.Second).
			WithPollInterval(1 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start server container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get server host: %w", err)
	}

	port, err := container.MappedPort(ctx, "3030")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get server port: %w", err)
	}

	serverAddr := fmt.Sprintf("%s:%s", host, port.Port())

	return container, serverAddr, nil
}

// waitForServer waits for the server to be ready to accept connections
func waitForServer(ctx context.Context, addr string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for server to be ready")
		case <-ticker.C:
			conn, err := grpc.NewClient(
				addr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err == nil {
				conn.Close()
				return nil
			}
		}
	}
}

// getTestServerAddr returns the test server address set up by TestMain
func getTestServerAddr() string {
	if testServerAddr != "" {
		return testServerAddr
	}
	// Fallback for when tests run individually (not via TestMain)
	return defaultServerAddr
}
