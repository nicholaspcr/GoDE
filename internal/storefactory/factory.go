package storefactory

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/glebarez/sqlite"
	// pgx driver for PostgreSQL
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/migrations"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/internal/store/composite"
	"github.com/nicholaspcr/GoDE/internal/store/gorm"
	redisstore "github.com/nicholaspcr/GoDE/internal/store/redis"
	"gorm.io/driver/postgres"
)

// Config extends store.Config with Redis and TTL settings.
type Config struct {
	store.Config `json:",inline" yaml:",inline" mapstructure:",squash"`
	Redis        redis.Config  `json:"redis" yaml:"redis" mapstructure:"redis"`
	ExecutionTTL time.Duration `json:"execution_ttl" yaml:"execution_ttl" mapstructure:"execution_ttl"`
	ResultTTL    time.Duration `json:"result_ttl" yaml:"result_ttl" mapstructure:"result_ttl"`
	ProgressTTL  time.Duration `json:"progress_ttl" yaml:"progress_ttl" mapstructure:"progress_ttl"`
}

// New returns a new Store instance that combines database and Redis.
func New(ctx context.Context, cfg Config) (store.Store, error) {
	// Run SQL migrations for PostgreSQL only (GORM AutoMigrate handles SQLite)
	if cfg.Type == "postgres" {
		connStr := cfg.ConnectionString()
		if connStr != "" {
			slog.Info("Running database migrations before connecting...")
			if err := migrations.Run(connStr); err != nil {
				slog.Error("Migration failed", slog.String("error", err.Error()))
				return nil, err
			}
		}
	}

	var dialector gorm.Dialector

	switch cfg.Type {
	case "memory":
		dialector = sqlite.Open(":memory:")
	case "sqlite":
		dialector = sqlite.Open(cfg.Sqlite.Filepath)
	case "postgres":
		dialector = postgres.Open(cfg.Postgresql.DNS)
	default:
		return nil, errors.New("invalid store type")
	}

	dbStore, err := gorm.New(dialector)
	if err != nil {
		return nil, err
	}

	// Run AutoMigrate only for SQLite/memory stores
	// PostgreSQL schema is fully managed by SQL migrations
	if cfg.Type != "postgres" {
		if err := dbStore.AutoMigrate(); err != nil {
			return nil, err
		}
	}

	// Initialize Redis client
	slog.Info("Connecting to Redis",
		slog.String("host", cfg.Redis.Host),
		slog.Int("port", cfg.Redis.Port))

	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		return nil, err
	}

	slog.Info("Redis connection established")

	// Create Redis execution store
	redisExecStore := redisstore.NewExecutionStore(
		redisClient,
		cfg.ExecutionTTL,
		cfg.ProgressTTL,
	)

	// Create and return composite store
	compositeStore := composite.New(dbStore, redisClient, redisExecStore)

	slog.Info("Composite store initialized with database and Redis")

	return compositeStore, nil
}
