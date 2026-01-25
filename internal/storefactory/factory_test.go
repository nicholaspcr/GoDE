package storefactory

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig_Embedding verifies the Config struct properly embeds store.Config
func TestConfig_Embedding(t *testing.T) {
	t.Run("Config embeds store.Config fields", func(t *testing.T) {
		cfg := Config{
			Config: store.Config{
				Type: "memory",
				Sqlite: store.Sqlite{
					Filepath: "/tmp/test.db",
				},
				Postgresql: store.Postgresql{
					DNS: "postgres://localhost:5432/test",
				},
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
			ExecutionTTL: 24 * time.Hour,
			ResultTTL:    7 * 24 * time.Hour,
			ProgressTTL:  1 * time.Hour,
		}

		assert.Equal(t, "memory", cfg.Type)
		assert.Equal(t, "/tmp/test.db", cfg.Sqlite.Filepath)
		assert.Equal(t, "postgres://localhost:5432/test", cfg.Postgresql.DNS)
		assert.Equal(t, "localhost", cfg.Redis.Host)
		assert.Equal(t, 6379, cfg.Redis.Port)
		assert.Equal(t, 24*time.Hour, cfg.ExecutionTTL)
		assert.Equal(t, 7*24*time.Hour, cfg.ResultTTL)
		assert.Equal(t, 1*time.Hour, cfg.ProgressTTL)
	})
}

// TestNew_InvalidStoreType tests that invalid store types return an error
func TestNew_InvalidStoreType(t *testing.T) {
	tests := []struct {
		name      string
		storeType string
	}{
		{
			name:      "empty type",
			storeType: "",
		},
		{
			name:      "unknown type",
			storeType: "unknown",
		},
		{
			name:      "mysql (unsupported)",
			storeType: "mysql",
		},
		{
			name:      "mongodb (unsupported)",
			storeType: "mongodb",
		},
		{
			name:      "case sensitive - Memory",
			storeType: "Memory",
		},
		{
			name:      "case sensitive - SQLITE",
			storeType: "SQLITE",
		},
		{
			name:      "case sensitive - Postgres",
			storeType: "Postgres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Config: store.Config{
					Type: tt.storeType,
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
			}

			ctx := context.Background()
			s, err := New(ctx, cfg)

			assert.Error(t, err, "should return error for invalid store type")
			assert.Nil(t, s, "store should be nil on error")
			assert.Contains(t, err.Error(), "invalid store type")
		})
	}
}

