// Package variants provides mutation strategy implementations for Differential Evolution algorithms.
package variants

import "errors"

var (
	// ErrInsufficientPopulation indicates the population is too small for the requested mutation operation.
	ErrInsufficientPopulation = errors.New("insufficient population size")
	// ErrInvalidVector indicates a vector has nil or empty elements.
	ErrInvalidVector          = errors.New("vector has nil or empty elements")
)
