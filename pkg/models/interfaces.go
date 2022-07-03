package models

import (
	"os"

	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type (

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
			problems.Problem,
			Variant,
			Population,
			*os.File,
		)
	}
)
