package models

import "os"

type (
	// ProblemInterface contains the definition of what a problem should have
	ProblemInterface interface {
		Name() string
		// Evaluate is the function responsible for altering the objective
		// slice of a vector, therefore is assumed that the Vector passed will
		// be modified by this func
		Evaluate(e *Vector, M int) error
	}

	VariantInterface interface {
		Name() string
		// Mutate is the funtion responsible for creating the trial vector
		Mutate(elems, rankZero Population, p VariantParams) (Vector, error)
	}

	ModeInterface interface {
		Execute(
			chan<- Population,
			chan<- []float64,
			AlgorithmParams,
			ProblemInterface,
			VariantInterface,
			Population,
			*os.File,
		)
	}
)
