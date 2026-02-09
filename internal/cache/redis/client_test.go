package redis

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestClient creates a Client backed by miniredis for testing.
func newTestClient(t *testing.T) (*Client, *miniredis.Miniredis) {
	t.Helper()
	s := miniredis.RunT(t)
	port, err := strconv.Atoi(s.Port())
	require.NoError(t, err)

	client, err := NewClient(Config{
		Host: s.Host(),
		Port: port,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	return client, s
}

func TestNewClient_InvalidConnection(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "Invalid host",
			config: Config{
				Host: "invalid-host-that-does-not-exist",
				Port: 6379,
			},
		},
		{
			name: "Invalid port",
			config: Config{
				Host: "localhost",
				Port: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			assert.Error(t, err)
			assert.Nil(t, client)
			assert.Contains(t, err.Error(), "failed to connect to Redis")
		})
	}
}

func TestNewClient_ValidConnection(t *testing.T) {
	client, _ := newTestClient(t)
	assert.NotNil(t, client)
}

func TestNewClient_CustomBreakerConfig(t *testing.T) {
	s := miniredis.RunT(t)
	port, err := strconv.Atoi(s.Port())
	require.NoError(t, err)

	client, err := NewClient(Config{
		Host: s.Host(),
		Port: port,
		BreakerConfig: BreakerConfig{
			MaxFailures: 10,
			Timeout:     30 * time.Second,
			MaxRequests: 5,
		},
	})
	require.NoError(t, err)
	defer func() { _ = client.Close() }()
	assert.NotNil(t, client)
}

func TestClient_HealthCheck(t *testing.T) {
	client, _ := newTestClient(t)
	err := client.HealthCheck(context.Background())
	assert.NoError(t, err)
}

func TestClient_Close(t *testing.T) {
	client, _ := newTestClient(t)
	err := client.Close()
	assert.NoError(t, err)
}

func TestClient_GetSet(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	t.Run("set and get value", func(t *testing.T) {
		err := client.Set(ctx, "key1", "value1", 0)
		require.NoError(t, err)

		val, err := client.Get(ctx, "key1")
		require.NoError(t, err)
		assert.Equal(t, "value1", val)
	})

	t.Run("get non-existent key", func(t *testing.T) {
		_, err := client.Get(ctx, "nonexistent")
		assert.Error(t, err)
	})

	t.Run("set with TTL", func(t *testing.T) {
		err := client.Set(ctx, "ttl-key", "ttl-value", 1*time.Hour)
		require.NoError(t, err)

		val, err := client.Get(ctx, "ttl-key")
		require.NoError(t, err)
		assert.Equal(t, "ttl-value", val)
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		err := client.Set(ctx, "overwrite", "v1", 0)
		require.NoError(t, err)

		err = client.Set(ctx, "overwrite", "v2", 0)
		require.NoError(t, err)

		val, err := client.Get(ctx, "overwrite")
		require.NoError(t, err)
		assert.Equal(t, "v2", val)
	})
}

func TestClient_Delete(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	t.Run("delete existing key", func(t *testing.T) {
		err := client.Set(ctx, "del-key", "del-value", 0)
		require.NoError(t, err)

		err = client.Delete(ctx, "del-key")
		require.NoError(t, err)

		_, err = client.Get(ctx, "del-key")
		assert.Error(t, err) // Should be gone
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		err := client.Delete(ctx, "nonexistent-key")
		assert.NoError(t, err) // Redis DEL returns 0 but no error
	})
}