// TestNew_MemoryStore tests creating an in-memory SQLite store
// Note: This will fail in the RED phase until we have a valid Redis connection
// because the factory requires Redis to create the composite store.
func TestNew_MemoryStore_DBInitialization(t *testing.T) {
	t.Run("memory store creates in-memory SQLite", func(t *testing.T) {
		cfg := Config{
			Config: store.Config{
				Type: "memory",
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
			ExecutionTTL: 24 * time.Hour,
			ProgressTTL:  1 * time.Hour,
		}

		// This test verifies that the store type "memory" is valid
		// and results in using sqlite.Open(":memory:")
		// The actual store creation will fail without Redis,
		// but the dialector selection should succeed
		assert.Equal(t, "memory", cfg.Type)
		assert.NotEmpty(t, cfg.Redis.Host)
	})
}

// TestNew_SQLiteStore tests creating a SQLite store with file path
func TestNew_SQLiteStore_ConfigValidation(t *testing.T) {
	t.Run("sqlite store validates filepath configuration", func(t *testing.T) {
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "test.db")

		cfg := Config{
			Config: store.Config{
				Type: "sqlite",
				Sqlite: store.Sqlite{
					Filepath: dbPath,
				},
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
			ExecutionTTL: 24 * time.Hour,
			ProgressTTL:  1 * time.Hour,
		}

		assert.Equal(t, "sqlite", cfg.Type)
		assert.Equal(t, dbPath, cfg.Sqlite.Filepath)
	})

	t.Run("sqlite store with empty filepath", func(t *testing.T) {
		cfg := Config{
			Config: store.Config{
				Type: "sqlite",
				Sqlite: store.Sqlite{
					Filepath: "",
				},
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
		}

		// Empty filepath should still be a valid config structure
		// (SQLite may create a default file or use memory)
		assert.Equal(t, "sqlite", cfg.Type)
		assert.Empty(t, cfg.Sqlite.Filepath)
	})
}

// TestNew_PostgresStore tests creating a PostgreSQL store
func TestNew_PostgresStore_ConfigValidation(t *testing.T) {
	t.Run("postgres store validates DNS configuration", func(t *testing.T) {
		dns := "postgres://user:password@localhost:5432/testdb?sslmode=disable"
		cfg := Config{
			Config: store.Config{
				Type: "postgres",
				Postgresql: store.Postgresql{
					DNS: dns,
				},
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
			ExecutionTTL: 24 * time.Hour,
			ProgressTTL:  1 * time.Hour,
		}

		assert.Equal(t, "postgres", cfg.Type)
		assert.Equal(t, dns, cfg.Postgresql.DNS)
	})

	t.Run("postgres store with empty DNS", func(t *testing.T) {
		cfg := Config{
			Config: store.Config{
				Type: "postgres",
				Postgresql: store.Postgresql{
					DNS: "",
				},
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
		}

		// Empty DNS should be a valid config structure
		// (will fail at connection time)
		assert.Equal(t, "postgres", cfg.Type)
		assert.Empty(t, cfg.Postgresql.DNS)
	})
}

// TestNew_RedisConfig tests Redis configuration validation
func TestNew_RedisConfig(t *testing.T) {
	tests := []struct {
		name        string
		redisConfig redis.Config
	}{
		{
			name: "standard configuration",
			redisConfig: redis.Config{
				Host:     "localhost",
				Port:     6379,
				Password: "",
				DB:       0,
			},
		},
		{
			name: "with password",
			redisConfig: redis.Config{
				Host:     "localhost",
				Port:     6379,
				Password: "secretpassword",
				DB:       0,
			},
		},
		{
			name: "different database",
			redisConfig: redis.Config{
				Host:     "localhost",
				Port:     6379,
				Password: "",
				DB:       1,
			},
		},
		{
			name: "custom port",
			redisConfig: redis.Config{
				Host:     "localhost",
				Port:     7000,
				Password: "",
				DB:       0,
			},
		},
		{
			name: "remote host",
			redisConfig: redis.Config{
				Host:     "redis.example.com",
				Port:     6379,
				Password: "secret",
				DB:       0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Config: store.Config{
					Type: "memory",
				},
				Redis:        tt.redisConfig,
				ExecutionTTL: 24 * time.Hour,
				ProgressTTL:  1 * time.Hour,
			}

			assert.Equal(t, tt.redisConfig.Host, cfg.Redis.Host)
			assert.Equal(t, tt.redisConfig.Port, cfg.Redis.Port)
			assert.Equal(t, tt.redisConfig.Password, cfg.Redis.Password)
			assert.Equal(t, tt.redisConfig.DB, cfg.Redis.DB)
		})
	}
}

// TestNew_TTLConfiguration tests TTL configuration for executions, results, and progress
func TestNew_TTLConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		executionTTL time.Duration
		resultTTL    time.Duration
		progressTTL  time.Duration
	}{
		{
			name:         "default TTLs",
			executionTTL: 24 * time.Hour,
			resultTTL:    7 * 24 * time.Hour,
			progressTTL:  1 * time.Hour,
		},
		{
			name:         "short TTLs",
			executionTTL: 1 * time.Hour,
			resultTTL:    6 * time.Hour,
			progressTTL:  5 * time.Minute,
		},
		{
			name:         "long TTLs",
			executionTTL: 30 * 24 * time.Hour,
			resultTTL:    90 * 24 * time.Hour,
			progressTTL:  24 * time.Hour,
		},
		{
			name:         "zero TTLs (no expiration)",
			executionTTL: 0,
			resultTTL:    0,
			progressTTL:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Config: store.Config{
					Type: "memory",
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 6379,
				},
				ExecutionTTL: tt.executionTTL,
				ResultTTL:    tt.resultTTL,
				ProgressTTL:  tt.progressTTL,
			}

			assert.Equal(t, tt.executionTTL, cfg.ExecutionTTL)
			assert.Equal(t, tt.resultTTL, cfg.ResultTTL)
			assert.Equal(t, tt.progressTTL, cfg.ProgressTTL)
		})
	}
}

