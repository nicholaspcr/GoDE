package variants

import "gitlab.com/nicholaspcr/go-de/pkg/problems/models"

// shared variables and definitions

// Params are the necessary values that a variant uses
type Params struct {
	DIM     int
	F       float64
	CurrPos int
	P       float64
}

// VariantFn function type of the multiple variants
type VariantFn struct {
	Fn   func(elems, rankZero models.Elements, p Params) (models.Elem, error)
	Name string
}
