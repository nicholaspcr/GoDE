package config

import (
	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/internal/store"
)

// Default configuration of the deserver binary.
func Default() *Config {
	return &Config{
		Log:    log.DefaultConfig(),
		Server: server.DefaultConfig(),
		Store:  store.DefaultConfig(),
	}
}
