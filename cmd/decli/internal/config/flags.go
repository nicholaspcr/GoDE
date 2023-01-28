package config

import (
	"github.com/spf13/pflag"
)

var CLI = &Config{
	PopulationSize: new(int),
	Generations:    new(int),
	Executions:     new(int),
	SaveEachGen:    new(bool),
	Dimensions: &Dimensions{
		Size:   new(int),
		Floors: new([]float64),
		Ceils:  new([]float64),
	},
	Constants: &Constants{
		M:  new(int),
		CR: new(float64),
		F:  new(float64),
		P:  new(float64),
	},
}

func globalFlags(set *pflag.FlagSet) {
	// General
	set.IntVarP(
		CLI.PopulationSize,
		"population",
		"p",
		50,
		"Determines size of population",
	)
	set.IntVarP(
		CLI.Generations,
		"generations",
		"g",
		100,
		"Determines amount of generations",
	)
	set.IntVarP(
		CLI.Executions,
		"executions",
		"e",
		1,
		"Determines amount of executions",
	)
	set.BoolVar(
		CLI.SaveEachGen,
		"save-each-gen",
		false,
		"Determines if population will be save on each generation",
	)

	// Dimensions
	set.IntVarP(
		CLI.Dimensions.Size,
		"dim-size",
		"d",
		7,
		"Determines size of the dimensions within a Population Element's vector",
	)
	set.Float64SliceVar(
		CLI.Dimensions.Ceils,
		"ceils",
		[]float64{1, 1, 1, 1, 1, 1, 1},
		"Ceil value for each dimension of a vector",
	)
	set.Float64SliceVar(
		CLI.Dimensions.Floors,
		"floors",
		[]float64{0, 0, 0, 0, 0, 0, 0},
		"Floor value for each dimension of a vector",
	)

	// Constants
	set.IntVar(CLI.Constants.M, "const-M", 3, "DE Constant")
	set.Float64Var(CLI.Constants.CR, "const-CR", 0.9, "DE Constant")
	set.Float64Var(CLI.Constants.CR, "const-F", 0.5, "DE Constant")
	set.Float64Var(CLI.Constants.CR, "const-P", 0.2, "DE Constant")
}

var ProblemName = new(string)
var VariantName = new(string)

// ModeFlags contains flags that should only be included in the Mode command.
func ModeFlags(set *pflag.FlagSet) {
	globalFlags(set)
	set.StringVar(
		ProblemName,
		"problem",
		"DTLZ1",
		"Selects what problem algorithm to run",
	)
	set.StringVar(
		VariantName,
		"variant",
		"rand1",
		"Selects what variant to user",
	)
}
