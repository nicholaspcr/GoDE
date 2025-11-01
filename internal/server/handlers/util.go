package handlers

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := crypto_rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generatePopulation creates a population without objective values calculated
func generatePopulation(p models.PopulationParams, random *rand.Rand) (models.Population, error) {
	return models.GeneratePopulation(p, random)
}
