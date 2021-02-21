package variants

import "gitlab.com/nicholaspcr/go-de/pkg/problems/models"

// shared variables and definitions

type varParams struct {
	DIM     int
	F       float64
	currPos int
	P       float64
}

// VariantFn function type of the multiple variants
type VariantFn struct {
	fn   func(elems, rankZero models.Elements, p varParams) (models.Elem, error)
	Name string
}
