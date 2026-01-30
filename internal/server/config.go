package server

import (
	"fmt"
	"os"
	"time"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/spf13/viper"
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
	MaxWorkers           int
	QueueSize            int
	MaxVectorsInProgress int // Maximum vectors to include in progress updates (default: 100)
	ExecutionTTL         time.Duration
	ResultTTL            time.Duration
	ProgressTTL          time.Duration
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

// LoadConfig loads configuration from environment variables and optional config file.
// If configPath is empty, it will look for config.yaml in common locations.
// Environment variables take precedence over config file values.
func LoadConfig(configPath string) (Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read config file if specified
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return Config{}, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Search for config in common locations
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/gode")
		// Ignore error if config file doesn't exist
		_ = v.ReadInConfig()
	}

	// Environment variables override config file
	v.SetEnvPrefix("GODE")
	v.AutomaticEnv()

	// Bind specific environment variables (for backward compatibility)
	bindEnvVars(v)

	// Build config struct
	cfg := Config{
		LisAddr:        v.GetString("lis_addr"),
		HTTPPort:       v.GetString("http_port"),
		JWTSecret:      v.GetString("jwt_secret"),
		JWTExpiry:      v.GetDuration("jwt_expiry"),
		MetricsEnabled: v.GetBool("metrics_enabled"),
		TracingEnabled: v.GetBool("tracing_enabled"),
		SLOEnabled:     v.GetBool("slo_enabled"),
		PprofEnabled:   v.GetBool("pprof_enabled"),
		PprofPort:      v.GetString("pprof_port"),
		TLS: TLSConfig{
			Enabled:  v.GetBool("tls.enabled"),
			CertFile: v.GetString("tls.cert_file"),
			KeyFile:  v.GetString("tls.key_file"),
		},
		RateLimit: RateLimitConfig{
			LoginRequestsPerMinute:    v.GetInt("rate_limit.login_requests_per_minute"),
			RegisterRequestsPerMinute: v.GetInt("rate_limit.register_requests_per_minute"),
			DEExecutionsPerUser:       v.GetInt("rate_limit.de_executions_per_user"),
			MaxConcurrentDEPerUser:    v.GetInt("rate_limit.max_concurrent_de_per_user"),
			MaxRequestsPerSecond:      v.GetInt("rate_limit.max_requests_per_second"),
			MaxMessageSizeBytes:       v.GetInt("rate_limit.max_message_size_bytes"),
		},
		Redis: redis.Config{
			Host:     v.GetString("redis.host"),
			Port:     v.GetInt("redis.port"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		Executor: ExecutorConfig{
			MaxWorkers:           v.GetInt("executor.max_workers"),
			QueueSize:            v.GetInt("executor.queue_size"),
			MaxVectorsInProgress: v.GetInt("executor.max_vectors_in_progress"),
			ExecutionTTL:         v.GetDuration("executor.execution_ttl"),
			ResultTTL:            v.GetDuration("executor.result_ttl"),
			ProgressTTL:          v.GetDuration("executor.progress_ttl"),
		},
		DE: de.Config{
			ParetoChannelLimiter: v.GetInt("de.pareto_channel_limiter"),
			MaxChannelLimiter:    v.GetInt("de.max_channel_limiter"),
			ResultLimiter:        v.GetInt("de.result_limiter"),
		},
		CORS: middleware.DefaultCORSConfig(),
	}

	// Handle metrics type enum
	if v.GetString("metrics_type") == "stdout" {
		cfg.MetricsType = telemetry.MetricsExporterStdout
	} else {
		cfg.MetricsType = telemetry.MetricsExporterPrometheus
	}

	// Handle tracing configuration
	cfg.TracingConfig = telemetry.TracingConfig{
		OTLPEndpoint: v.GetString("tracing.otlp_endpoint"),
		SampleRatio:  v.GetFloat64("tracing.sample_ratio"),
		FilePath:     v.GetString("tracing.file_path"),
	}

	switch v.GetString("tracing.exporter") {
	case "stdout":
		cfg.TracingConfig.ExporterType = telemetry.TracingExporterStdout
	case "otlp":
		cfg.TracingConfig.ExporterType = telemetry.TracingExporterOTLP
	case "file":
		cfg.TracingConfig.ExporterType = telemetry.TracingExporterFile
	default:
		cfg.TracingConfig.ExporterType = telemetry.TracingExporterNone
	}

	return cfg, nil
}

// setDefaults sets default values for all configuration options.
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("lis_addr", "localhost:3030")
	v.SetDefault("http_port", ":8081")
	v.SetDefault("jwt_expiry", 24*time.Hour)
	v.SetDefault("metrics_enabled", true)
	v.SetDefault("metrics_type", "prometheus")
	v.SetDefault("tracing_enabled", true)
	v.SetDefault("slo_enabled", true)
	v.SetDefault("pprof_enabled", false)
	v.SetDefault("pprof_port", ":6060")

	// TLS defaults
	v.SetDefault("tls.enabled", false)

	// Rate limit defaults
	v.SetDefault("rate_limit.login_requests_per_minute", 5)
	v.SetDefault("rate_limit.register_requests_per_minute", 3)
	v.SetDefault("rate_limit.de_executions_per_user", 10)
	v.SetDefault("rate_limit.max_concurrent_de_per_user", 3)
	v.SetDefault("rate_limit.max_requests_per_second", 100)
	v.SetDefault("rate_limit.max_message_size_bytes", 4*1024*1024)

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// Executor defaults
	v.SetDefault("executor.max_workers", 10)
	v.SetDefault("executor.queue_size", 100)
	v.SetDefault("executor.max_vectors_in_progress", 100)
	v.SetDefault("executor.execution_ttl", 24*time.Hour)
	v.SetDefault("executor.result_ttl", 7*24*time.Hour)
	v.SetDefault("executor.progress_ttl", 1*time.Hour)

	// DE algorithm defaults
	v.SetDefault("de.pareto_channel_limiter", 100)
	v.SetDefault("de.max_channel_limiter", 100)
	v.SetDefault("de.result_limiter", 1000)

	// Tracing defaults
	v.SetDefault("tracing.exporter", "none")
	v.SetDefault("tracing.otlp_endpoint", "localhost:4317")
	v.SetDefault("tracing.sample_ratio", 1.0)
	v.SetDefault("tracing.file_path", "traces.json")
}

// bindEnvVars binds specific environment variables for backward compatibility.
func bindEnvVars(v *viper.Viper) {
	// Direct environment variable names (without GODE_ prefix) for backward compatibility
	envBindings := map[string]string{
		"JWT_SECRET":         "jwt_secret",
		"GRPC_PORT":          "lis_addr",
		"HTTP_PORT":          "http_port",
		"METRICS_TYPE":       "metrics_type",
		"METRICS_ENABLED":    "metrics_enabled",
		"TRACING_ENABLED":    "tracing_enabled",
		"TRACING_EXPORTER":   "tracing.exporter",
		"TRACE_FILE_PATH":    "tracing.file_path",
		"TRACE_SAMPLE_RATIO": "tracing.sample_ratio",
		"OTLP_ENDPOINT":      "tracing.otlp_endpoint",
		"SLO_ENABLED":        "slo_enabled",
		"PPROF_ENABLED":      "pprof_enabled",
		"PPROF_PORT":         "pprof_port",
		"REDIS_HOST":         "redis.host",
		"REDIS_PORT":         "redis.port",
		"REDIS_PASSWORD":     "redis.password",
		"REDIS_DB":           "redis.db",
	}

	for envVar, configKey := range envBindings {
		_ = v.BindEnv(configKey, envVar)
	}
}

// DefaultConfig returns the default configuration of the server.
// For production use, prefer LoadConfig() which supports config files and environment variables.
func DefaultConfig() Config {
	cfg, err := LoadConfig("")
	if err != nil {
		// Fallback to basic defaults if config loading fails
		return fallbackConfig()
	}
	return cfg
}

// fallbackConfig provides minimal defaults when config loading fails.
func fallbackConfig() Config {
	return Config{
		LisAddr:        "localhost:3030",
		HTTPPort:       ":8081",
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiry:      24 * time.Hour,
		MetricsEnabled: true,
		MetricsType:    telemetry.MetricsExporterPrometheus,
		TracingEnabled: true,
		TracingConfig: telemetry.TracingConfig{
			ExporterType: telemetry.TracingExporterNone,
			OTLPEndpoint: "localhost:4317",
			SampleRatio:  1.0,
			FilePath:     "traces.json",
		},
		SLOEnabled:   true,
		PprofEnabled: false,
		PprofPort:    ":6060",
		Redis: redis.Config{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
		Executor: ExecutorConfig{
			MaxWorkers:           10,
			QueueSize:            100,
			MaxVectorsInProgress: 100,
			ExecutionTTL:         24 * time.Hour,
			ResultTTL:            7 * 24 * time.Hour,
			ProgressTTL:          1 * time.Hour,
		},
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
	if c.Executor.MaxVectorsInProgress < 1 {
		return fmt.Errorf("executor max_vectors_in_progress must be at least 1")
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
