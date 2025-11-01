package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTService(t *testing.T) {
	service := NewJWTService("test-secret", 24*time.Hour)
	assert.NotNil(t, service)
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		username string
		secret   string
		expiry   time.Duration
		wantErr  bool
	}{
		{
			name:     "valid token generation",
			username: "testuser",
			secret:   "test-secret",
			expiry:   24 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			secret:   "test-secret",
			expiry:   24 * time.Hour,
			wantErr:  false, // JWT allows empty claims
		},
		{
			name:     "short expiry",
			username: "testuser",
			secret:   "test-secret",
			expiry:   1 * time.Second,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewJWTService(tt.secret, tt.expiry)
			token, err := service.GenerateToken(tt.username)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	service := NewJWTService(secret, 24*time.Hour)

	tests := []struct {
		name      string
		setupFn   func() string
		wantErr   bool
		wantClaim string
	}{
		{
			name: "valid token",
			setupFn: func() string {
				token, _ := service.GenerateToken("testuser")
				return token
			},
			wantErr:   false,
			wantClaim: "testuser",
		},
		{
			name: "invalid token format",
			setupFn: func() string {
				return "invalid.token.format"
			},
			wantErr: true,
		},
		{
			name: "empty token",
			setupFn: func() string {
				return ""
			},
			wantErr: true,
		},
		{
			name: "malformed token",
			setupFn: func() string {
				return "not-a-jwt-token"
			},
			wantErr: true,
		},
		{
			name: "token with wrong secret",
			setupFn: func() string {
				wrongService := NewJWTService("wrong-secret", 24*time.Hour)
				token, _ := wrongService.GenerateToken("testuser")
				return token
			},
			wantErr: true,
		},
		{
			name: "expired token",
			setupFn: func() string {
				expiredService := NewJWTService(secret, -1*time.Hour)
				token, _ := expiredService.GenerateToken("testuser")
				return token
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupFn()
			claims, err := service.ValidateToken(token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, claims)
				assert.Equal(t, tt.wantClaim, claims.Username)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	service := NewJWTService("test-secret", 100*time.Millisecond)

	// Generate token
	token, err := service.GenerateToken("testuser")
	require.NoError(t, err)

	// Should be valid immediately
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)

	// Wait for expiration
	time.Sleep(200 * time.Millisecond)

	// Should be expired now
	_, err = service.ValidateToken(token)
	assert.Error(t, err)
	assert.Equal(t, ErrExpiredToken, err)
}

func TestTokenClaims(t *testing.T) {
	service := NewJWTService("test-secret", 24*time.Hour)

	username := "testuser"
	token, err := service.GenerateToken(username)
	require.NoError(t, err)

	claims, err := service.ValidateToken(token)
	require.NoError(t, err)

	// Verify claims
	assert.Equal(t, username, claims.Username)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)

	// Verify timestamps make sense
	now := time.Now()
	assert.True(t, claims.IssuedAt.Before(now.Add(1*time.Second)))
	assert.True(t, claims.ExpiresAt.After(now))
}

func TestMultipleTokens(t *testing.T) {
	service := NewJWTService("test-secret", 24*time.Hour)

	// Generate multiple tokens
	token1, err1 := service.GenerateToken("user1")
	token2, err2 := service.GenerateToken("user2")

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, token1, token2)

	// Validate each token
	claims1, err := service.ValidateToken(token1)
	require.NoError(t, err)
	assert.Equal(t, "user1", claims1.Username)

	claims2, err := service.ValidateToken(token2)
	require.NoError(t, err)
	assert.Equal(t, "user2", claims2.Username)
}
