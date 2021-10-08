package models

import "os"

type (
	// Problem contains the definition of what a problem should have
	Problem interface {
		Name() string
		// Evaluate is the function responsible for altering the objective
		// slice of a vector, therefore is assumed that the Vector passed will
		// be modified by this func
		Evaluate(*Vector, int) error
	}

	// Var
	Variant interface {
		Name() string
		// Mutate is the funtion responsible for creating the trial vector
		Mutate(elems, rankZero Population, p VariantParams) (Vector, error)
	}

	Mode interface {
		Execute(
			chan<- Population,
			chan<- []float64,
			AlgorithmParams,
			Problem,
			Variant,
			Population,
			*os.File,
		)
	}
)
