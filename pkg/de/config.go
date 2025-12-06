// Package de provides the core Differential Evolution algorithm framework and execution utilities.
package de

import (
	"errors"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// Config of the DE implementation
type Config struct {
	ParetoChannelLimiter int
	MaxChannelLimiter    int
	ResultLimiter        int
}

// ProgressCallback is called periodically during execution to report progress.
type ProgressCallback func(generation int, totalGenerations int, paretoSize int, currentPareto []models.Vector)

// Constants are the set of values that determine the behaviour of the Mode
// execution.
type Constants struct {
	// Executions is the amount of times the algorithm will run.
	// All executions start with the same initial population.
	Executions int `json:"executions" yaml:"executions" name:"executions"`

	// Generations of an execution.
	Generations int `json:"generations" yaml:"generations" name:"generations"`

	// Dimensions is the size of the dimensions of each vector element.
	Dimensions int `json:"dimensions" yaml:"dimensions" name:"dimensions"`

	// ObjFuncAmount represents the amount of objective functions.
	ObjFuncAmount int `json:"obj_func_amount" yaml:"obj_func_amount" name:"obj_func_amount"`
}

// Validate checks that all Constants fields have valid values.
func (c *Constants) Validate() error {
	if c.Executions <= 0 {
		return errors.New("executions must be positive")
	}
	if c.Generations <= 0 {
		return errors.New("generations must be positive")
	}
	if c.Dimensions <= 0 {
		return errors.New("dimensions must be positive")
	}
	if c.ObjFuncAmount <= 0 {
		return errors.New("objective function amount must be positive")
	}
	return nil
}
