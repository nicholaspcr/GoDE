package config

import (
	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/internal/store"
)

// Default configuration of the deserver binary.
func Default() *DeServer {
	return &DeServer{
		Log:    log.DefaultConfig(),
		Store:  store.Config{},
		Server: server.Config{},
	}
}
