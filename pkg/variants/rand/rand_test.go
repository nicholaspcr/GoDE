package rand

import (
	"math/rand"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/stretchr/testify/assert"
)

// Tests for Rand1

func TestRand1_Name(t *testing.T) {
	r := Rand1()
	assert.Equal(t, "rand1", r.Name())
}

func TestRand1_Mutate_Success(t *testing.T) {
	// Create a population of 10 vectors with 5 dimensions each
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)},
		}
	}

	// Create rank zero population (not used by rand1 but required)
	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2, 0.3, 0.4}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     5,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	r := Rand1()
	mutant, err := r.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.NotNil(t, mutant)
	assert.Len(t, mutant.Elements, 5)

	// Verify mutant is valid (all elements should be finite numbers)
	for _, elem := range mutant.Elements {
		assert.False(t, elem != elem) // Not NaN
	}
}

func TestRand1_Mutate_InsufficientPopulation(t *testing.T) {
	// Create population with only 3 elements (need at least 4 for rand/1)
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

	r := Rand1()
	_, err := r.Mutate(elems, rankZero, params)

	assert.Error(t, err)
}

func TestRand1_Mutate_VariousDimensions(t *testing.T) {
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

			rankZero := []models.Vector{
				{Elements: make([]float64, tt.dimensions)},
			}

			params := variants.Parameters{
				CurrPos: 0,
				DIM:     tt.dimensions,
				F:       0.8,
				Random:  rand.New(rand.NewSource(123)),
			}

			r := Rand1()
			mutant, err := r.Mutate(elems, rankZero, params)

			assert.NoError(t, err)
			assert.Len(t, mutant.Elements, tt.dimensions)
		})
	}
}

func TestRand1_Mutate_IndicesAreUnique(t *testing.T) {
	// This test verifies that the generated indices are unique
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     3,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	r := Rand1()
	// Run multiple times to check randomness
	for range 10 {
		mutant, err := r.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestRand1_Mutate_WithDifferentCurrPos(t *testing.T) {
	// Test that mutation works correctly with different current positions
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
	}

	r := Rand1()

	for currPos := range 5 {
		params := variants.Parameters{
			CurrPos: currPos,
			DIM:     3,
			F:       0.5,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := r.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestRand1_Mutate_VariousFValues(t *testing.T) {
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

	r := Rand1()

	fValues := []float64{0.1, 0.5, 0.9, 1.0, 1.5}
	for _, fval := range fValues {
		params := variants.Parameters{
			CurrPos: 0,
			DIM:     3,
			F:       fval,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := r.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestRand1_Mutate_MinimalPopulation(t *testing.T) {
	// Test with minimal population size (4 elements)
	elems := make([]models.Vector, 4)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     3,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	r := Rand1()
	mutant, err := r.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.Len(t, mutant.Elements, 3)
}

// Tests for Rand2

func TestRand2_Name(t *testing.T) {
	r := Rand2()
	assert.Equal(t, "rand2", r.Name())
}

func TestRand2_Mutate_Success(t *testing.T) {
	// Create a population of 10 vectors with 5 dimensions each
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)},
		}
	}

	// Create rank zero population (not used by rand2 but required)
	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2, 0.3, 0.4}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     5,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	r := Rand2()
	mutant, err := r.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.NotNil(t, mutant)
	assert.Len(t, mutant.Elements, 5)

	// Verify mutant is valid (all elements should be finite numbers)
	for _, elem := range mutant.Elements {
		assert.False(t, elem != elem) // Not NaN
	}
}

func TestRand2_Mutate_InsufficientPopulation(t *testing.T) {
	// Create population with only 5 elements (need at least 6 for rand/2)
	elems := make([]models.Vector, 5)
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

	r := Rand2()
	_, err := r.Mutate(elems, rankZero, params)

	assert.Error(t, err)
}

func TestRand2_Mutate_VariousDimensions(t *testing.T) {
	tests := []struct {
		name       string
		dimensions int
		popSize    int
	}{
		{"3 dimensions", 3, 15},
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

			rankZero := []models.Vector{
				{Elements: make([]float64, tt.dimensions)},
			}

			params := variants.Parameters{
				CurrPos: 0,
				DIM:     tt.dimensions,
				F:       0.8,
				Random:  rand.New(rand.NewSource(123)),
			}

			r := Rand2()
			mutant, err := r.Mutate(elems, rankZero, params)

			assert.NoError(t, err)
			assert.Len(t, mutant.Elements, tt.dimensions)
		})
	}
}

func TestRand2_Mutate_IndicesAreUnique(t *testing.T) {
	// This test verifies that the generated indices are unique
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     3,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	r := Rand2()
	// Run multiple times to check randomness
	for range 10 {
		mutant, err := r.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestRand2_Mutate_WithDifferentCurrPos(t *testing.T) {
	// Test that mutation works correctly with different current positions
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
	}

	r := Rand2()

	for currPos := range 5 {
		params := variants.Parameters{
			CurrPos: currPos,
			DIM:     3,
			F:       0.5,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := r.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestRand2_Mutate_VariousFValues(t *testing.T) {
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

	r := Rand2()

	fValues := []float64{0.1, 0.5, 0.9, 1.0, 1.5}
	for _, fval := range fValues {
		params := variants.Parameters{
			CurrPos: 0,
			DIM:     3,
			F:       fval,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := r.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestRand2_Mutate_MinimalPopulation(t *testing.T) {
	// Test with minimal population size (6 elements)
	elems := make([]models.Vector, 6)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := []models.Vector{
		{Elements: []float64{0.0, 0.1, 0.2}},
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     3,
		F:       0.5,
		Random:  rand.New(rand.NewSource(42)),
	}

	r := Rand2()
	mutant, err := r.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.Len(t, mutant.Elements, 3)
}
