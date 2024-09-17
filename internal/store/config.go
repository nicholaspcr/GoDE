package store

// Config contains options related to the Store implementation.
type Config struct {
	// Type supported are 'memory', 'sqlite', 'postgresql'.
	Type       string `json:"type" yaml:"type"`
	Memory     Memory
	Sqlite     Sqlite
	Postgresql Postgresql
}

type Memory struct{}
type Sqlite struct {
	Filepath string `json:"filepath" yaml:"filepath"`
}
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
