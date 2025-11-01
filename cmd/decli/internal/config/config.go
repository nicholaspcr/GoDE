package config

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state/sqlite"
	"github.com/nicholaspcr/GoDE/internal/log"
)

type (
	// Config is a set of values that are necessary to execute an Differential
	// Evolutionary algorithm.
	Config struct {
		Log    LogConfig     `json:"log" yaml:"log"`
		Server ServerConfig  `json:"server" yaml:"server"`
		State  sqlite.Config `json:"state" yaml:"state"`
		Run    RunConfig     `json:"run" yaml:"run"`
	}

	RunConfig struct {
		Algorithm string   `json:"algorithm" yaml:"algorithm"`
		Variant   string   `json:"variant" yaml:"variant"`
		Problem   string   `json:"problem" yaml:"problem"`
		DeConfig  DEConfig `json:"de" yaml:"de"`
	}

	DEConfig struct {
		Executions     int64      `json:"executions" yaml:"executions"`
		Generations    int64      `json:"generations" yaml:"generations"`
		PopulationSize int64      `json:"population_size" yaml:"population_size"`
		DimensionsSize int64      `json:"dimensions_size" yaml:"dimensions_size"`
		ObjectivesSize int64      `json:"objectives_size" yaml:"objectives_size"`
		FloorLimiter   float32    `json:"floor_limiter" yaml:"floor_limiter"`
		CeilLimiter    float32    `json:"ceil_limiter" yaml:"ceil_limiter"`
		GDE3           GDE3Config `json:"gde3" yaml:"gde3"`
	}

	GDE3Config struct {
		CR float32 `json:"cr" yaml:"cr"`
		F  float32 `json:"f" yaml:"f"`
		P  float32 `json:"p" yaml:"p"`
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
		Floors []float64 `json:"floors" yaml:"floors"`
		Ceils  []float64 `json:"ceils"  yaml:"ceils"`
		Size   int       `json:"size"   yaml:"size"`
	}

	// Constants is a set of values to define the behaviour of a DE.
	Constants struct {
		M  int     `json:"m"  yaml:"m"`
		CR float64 `json:"cr" yaml:"cr"`
		F  float64 `json:"f"  yaml:"f"`
		P  float64 `json:"p"  yaml:"p"`
	}
)
