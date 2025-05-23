package server

import "github.com/nicholaspcr/GoDE/pkg/de"

// Config contains all the necessary configuration options for the server.
type Config struct {
	LisAddr  string
	HTTPPort string
	DE       de.Config
}

// DefaultConfig returns the default configuration of the server.
func DefaultConfig() Config {
	return Config{
		LisAddr:  "localhost:3030",
		HTTPPort: ":8081",
		DE: de.Config{
			ParetoChannelLimiter: 100,
			MaxChannelLimiter:    100,
			ResultLimiter:        1000,
		},
	}
}
