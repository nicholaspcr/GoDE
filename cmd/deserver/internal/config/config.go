package config

import (
	"encoding/json"

	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/internal/server"
	"github.com/nicholaspcr/GoDE/internal/store"
	"gopkg.in/yaml.v3"
)

// Config configuration.
type Config struct {
	Log    log.Config    `json:"log" yaml:"log"`
	Store  store.Config  `json:"store" yaml:"store"`
	Server server.Config `json:"server" yaml:"server"`
}

// StringifyJSON returns a string with the JSON object of the configuration.
func (cfg *Config) StringifyJSON() (string, error) {
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// StringifyYAML returns a string block with the yaml configuration contents.
func (cfg *Config) StringifyYAML() (string, error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
