package config

import (
	"encoding/json"

	"github.com/nicholaspcr/GoDE/internal/log"
	"gopkg.in/yaml.v3"
)

type (
	// Config is a set of values that are necessary to execute an Differential
	// Evolutionary algorithm.
	Config struct {
		Local  LocalConfig  `json:"local" yaml:"local"`
		Log    LogConfig    `json:"log" yaml:"log"`
		Server ServerConfig `json:"server" yaml:"server"`
	}

	LocalConfig struct {
		PopulationSize int        `json:"populationSize" yaml:"populationSize"`
		Generations    int        `json:"generations" yaml:"generations"`
		Executions     int        `json:"executions" yaml:"executions"`
		Dimensions     Dimensions `json:"dimensions" yaml:"dimensions"`
		Constants      Constants  `json:"constants" yaml:"constants"`
		Problem        string     `json:"problem" yaml:"problem"`
		Variant        string     `json:"variant" yaml:"variant"`
	}

	// LogConfig is a set of values that are necessary to configure the logger.
	LogConfig struct {
		log.Config `json:"config" yaml:"config" mapstructure:",squash"`
		Filename   string `json:"filename" yaml:"filename"`
	}

	// ServerConfig is a set of values that are necessary for making requests to
	// the DE server.
	ServerConfig struct {
		HTTPAddr string `json:"http-addr" yaml:"http-addr"`
		GRPCAddr string `json:"grpc-addr" yaml:"grpc-addr"`
	}

	// Dimensions is a set of values to define the behaviour that happens in
	// each dimension of the DE.
	Dimensions struct {
		Size   int       `json:"size"   yaml:"size"`
		Floors []float64 `json:"floors" yaml:"floors"`
		Ceils  []float64 `json:"ceils"  yaml:"ceils"`
	}

	// Constants is a set of values to define the behaviour of a DE.
	Constants struct {
		M  int     `json:"m"  yaml:"m"`
		CR float64 `json:"cr" yaml:"cr"`
		F  float64 `json:"f"  yaml:"f"`
		P  float64 `json:"p"  yaml:"p"`
	}
)

// Default configuration of the decli binary.
func Default() *Config {
	return &Config{
		Local: LocalConfig{
			PopulationSize: 50,
			Generations:    100,
			Executions:     1,
			Dimensions: Dimensions{
				Size:   7,
				Floors: []float64{0, 0, 0, 0, 0, 0, 0},
				Ceils:  []float64{1, 1, 1, 1, 1, 1, 1},
			},
			Constants: Constants{
				M:  int(3),
				CR: float64(0.9),
				F:  float64(0.5),
				P:  float64(0.2),
			},
			Problem: "dtlz1",
			Variant: "rand1",
		},
		Log: LogConfig{
			Config: log.DefaultConfig(),
		},
		Server: ServerConfig{
			GRPCAddr: "localhost:3030",
			HTTPAddr: "http://localhost:8081",
		},
	}
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
