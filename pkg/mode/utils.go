package mode

import (
	"math/rand"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// GeneratePopulation - creates a population without objs calculates
func GeneratePopulation(p models.AlgorithmParams) models.Population {
	ret := make(models.Population, p.NP)
	for i := 0; i < p.NP; i++ {
		ret[i].X = make([]float64, p.DIM)

		for j := 0; j < p.DIM; j++ {
			// range between floor and ceiling
			constant := p.CEIL[j] - p.FLOOR[j]
			// value varies within [ceil,upper]
			ret[i].X[j] = rand.Float64()*constant + p.FLOOR[j]
		}
	}
	return ret
}
