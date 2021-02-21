package variants

import (
	"errors"
	"math/rand"
)

// generates random indices in the int slice, r -> it's a pointer
func generateIndices(startInd, NP int, r []int) error {
	if len(r) > NP {
		return errors.New("insufficient elements in population to generate random indices")
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
