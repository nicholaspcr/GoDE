// Package store defines storage interfaces and implementations for persistence.
package store

import "fmt"

// Config contains options related to the Store implementation.
type Config struct {
	// Type supported are 'memory', 'sqlite', 'postgresql'.
	Type       string `json:"type" yaml:"type"`
	Memory     Memory
	Sqlite     Sqlite
	Postgresql Postgresql
}

// Memory represents in-memory storage configuration.
type Memory struct{}

// Sqlite represents SQLite storage configuration.
type Sqlite struct {
	Filepath string `json:"filepath" yaml:"filepath"`
}

// Postgresql represents PostgreSQL storage configuration.
type Postgresql struct {
	DNS string `json:"dns" yaml:"dns"`
}

// DefaultConfig returns the standard configuration for the Store package.
func DefaultConfig() Config {
	return Config{
		Type:   "sqlite",
		Sqlite: Sqlite{Filepath: ".dev/server/sqlite.db"},
	}
}

// ConnectionString returns the database connection string for migrations.
func (c *Config) ConnectionString() string {
	switch c.Type {
	case "sqlite":
		return fmt.Sprintf("sqlite3://%s", c.Sqlite.Filepath)
	case "postgresql":
		return c.Postgresql.DNS
	case "memory":
		return "" // Memory store doesn't support migrations
	default:
		return ""
	}
}
