package variants

import (
	"errors"
	"math/rand"
	"strings"
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

// GetVariantByName -> Returns the variant function
func GetVariantByName(name string) VariantFn {
	name = strings.ToLower(name)
	variants := map[string]VariantFn{
		"rand1":       rand1,
		"rand2":       rand2,
		"best1":       best1,
		"best2":       best2,
		"currtobest1": currToBest1,
		"pbest":       pbest,
	}
	for k, v := range variants {
		if name == k {
			return v
		}
	}
	return VariantFn{}
}

// GetAllVariants returns all the variants implemented in this package
func GetAllVariants() []VariantFn {
	variants := []VariantFn{
		rand1,
		rand2,
		best1,
		best2,
		currToBest1,
		pbest,
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
