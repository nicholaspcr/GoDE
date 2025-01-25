package config

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

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