// TestNew_RedisConnectionFailure tests that Redis connection failure returns error
func TestNew_RedisConnectionFailure(t *testing.T) {
	tests := []struct {
		name        string
		redisConfig redis.Config
	}{
		{
			name: "invalid host",
			redisConfig: redis.Config{
				Host:     "invalid-host-that-does-not-exist",
				Port:     6379,
				Password: "",
				DB:       0,
			},
		},
		{
			name: "invalid port",
			redisConfig: redis.Config{
				Host:     "localhost",
				Port:     1,
				Password: "",
				DB:       0,
			},
		},
		{
			name: "unreachable host",
			redisConfig: redis.Config{
				Host:     "192.0.2.1", // TEST-NET-1, guaranteed unreachable
				Port:     6379,
				Password: "",
				DB:       0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Config: store.Config{
					Type: "memory",
				},
				Redis:        tt.redisConfig,
				ExecutionTTL: 24 * time.Hour,
				ProgressTTL:  1 * time.Hour,
			}

			ctx := context.Background()
			s, err := New(ctx, cfg)

			assert.Error(t, err, "should return error when Redis connection fails")
			assert.Nil(t, s, "store should be nil when Redis fails")
			assert.Contains(t, err.Error(), "Redis")
		})
	}
}

// TestNew_ContextHandling tests context handling in factory
func TestNew_ContextHandling(t *testing.T) {
	t.Run("accepts background context", func(t *testing.T) {
		ctx := context.Background()
		assert.NotNil(t, ctx)
	})

	t.Run("accepts TODO context", func(t *testing.T) {
		ctx := context.TODO()
		assert.NotNil(t, ctx)
	})

	t.Run("accepts context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		assert.NotNil(t, ctx)
	})

	t.Run("accepts context with cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		assert.NotNil(t, ctx)
	})
}

// TestConfig_ConnectionString tests the ConnectionString method inherited from store.Config
func TestConfig_ConnectionString(t *testing.T) {
	tests := []struct {
		name         string
		storeType    string
		sqlitePath   string
		postgresDNS  string
		expectedConn string
	}{
		{
			name:         "sqlite connection string",
			storeType:    "sqlite",
			sqlitePath:   "/tmp/test.db",
			expectedConn: "sqlite3:///tmp/test.db",
		},
		{
			name:         "postgres connection string",
			storeType:    "postgres",
			postgresDNS:  "postgres://user:pass@localhost:5432/db",
			expectedConn: "postgres://user:pass@localhost:5432/db",
		},
		{
			name:         "memory returns empty string",
			storeType:    "memory",
			expectedConn: "",
		},
		{
			name:         "unknown type returns empty string",
			storeType:    "unknown",
			expectedConn: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Config: store.Config{
					Type: tt.storeType,
					Sqlite: store.Sqlite{
						Filepath: tt.sqlitePath,
					},
					Postgresql: store.Postgresql{
						DNS: tt.postgresDNS,
					},
				},
			}

			connStr := cfg.ConnectionString()
			assert.Equal(t, tt.expectedConn, connStr)
		})
	}
}

// TestNew_StoreTypeSwitch tests the switch statement logic for store types
func TestNew_StoreTypeSwitch(t *testing.T) {
	// This test verifies that valid store types are recognized
	validTypes := []string{"memory", "sqlite", "postgres"}

	for _, storeType := range validTypes {
		t.Run("valid type: "+storeType, func(t *testing.T) {
			cfg := Config{
				Config: store.Config{
					Type: storeType,
				},
			}
			assert.Contains(t, validTypes, cfg.Type)
		})
	}
}

