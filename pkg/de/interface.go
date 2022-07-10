package de

import (
	"os"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)


type Mode interface {
	Execute(
		chan<- models.Population,
		chan<- []float64,
		AlgorithmParams,
		problems.Interface,
		variants.Interface,
		models.Population,
		*os.File,
	)
}
