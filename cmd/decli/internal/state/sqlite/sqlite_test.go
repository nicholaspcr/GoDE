package sqlite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	t.Run("memory provider", func(t *testing.T) {
		ops, err := New(ctx, Config{Provider: "memory"})
		require.NoError(t, err)
		assert.NotNil(t, ops)
	})

	t.Run("file provider with temp path", func(t *testing.T) {
		tmpFile := t.TempDir() + "/test.db"
		ops, err := New(ctx, Config{Provider: "file", Filepath: tmpFile})
		require.NoError(t, err)
		assert.NotNil(t, ops)
	})

	t.Run("invalid provider", func(t *testing.T) {
		ops, err := New(ctx, Config{Provider: "invalid"})
		assert.Error(t, err)
		assert.Nil(t, ops)
		assert.Contains(t, err.Error(), "invalid store type")
	})
}

func TestAuthTokenStore_SaveAndGet(t *testing.T) {
	ctx := context.Background()
	ops, err := New(ctx, Config{Provider: "memory"})
	require.NoError(t, err)

	t.Run("get token when none saved", func(t *testing.T) {
		_, err := ops.GetAuthToken()
		assert.Error(t, err)
	})

	t.Run("save and retrieve token", func(t *testing.T) {
		err := ops.SaveAuthToken("test-token-123")
		require.NoError(t, err)

		token, err := ops.GetAuthToken()
		require.NoError(t, err)
		assert.Equal(t, "test-token-123", token)
	})

	t.Run("get latest token after multiple saves", func(t *testing.T) {
		err := ops.SaveAuthToken("token-1")
		require.NoError(t, err)

		err = ops.SaveAuthToken("token-2")
		require.NoError(t, err)

		token, err := ops.GetAuthToken()
		require.NoError(t, err)
		assert.Equal(t, "token-2", token)
	})
}

func TestAuthTokenStore_Invalidate(t *testing.T) {
	ctx := context.Background()
	ops, err := New(ctx, Config{Provider: "memory"})
	require.NoError(t, err)

	t.Run("invalidate when no tokens exist", func(t *testing.T) {
		err := ops.InvalidateAuthToken()
		assert.NoError(t, err)
	})

	t.Run("save then invalidate", func(t *testing.T) {
		err := ops.SaveAuthToken("token-to-invalidate")
		require.NoError(t, err)

		err = ops.InvalidateAuthToken()
		assert.NoError(t, err)
	})
}

func TestAuthTokenStore_FullLifecycle(t *testing.T) {
	ctx := context.Background()
	ops, err := New(ctx, Config{Provider: "memory"})
	require.NoError(t, err)

	// Save a token
	err = ops.SaveAuthToken("lifecycle-token")
	require.NoError(t, err)

	// Retrieve it
	token, err := ops.GetAuthToken()
	require.NoError(t, err)
	assert.Equal(t, "lifecycle-token", token)

	// Invalidate all tokens
	err = ops.InvalidateAuthToken()
	require.NoError(t, err)

	// Save a new token
	err = ops.SaveAuthToken("new-token")
	require.NoError(t, err)

	// Should get the new token
	token, err = ops.GetAuthToken()
	require.NoError(t, err)
	assert.Equal(t, "new-token", token)
}
