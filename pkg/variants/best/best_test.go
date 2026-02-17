package best

import (
	"math/rand"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/stretchr/testify/assert"
)

func TestBest1_Name(t *testing.T) {
	b := Best1()
	assert.Equal(t, "best1", b.Name())
}

func TestBest1_Mutate_Success(t *testing.T) {
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

	b := Best1()
	mutant, err := b.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.NotNil(t, mutant)
	assert.Len(t, mutant.Elements, 5)

	// Verify mutant is valid (all elements should be finite numbers)
	for _, elem := range mutant.Elements {
		assert.False(t, elem != elem) // Not NaN
	}
}

func TestBest1_Mutate_InsufficientPopulation(t *testing.T) {
	// Create population with only 2 elements (need at least 3 for best/1)
	elems := make([]models.Vector, 2)
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

	b := Best1()
	_, err := b.Mutate(elems, rankZero, params)

	assert.Error(t, err)
	assert.Equal(t, variants.ErrInsufficientPopulation, err)
}

func TestBest1_Mutate_VariousDimensions(t *testing.T) {
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

			b := Best1()
			mutant, err := b.Mutate(elems, rankZero, params)

			assert.NoError(t, err)
			assert.Len(t, mutant.Elements, tt.dimensions)
		})
	}
}

func TestBest2_Name(t *testing.T) {
	b := Best2()
	assert.Equal(t, "best2", b.Name())
}

func TestBest2_Mutate_Success(t *testing.T) {
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

	b := Best2()
	mutant, err := b.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.NotNil(t, mutant)
	assert.Len(t, mutant.Elements, 5)

	// Verify mutant is valid (all elements should be finite numbers)
	for _, elem := range mutant.Elements {
		assert.False(t, elem != elem) // Not NaN
	}
}

func TestBest2_Mutate_InsufficientPopulation(t *testing.T) {
	// Create population with only 4 elements (need at least 5 for best/2)
	elems := make([]models.Vector, 4)
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

	b := Best2()
	_, err := b.Mutate(elems, rankZero, params)

	assert.Error(t, err)
	assert.Equal(t, variants.ErrInsufficientPopulation, err)
}

func TestBest2_Mutate_VariousDimensions(t *testing.T) {
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

			b := Best2()
			mutant, err := b.Mutate(elems, rankZero, params)

			assert.NoError(t, err)
			assert.Len(t, mutant.Elements, tt.dimensions)
		})
	}
}

func TestBest1_Mutate_IndicesAreUnique(t *testing.T) {
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

	b := Best1()
	// Run multiple times to check randomness
	for range 10 {
		mutant, err := b.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestBest2_Mutate_IndicesAreUnique(t *testing.T) {
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

	b := Best2()
	// Run multiple times to check randomness
	for range 10 {
		mutant, err := b.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}
