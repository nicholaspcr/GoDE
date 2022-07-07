package models

type (
	// AlgorithmParams are the set of parameters that can be used in the gde3
	// cli
	AlgorithmParams struct {
		// EXECS is the amount of times the algorithm will run on top of the
		// same
		// initial population
		EXECS int `json:"execs" yaml:"execs"`
		// DIM is the size of the dimensions of each vector element
		DIM int `json:"dim"   yaml:"dim"`
		// GEN is the amount of generations that the algorithm will execute
		GEN int `json:"gen"   yaml:"gen"`

		// NP represents the size of the population
		NP int `json:"np" yaml:"np"`
		// M represents the amount of objective functions
		M int `json:"m"  yaml:"m"`

		// limits the search area of the vectors in the population
		FLOOR []float64 `json:"floor" yaml:"floor"`
		CEIL  []float64 `json:"ceil"  yaml:"ceil"`

		// constants used for the gde3 algorithm
		CR float64 `json:"cr" yaml:"cr"`
		F  float64 `json:"f"  yaml:"f"`
		P  float64 `json:"p"  yaml:"p"`

		// Disables the data generation for the plot
		DisablePlot bool `json:"disable_plot" yaml:"disable_plot"`
	}
)
