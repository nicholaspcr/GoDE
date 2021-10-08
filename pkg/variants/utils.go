package variants

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// generates random indices in the int slice, r -> it's a pointer
func generateIndices(startInd, NP int, r []int) error {
	if len(r) > NP {
		return errors.New(
			"insufficient elements in population to generate random indices",
		)
	}
	for i := startInd; i < len(r); i++ {
		for done := false; !done; {
			r[i] = rand.Int() % NP
			done = true
			for j := 0; j < i; j++ {
				done = done && r[j] != r[i]
			}
		}
	}
	return nil
}

// GetVariantByVariantName -> Returns the variant function
func GetVariantByName(name string) models.Variant {
	name = strings.ToLower(name)
	variants := map[string]models.Variant{
		Rand1().Name():       Rand1(),
		Rand2().Name():       Rand2(),
		Best1().Name():       Best1(),
		Best2().Name():       Best2(),
		CurrToBest1().Name(): CurrToBest1(),
		Pbest().Name():       Pbest(),
	}
	for k, v := range variants {
		if name == k {
			return v
		}
	}
	return nil
}

// GetAllVariants returns all the variants implemented in this package
func GetAllVariants() []models.Variant {
	variants := []models.Variant{
		Rand1(),
		Rand2(),
		Best1(),
		Best2(),
		CurrToBest1(),
		Pbest(),
	}
	return variants
}

// GetStandardPValues returns the default values used for Pbest variants
func GetStandardPValues() []float64 {
	return []float64{
		0.05,
		0.10,
		0.15,
		0.20,
	}
}
