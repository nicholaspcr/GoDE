package config

import (
	"time"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/internal/server"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/internal/storefactory"
)

// Default configuration of the deserver binary.
func Default() *Config {
	return &Config{
		Log:    log.DefaultConfig(),
		Server: server.DefaultConfig(),
		Store: storefactory.Config{
			Config: store.DefaultConfig(),
			Redis: redis.Config{
				Host:     "localhost",
				Port:     6379,
				Password: "",
				DB:       0,
			},
			ExecutionTTL: 24 * time.Hour,
			ResultTTL:    7 * 24 * time.Hour,
			ProgressTTL:  time.Hour,
		},
	}
}
