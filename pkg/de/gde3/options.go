package gde3

import (
	"github.com/nicholaspcr/GoDE/internal/store"
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

// WithStore gde3termines the store to be used for storing the algorithm results.
func WithStore(s store.Store) Option {
	return func(m *gde3) *gde3 {
		m.store = s
		return m
	}
}

// WithPopulationParams determines the contants used to create the initial
// population of an execution.
func WithPopulationParams(params models.PopulationParams) Option {
	return func(m *gde3) *gde3 {
		m.population_params = params
		return m
	}
}
