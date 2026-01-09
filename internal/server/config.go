package server

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
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
	Redis          redis.Config
	Executor       ExecutorConfig
	MetricsEnabled bool
	MetricsType    telemetry.MetricsExporterType
	TracingEnabled bool
	TracingConfig  telemetry.TracingConfig
	SLOEnabled     bool
	PprofEnabled   bool
	PprofPort      string
	CORS           middleware.CORSConfig
}

// ExecutorConfig contains configuration for the background execution executor.
type ExecutorConfig struct {
	MaxWorkers    int
	QueueSize     int
	ExecutionTTL  time.Duration
	ResultTTL     time.Duration
	ProgressTTL   time.Duration
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

	tracingEnabled := os.Getenv("TRACING_ENABLED") != "false" // Enabled by default
	tracingExporterType := telemetry.TracingExporterNone      // Default to none to avoid stdout noise
	switch os.Getenv("TRACING_EXPORTER") {
	case "stdout":
		tracingExporterType = telemetry.TracingExporterStdout
	case "otlp":
		tracingExporterType = telemetry.TracingExporterOTLP
	case "file":
		tracingExporterType = telemetry.TracingExporterFile
	}

	traceFilePath := os.Getenv("TRACE_FILE_PATH")
	if traceFilePath == "" {
		traceFilePath = "traces.json"
	}

	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317" // Default Jaeger/OTLP endpoint
	}

	sampleRatio := 1.0 // Sample all traces by default in development
	if ratio := os.Getenv("TRACE_SAMPLE_RATIO"); ratio != "" {
		if parsed, err := strconv.ParseFloat(ratio, 64); err == nil {
			sampleRatio = parsed
		}
	}

	pprofEnabled := os.Getenv("PPROF_ENABLED") == "true"
	pprofPort := os.Getenv("PPROF_PORT")
	if pprofPort == "" {
		pprofPort = ":6060" // Default pprof port
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := 6379
	if portStr := os.Getenv("REDIS_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			redisPort = port
		}
	}

	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			redisDB = db
		}
	}

	sloEnabled := os.Getenv("SLO_ENABLED") != "false" // Enabled by default

	return Config{
		LisAddr:        "localhost:3030",
		HTTPPort:       ":8081",
		JWTSecret:      os.Getenv("JWT_SECRET"), // No default - must be set via env var
		JWTExpiry:      24 * time.Hour,
		MetricsEnabled: true, // Metrics enabled by default
		MetricsType:    metricsType,
		TracingEnabled: tracingEnabled,
		TracingConfig: telemetry.TracingConfig{
			ExporterType: tracingExporterType,
			OTLPEndpoint: otlpEndpoint,
			SampleRatio:  sampleRatio,
			FilePath:     traceFilePath,
		},
		SLOEnabled:     sloEnabled,
		PprofEnabled:   pprofEnabled,
		PprofPort:      pprofPort,
		Redis: redis.Config{
			Host:     redisHost,
			Port:     redisPort,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		Executor: ExecutorConfig{
			MaxWorkers:    10,
			QueueSize:     100,
			ExecutionTTL:  24 * time.Hour,
			ResultTTL:     7 * 24 * time.Hour,
			ProgressTTL:   1 * time.Hour,
		},
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
		CORS: middleware.DefaultCORSConfig(),
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

	// Redis validation
	if c.Redis.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}
	if c.Redis.Port < 1 || c.Redis.Port > 65535 {
		return fmt.Errorf("redis port must be between 1 and 65535")
	}

	// Executor validation
	if c.Executor.MaxWorkers < 1 {
		return fmt.Errorf("executor max_workers must be at least 1")
	}
	if c.Executor.QueueSize < 1 {
		return fmt.Errorf("executor queue_size must be at least 1")
	}
	if c.Executor.ExecutionTTL < time.Minute {
		return fmt.Errorf("executor execution_ttl must be at least 1 minute")
	}
	if c.Executor.ResultTTL < time.Hour {
		return fmt.Errorf("executor result_ttl must be at least 1 hour")
	}
	if c.Executor.ProgressTTL < time.Minute {
		return fmt.Errorf("executor progress_ttl must be at least 1 minute")
	}

	return nil
}
