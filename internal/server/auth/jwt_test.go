package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		setupFn   func() string
		name      string
		wantClaim string
		wantErr   bool
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
	service := NewJWTService("test-secret", 10*time.Millisecond)

	// This sets the precision for the JWT package used for creating tokens.
	// We set this to a lower value in tests as the default is 1s.
	jwt.TimePrecision = time.Millisecond

	// Generate token
	token, err := service.GenerateToken("testuser")
	require.NoError(t, err)

	// Should be valid immediately
	claims, err := service.ValidateToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, "testuser", claims.Username)

	// Wait for expiration
	time.Sleep(15 * time.Millisecond)

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

func TestGenerateTokenPair(t *testing.T) {
	service := NewJWTService("test-secret", 15*time.Minute)

	accessToken, refreshToken, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.NotEqual(t, accessToken, refreshToken)

	// Validate access token
	accessClaims, err := service.ValidateToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, "testuser", accessClaims.Username)
	assert.Equal(t, AccessToken, accessClaims.TokenType)

	// Validate refresh token
	refreshClaims, err := service.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, "testuser", refreshClaims.Username)
	assert.Equal(t, RefreshToken, refreshClaims.TokenType)
}

func TestValidateRefreshToken(t *testing.T) {
	service := NewJWTService("test-secret", 15*time.Minute)

	// Generate token pair
	accessToken, refreshToken, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)

	// Refresh token should be validated successfully
	claims, err := service.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, RefreshToken, claims.TokenType)

	// Access token should fail validation as refresh token
	_, err = service.ValidateRefreshToken(accessToken)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTokenType, err)
}

func TestValidateAccessToken_WithRefreshToken(t *testing.T) {
	service := NewJWTService("test-secret", 15*time.Minute)

	// Generate token pair
	accessToken, refreshToken, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)

	// Access token should be validated successfully
	claims, err := service.ValidateToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, AccessToken, claims.TokenType)

	// Refresh token should fail validation as access token
	_, err = service.ValidateToken(refreshToken)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTokenType, err)
}

func TestRefreshAccessToken(t *testing.T) {
	service := NewJWTService("test-secret", 15*time.Minute)

	// Generate initial token pair
	oldAccessToken, oldRefreshToken, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)

	// Small delay to ensure new tokens have different timestamps
	time.Sleep(10 * time.Millisecond)

	// Refresh the access token
	newAccessToken, newRefreshToken, err := service.RefreshAccessToken(oldRefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)

	// New tokens should be different from old ones
	assert.NotEqual(t, oldAccessToken, newAccessToken)
	assert.NotEqual(t, oldRefreshToken, newRefreshToken)

	// Validate new access token
	accessClaims, err := service.ValidateToken(newAccessToken)
	require.NoError(t, err)
	assert.Equal(t, "testuser", accessClaims.Username)
	assert.Equal(t, AccessToken, accessClaims.TokenType)

	// Validate new refresh token
	refreshClaims, err := service.ValidateRefreshToken(newRefreshToken)
	require.NoError(t, err)
	assert.Equal(t, "testuser", refreshClaims.Username)
	assert.Equal(t, RefreshToken, refreshClaims.TokenType)
}

func TestRefreshAccessToken_WithExpiredRefreshToken(t *testing.T) {
	// Create service with very short refresh token expiry
	service := NewJWTServiceWithRefreshExpiry("test-secret", 15*time.Minute, 10*time.Millisecond)
	jwt.TimePrecision = time.Millisecond

	// Generate token pair
	_, refreshToken, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)

	// Wait for refresh token to expire
	time.Sleep(15 * time.Millisecond)

	// Attempt to refresh should fail
	_, _, err = service.RefreshAccessToken(refreshToken)
	assert.Error(t, err)
	assert.Equal(t, ErrExpiredToken, err)
}

func TestRefreshAccessToken_WithAccessToken(t *testing.T) {
	service := NewJWTService("test-secret", 15*time.Minute)

	// Generate token pair
	accessToken, _, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)

	// Trying to use access token to refresh should fail
	_, _, err = service.RefreshAccessToken(accessToken)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTokenType, err)
}

func TestTokenTypeField(t *testing.T) {
	service := NewJWTService("test-secret", 15*time.Minute)

	// Generate token pair
	accessToken, refreshToken, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)

	// Parse access token and check type
	accessClaims, err := service.ValidateToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, AccessToken, accessClaims.TokenType)

	// Parse refresh token and check type
	refreshClaims, err := service.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, RefreshToken, refreshClaims.TokenType)
}

func TestRefreshTokenRotation(t *testing.T) {
	service := NewJWTService("test-secret", 15*time.Minute)

	// Generate initial tokens
	_, refreshToken1, err := service.GenerateTokenPair("testuser")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	// First refresh
	_, refreshToken2, err := service.RefreshAccessToken(refreshToken1)
	require.NoError(t, err)
	assert.NotEqual(t, refreshToken1, refreshToken2)

	time.Sleep(10 * time.Millisecond)

	// Second refresh with new refresh token
	_, refreshToken3, err := service.RefreshAccessToken(refreshToken2)
	require.NoError(t, err)
	assert.NotEqual(t, refreshToken2, refreshToken3)
	assert.NotEqual(t, refreshToken1, refreshToken3)

	// All refresh tokens should be different
	tokens := []string{refreshToken1, refreshToken2, refreshToken3}
	for i := 0; i < len(tokens); i++ {
		for j := i + 1; j < len(tokens); j++ {
			assert.NotEqual(t, tokens[i], tokens[j])
		}
	}
}