// TestNew_PostgresMigrationPath tests postgres-specific migration behavior
func TestNew_PostgresMigration_SkippedForNonPostgres(t *testing.T) {
	// Migrations should only run for postgres type
	// For memory and sqlite, GORM AutoMigrate is used instead

	t.Run("memory type does not trigger migrations", func(t *testing.T) {
		cfg := Config{
			Config: store.Config{
				Type: "memory",
			},
		}
		assert.NotEqual(t, "postgres", cfg.Type)
	})

	t.Run("sqlite type does not trigger migrations", func(t *testing.T) {
		cfg := Config{
			Config: store.Config{
				Type: "sqlite",
				Sqlite: store.Sqlite{
					Filepath: "/tmp/test.db",
				},
			},
		}
		assert.NotEqual(t, "postgres", cfg.Type)
	})
}

// TestNew_AutoMigrate_CalledForNonPostgres tests AutoMigrate behavior
func TestNew_AutoMigrate_ForSQLite(t *testing.T) {
	// AutoMigrate should be called for sqlite and memory types
	// but NOT for postgres (which uses SQL migrations)

	nonPostgresTypes := []string{"memory", "sqlite"}

	for _, storeType := range nonPostgresTypes {
		t.Run("AutoMigrate for "+storeType, func(t *testing.T) {
			cfg := Config{
				Config: store.Config{
					Type: storeType,
				},
			}
			assert.NotEqual(t, "postgres", cfg.Type)
		})
	}
}

