package config

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state/sqlite"
	"github.com/nicholaspcr/GoDE/internal/log"
)

// Default configuration of the decli binary.
func Default() *Config {
	return &Config{
		Log: LogConfig{
			Config: log.DefaultConfig(),
		},
		Server: ServerConfig{
			GRPCAddr: "localhost:3030",
			HTTPAddr: "http://localhost:8081",
		},
		State: sqlite.Config{
			Provider: "file",
			Filepath: ".dev/cli/state",
		},
	}
}
