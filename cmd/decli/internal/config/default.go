package config

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state/sqlite"
	"github.com/nicholaspcr/GoDE/internal/log"
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
		State: sqlite.Config{
			Provider: "file",
			Filepath: ".dev/cli/state",
		},
	}
}
