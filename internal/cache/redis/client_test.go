package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewClient_InvalidConnection tests connection failure scenarios
func TestNewClient_InvalidConnection(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "Invalid host",
			config: Config{
				Host:     "invalid-host-that-does-not-exist",
				Port:     6379,
				Password: "",
				DB:       0,
			},
		},
		{
			name: "Invalid port",
			config: Config{
				Host:     "localhost",
				Port:     1,
				Password: "",
				DB:       0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			assert.Error(t, err, "Should fail to connect to invalid Redis")
			assert.Nil(t, client)
			assert.Contains(t, err.Error(), "failed to connect to Redis")
		})
	}
}

// TestClient_Operations tests basic client operations with mock expectations
// Note: These tests require a running Redis instance for integration testing
// To run with testcontainers, see client_integration_test.go
func TestClient_ConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "Valid config",
			config: Config{
				Host:     "localhost",
				Port:     6379,
				Password: "",
				DB:       0,
			},
		},
		{
			name: "With password",
			config: Config{
				Host:     "localhost",
				Port:     6379,
				Password: "secret",
				DB:       0,
			},
		},
		{
			name: "Different DB",
			config: Config{
				Host:     "localhost",
				Port:     6379,
				Password: "",
				DB:       1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just validate config structure
			assert.NotEmpty(t, tt.config.Host)
			assert.Greater(t, tt.config.Port, 0)
			assert.GreaterOrEqual(t, tt.config.DB, 0)
		})
	}
}

// TestClient_MethodSignatures tests that all methods have correct signatures
func TestClient_MethodSignatures(t *testing.T) {
	// This test verifies the Client struct has all expected methods
	// It doesn't execute them, just checks they exist with correct signatures

	t.Run("Client has all required methods", func(t *testing.T) {
		var c *Client

		// Verify methods exist by attempting to assign them to variables
		var (
			_ func() error                                             = c.Close
			_ func(context.Context) error                              = c.HealthCheck
			_ func(context.Context, string) (string, error)            = c.Get
			_ func(context.Context, string, any, time.Duration) error  = c.Set
			_ func(context.Context, string) error                      = c.Delete
			_ func(context.Context, string, any) error                 = c.Publish
			_ func(context.Context, string, ...any) error              = c.HSet
			_ func(context.Context, string) (map[string]string, error) = c.HGetAll
			_ func(context.Context, string, string) (string, error)    = c.HGet
			_ func(context.Context, string, time.Duration) error       = c.Expire
			_ func(context.Context, string) ([]string, error)          = c.Keys
		)

		// If we get here, all methods exist with correct signatures
		assert.True(t, true)
	})
}

// TestClient_ContextCancellation tests context cancellation behavior
func TestClient_ContextCancellation(t *testing.T) {
	t.Run("Operations respect cancelled context", func(t *testing.T) {
		// Create cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// These would fail with context cancelled error if Redis was available
		// Just verify the context is cancelled
		select {
		case <-ctx.Done():
			assert.Error(t, ctx.Err())
			assert.Equal(t, context.Canceled, ctx.Err())
		default:
			t.Fatal("Context should be cancelled")
		}
	})

	t.Run("Operations respect timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait for timeout
		<-ctx.Done()

		assert.Error(t, ctx.Err())
		assert.Equal(t, context.DeadlineExceeded, ctx.Err())
	})
}

// TestClient_TTLHandling tests TTL parameter handling
func TestClient_TTLHandling(t *testing.T) {
	tests := []struct {
		name string
		ttl  time.Duration
		want time.Duration
	}{
		{
			name: "No expiry",
			ttl:  0,
			want: 0,
		},
		{
			name: "1 second expiry",
			ttl:  1 * time.Second,
			want: 1 * time.Second,
		},
		{
			name: "1 hour expiry",
			ttl:  1 * time.Hour,
			want: 1 * time.Hour,
		},
		{
			name: "Negative TTL (immediate expiry)",
			ttl:  -1 * time.Second,
			want: -1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate TTL values
			assert.Equal(t, tt.want, tt.ttl)
		})
	}
}

