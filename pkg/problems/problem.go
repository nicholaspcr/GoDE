package problems

import "github.com/nicholaspcr/GoDE/pkg/api"

// Problem contains the definition of what a problem should have
type Interface interface {
	Name() string
	// Evaluate is the function responsible for altering the objective
	// slice of a vector, therefore is assumed that the Vector passed will
	// be modified by this func
	Evaluate(*api.Vector, int) error
}
