package models

import (
	"os"

	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

type (
	Mode interface {
		Execute(
			chan<- Population,
			chan<- []float64,
			AlgorithmParams,
			problems.Problem,
			variants.Interface,
			Population,
			*os.File,
		)
	}
)
