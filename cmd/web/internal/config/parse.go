package config

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

func (c *Config) JSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *Config) YAML() ([]byte, error) {
	return yaml.Marshal(c)
}
