// Package redis provides a Redis client wrapper for caching and pub/sub operations.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Client wraps the Redis client with application-specific operations and circuit breaker.
type Client struct {
	rdb     *redis.Client
	breaker *gobreaker.CircuitBreaker
}

// Config holds Redis connection configuration.
type Config struct {
	Host          string
	Port          int
	Password      string
	DB            int
	BreakerConfig BreakerConfig // Circuit breaker configuration
}

// NewClient creates a new Redis client wrapper with circuit breaker.
func NewClient(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Use provided breaker config or defaults
	breakerCfg := cfg.BreakerConfig
	if breakerCfg.MaxFailures == 0 {
		breakerCfg = DefaultBreakerConfig()
	}

	breaker := newCircuitBreaker("redis-client", breakerCfg)

	return &Client{
		rdb:     rdb,
		breaker: breaker,
	}, nil
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// HealthCheck verifies the Redis connection is healthy.
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Get retrieves a value by key with circuit breaker protection.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	tracer := otel.Tracer("redis")
	ctx, span := tracer.Start(ctx, "redis.Get",
		trace.WithAttributes(attribute.String("redis.key", key)),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.rdb.Get(ctx, key).Result()
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	span.SetStatus(codes.Ok, "")
	str, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type from Redis Get: %T", result)
	}
	return str, nil
}

// Set stores a value with an optional TTL with circuit breaker protection.
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	tracer := otel.Tracer("redis")
	ctx, span := tracer.Start(ctx, "redis.Set",
		trace.WithAttributes(
			attribute.String("redis.key", key),
			attribute.String("redis.ttl", ttl.String()),
		),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.rdb.Set(ctx, key, value, ttl).Err()
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

// Delete removes a key with circuit breaker protection.
func (c *Client) Delete(ctx context.Context, key string) error {
	tracer := otel.Tracer("redis")
	ctx, span := tracer.Start(ctx, "redis.Delete",
		trace.WithAttributes(attribute.String("redis.key", key)),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.rdb.Del(ctx, key).Err()
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

// Publish publishes a message to a channel with circuit breaker protection.
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	tracer := otel.Tracer("redis")
	ctx, span := tracer.Start(ctx, "redis.Publish",
		trace.WithAttributes(attribute.String("redis.channel", channel)),
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer span.End()

	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.rdb.Publish(ctx, channel, message).Err()
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "")
	return nil
}

// Subscribe subscribes to a channel and returns a PubSub instance.
func (c *Client) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return c.rdb.Subscribe(ctx, channel)
}

// HSet sets a hash field with circuit breaker protection.
func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.rdb.HSet(ctx, key, values...).Err()
	})
	return err
}

// HGetAll retrieves all fields and values from a hash with circuit breaker protection.
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.rdb.HGetAll(ctx, key).Result()
	})
	if err != nil {
		return nil, err
	}
	m, ok := result.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("unexpected type from Redis HGetAll: %T", result)
	}
	return m, nil
}

// HGet retrieves a specific field from a hash with circuit breaker protection.
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.rdb.HGet(ctx, key, field).Result()
	})
	if err != nil {
		return "", err
	}
	str, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type from Redis HGet: %T", result)
	}
	return str, nil
}

// Expire sets a TTL on a key with circuit breaker protection.
func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	_, err := c.breaker.Execute(func() (interface{}, error) {
		return nil, c.rdb.Expire(ctx, key, ttl).Err()
	})
	return err
}

// Keys returns all keys matching a pattern with circuit breaker protection.
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		return c.rdb.Keys(ctx, pattern).Result()
	})
	if err != nil {
		return nil, err
	}
	return result.([]string), nil
}