func TestClient_HSet_HGetAll(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	t.Run("set and get hash fields", func(t *testing.T) {
		err := client.HSet(ctx, "hash1", "field1", "val1", "field2", "val2")
		require.NoError(t, err)

		result, err := client.HGetAll(ctx, "hash1")
		require.NoError(t, err)
		assert.Equal(t, "val1", result["field1"])
		assert.Equal(t, "val2", result["field2"])
	})

	t.Run("get all from empty hash", func(t *testing.T) {
		result, err := client.HGetAll(ctx, "empty-hash")
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestClient_HGet(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	err := client.HSet(ctx, "hash2", "name", "alice", "age", "30")
	require.NoError(t, err)

	t.Run("get existing field", func(t *testing.T) {
		val, err := client.HGet(ctx, "hash2", "name")
		require.NoError(t, err)
		assert.Equal(t, "alice", val)
	})

	t.Run("get non-existent field", func(t *testing.T) {
		_, err := client.HGet(ctx, "hash2", "nonexistent")
		assert.Error(t, err)
	})
}

func TestClient_HDel(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	err := client.HSet(ctx, "hash3", "f1", "v1", "f2", "v2", "f3", "v3")
	require.NoError(t, err)

	t.Run("delete single field", func(t *testing.T) {
		err := client.HDel(ctx, "hash3", "f1")
		require.NoError(t, err)

		_, err = client.HGet(ctx, "hash3", "f1")
		assert.Error(t, err)

		// Other fields should still exist
		val, err := client.HGet(ctx, "hash3", "f2")
		require.NoError(t, err)
		assert.Equal(t, "v2", val)
	})

	t.Run("delete multiple fields", func(t *testing.T) {
		err := client.HDel(ctx, "hash3", "f2", "f3")
		require.NoError(t, err)

		result, err := client.HGetAll(ctx, "hash3")
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestClient_HLen(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	t.Run("empty hash", func(t *testing.T) {
		length, err := client.HLen(ctx, "empty-hlen")
		require.NoError(t, err)
		assert.Equal(t, int64(0), length)
	})

	t.Run("hash with fields", func(t *testing.T) {
		err := client.HSet(ctx, "hlen-hash", "a", "1", "b", "2", "c", "3")
		require.NoError(t, err)

		length, err := client.HLen(ctx, "hlen-hash")
		require.NoError(t, err)
		assert.Equal(t, int64(3), length)
	})
}

func TestClient_HScan(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	// Populate hash with multiple fields
	for i := range 10 {
		err := client.HSet(ctx, "scan-hash",
			fmt.Sprintf("field-%d", i), fmt.Sprintf("value-%d", i))
		require.NoError(t, err)
	}

	t.Run("scan all fields", func(t *testing.T) {
		keys, _, err := client.HScan(ctx, "scan-hash", 0, "*", 100)
		require.NoError(t, err)
		// HScan returns alternating field/value pairs
		assert.Equal(t, 20, len(keys)) // 10 fields * 2 (field + value)
	})

	t.Run("scan with pattern", func(t *testing.T) {
		keys, _, err := client.HScan(ctx, "scan-hash", 0, "field-1*", 100)
		require.NoError(t, err)
		assert.NotEmpty(t, keys)
	})

	t.Run("scan empty hash", func(t *testing.T) {
		keys, _, err := client.HScan(ctx, "empty-scan", 0, "*", 100)
		require.NoError(t, err)
		assert.Empty(t, keys)
	})
}

func TestClient_Expire(t *testing.T) {
	client, s := newTestClient(t)
	ctx := context.Background()

	t.Run("set expiry on key", func(t *testing.T) {
		err := client.Set(ctx, "expire-key", "expire-val", 0)
		require.NoError(t, err)

		err = client.Expire(ctx, "expire-key", 1*time.Hour)
		require.NoError(t, err)

		// Key should still exist
		val, err := client.Get(ctx, "expire-key")
		require.NoError(t, err)
		assert.Equal(t, "expire-val", val)

		// Fast-forward time in miniredis
		s.FastForward(2 * time.Hour)

		// Key should be expired
		_, err = client.Get(ctx, "expire-key")
		assert.Error(t, err)
	})
}

func TestClient_TTLExpiration(t *testing.T) {
	client, s := newTestClient(t)
	ctx := context.Background()

	err := client.Set(ctx, "ttl-test", "value", 30*time.Minute)
	require.NoError(t, err)

	// Key should exist before TTL
	val, err := client.Get(ctx, "ttl-test")
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	// Fast-forward past TTL
	s.FastForward(31 * time.Minute)

	// Key should be expired
	_, err = client.Get(ctx, "ttl-test")
	assert.Error(t, err)
}

func TestClient_Publish(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	t.Run("publish message", func(t *testing.T) {
		err := client.Publish(ctx, "test-channel", "test-message")
		assert.NoError(t, err)
	})

	t.Run("publish to different channels", func(t *testing.T) {
		err := client.Publish(ctx, "channel-1", "msg-1")
		assert.NoError(t, err)

		err = client.Publish(ctx, "channel-2", "msg-2")
		assert.NoError(t, err)
	})
}

func TestClient_Subscribe(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	t.Run("subscribe returns PubSub", func(t *testing.T) {
		pubsub := client.Subscribe(ctx, "sub-channel")
		require.NotNil(t, pubsub)
		defer func() { _ = pubsub.Close() }()
	})
}

func TestClient_Keys(t *testing.T) {
	client, _ := newTestClient(t)
	ctx := context.Background()

	// Set up test keys
	for i := range 5 {
		err := client.Set(ctx, fmt.Sprintf("prefix:%d", i), "val", 0)
		require.NoError(t, err)
	}
	err := client.Set(ctx, "other:key", "val", 0)
	require.NoError(t, err)

	t.Run("match all keys", func(t *testing.T) {
		keys, err := client.Keys(ctx, "*")
		require.NoError(t, err)
		assert.Len(t, keys, 6)
	})

	t.Run("match prefix pattern", func(t *testing.T) {
		keys, err := client.Keys(ctx, "prefix:*")
		require.NoError(t, err)
		assert.Len(t, keys, 5)
	})

	t.Run("no matches", func(t *testing.T) {
		keys, err := client.Keys(ctx, "nonexistent:*")
		require.NoError(t, err)
		assert.Empty(t, keys)
	})
}

func TestClient_CircuitBreaker_Integration(t *testing.T) {
	s := miniredis.RunT(t)
	port, err := strconv.Atoi(s.Port())
	require.NoError(t, err)

	rdb := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%d", s.Host(), port),
	})

	breaker := newCircuitBreaker("test", BreakerConfig{
		MaxFailures: 3,
		Timeout:     100 * time.Millisecond,
		MaxRequests: 2,
	})

	client := &Client{rdb: rdb, breaker: breaker}
	defer func() { _ = client.Close() }()
	ctx := context.Background()

	// Normal operations should work
	err = client.Set(ctx, "cb-key", "cb-val", 0)
	require.NoError(t, err)

	val, err := client.Get(ctx, "cb-key")
	require.NoError(t, err)
	assert.Equal(t, "cb-val", val)

	// Stop miniredis to simulate failures
	s.Close()

	// Operations should start failing
	for range 4 {
		_ = client.Set(ctx, "fail-key", "val", 0)
	}

	// Circuit should be open now
	assert.Equal(t, gobreaker.StateOpen, breaker.State())

	// Subsequent calls should fail fast with circuit breaker error
	err = client.Set(ctx, "blocked-key", "val", 0)
	assert.Error(t, err)
}

func TestNewClient_ConfigFormatting(t *testing.T) {
	tests := []struct {
		name         string
		config       Config
		expectedAddr string
	}{
		{
			name:         "Standard port",
			config:       Config{Host: "localhost", Port: 6379},
			expectedAddr: "localhost:6379",
		},
		{
			name:         "Custom port",
			config:       Config{Host: "redis.example.com", Port: 7000},
			expectedAddr: "redis.example.com:7000",
		},
		{
			name:         "IP address",
			config:       Config{Host: "192.168.1.100", Port: 6379},
			expectedAddr: "192.168.1.100:6379",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := fmt.Sprintf("%s:%d", tt.config.Host, tt.config.Port)
			assert.Equal(t, tt.expectedAddr, addr)
		})
	}
}

func TestClient_InterfaceCompliance(t *testing.T) {
	// Verify Client implements ClientInterface at compile time
	var _ ClientInterface = (*Client)(nil)
}
