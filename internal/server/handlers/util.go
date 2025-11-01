package handlers

import (
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// generatePopulation creates a population without objective values calculated
func generatePopulation(p models.PopulationParams, random *rand.Rand) (models.Population, error) {
	return models.GeneratePopulation(p, random)
}
