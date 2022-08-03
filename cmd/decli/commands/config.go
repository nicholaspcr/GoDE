package commands

import "github.com/nicholaspcr/GoDE/internal/config"

type (
	// Config is a set of values that are necessary to execute an Differential
	// Evolutionary algorithm.
	Config struct {
		PopulationSize int        `name:"population_size" json:"population_size" yaml:"population_size"`
		Generations    int        `name:"generations"     json:"generations"     yaml:"generations"`
		Executions     int        `name:"executions"      json:"executions"      yaml:"executions"`
		Dimensions     Dimensions `name:"dimensions"      json:"dimensions"      yaml:"dimensions"`
		Constants      Constants  `name:"constants"       json:"constants"       yaml:"constants"`
	}

	// Dimensions is a set of values to define the behaviour that happens in each
	// Dimension of the DE.
	Dimensions struct {
		Size        int       `name:"size"          json:"size"          yaml:"size"`
		Floors      []float64 `name:"floors"        json:"floors"        yaml:"floors"`
		Ceils       []float64 `name:"ceils"         json:"ceils"         yaml:"ceils"`
		SaveEachGen bool      `name:"save_each_gen" json:"save_each_gen" yaml:"save_each_gen"`
	}

	// Constants is a set of values to define the behaviour of a DE.
	Constants struct {
		M  int     `name:"m"  json:"m"  yaml:"m"`
		CR float64 `name:"cr" json:"cr" yaml:"cr"`
		F  float64 `name:"f"  json:"f"  yaml:"f"`
		P  float64 `name:"p"  json:"p"  yaml:"p"`
	}
)

var (
	// DefaultConfig is a predefined set of values to run `decli` without having to
	// config anything.
	DefaultConfig = Config{
		PopulationSize: 100,
		Generations:    100,
		Executions:     1,
		Dimensions:     DefaultDimensions,
		Constants:      DefaultConstants,
	}

	// DefaultDimensions contains values that define a default behaviour for the
	// dimensions of a DE.
	DefaultDimensions = Dimensions{
		Size:   7,
		Floors: []float64{0, 0, 0, 0, 0, 0, 0},
		Ceils:  []float64{1, 1, 1, 1, 1, 1, 1},
	}

	// DefaultConstants constaints values to determine a set way to execute the DE.
	DefaultConstants = Constants{
		M:  3,
		CR: 0.9,
		F:  0.5,
		P:  0.5,
	}

	np, dim                 int
	gen, execs              int
	floor, ceil             []float64
	crConst, fConst, pConst float64
	mConst                  int
	functionName            string
	disablePlot             bool
  filename string
)

func setDefaults(conf *config.Config) {
	// population
	conf.SetDefault("np", DefaultConfig.PopulationSize)
	conf.SetDefault("gen", DefaultConfig.Generations)
	conf.SetDefault("execs", DefaultConfig.Executions)
	conf.SetDefault("dim", DefaultConfig.Dimensions)

	// constants
	conf.SetDefault("CR", DefaultConfig.Constants.CR)
	conf.SetDefault("F", DefaultConfig.Constants.F)
	conf.SetDefault("P", DefaultConfig.Constants.P)
	conf.SetDefault("M", DefaultConfig.Constants.M)
}
