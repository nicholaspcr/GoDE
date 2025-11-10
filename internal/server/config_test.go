package server

import (
	"os"
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
					MaxWorkers:   10,
					QueueSize:    100,
					ExecutionTTL: 24 * time.Hour,
					ResultTTL:    7 * 24 * time.Hour,
					ProgressTTL:  time.Hour,
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
					DEExecutionsPerUser:    10,
					MaxConcurrentDEPerUser: 3,
					MaxRequestsPerSecond:   100,
					MaxMessageSizeBytes:    4 * 1024 * 1024,
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
					DEExecutionsPerUser:    0,
					MaxConcurrentDEPerUser: 3,
					MaxRequestsPerSecond:   100,
					MaxMessageSizeBytes:    4 * 1024 * 1024,
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
					DEExecutionsPerUser:    10,
					MaxConcurrentDEPerUser: 0,
					MaxRequestsPerSecond:   100,
					MaxMessageSizeBytes:    4 * 1024 * 1024,
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
					DEExecutionsPerUser:    10,
					MaxConcurrentDEPerUser: 3,
					MaxRequestsPerSecond:   0,
					MaxMessageSizeBytes:    4 * 1024 * 1024,
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
					DEExecutionsPerUser:    10,
					MaxConcurrentDEPerUser: 3,
					MaxRequestsPerSecond:   100,
					MaxMessageSizeBytes:    512,
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
				DEExecutionsPerUser:    10,
				MaxConcurrentDEPerUser: 3,
				MaxRequestsPerSecond:   100,
				MaxMessageSizeBytes:    4 * 1024 * 1024,
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
			Executor: ExecutorConfig{
				MaxWorkers:   10,
				QueueSize:    100,
				ExecutionTTL: 24 * time.Hour,
				ResultTTL:    7 * 24 * time.Hour,
				ProgressTTL:  time.Hour,
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
				DEExecutionsPerUser:    10,
				MaxConcurrentDEPerUser: 3,
				MaxRequestsPerSecond:   100,
				MaxMessageSizeBytes:    4 * 1024 * 1024,
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
				DEExecutionsPerUser:    10,
				MaxConcurrentDEPerUser: 3,
				MaxRequestsPerSecond:   100,
				MaxMessageSizeBytes:    4 * 1024 * 1024,
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
	assert.Equal(t, 24*time.Hour, config.JWTExpiry)
	assert.False(t, config.TLS.Enabled)
	assert.Equal(t, 5, config.RateLimit.LoginRequestsPerMinute)
	assert.Equal(t, 10, config.RateLimit.DEExecutionsPerUser)
	assert.Equal(t, 3, config.RateLimit.MaxConcurrentDEPerUser)
	assert.Equal(t, 100, config.RateLimit.MaxRequestsPerSecond)
	assert.Equal(t, 4*1024*1024, config.RateLimit.MaxMessageSizeBytes)
}
