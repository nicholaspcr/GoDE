package server

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr string
	}{
		{
			name: "valid configuration",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           10,
					QueueSize:            100,
					MaxVectorsInProgress: 100,
					ExecutionTTL:         24 * time.Hour,
					ResultTTL:            7 * 24 * time.Hour,
					ProgressTTL:          time.Hour,
				},
				DE: de.Config{},
			},
			wantErr: "",
		},
		{
			name: "empty JWT secret",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "",
				JWTExpiry: 24 * time.Hour,
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "JWT_SECRET environment variable is required",
		},
		{
			name: "insecure default JWT secret",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "change-me-in-production",
				JWTExpiry: 24 * time.Hour,
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "insecure default value",
		},
		{
			name: "JWT secret too short",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "short",
				JWTExpiry: 24 * time.Hour,
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "must be at least 32 characters",
		},
		{
			name: "TLS enabled but missing cert file",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled:  true,
					CertFile: "",
					KeyFile:  "/path/to/key.pem",
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "cert_file is not specified",
		},
		{
			name: "TLS enabled but missing key file",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled:  true,
					CertFile: "/path/to/cert.pem",
					KeyFile:  "",
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "key_file is not specified",
		},
		{
			name: "invalid auth requests per minute",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    0,
					RegisterRequestsPerMinute: 0,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "login_requests_per_minute must be at least 1",
		},
		{
			name: "invalid DE executions per user",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       0,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "de_executions_per_user must be at least 1",
		},
		{
			name: "invalid max concurrent DE per user",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    0,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "max_concurrent_de_per_user must be at least 1",
		},
		{
			name: "invalid max requests per second",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      0,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "max_requests_per_second must be at least 1",
		},
		{
			name: "invalid max message size",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       512,
				},
			},
			wantErr: "max_message_size_bytes must be at least 1024 bytes",
		},
		{
			name: "empty listen address",
			config: Config{
				LisAddr:   "",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "listen address cannot be empty",
		},
		{
			name: "empty HTTP port",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  "",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "HTTP port cannot be empty",
		},
		{
			name: "JWT expiry too short",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 30 * time.Second,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
			},
			wantErr: "JWT expiry must be at least 1 minute",
		},
		{
			name: "invalid executor max workers",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           0,
					QueueSize:            100,
					MaxVectorsInProgress: 100,
					ExecutionTTL:         24 * time.Hour,
					ResultTTL:            7 * 24 * time.Hour,
					ProgressTTL:          time.Hour,
				},
			},
			wantErr: "executor max_workers must be at least 1",
		},
		{
			name: "invalid executor queue size",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           10,
					QueueSize:            0,
					MaxVectorsInProgress: 100,
					ExecutionTTL:         24 * time.Hour,
					ResultTTL:            7 * 24 * time.Hour,
					ProgressTTL:          time.Hour,
				},
			},
			wantErr: "executor queue_size must be at least 1",
		},
		{
			name: "invalid executor max vectors in progress",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           10,
					QueueSize:            100,
					MaxVectorsInProgress: 0,
					ExecutionTTL:         24 * time.Hour,
					ResultTTL:            7 * 24 * time.Hour,
					ProgressTTL:          time.Hour,
				},
			},
			wantErr: "executor max_vectors_in_progress must be at least 1",
		},
		{
			name: "invalid executor execution TTL",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           10,
					QueueSize:            100,
					MaxVectorsInProgress: 100,
					ExecutionTTL:         30 * time.Second,
					ResultTTL:            7 * 24 * time.Hour,
					ProgressTTL:          time.Hour,
				},
			},
			wantErr: "executor execution_ttl must be at least 1 minute",
		},
		{
			name: "invalid executor result TTL",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           10,
					QueueSize:            100,
					MaxVectorsInProgress: 100,
					ExecutionTTL:         24 * time.Hour,
					ResultTTL:            30 * time.Minute,
					ProgressTTL:          time.Hour,
				},
			},
			wantErr: "executor result_ttl must be at least 1 hour",
		},
		{
			name: "invalid executor progress TTL",
			config: Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
				JWTExpiry: 24 * time.Hour,
				TLS: TLSConfig{
					Enabled: false,
				},
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           10,
					QueueSize:            100,
					MaxVectorsInProgress: 100,
					ExecutionTTL:         24 * time.Hour,
					ResultTTL:            7 * 24 * time.Hour,
					ProgressTTL:          30 * time.Second,
				},
			},
			wantErr: "executor progress_ttl must be at least 1 minute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_ValidateTLSFiles(t *testing.T) {
	// Create temporary cert and key files for testing
	certFile, err := os.CreateTemp("", "cert-*.pem")
	require.NoError(t, err)
	defer func() { _ = os.Remove(certFile.Name()) }()
	_, _ = certFile.WriteString("fake cert")
	defer func() { _ = certFile.Close() }()

	keyFile, err := os.CreateTemp("", "key-*.pem")
	require.NoError(t, err)
	defer func() { _ = os.Remove(keyFile.Name()) }()
	_, _ = keyFile.WriteString("fake key")
	defer func() { _ = keyFile.Close() }()

	t.Run("TLS with existing files", func(t *testing.T) {
		config := Config{
			LisAddr:   "localhost:3030",
			HTTPPort:  ":8081",
			JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
			JWTExpiry: 24 * time.Hour,
			TLS: TLSConfig{
				Enabled:  true,
				CertFile: certFile.Name(),
				KeyFile:  keyFile.Name(),
			},
			RateLimit: RateLimitConfig{
				LoginRequestsPerMinute:    5,
				RegisterRequestsPerMinute: 3,
				DEExecutionsPerUser:       10,
				MaxConcurrentDEPerUser:    3,
				MaxRequestsPerSecond:      100,
				MaxMessageSizeBytes:       4 * 1024 * 1024,
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
			Executor: ExecutorConfig{
				MaxWorkers:           10,
				QueueSize:            100,
				MaxVectorsInProgress: 100,
				ExecutionTTL:         24 * time.Hour,
				ResultTTL:            7 * 24 * time.Hour,
				ProgressTTL:          time.Hour,
			},
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("TLS with non-existent cert file", func(t *testing.T) {
		config := Config{
			LisAddr:   "localhost:3030",
			HTTPPort:  ":8081",
			JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
			JWTExpiry: 24 * time.Hour,
			TLS: TLSConfig{
				Enabled:  true,
				CertFile: "/nonexistent/cert.pem",
				KeyFile:  keyFile.Name(),
			},
			RateLimit: RateLimitConfig{
				LoginRequestsPerMinute:    5,
				RegisterRequestsPerMinute: 3,
				DEExecutionsPerUser:       10,
				MaxConcurrentDEPerUser:    3,
				MaxRequestsPerSecond:      100,
				MaxMessageSizeBytes:       4 * 1024 * 1024,
			},
		}

		err := config.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cert file does not exist")
	})

	t.Run("TLS with non-existent key file", func(t *testing.T) {
		config := Config{
			LisAddr:   "localhost:3030",
			HTTPPort:  ":8081",
			JWTSecret: "this-is-a-very-secure-secret-with-more-than-32-characters",
			JWTExpiry: 24 * time.Hour,
			TLS: TLSConfig{
				Enabled:  true,
				CertFile: certFile.Name(),
				KeyFile:  "/nonexistent/key.pem",
			},
			RateLimit: RateLimitConfig{
				LoginRequestsPerMinute:    5,
				RegisterRequestsPerMinute: 3,
				DEExecutionsPerUser:       10,
				MaxConcurrentDEPerUser:    3,
				MaxRequestsPerSecond:      100,
				MaxMessageSizeBytes:       4 * 1024 * 1024,
			},
		}

		err := config.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "key file does not exist")
	})
}

func TestDefaultConfig(t *testing.T) {
	// Set JWT_SECRET for this test
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() { _ = os.Setenv("JWT_SECRET", originalSecret) }()

	_ = os.Setenv("JWT_SECRET", "test-secret-for-default-config-testing-123456")

	config := DefaultConfig()

	assert.Equal(t, "localhost:3030", config.LisAddr)
	assert.Equal(t, ":8081", config.HTTPPort)
	assert.Equal(t, "test-secret-for-default-config-testing-123456", config.JWTSecret)
	assert.Equal(t, 15*time.Minute, config.JWTExpiry)
	assert.False(t, config.TLS.Enabled)
	assert.Equal(t, 5, config.RateLimit.LoginRequestsPerMinute)
	assert.Equal(t, 10, config.RateLimit.DEExecutionsPerUser)
	assert.Equal(t, 3, config.RateLimit.MaxConcurrentDEPerUser)
	assert.Equal(t, 100, config.RateLimit.MaxRequestsPerSecond)
	assert.Equal(t, 4*1024*1024, config.RateLimit.MaxMessageSizeBytes)

	// Executor config defaults
	assert.Equal(t, 10, config.Executor.MaxWorkers)
	assert.Equal(t, 100, config.Executor.QueueSize)
	assert.Equal(t, 100, config.Executor.MaxVectorsInProgress)
	assert.Equal(t, 24*time.Hour, config.Executor.ExecutionTTL)
	assert.Equal(t, 7*24*time.Hour, config.Executor.ResultTTL)
	assert.Equal(t, 1*time.Hour, config.Executor.ProgressTTL)
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Set required JWT secret
	t.Setenv("JWT_SECRET", "test-jwt-secret-with-at-least-32-characters-long")

	cfg, err := LoadConfig("")
	require.NoError(t, err)

	// Verify default values
	assert.Equal(t, "localhost:3030", cfg.LisAddr)
	assert.Equal(t, ":8081", cfg.HTTPPort)
	assert.Equal(t, 15*time.Minute, cfg.JWTExpiry)
	assert.True(t, cfg.MetricsEnabled)
	assert.True(t, cfg.TracingEnabled)
	assert.True(t, cfg.SLOEnabled)
	assert.False(t, cfg.PprofEnabled)
	assert.Equal(t, "localhost:6060", cfg.PprofPort)

	// TLS defaults
	assert.False(t, cfg.TLS.Enabled)

	// Rate limit defaults
	assert.Equal(t, 5, cfg.RateLimit.LoginRequestsPerMinute)
	assert.Equal(t, 3, cfg.RateLimit.RegisterRequestsPerMinute)
	assert.Equal(t, 10, cfg.RateLimit.DEExecutionsPerUser)
	assert.Equal(t, 3, cfg.RateLimit.MaxConcurrentDEPerUser)
	assert.Equal(t, 100, cfg.RateLimit.MaxRequestsPerSecond)
	assert.Equal(t, 4*1024*1024, cfg.RateLimit.MaxMessageSizeBytes)

	// Redis defaults
	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Empty(t, cfg.Redis.Password)

	// Executor defaults
	assert.Equal(t, 10, cfg.Executor.MaxWorkers)
	assert.Equal(t, 100, cfg.Executor.QueueSize)
	assert.Equal(t, 100, cfg.Executor.MaxVectorsInProgress)
	assert.Equal(t, 24*time.Hour, cfg.Executor.ExecutionTTL)
	assert.Equal(t, 7*24*time.Hour, cfg.Executor.ResultTTL)
	assert.Equal(t, 1*time.Hour, cfg.Executor.ProgressTTL)

	// DE defaults
	assert.Equal(t, 100, cfg.DE.ParetoChannelLimiter)
	assert.Equal(t, 100, cfg.DE.MaxChannelLimiter)
	assert.Equal(t, 1000, cfg.DE.ResultLimiter)
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	t.Setenv("JWT_SECRET", "env-jwt-secret-at-least-32-chars-long-test")
	t.Setenv("GRPC_PORT", "localhost:4040")
	t.Setenv("HTTP_PORT", ":9091")
	t.Setenv("REDIS_HOST", "redis.example.com")
	t.Setenv("REDIS_PORT", "6380")
	t.Setenv("REDIS_PASSWORD", "secret-password")
	t.Setenv("REDIS_DB", "2")
	t.Setenv("METRICS_ENABLED", "false")
	t.Setenv("TRACING_ENABLED", "false")
	t.Setenv("SLO_ENABLED", "false")
	t.Setenv("PPROF_ENABLED", "true")
	t.Setenv("PPROF_PORT", ":7070")

	cfg, err := LoadConfig("")
	require.NoError(t, err)

	// Verify environment variables override defaults
	assert.Equal(t, "env-jwt-secret-at-least-32-chars-long-test", cfg.JWTSecret)
	assert.Equal(t, "localhost:4040", cfg.LisAddr)
	assert.Equal(t, ":9091", cfg.HTTPPort)
	assert.Equal(t, "redis.example.com", cfg.Redis.Host)
	assert.Equal(t, 6380, cfg.Redis.Port)
	assert.Equal(t, "secret-password", cfg.Redis.Password)
	assert.Equal(t, 2, cfg.Redis.DB)
	assert.False(t, cfg.MetricsEnabled)
	assert.False(t, cfg.TracingEnabled)
	assert.False(t, cfg.SLOEnabled)
	assert.True(t, cfg.PprofEnabled)
	assert.Equal(t, ":7070", cfg.PprofPort)
}

func TestLoadConfig_GODEPrefixedEnvVars(t *testing.T) {
	// Viper with GODE_ prefix should also work
	t.Setenv("JWT_SECRET", "base-jwt-secret-at-least-32-chars-test")
	t.Setenv("GODE_LIS_ADDR", "localhost:5050")
	t.Setenv("GODE_HTTP_PORT", ":8082")
	t.Setenv("GODE_METRICS_ENABLED", "true")

	cfg, err := LoadConfig("")
	require.NoError(t, err)

	// GODE_ prefixed vars should be read by Viper
	assert.Equal(t, "localhost:5050", cfg.LisAddr)
	assert.Equal(t, ":8082", cfg.HTTPPort)
	assert.True(t, cfg.MetricsEnabled)
}

func TestLoadConfig_ConfigFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
lis_addr: "localhost:7070"
http_port: ":8888"
jwt_secret: "file-jwt-secret-at-least-32-chars-long-test"
jwt_expiry: "48h"
metrics_enabled: false
tracing_enabled: true
slo_enabled: false
pprof_enabled: true
pprof_port: ":8080"

tls:
  enabled: false

rate_limit:
  login_requests_per_minute: 10
  register_requests_per_minute: 5
  de_executions_per_user: 20
  max_concurrent_de_per_user: 5
  max_requests_per_second: 200
  max_message_size_bytes: 8388608

redis:
  host: "redis-server"
  port: 6380
  password: "file-password"
  db: 1

executor:
  max_workers: 20
  queue_size: 200
  max_vectors_in_progress: 200
  execution_ttl: "48h"
  result_ttl: "336h"
  progress_ttl: "2h"

de:
  pareto_channel_limiter: 200
  max_channel_limiter: 200
  result_limiter: 2000
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)

	// Verify config file values
	assert.Equal(t, "localhost:7070", cfg.LisAddr)
	assert.Equal(t, ":8888", cfg.HTTPPort)
	assert.Equal(t, "file-jwt-secret-at-least-32-chars-long-test", cfg.JWTSecret)
	assert.Equal(t, 48*time.Hour, cfg.JWTExpiry)
	assert.False(t, cfg.MetricsEnabled)
	assert.True(t, cfg.TracingEnabled)
	assert.False(t, cfg.SLOEnabled)
	assert.True(t, cfg.PprofEnabled)
	assert.Equal(t, ":8080", cfg.PprofPort)

	// Rate limits
	assert.Equal(t, 10, cfg.RateLimit.LoginRequestsPerMinute)
	assert.Equal(t, 5, cfg.RateLimit.RegisterRequestsPerMinute)
	assert.Equal(t, 20, cfg.RateLimit.DEExecutionsPerUser)
	assert.Equal(t, 5, cfg.RateLimit.MaxConcurrentDEPerUser)
	assert.Equal(t, 200, cfg.RateLimit.MaxRequestsPerSecond)
	assert.Equal(t, 8388608, cfg.RateLimit.MaxMessageSizeBytes)

	// Redis
	assert.Equal(t, "redis-server", cfg.Redis.Host)
	assert.Equal(t, 6380, cfg.Redis.Port)
	assert.Equal(t, "file-password", cfg.Redis.Password)
	assert.Equal(t, 1, cfg.Redis.DB)

	// Executor
	assert.Equal(t, 20, cfg.Executor.MaxWorkers)
	assert.Equal(t, 200, cfg.Executor.QueueSize)
	assert.Equal(t, 200, cfg.Executor.MaxVectorsInProgress)
	assert.Equal(t, 48*time.Hour, cfg.Executor.ExecutionTTL)
	assert.Equal(t, 336*time.Hour, cfg.Executor.ResultTTL)
	assert.Equal(t, 2*time.Hour, cfg.Executor.ProgressTTL)

	// DE config
	assert.Equal(t, 200, cfg.DE.ParetoChannelLimiter)
	assert.Equal(t, 200, cfg.DE.MaxChannelLimiter)
	assert.Equal(t, 2000, cfg.DE.ResultLimiter)
}

