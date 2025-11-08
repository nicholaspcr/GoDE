// Package problems defines the optimization problem interface and common utilities.
package problems

import "github.com/nicholaspcr/GoDE/pkg/models"

// Interface defines the contract for optimization problems.
// Implementations must provide an evaluation function that computes objective values
// for a given solution vector. The Evaluate method modifies the vector's Objectives
// slice in-place for performance.
type Interface interface {
	Name() string
	// Evaluate is the function responsible for altering the objective
	// slice of a vector, therefore is assumed that the Vector passed will
	// be modified by this func
	Evaluate(*models.Vector, int) error
}
