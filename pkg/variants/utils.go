package variants

import (
	"errors"

	"github.com/nicholaspcr/GoDE/internal/fastrand"
)

// GenerateIndices returns random indices into the r slice.
func GenerateIndices(startInd, NP int, r []int) error {
	random := fastrand.NewRand()
	if len(r) > NP {
		return errors.New(
			"insufficient elements in population to generate random indices",
		)
	}
	for i := startInd; i < len(r); i++ {
		for done := false; !done; {
			r[i] = random.Int() % NP
			done = true
			for j := 0; j < i; j++ {
				done = done && r[j] != r[i]
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
