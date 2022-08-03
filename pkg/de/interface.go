package de

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// InjectConstants allows for each implementation of a differential evolution
// to inject its own necessary contants into the context.
type InjectConstants func(context.Context) context.Context

// Mode defines the methods that a Differential Evolution algorihtm should
// implement, this method will be executed in each generation.
type Algorithm interface {
	Execute(
		context.Context,
		models.Population,
		problems.Interface,
		variants.Interface,
		Store,
		chan<- []models.Vector,
	) error
}

// Store TODO: move this to its own packer
type Store interface {
	Header(...string) error
	Population(models.Population) error
}
