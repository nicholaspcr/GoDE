// Package store defines storage interfaces and implementations for persistence.
package store

import (
	"fmt"
	"time"
)

// Config contains options related to the Store implementation.
type Config struct {
	// Type supported are 'memory', 'sqlite', 'postgresql'.
	Type       string     `json:"type" yaml:"type" mapstructure:"type"`
	Memory     Memory     `json:"memory" yaml:"memory" mapstructure:"memory"`
	Sqlite     Sqlite     `json:"sqlite" yaml:"sqlite" mapstructure:"sqlite"`
	Postgresql Postgresql `json:"postgresql" yaml:"postgresql" mapstructure:"postgresql"`
}

// ConnectionPool configures database connection pool settings.
type ConnectionPool struct {
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`       // Maximum number of idle connections in the pool
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns" mapstructure:"max_open_conns"`       // Maximum number of open connections
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"` // Maximum lifetime of a connection
}

// DefaultConnectionPool returns sensible defaults for connection pooling.
func DefaultConnectionPool() ConnectionPool {
	return ConnectionPool{
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
	}
}

// Memory represents in-memory storage configuration.
type Memory struct{}

// Sqlite represents SQLite storage configuration.
type Sqlite struct {
	Filepath string         `json:"filepath" yaml:"filepath" mapstructure:"filepath"`
	Pool     ConnectionPool `json:"pool" yaml:"pool" mapstructure:"pool"`
}

// Postgresql represents PostgreSQL storage configuration.
type Postgresql struct {
	DNS  string         `json:"dns" yaml:"dns" mapstructure:"dns"`
	Pool ConnectionPool `json:"pool" yaml:"pool" mapstructure:"pool"`
}

// DefaultConfig returns the standard configuration for the Store package.
func DefaultConfig() Config {
	return Config{
		Type: "sqlite",
		Sqlite: Sqlite{
			Filepath: ".dev/server/sqlite.db",
			Pool:     DefaultConnectionPool(),
		},
		Postgresql: Postgresql{
			Pool: DefaultConnectionPool(),
		},
	}
}

// ConnectionString returns the database connection string for migrations.
func (c *Config) ConnectionString() string {
	switch c.Type {
	case "sqlite":
		return fmt.Sprintf("sqlite3://%s", c.Sqlite.Filepath)
	case "postgres":
		return c.Postgresql.DNS
	case "memory":
		return "" // Memory store doesn't support migrations
	default:
		return ""
	}
}
