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
func GetVariantByName(name string) models.VariantInterface {
	name = strings.ToLower(name)
	variants := map[string]models.VariantInterface{
		rand1.VariantName:       &rand1,
		rand2.VariantName:       &rand2,
		best1.VariantName:       &best1,
		best2.VariantName:       &best2,
		currToBest1.VariantName: &currToBest1,
		pbest.VariantName:       &pbest,
	}
	for k, v := range variants {
		if name == k {
			return v
		}
	}
	return &models.Variant{}
}

// GetAllVariants returns all the variants implemented in this package
func GetAllVariants() []models.VariantInterface {
	variants := []models.VariantInterface{
		&rand1,
		&rand2,
		&best1,
		&best2,
		&currToBest1,
		&pbest,
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
