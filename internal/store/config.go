package store

// Config contains options related to the Store implementation.
type Config struct {
	// Type supported are 'memory', 'sqlite', 'postgresql'.
	Type string

	Memory     struct{}
	Sqlite     struct{ Filepath string }
	Postgresql struct{ DNS string }
}

// DefaultConfig returns the standard configuration for the Store package.
func DefaultConfig() Config {
	return Config{
		Type: "memory",
	}
}
