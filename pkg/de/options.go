package de

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// Constants are the set of values that determine the behaviour of the Mode execution.
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

	// constants used for the gde3 algorithm.
	CR float64 `json:"cr" yaml:"cr" name:"cr"`
	F  float64 `json:"f" yaml:"f" name:"f"`
	P  float64 `json:"p" yaml:"p" name:"p"`
}

// ModeOption defines a configuration method for the de struct
type ModeOptions func(*de) *de

// WithProblem attaches the Problem interface
// implementation that will be ran on DE execution.
func WithProblem(p problems.Interface) ModeOptions {
	return func(m *de) *de {
		m.problem = p
		return m
	}
}

// WithVariant attaches the Variant interface
// implementation that will be ran on DE execution.
func WithVariant(v variants.Interface) ModeOptions {
	return func(m *de) *de {
		m.variant = v
		return m
	}
}

// WithAlgorithm specifies which algorithm is to be runned on top of the problem.
func WithAlgorithm(alg Algorithm) ModeOptions {
	return func(m *de) *de {
		m.algorithm = alg
		return m
	}

}

// WithStore determines the store to be used for storing the algorithm results.
func WithStore(s Store) ModeOptions {
	return func(m *de) *de {
		m.store = s
		return m
	}
}

// WithPopulation determines the initial population of the DE.
func WithPopulation(pop models.Population) ModeOptions {
	return func(m *de) *de {
		m.population = pop
		return m
	}
}

// WithExecs determines the amount of times the DE will be
// executed, all executions start with the same population.
func WithExecs(n int) ModeOptions {
	return func(m *de) *de {
		m.constants.Executions = n
		return m
	}
}

// WithDimensions determines the size of each vector in the population element.
func WithDimensions(dim int) ModeOptions {
	return func(m *de) *de {
		m.constants.Dimensions = dim
		return m
	}
}

// WithGenerations determines the amount of generations of the DE.
func WithGenerations(gen int) ModeOptions {
	return func(m *de) *de {
		m.constants.Generations = gen
		return m
	}
}

// WithObjFuncAmount determines the amount of objective functions to be executed.
func WithObjFuncAmount(n int) ModeOptions {
	return func(m *de) *de {
		m.constants.ObjFuncAmount = n
		return m
	}
}

// WithCRConstant sets the CR constant value.
func WithCRConstant(cr float64) ModeOptions {
	return func(m *de) *de {
		m.constants.CR = cr
		return m
	}
}

// WithFConstant sets the F constant value.
func WithFConstant(f float64) ModeOptions {
	return func(m *de) *de {
		m.constants.F = f
		return m
	}
}

// WithPConstant sets the P constant value.
func WithPConstant(p float64) ModeOptions {
	return func(m *de) *de {
		m.constants.P = p
		return m
	}
}
