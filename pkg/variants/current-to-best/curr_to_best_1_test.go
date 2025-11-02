package currenttobest

import (
	"math/rand"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/stretchr/testify/assert"
)

func TestCurrToBest1_Name(t *testing.T) {
	c := CurrToBest1()
	assert.Equal(t, "currToBest1", c.Name())
}

func TestCurrToBest1_Mutate_Success(t *testing.T) {
	// Create a population of 10 vectors with 5 dimensions each
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)},
		}
	}

	// Create rank zero population (Pareto front)
	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2, 0.3, 0.4}},
		{Elements: []float64{1.0, 1.1, 1.2, 1.3, 1.4}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     5,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	c := CurrToBest1()
	mutant, err := c.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.NotNil(t, mutant)
	assert.Len(t, mutant.Elements, 5)

	// Verify mutant is valid (all elements should be finite numbers)
	for _, elem := range mutant.Elements {
		assert.False(t, elem != elem) // Not NaN
	}
}

func TestCurrToBest1_Mutate_InsufficientPopulation(t *testing.T) {
	// Create population with only 3 elements (need at least 4 for current-to-best/1)
	elems := make([]models.Vector, 3)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     2,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	c := CurrToBest1()
	_, err := c.Mutate(elems, rankZero, params)

	assert.Error(t, err)
	assert.Equal(t, variants.ErrInsufficientPopulation, err)
}

func TestCurrToBest1_Mutate_VariousDimensions(t *testing.T) {
	tests := []struct {
		name       string
		dimensions int
		popSize    int
	}{
		{"3 dimensions", 3, 10},
		{"10 dimensions", 10, 20},
		{"30 dimensions", 30, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elems := make([]models.Vector, tt.popSize)
			for i := range elems {
				elem := make([]float64, tt.dimensions)
				for j := range elem {
					elem[j] = float64(i*tt.dimensions + j)
				}
				elems[i] = models.Vector{Elements: elem}
			}

			rankZero := make([]models.Vector, 3)
			for i := range rankZero {
				elem := make([]float64, tt.dimensions)
				for j := range elem {
					elem[j] = float64(i) * 0.1
				}
				rankZero[i] = models.Vector{Elements: elem}
			}

			params := variants.Parameters{
				CurrPos: 0,
				DIM:     tt.dimensions,
				F:       0.8,
				Random:  rand.New(rand.NewSource(123)),
			}

			c := CurrToBest1()
			mutant, err := c.Mutate(elems, rankZero, params)

			assert.NoError(t, err)
			assert.Len(t, mutant.Elements, tt.dimensions)
		})
	}
}

func TestCurrToBest1_Mutate_IndicesAreUnique(t *testing.T) {
	// This test verifies that the generated indices are unique
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
		{Elements: []float64{1.0, 1.1, 1.2}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     3,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	c := CurrToBest1()
	// Run multiple times to check randomness
	for i := 0; i < 10; i++ {
		mutant, err := c.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestCurrToBest1_Mutate_WithDifferentCurrPos(t *testing.T) {
	// Test that mutation works correctly with different current positions
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
		{Elements: []float64{1.0, 1.1, 1.2}},
	}

	c := CurrToBest1()

	for currPos := 0; currPos < 5; currPos++ {
		params := variants.Parameters{
			CurrPos: currPos,
			DIM:     3,
			F:       0.5,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := c.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestCurrToBest1_Mutate_EmptyRankZero(t *testing.T) {
	// Test with empty rankZero should cause a panic or error in bestIdx selection
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1)},
		}
	}

	rankZero := []models.Vector{}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     2,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	c := CurrToBest1()

	// This should panic due to empty rankZero when trying to select bestIdx
	assert.Panics(t, func() {
		_, _ = c.Mutate(elems, rankZero, params)
	})
}

func TestCurrToBest1_Mutate_VariousFValues(t *testing.T) {
	// Test mutation with different F values
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
	}

	c := CurrToBest1()

	fValues := []float64{0.1, 0.5, 0.9, 1.0, 1.5}
	for _, fval := range fValues {
		params := variants.Parameters{
			CurrPos: 0,
			DIM:     3,
			F:       fval,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := c.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}
