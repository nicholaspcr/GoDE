package gde3

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// WithProblem attaches the Problem interface
// implementation that will be ran on DE execution.
func WithProblem(p problems.Interface) Option {
	return func(m *gde3) *gde3 {
		m.problem = p
		return m
	}
}

// WithVariant attaches the Variant interface
// implementation that will be ran on DE execution.
func WithVariant(v variants.Interface) Option {
	return func(m *gde3) *gde3 {
		m.variant = v
		return m
	}
}

// WithPopulationParams determines the contants used to create the initial
// population of an execution.
func WithPopulationParams(params models.PopulationParams) Option {
	return func(m *gde3) *gde3 {
		m.populationParams = params
		return m
	}
}

// WithConstants sets the constants used on DE execution.
func WithConstants(c Constants) Option {
	return func(m *gde3) *gde3 {
		m.constants = c
		return m
	}
}

// WithInitialPopulation determines the initial population of an execution.
func WithInitialPopulation(p models.Population) Option {
	return func(m *gde3) *gde3 {
		m.initialPopulation = p
		return m
	}
}
