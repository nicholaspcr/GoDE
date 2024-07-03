package config

import (
	"encoding/json"

	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/internal/store"
)

// DeServer configuration.
type DeServer struct {
	Log    log.Config    `json:"log" yaml:"log"`
	Store  store.Config  `json:"store" yaml:"store"`
	Server server.Config `json:"server" yaml:"server"`
}

func (cfg *DeServer) StringifyJSON() (string, error) {
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