func TestLoadConfig_EnvironmentOverridesFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
lis_addr: "localhost:7070"
jwt_secret: "file-jwt-secret-at-least-32-chars-long-test"
redis:
  host: "file-redis"
  port: 6380
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set environment variables that should override file
	t.Setenv("GRPC_PORT", "localhost:9999")
	t.Setenv("REDIS_HOST", "env-redis")

	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)

	// Environment should override file
	assert.Equal(t, "localhost:9999", cfg.LisAddr, "env var should override file")
	assert.Equal(t, "env-redis", cfg.Redis.Host, "env var should override file")
	// File value should be used when no env var
	assert.Equal(t, 6380, cfg.Redis.Port, "file value should be used when no env var")
}

func TestLoadConfig_InvalidConfigFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoadConfig_MetricsTypeEnum(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "prometheus (default)",
			envValue: "prometheus",
			expected: "prometheus",
		},
		{
			name:     "stdout",
			envValue: "stdout",
			expected: "stdout",
		},
		{
			name:     "invalid defaults to prometheus",
			envValue: "invalid",
			expected: "prometheus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", "test-jwt-secret-at-least-32-characters-test")
			t.Setenv("METRICS_TYPE", tt.envValue)

			cfg, err := LoadConfig("")
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(cfg.MetricsType))
		})
	}
}

