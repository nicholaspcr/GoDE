package server

import (
	"fmt"
	"os"
	"time"

	"github.com/nicholaspcr/GoDE/internal/telemetry"
	"github.com/nicholaspcr/GoDE/pkg/de"
)

// Config contains all the necessary configuration options for the server.
type Config struct {
	LisAddr        string
	HTTPPort       string
	JWTSecret      string
	DE             de.Config
	JWTExpiry      time.Duration
	TLS            TLSConfig
	RateLimit      RateLimitConfig
	MetricsEnabled bool
	MetricsType    telemetry.MetricsExporterType
	PprofEnabled   bool
	PprofPort      string
}

// TLSConfig contains TLS/HTTPS configuration.
type TLSConfig struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	CertFile string `json:"cert_file" yaml:"cert_file"`
	KeyFile  string `json:"key_file" yaml:"key_file"`
}

// RateLimitConfig contains rate limiting configuration.
type RateLimitConfig struct {
	// LoginRequestsPerMinute limits login requests per IP
	LoginRequestsPerMinute int `json:"login_requests_per_minute" yaml:"login_requests_per_minute"`
	// RegisterRequestsPerMinute limits registration requests per IP (stricter than login)
	RegisterRequestsPerMinute int `json:"register_requests_per_minute" yaml:"register_requests_per_minute"`
	// DEExecutionsPerUser limits DE executions per user
	DEExecutionsPerUser int `json:"de_executions_per_user" yaml:"de_executions_per_user"`
	// MaxConcurrentDEPerUser limits concurrent DE executions per user
	MaxConcurrentDEPerUser int `json:"max_concurrent_de_per_user" yaml:"max_concurrent_de_per_user"`
	// MaxRequestsPerSecond is the global API rate limit per second
	MaxRequestsPerSecond int `json:"max_requests_per_second" yaml:"max_requests_per_second"`
	// MaxMessageSizeBytes is the maximum gRPC message size
	MaxMessageSizeBytes int `json:"max_message_size_bytes" yaml:"max_message_size_bytes"`
}

// DefaultConfig returns the default configuration of the server.
func DefaultConfig() Config {
	metricsType := telemetry.MetricsExporterPrometheus
	if os.Getenv("METRICS_TYPE") == "stdout" {
		metricsType = telemetry.MetricsExporterStdout
	}

	pprofEnabled := os.Getenv("PPROF_ENABLED") == "true"
	pprofPort := os.Getenv("PPROF_PORT")
	if pprofPort == "" {
		pprofPort = ":6060" // Default pprof port
	}

	return Config{
		LisAddr:        "localhost:3030",
		HTTPPort:       ":8081",
		JWTSecret:      os.Getenv("JWT_SECRET"), // No default - must be set via env var
		JWTExpiry:      24 * time.Hour,
		MetricsEnabled: true, // Metrics enabled by default
		MetricsType:    metricsType,
		PprofEnabled:   pprofEnabled,
		PprofPort:      pprofPort,
		TLS: TLSConfig{
			Enabled:  false, // TLS disabled by default for development
			CertFile: "",
			KeyFile:  "",
		},
		RateLimit: RateLimitConfig{
			LoginRequestsPerMinute:    5,  // 5 login attempts per minute per IP
			RegisterRequestsPerMinute: 3,  // 3 registrations per minute per IP (stricter)
			DEExecutionsPerUser:       10,
			MaxConcurrentDEPerUser:    3,
			MaxRequestsPerSecond:      100,
			MaxMessageSizeBytes:       4 * 1024 * 1024, // 4MB
		},
		DE: de.Config{
			ParetoChannelLimiter: 100,
			MaxChannelLimiter:    100,
			ResultLimiter:        1000,
		},
	}
}

// Validate validates the server configuration and returns an error if invalid.
func (c *Config) Validate() error {
	// JWT Secret validation
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required and must not be empty")
	}
	if c.JWTSecret == "change-me-in-production" {
		return fmt.Errorf("JWT_SECRET is set to the insecure default value 'change-me-in-production' - please use a secure random secret")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long for security")
	}

	// TLS validation
	if c.TLS.Enabled {
		if c.TLS.CertFile == "" {
			return fmt.Errorf("TLS is enabled but cert_file is not specified")
		}
		if c.TLS.KeyFile == "" {
			return fmt.Errorf("TLS is enabled but key_file is not specified")
		}
		// Check if files exist
		if _, err := os.Stat(c.TLS.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS cert file does not exist: %s", c.TLS.CertFile)
		}
		if _, err := os.Stat(c.TLS.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS key file does not exist: %s", c.TLS.KeyFile)
		}
	}

	// Rate limit validation
	if c.RateLimit.LoginRequestsPerMinute < 1 {
		return fmt.Errorf("login_requests_per_minute must be at least 1")
	}
	if c.RateLimit.RegisterRequestsPerMinute < 1 {
		return fmt.Errorf("register_requests_per_minute must be at least 1")
	}
	if c.RateLimit.DEExecutionsPerUser < 1 {
		return fmt.Errorf("de_executions_per_user must be at least 1")
	}
	if c.RateLimit.MaxConcurrentDEPerUser < 1 {
		return fmt.Errorf("max_concurrent_de_per_user must be at least 1")
	}
	if c.RateLimit.MaxRequestsPerSecond < 1 {
		return fmt.Errorf("max_requests_per_second must be at least 1")
	}
	if c.RateLimit.MaxMessageSizeBytes < 1024 {
		return fmt.Errorf("max_message_size_bytes must be at least 1024 bytes")
	}

	// Port validation
	if c.LisAddr == "" {
		return fmt.Errorf("listen address cannot be empty")
	}
	if c.HTTPPort == "" {
		return fmt.Errorf("HTTP port cannot be empty")
	}

	// JWT expiry validation
	if c.JWTExpiry < time.Minute {
		return fmt.Errorf("JWT expiry must be at least 1 minute")
	}

	return nil
}
