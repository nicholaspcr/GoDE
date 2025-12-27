// Package auth provides JWT-based authentication for the DE server.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken indicates a JWT token is malformed or has an invalid signature.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken indicates a JWT token has passed its expiration time.
	ErrExpiredToken = errors.New("token has expired")
	// ErrInvalidTokenType indicates the token type doesn't match what was expected.
	ErrInvalidTokenType = errors.New("invalid token type")
)

// TokenType represents the type of JWT token
type TokenType string

const (
	// AccessToken is a short-lived token for API access
	AccessToken TokenType = "access"
	// RefreshToken is a long-lived token for obtaining new access tokens
	RefreshToken TokenType = "refresh"
)

// Scope represents a permission scope for authorization
type Scope string

const (
	// ScopeUserRead allows reading user information
	ScopeUserRead Scope = "user:read"
	// ScopeUserWrite allows modifying user information
	ScopeUserWrite Scope = "user:write"
	// ScopeDERun allows running DE executions
	ScopeDERun Scope = "de:run"
	// ScopeDERead allows reading DE executions and results
	ScopeDERead Scope = "de:read"
	// ScopeParetoRead allows reading Pareto sets
	ScopeParetoRead Scope = "pareto:read"
	// ScopeParetoWrite allows modifying Pareto sets
	ScopeParetoWrite Scope = "pareto:write"
	// ScopeAdmin allows all operations
	ScopeAdmin Scope = "admin"
)

// DefaultUserScopes returns the default scopes for regular users
func DefaultUserScopes() []Scope {
	return []Scope{ScopeUserRead, ScopeUserWrite, ScopeDERun, ScopeDERead, ScopeParetoRead, ScopeParetoWrite}
}

// Claims represents the JWT claims
type Claims struct {
	Username  string    `json:"username"`
	TokenType TokenType `json:"token_type"`
	Scopes    []Scope   `json:"scopes,omitempty"`
	jwt.RegisteredClaims
}

// HasScope checks if the claims include a specific scope
func (c *Claims) HasScope(scope Scope) bool {
	for _, s := range c.Scopes {
		if s == ScopeAdmin || s == scope {
			return true
		}
	}
	return false
}

// HasAnyScope checks if the claims include any of the specified scopes
func (c *Claims) HasAnyScope(scopes ...Scope) bool {
	for _, scope := range scopes {
		if c.HasScope(scope) {
			return true
		}
	}
	return false
}

// JWTService defines methods for JWT token operations
type JWTService interface {
	GenerateToken(username string) (string, error)
	GenerateTokenPair(username string) (accessToken, refreshToken string, err error)
	ValidateToken(tokenString string) (*Claims, error)
	ValidateRefreshToken(tokenString string) (*Claims, error)
	RefreshAccessToken(refreshTokenString string) (accessToken, newRefreshToken string, err error)
}

type jwtService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
	audience      string
}

// JWTConfig holds configuration for JWT service.
type JWTConfig struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string // Token issuer (e.g., "gode-server")
	Audience      string // Intended audience (e.g., "gode-api")
}

// DefaultJWTConfig returns a JWTConfig with sensible defaults.
func DefaultJWTConfig(secretKey string, accessExpiry time.Duration) JWTConfig {
	return JWTConfig{
		SecretKey:     secretKey,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: 7 * 24 * time.Hour,
		Issuer:        "gode-server",
		Audience:      "gode-api",
	}
}

// NewJWTService creates a new JWT service instance
// accessExpiry is the duration for access tokens (typically short-lived, e.g., 15 minutes)
// refreshExpiry is the duration for refresh tokens (typically long-lived, e.g., 7 days)
func NewJWTService(secretKey string, accessExpiry time.Duration) JWTService {
	cfg := DefaultJWTConfig(secretKey, accessExpiry)
	return NewJWTServiceWithConfig(cfg)
}

// NewJWTServiceWithRefreshExpiry creates a new JWT service with custom refresh token expiry
func NewJWTServiceWithRefreshExpiry(secretKey string, accessExpiry, refreshExpiry time.Duration) JWTService {
	cfg := DefaultJWTConfig(secretKey, accessExpiry)
	cfg.RefreshExpiry = refreshExpiry
	return NewJWTServiceWithConfig(cfg)
}

// NewJWTServiceWithConfig creates a new JWT service with full configuration.
func NewJWTServiceWithConfig(cfg JWTConfig) JWTService {
	return &jwtService{
		secretKey:     []byte(cfg.SecretKey),
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
		issuer:        cfg.Issuer,
		audience:      cfg.Audience,
	}
}

// GenerateToken creates a new JWT access token for a user (legacy method for backward compatibility)
func (j *jwtService) GenerateToken(username string) (string, error) {
	return j.generateToken(username, AccessToken, j.accessExpiry)
}

// GenerateTokenPair creates both access and refresh tokens for a user
func (j *jwtService) GenerateTokenPair(username string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.generateToken(username, AccessToken, j.accessExpiry)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.generateToken(username, RefreshToken, j.refreshExpiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken is the internal method that creates a token with specific type and expiry
func (j *jwtService) generateToken(username string, tokenType TokenType, expiry time.Duration) (string, error) {
	return j.generateTokenWithScopes(username, tokenType, expiry, DefaultUserScopes())
}

// generateTokenWithScopes creates a token with specific scopes
func (j *jwtService) generateTokenWithScopes(username string, tokenType TokenType, expiry time.Duration, scopes []Scope) (string, error) {
	now := time.Now()
	claims := &Claims{
		Username:  username,
		TokenType: tokenType,
		Scopes:    scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{j.audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken parses and validates a JWT access token
func (j *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := j.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Verify it's an access token
	if claims.TokenType != AccessToken {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// ValidateRefreshToken parses and validates a JWT refresh token
func (j *jwtService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Verify it's a refresh token
	if claims.TokenType != RefreshToken {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// RefreshAccessToken validates a refresh token and generates a new access token and refresh token
func (j *jwtService) RefreshAccessToken(refreshTokenString string) (accessToken, newRefreshToken string, err error) {
	// Validate the refresh token
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	// Generate new token pair
	return j.GenerateTokenPair(claims.Username)
}

// parseToken is the internal method that parses and validates a token
func (j *jwtService) parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Build parser options for validation
	parserOpts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	}

	// Add issuer validation if configured
	if j.issuer != "" {
		parserOpts = append(parserOpts, jwt.WithIssuer(j.issuer))
	}

	// Add audience validation if configured
	if j.audience != "" {
		parserOpts = append(parserOpts, jwt.WithAudience(j.audience))
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Verify signing method (additional check beyond WithValidMethods)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	}, parserOpts...)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