func TestLoadConfig_TracingExporterEnum(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "none (default)",
			envValue: "none",
			expected: "none",
		},
		{
			name:     "stdout",
			envValue: "stdout",
			expected: "stdout",
		},
		{
			name:     "otlp",
			envValue: "otlp",
			expected: "otlp",
		},
		{
			name:     "file",
			envValue: "file",
			expected: "file",
		},
		{
			name:     "invalid defaults to none",
			envValue: "invalid",
			expected: "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", "test-jwt-secret-at-least-32-characters-test")
			t.Setenv("TRACING_EXPORTER", tt.envValue)

			cfg, err := LoadConfig("")
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(cfg.TracingConfig.ExporterType))
		})
	}
}

func TestLoadConfig_TracingConfig(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-jwt-secret-at-least-32-characters-test")
	t.Setenv("TRACING_EXPORTER", "otlp")
	t.Setenv("OTLP_ENDPOINT", "otel-collector:4317")
	t.Setenv("TRACE_SAMPLE_RATIO", "0.5")
	t.Setenv("TRACE_FILE_PATH", "custom-traces.json")

	cfg, err := LoadConfig("")
	require.NoError(t, err)

	assert.Equal(t, "otlp", string(cfg.TracingConfig.ExporterType))
	assert.Equal(t, "otel-collector:4317", cfg.TracingConfig.OTLPEndpoint)
	assert.Equal(t, 0.5, cfg.TracingConfig.SampleRatio)
	assert.Equal(t, "custom-traces.json", cfg.TracingConfig.FilePath)
}