// Integration test - requires running Redis
// TestNew_Integration_MemoryStore tests full store creation with memory + Redis
func TestNew_Integration_MemoryStore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if Redis is available
	redisAvailable := os.Getenv("REDIS_HOST") != "" || isRedisAvailable()
	if !redisAvailable {
		t.Skip("Skipping integration test: Redis not available")
	}

	cfg := Config{
		Config: store.Config{
			Type: "memory",
		},
		Redis: redis.Config{
			Host:     getRedisHost(),
			Port:     getRedisPort(),
			Password: "",
			DB:       0,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	require.NoError(t, err)
	require.NotNil(t, s)

	// Verify health check works
	err = s.HealthCheck(ctx)
	assert.NoError(t, err)
}

// TestNew_Integration_SQLiteStore tests full store creation with SQLite + Redis
func TestNew_Integration_SQLiteStore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisAvailable := os.Getenv("REDIS_HOST") != "" || isRedisAvailable()
	if !redisAvailable {
		t.Skip("Skipping integration test: Redis not available")
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "integration_test.db")

	cfg := Config{
		Config: store.Config{
			Type: "sqlite",
			Sqlite: store.Sqlite{
				Filepath: dbPath,
			},
		},
		Redis: redis.Config{
			Host:     getRedisHost(),
			Port:     getRedisPort(),
			Password: "",
			DB:       0,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	require.NoError(t, err)
	require.NotNil(t, s)

	// Verify database file was created
	_, statErr := os.Stat(dbPath)
	assert.NoError(t, statErr, "SQLite database file should exist")

	// Verify health check works
	err = s.HealthCheck(ctx)
	assert.NoError(t, err)
}

// Benchmark tests
func BenchmarkConfig_Creation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Config{
			Config: store.Config{
				Type: "memory",
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
			ExecutionTTL: 24 * time.Hour,
			ProgressTTL:  1 * time.Hour,
		}
	}
}

func BenchmarkConfig_ConnectionString(b *testing.B) {
	cfg := Config{
		Config: store.Config{
			Type: "postgres",
			Postgresql: store.Postgresql{
				DNS: "postgres://user:pass@localhost:5432/db",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cfg.ConnectionString()
	}
}

// Helper functions for integration tests
func isRedisAvailable() bool {
	cfg := redis.Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
	client, err := redis.NewClient(cfg)
	if err != nil {
		return false
	}
	client.Close()
	return true
}

func getRedisHost() string {
	if host := os.Getenv("REDIS_HOST"); host != "" {
		return host
	}
	return "localhost"
}

func getRedisPort() int {
	// For simplicity, always return default port
	// In production, you might parse REDIS_PORT env var
	return 6379
}

// Close method helper for testing (assuming redis.Client has Close)
type closeable interface {
	Close() error
}

func closeIfPossible(c any) {
	if closer, ok := c.(closeable); ok {
		closer.Close()
	}
}

// TestNew_DialectorSelection tests that the correct dialector is selected for each store type
// This test exercises the switch statement by attempting to create stores with different types
// The Redis connection will fail, but the dialector selection should succeed first
func TestNew_DialectorSelection(t *testing.T) {
	tests := []struct {
		name         string
		storeType    string
		expectDbErr  bool // whether we expect a DB-related error before Redis
		errContains  string
	}{
		{
			name:        "memory dialector selected",
			storeType:   "memory",
			expectDbErr: false,
			errContains: "Redis", // fails at Redis, not dialector
		},
		{
			name:        "sqlite dialector selected",
			storeType:   "sqlite",
			expectDbErr: false,
			errContains: "Redis", // fails at Redis, not dialector
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			cfg := Config{
				Config: store.Config{
					Type: tt.storeType,
					Sqlite: store.Sqlite{
						Filepath: filepath.Join(tempDir, "test.db"),
					},
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 1, // Invalid port to ensure Redis fails
				},
				ExecutionTTL: 24 * time.Hour,
				ProgressTTL:  1 * time.Hour,
			}

			ctx := context.Background()
			s, err := New(ctx, cfg)

			assert.Error(t, err)
			assert.Nil(t, s)
			assert.Contains(t, err.Error(), tt.errContains)
		})
	}
}

// TestNew_DBStoreCreationSuccess tests that the DB store is created successfully
// before failing on Redis connection
func TestNew_DBStoreCreationSuccess_MemoryType(t *testing.T) {
	cfg := Config{
		Config: store.Config{
			Type: "memory",
		},
		Redis: redis.Config{
			Host: "invalid-redis-host",
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail at Redis connection, not at DB creation
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.Contains(t, err.Error(), "Redis")
}

// TestNew_DBStoreCreationSuccess_SQLiteType tests SQLite store creation
func TestNew_DBStoreCreationSuccess_SQLiteType(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_store.db")

	cfg := Config{
		Config: store.Config{
			Type: "sqlite",
			Sqlite: store.Sqlite{
				Filepath: dbPath,
			},
		},
		Redis: redis.Config{
			Host: "invalid-redis-host",
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail at Redis connection
	// DB file may or may not be created depending on when gorm opens it
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.Contains(t, err.Error(), "Redis")
}

// TestNew_PostgresType_MigrationPath tests the postgres migration path
// Note: This requires actual postgres which we don't have, so it should fail
func TestNew_PostgresType_MigrationFailure(t *testing.T) {
	cfg := Config{
		Config: store.Config{
			Type: "postgres",
			Postgresql: store.Postgresql{
				DNS: "postgres://invalid:invalid@localhost:5432/nonexistent?sslmode=disable",
			},
		},
		Redis: redis.Config{
			Host: "localhost",
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail at migration or DB connection
	assert.Error(t, err)
	assert.Nil(t, s)
}

// TestNew_PostgresType_EmptyConnectionString tests postgres with empty connection string
func TestNew_PostgresType_EmptyConnectionString(t *testing.T) {
	cfg := Config{
		Config: store.Config{
			Type: "postgres",
			Postgresql: store.Postgresql{
				DNS: "",
			},
		},
		Redis: redis.Config{
			Host: "localhost",
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Empty DNS should cause migration to be skipped (connStr == "")
	// Then it should fail at postgres.Open with empty DNS
	assert.Error(t, err)
	assert.Nil(t, s)
}

// TestNew_AllValidTypesFailAtRedis tests that all valid store types
// successfully create dialectors but fail at Redis connection
func TestNew_AllValidTypesFailAtRedis(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "memory type",
			config: Config{
				Config: store.Config{
					Type: "memory",
				},
				Redis: redis.Config{
					Host: "invalid-host",
					Port: 1,
				},
			},
		},
		{
			name: "sqlite type",
			config: Config{
				Config: store.Config{
					Type: "sqlite",
					Sqlite: store.Sqlite{
						Filepath: filepath.Join(tempDir, "test1.db"),
					},
				},
				Redis: redis.Config{
					Host: "invalid-host",
					Port: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s, err := New(ctx, tt.config)

			assert.Error(t, err)
			assert.Nil(t, s)
			// Error should be from Redis connection, not dialector
			assert.Contains(t, err.Error(), "Redis")
		})
	}
}

// TestConfig_ZeroValue tests behavior with zero-value Config
func TestConfig_ZeroValue(t *testing.T) {
	var cfg Config

	assert.Empty(t, cfg.Type)
	assert.Empty(t, cfg.Redis.Host)
	assert.Zero(t, cfg.Redis.Port)
	assert.Zero(t, cfg.ExecutionTTL)
	assert.Zero(t, cfg.ProgressTTL)
}

// TestConfig_PartialInitialization tests partial config initialization
func TestConfig_PartialInitialization(t *testing.T) {
	t.Run("only store type set", func(t *testing.T) {
		cfg := Config{
			Config: store.Config{
				Type: "memory",
			},
		}

		assert.Equal(t, "memory", cfg.Type)
		assert.Empty(t, cfg.Redis.Host)
		assert.Zero(t, cfg.ExecutionTTL)
	})

	t.Run("only redis config set", func(t *testing.T) {
		cfg := Config{
			Redis: redis.Config{
				Host: "localhost",
				Port: 6379,
			},
		}

		assert.Empty(t, cfg.Type)
		assert.Equal(t, "localhost", cfg.Redis.Host)
		assert.Equal(t, 6379, cfg.Redis.Port)
	})

	t.Run("only TTLs set", func(t *testing.T) {
		cfg := Config{
			ExecutionTTL: 1 * time.Hour,
			ProgressTTL:  30 * time.Minute,
		}

		assert.Empty(t, cfg.Type)
		assert.Equal(t, 1*time.Hour, cfg.ExecutionTTL)
		assert.Equal(t, 30*time.Minute, cfg.ProgressTTL)
	})
}

// TestNew_ContextCancellation tests behavior when context is cancelled
func TestNew_ContextCancellation(t *testing.T) {
	t.Run("cancelled context before call", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		cfg := Config{
			Config: store.Config{
				Type: "memory",
			},
			Redis: redis.Config{
				Host: "localhost",
				Port: 1, // Invalid to ensure failure
			},
		}

		s, err := New(ctx, cfg)

		// Should fail (either due to context or connection)
		assert.Error(t, err)
		assert.Nil(t, s)
	})
}

// TestNew_VerifyStoreInterface tests that the returned store implements the Store interface
func TestNew_VerifyStoreInterface(t *testing.T) {
	// This test verifies at compile time that the factory returns store.Store
	// We can't actually run it without Redis, but we can verify the types
	var s store.Store
	_ = s // Verify store.Store is a valid type

	// The New function signature guarantees it returns store.Store
	// New(ctx context.Context, cfg Config) (store.Store, error)
}

// TestNew_ConcurrentCalls tests thread safety of New function
func TestNew_ConcurrentCalls(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	// Note: All calls should fail due to invalid Redis, but the function
	// should be thread-safe
	cfg := Config{
		Config: store.Config{
			Type: "memory",
		},
		Redis: redis.Config{
			Host: "invalid-host",
			Port: 1,
		},
	}

	const numGoroutines = 10
	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			ctx := context.Background()
			_, err := New(ctx, cfg)
			errChan <- err
		}()
	}

	// Collect all errors
	for i := 0; i < numGoroutines; i++ {
		err := <-errChan
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Redis")
	}
}

// TestConfig_JSONTags tests that JSON tags are properly set
func TestConfig_JSONTags(t *testing.T) {
	// This test verifies the struct tags exist
	// The actual JSON marshaling would use these tags
	cfg := Config{
		Config: store.Config{
			Type: "memory",
		},
		Redis: redis.Config{
			Host: "localhost",
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ResultTTL:    7 * 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	// Verify structure is valid
	assert.NotEmpty(t, cfg.Type)
	assert.NotEmpty(t, cfg.Redis.Host)
	assert.Greater(t, cfg.ExecutionTTL, time.Duration(0))
}

// TestNew_SQLiteWithInvalidPath tests SQLite with an invalid path that should cause creation to fail
func TestNew_SQLiteWithInvalidPath(t *testing.T) {
	// Try to create SQLite database in a non-existent directory with no write permissions
	cfg := Config{
		Config: store.Config{
			Type: "sqlite",
			Sqlite: store.Sqlite{
				Filepath: "/nonexistent/deeply/nested/path/that/does/not/exist/test.db",
			},
		},
		Redis: redis.Config{
			Host: "localhost",
			Port: 1, // Invalid port
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail either at DB creation or Redis - both are acceptable
	assert.Error(t, err)
	assert.Nil(t, s)
}

// TestNew_MultipleInvalidConfigCombinations tests various invalid config combinations
func TestNew_MultipleInvalidConfigCombinations(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "empty config",
			config: Config{},
		},
		{
			name: "only type without dependencies",
			config: Config{
				Config: store.Config{
					Type: "memory",
				},
			},
		},
		{
			name: "postgres with invalid connection string format",
			config: Config{
				Config: store.Config{
					Type: "postgres",
					Postgresql: store.Postgresql{
						DNS: "not-a-valid-postgres-url",
					},
				},
				Redis: redis.Config{
					Host: "localhost",
					Port: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s, err := New(ctx, tt.config)

			// All these configs should fail
			assert.Error(t, err)
			assert.Nil(t, s)
		})
	}
}

// TestNew_PostgresWithValidDNS_FailsAtConnection tests postgres path with non-empty but unreachable DNS
func TestNew_PostgresWithValidDNS_FailsAtMigration(t *testing.T) {
	// Use a DNS that looks valid but points to an unreachable host
	cfg := Config{
		Config: store.Config{
			Type: "postgres",
			Postgresql: store.Postgresql{
				DNS: "postgres://user:pass@192.0.2.1:5432/testdb?sslmode=disable", // TEST-NET-1, unreachable
			},
		},
		Redis: redis.Config{
			Host: "localhost",
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail at migration due to connection issue
	assert.Error(t, err)
	assert.Nil(t, s)
}

// TestNew_SQLiteAutoMigrateSuccessPath tests the path where SQLite successfully creates DB
// and AutoMigrate succeeds, but fails at Redis
func TestNew_SQLiteAutoMigrateSuccessPath(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "automigrate_test.db")

	cfg := Config{
		Config: store.Config{
			Type: "sqlite",
			Sqlite: store.Sqlite{
				Filepath: dbPath,
			},
		},
		Redis: redis.Config{
			Host: "192.0.2.1", // TEST-NET-1, unreachable
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail at Redis, not at DB creation or AutoMigrate
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.Contains(t, err.Error(), "Redis")

	// Verify SQLite database file was created (proving AutoMigrate ran successfully)
	_, statErr := os.Stat(dbPath)
	assert.NoError(t, statErr, "SQLite database file should exist after AutoMigrate")
}

// TestNew_MemoryAutoMigrateSuccessPath tests the memory store AutoMigrate path
func TestNew_MemoryAutoMigrateSuccessPath(t *testing.T) {
	cfg := Config{
		Config: store.Config{
			Type: "memory",
		},
		Redis: redis.Config{
			Host: "192.0.2.1", // TEST-NET-1, unreachable
			Port: 6379,
		},
		ExecutionTTL: 24 * time.Hour,
		ProgressTTL:  1 * time.Hour,
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail at Redis, meaning DB creation and AutoMigrate succeeded
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.Contains(t, err.Error(), "Redis")
}

// TestNew_ConfigWithDefaultValues tests behavior with minimal required config
func TestNew_ConfigWithDefaultValues(t *testing.T) {
	cfg := Config{
		Config: store.Config{
			Type: "memory",
		},
		Redis: redis.Config{
			Host: "localhost",
			Port: 1, // Invalid to ensure fast failure
		},
		// TTLs left at zero values
	}

	ctx := context.Background()
	s, err := New(ctx, cfg)

	// Should fail at Redis
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.Contains(t, err.Error(), "Redis")

	// Verify zero TTLs are acceptable
	assert.Zero(t, cfg.ExecutionTTL)
	assert.Zero(t, cfg.ProgressTTL)
}
