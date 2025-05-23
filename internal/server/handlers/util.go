package handlers

import (
	crypto_rand "crypto/rand"
	"encoding/hex"
	"math/rand/v2"

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
func generatePopulation(p models.PopulationParams) models.Population {
	ret := make(models.Population, p.PopulationSize)
	for i := 0; i < p.PopulationSize; i++ {
		ret[i].Elements = make([]float64, p.DimensionSize)
		ret[i].Objectives = make([]float64, p.ObjectivesSize)

		for j := 0; j < p.DimensionSize; j++ {
			// range between floor and ceiling
			constant := p.CeilRange[j] - p.FloorRange[j]
			// value varies within [ceil,upper]
			ret[i].Elements[j] = rand.Float64()*constant + p.FloorRange[j]
		}
	}
	return ret
}