func TestConfig_Validate_Redis(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		port      int
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid redis config",
			host:      "localhost",
			port:      6379,
			wantError: false,
		},
		{
			name:      "empty host",
			host:      "",
			port:      6379,
			wantError: true,
			errorMsg:  "redis host cannot be empty",
		},
		{
			name:      "port too low",
			host:      "localhost",
			port:      0,
			wantError: true,
			errorMsg:  "redis port must be between 1 and 65535",
		},
		{
			name:      "port too high",
			host:      "localhost",
			port:      70000,
			wantError: true,
			errorMsg:  "redis port must be between 1 and 65535",
		},
		{
			name:      "minimum valid port",
			host:      "localhost",
			port:      1,
			wantError: false,
		},
		{
			name:      "maximum valid port",
			host:      "localhost",
			port:      65535,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				LisAddr:   "localhost:3030",
				HTTPPort:  ":8081",
				JWTSecret: "valid-jwt-secret-at-least-32-characters-long",
				JWTExpiry: 24 * time.Hour,
				RateLimit: RateLimitConfig{
					LoginRequestsPerMinute:    5,
					RegisterRequestsPerMinute: 3,
					DEExecutionsPerUser:       10,
					MaxConcurrentDEPerUser:    3,
					MaxRequestsPerSecond:      100,
					MaxMessageSizeBytes:       4 * 1024 * 1024,
				},
				Redis: redis.Config{
					Host: tt.host,
					Port: tt.port,
				},
				Executor: ExecutorConfig{
					MaxWorkers:           10,
					QueueSize:            100,
					MaxVectorsInProgress: 100,
					ExecutionTTL:         24 * time.Hour,
					ResultTTL:            7 * 24 * time.Hour,
					ProgressTTL:          time.Hour,
				},
			}

			err := cfg.Validate()
			if tt.wantError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
