package variants

import "errors"

var (
	ErrInsufficientPopulation = errors.New("insufficient population size")
	ErrInvalidVector          = errors.New("vector has nil or empty elements")
)
