package auth

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"

	rediscache "github.com/nicholaspcr/GoDE/internal/cache/redis"
)

// TokenRevoker manages a token revocation list.
// Implementations must be safe for concurrent use.
type TokenRevoker interface {
	// RevokeToken marks a token JTI as revoked until the given TTL expires.
	RevokeToken(ctx context.Context, jti string, ttl time.Duration) error
	// IsRevoked reports whether the given JTI has been revoked.
	IsRevoked(ctx context.Context, jti string) (bool, error)
}

// redisRevokerKeyPrefix is the Redis key prefix for revoked JTIs.
const redisRevokerKeyPrefix = "revoked:token:"

// RedisTokenRevoker implements TokenRevoker using Redis.
type RedisTokenRevoker struct {
	client rediscache.ClientInterface
}

// NewRedisTokenRevoker creates a new Redis-backed token revoker.
func NewRedisTokenRevoker(client rediscache.ClientInterface) *RedisTokenRevoker {
	return &RedisTokenRevoker{client: client}
}

// RevokeToken stores the JTI in Redis with the given TTL.
func (r *RedisTokenRevoker) RevokeToken(ctx context.Context, jti string, ttl time.Duration) error {
	if jti == "" || ttl <= 0 {
		return nil
	}
	return r.client.Set(ctx, redisRevokerKeyPrefix+jti, "1", ttl)
}

// IsRevoked checks whether the JTI exists in the revocation list.
func (r *RedisTokenRevoker) IsRevoked(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	_, err := r.client.Get(ctx, redisRevokerKeyPrefix+jti)
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
