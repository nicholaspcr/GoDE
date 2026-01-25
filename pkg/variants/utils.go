package variants

import (
	"fmt"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// ErrMaxRetriesExceeded is returned when the random index generation
// exceeds the maximum number of retries.
var ErrMaxRetriesExceeded = fmt.Errorf("exceeded maximum retries while generating unique indices")

// maxRetries is the maximum number of attempts to find a unique random index.
// This prevents excessive spinning when most values are already used.
const maxRetries = 1000

// GenerateIndices returns random indices into the r slice.
func GenerateIndices(startInd, np int, r []int, random *rand.Rand) error {
	if len(r) > np {
		return ErrInsufficientPopulation
	}
	used := make(map[int]bool, len(r))
	for i := 0; i < startInd; i++ {
		used[r[i]] = true
	}

	for i := startInd; i < len(r); i++ {
		retries := 0
		for {
			if retries >= maxRetries {
				return fmt.Errorf("%w: slot %d after %d attempts", ErrMaxRetriesExceeded, i, retries)
			}
			val := random.Intn(np)
			if !used[val] {
				r[i] = val
				used[val] = true
				break
			}
			retries++
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

// ValidateVectors validates that vectors at specified indices have non-nil elements
// with the correct dimensions. This reduces duplication across variant implementations.
func ValidateVectors(vectors []models.Vector, indices []int, expectedDim int) error {
	for _, idx := range indices {
		if idx < 0 || idx >= len(vectors) {
			return ErrInvalidVector
		}
		if vectors[idx].Elements == nil || len(vectors[idx].Elements) != expectedDim {
			return ErrInvalidVector
		}
	}
	return nil
}
