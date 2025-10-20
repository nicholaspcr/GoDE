package variants

import (
	"errors"
	"math/rand"
)

// GenerateIndices returns random indices into the r slice.
func GenerateIndices(startInd, NP int, r []int, random *rand.Rand) error {
	if len(r) > NP {
		return errors.New(
			"insufficient elements in population to generate random indices",
		)
	}
	used := make(map[int]bool, len(r))
	for i := 0; i < startInd; i++ {
		used[r[i]] = true
	}

	for i := startInd; i < len(r); i++ {
		for {
			val := random.Intn(NP)
			if !used[val] {
				r[i] = val
				used[val] = true
				break
			}
		}
	}
	return nil
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
