package server

// Config contains all the necessary configuration options for the server.
type Config struct {
	LisAddr  string
	HTTPPort string
}

// DefaultConfig returns the default configuration of the server.
func DefaultConfig() Config {
	return Config{
		LisAddr:  "localhost:3030",
		HTTPPort: ":8081",
	}
}
