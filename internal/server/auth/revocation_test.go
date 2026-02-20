package auth

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	rediscache "github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestRevoker creates a RedisTokenRevoker backed by miniredis for testing.
func newTestRevoker(t *testing.T) (*RedisTokenRevoker, *miniredis.Miniredis) {
	t.Helper()
	s := miniredis.RunT(t)
	port, err := strconv.Atoi(s.Port())
	require.NoError(t, err)

	client, err := rediscache.NewClient(rediscache.Config{
		Host: s.Host(),
		Port: port,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	return NewRedisTokenRevoker(client), s
}

func TestRedisTokenRevoker_RevokeToken(t *testing.T) {
	t.Run("revokes valid JTI", func(t *testing.T) {
		revoker, _ := newTestRevoker(t)
		ctx := context.Background()

		err := revoker.RevokeToken(ctx, "test-jti-1", 5*time.Minute)
		require.NoError(t, err)

		revoked, err := revoker.IsRevoked(ctx, "test-jti-1")
		require.NoError(t, err)
		assert.True(t, revoked)
	})

	t.Run("ignores empty JTI", func(t *testing.T) {
		revoker, _ := newTestRevoker(t)
		ctx := context.Background()

		err := revoker.RevokeToken(ctx, "", 5*time.Minute)
		require.NoError(t, err)
	})

	t.Run("ignores non-positive TTL", func(t *testing.T) {
		revoker, _ := newTestRevoker(t)
		ctx := context.Background()

		err := revoker.RevokeToken(ctx, "test-jti-2", 0)
		require.NoError(t, err)

		// Key should not have been set
		revoked, err := revoker.IsRevoked(ctx, "test-jti-2")
		require.NoError(t, err)
		assert.False(t, revoked)
	})

	t.Run("token expires after TTL", func(t *testing.T) {
		revoker, s := newTestRevoker(t)
		ctx := context.Background()

		err := revoker.RevokeToken(ctx, "expiring-jti", 1*time.Second)
		require.NoError(t, err)

		revoked, err := revoker.IsRevoked(ctx, "expiring-jti")
		require.NoError(t, err)
		assert.True(t, revoked, "token should be revoked before TTL")

		// Fast-forward time past TTL
		s.FastForward(2 * time.Second)

		revoked, err = revoker.IsRevoked(ctx, "expiring-jti")
		require.NoError(t, err)
		assert.False(t, revoked, "token should no longer be revoked after TTL")
	})
}

func TestRedisTokenRevoker_IsRevoked(t *testing.T) {
	t.Run("returns false for unknown JTI", func(t *testing.T) {
		revoker, _ := newTestRevoker(t)
		ctx := context.Background()

		revoked, err := revoker.IsRevoked(ctx, "unknown-jti")
		require.NoError(t, err)
		assert.False(t, revoked)
	})

	t.Run("returns false for empty JTI", func(t *testing.T) {
		revoker, _ := newTestRevoker(t)
		ctx := context.Background()

		revoked, err := revoker.IsRevoked(ctx, "")
		require.NoError(t, err)
		assert.False(t, revoked)
	})

	t.Run("returns true for revoked JTI", func(t *testing.T) {
		revoker, _ := newTestRevoker(t)
		ctx := context.Background()

		require.NoError(t, revoker.RevokeToken(ctx, "revoked-jti", 10*time.Minute))

		revoked, err := revoker.IsRevoked(ctx, "revoked-jti")
		require.NoError(t, err)
		assert.True(t, revoked)
	})

	t.Run("multiple independent JTIs", func(t *testing.T) {
		revoker, _ := newTestRevoker(t)
		ctx := context.Background()

		require.NoError(t, revoker.RevokeToken(ctx, "jti-a", 10*time.Minute))

		revokedA, err := revoker.IsRevoked(ctx, "jti-a")
		require.NoError(t, err)
		assert.True(t, revokedA)

		revokedB, err := revoker.IsRevoked(ctx, "jti-b")
		require.NoError(t, err)
		assert.False(t, revokedB)
	})
}
