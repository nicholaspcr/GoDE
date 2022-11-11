package store

import "github.com/nicholaspcr/GoDE/pkg/models"

// Store TODO
type Store interface {
	PopulationStore
}

// PopulationStore TODO
type PopulationStore interface {
	Create(models.Population) error
	Update(models.Population) error
	Read(models.Population) error
	Delete(models.Population) error
}