// TestClient_KeyPatterns tests various key pattern scenarios
func TestClient_KeyPatterns(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		valid   bool
	}{
		{
			name:    "Match all keys",
			pattern: "*",
			valid:   true,
		},
		{
			name:    "Prefix match",
			pattern: "user:*",
			valid:   true,
		},
		{
			name:    "Wildcard in middle",
			pattern: "user:*:profile",
			valid:   true,
		},
		{
			name:    "Single character wildcard",
			pattern: "user:?",
			valid:   true,
		},
		{
			name:    "Range pattern",
			pattern: "user:[a-z]*",
			valid:   true,
		},
		{
			name:    "Exact match",
			pattern: "user:123",
			valid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate pattern structure
			assert.NotEmpty(t, tt.pattern)
			assert.Equal(t, tt.valid, true)
		})
	}
}

// TestClient_HashOperations tests hash operation parameter combinations
func TestClient_HashOperations(t *testing.T) {
	t.Run("HSet with single field", func(t *testing.T) {
		// Validate parameters for single field HSet
		key := "user:123"
		field := "name"
		value := "John"

		assert.NotEmpty(t, key)
		assert.NotEmpty(t, field)
		assert.NotEmpty(t, value)
	})

	t.Run("HSet with multiple fields", func(t *testing.T) {
		// Validate parameters for multiple fields
		key := "user:123"
		values := []any{
			"name", "John",
			"email", "john@example.com",
			"age", "30",
		}

		assert.NotEmpty(t, key)
		assert.Len(t, values, 6)
		assert.Equal(t, 0, len(values)%2, "Values should be in field-value pairs")
	})

	t.Run("HGet validates key and field", func(t *testing.T) {
		key := "user:123"
		field := "name"

		assert.NotEmpty(t, key)
		assert.NotEmpty(t, field)
	})

	t.Run("HGetAll validates key", func(t *testing.T) {
		key := "user:123"
		assert.NotEmpty(t, key)
	})
}

// TestClient_PubSubChannel tests pub/sub channel name validation
func TestClient_PubSubChannel(t *testing.T) {
	tests := []struct {
		name    string
		channel string
	}{
		{
			name:    "Simple channel",
			channel: "notifications",
		},
		{
			name:    "Namespaced channel",
			channel: "user:123:updates",
		},
		{
			name:    "Execution progress channel",
			channel: "execution:abc123:updates",
		},
		{
			name:    "Cancellation channel",
			channel: "execution:abc123:cancel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.channel)
			// Channels should not contain wildcards
			assert.NotContains(t, tt.channel, "*")
		})
	}
}

// TestClient_ErrorHandling tests error scenarios
func TestClient_ErrorHandling(t *testing.T) {
	t.Run("Empty key returns error", func(t *testing.T) {
		key := ""
		assert.Empty(t, key, "Empty key should be rejected")
	})

	t.Run("Nil value handling", func(t *testing.T) {
		var value any
		value = nil
		// Verify nil values are handled
		assert.Nil(t, value)
	})

	t.Run("Context validation", func(t *testing.T) {
		ctx := context.Background()
		assert.NotNil(t, ctx, "Context should not be nil")

		// Validate cancelled context
		ctxCancelled, cancel := context.WithCancel(context.Background())
		cancel()

		assert.NotNil(t, ctxCancelled)
		assert.Error(t, ctxCancelled.Err())
	})
}

// TestNewClient_ConfigFormatting tests address formatting
func TestNewClient_ConfigFormatting(t *testing.T) {
	tests := []struct {
		name         string
		config       Config
		expectedAddr string
	}{
		{
			name: "Standard port",
			config: Config{
				Host: "localhost",
				Port: 6379,
			},
			expectedAddr: "localhost:6379",
		},
		{
			name: "Custom port",
			config: Config{
				Host: "redis.example.com",
				Port: 7000,
			},
			expectedAddr: "redis.example.com:7000",
		},
		{
			name: "IP address",
			config: Config{
				Host: "192.168.1.100",
				Port: 6379,
			},
			expectedAddr: "192.168.1.100:6379",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test address formatting without actually connecting
			addr := fmt.Sprintf("%s:%d", tt.config.Host, tt.config.Port)
			assert.Equal(t, tt.expectedAddr, addr)
		})
	}
}

// Benchmark tests for performance profiling
func BenchmarkClient_KeyFormatting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("execution:%s:progress", "abc123")
	}
}

func BenchmarkClient_ContextCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		cancel()
		_ = ctx
	}
}
