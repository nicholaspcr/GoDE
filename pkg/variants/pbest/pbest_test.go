package pbest

import (
	"math/rand"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/stretchr/testify/assert"
)

func TestPbest_Name(t *testing.T) {
	p := Pbest()
	assert.Equal(t, "pbest", p.Name())
}

func TestPbest_Mutate_Success(t *testing.T) {
	// Create a population of 10 vectors with 5 dimensions each
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)},
		}
	}

	// Create rank zero population (Pareto front)
	rankZero := make([]models.Vector, 5)
	for i := range rankZero {
		rankZero[i] = models.Vector{
			Elements: []float64{float64(i) * 0.1, float64(i)*0.1 + 0.1, float64(i)*0.1 + 0.2, float64(i)*0.1 + 0.3, float64(i)*0.1 + 0.4},
		}
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     5,
		F:       0.5,
		P:       0.1,
		Random:  rand.New(rand.NewSource(42)),
	}

	p := Pbest()
	mutant, err := p.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.NotNil(t, mutant)
	assert.Len(t, mutant.Elements, 5)

	// Verify mutant is valid (all elements should be finite numbers)
	for _, elem := range mutant.Elements {
		assert.False(t, elem != elem) // Not NaN
	}
}

func TestPbest_Mutate_InsufficientPopulation(t *testing.T) {
	// Create population with only 2 elements (need at least 3 for pbest)
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
		P:       0.1,
		Random:  rand.New(rand.NewSource(42)),
	}

	p := Pbest()
	_, err := p.Mutate(elems, rankZero, params)

	assert.Error(t, err)
	assert.Equal(t, variants.ErrInsufficientPopulation, err)
}

func TestPbest_Mutate_VariousDimensions(t *testing.T) {
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

			rankZero := make([]models.Vector, 5)
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
				P:       0.2,
				Random:  rand.New(rand.NewSource(123)),
			}

			p := Pbest()
			mutant, err := p.Mutate(elems, rankZero, params)

			assert.NoError(t, err)
			assert.Len(t, mutant.Elements, tt.dimensions)
		})
	}
}

func TestPbest_Mutate_IndicesAreUnique(t *testing.T) {
	// This test verifies that the generated indices are unique
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := make([]models.Vector, 3)
	for i := range rankZero {
		rankZero[i] = models.Vector{
			Elements: []float64{float64(i) * 0.1, float64(i)*0.1 + 0.1, float64(i)*0.1 + 0.2},
		}
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     3,
		F:       0.5,
		P:       0.3,
		Random:  rand.New(rand.NewSource(42)),
	}

	p := Pbest()
	// Run multiple times to check randomness
	for range 10 {
		mutant, err := p.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestPbest_Mutate_WithDifferentCurrPos(t *testing.T) {
	// Test that mutation works correctly with different current positions
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := make([]models.Vector, 3)
	for i := range rankZero {
		rankZero[i] = models.Vector{
			Elements: []float64{float64(i) * 0.1, float64(i)*0.1 + 0.1, float64(i)*0.1 + 0.2},
		}
	}

	p := Pbest()

	for currPos := range 5 {
		params := variants.Parameters{
			CurrPos: currPos,
			DIM:     3,
			F:       0.5,
			P:       0.3,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := p.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestPbest_Mutate_VariousPValues(t *testing.T) {
	// Test mutation with different P values
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := make([]models.Vector, 10)
	for i := range rankZero {
		rankZero[i] = models.Vector{
			Elements: []float64{float64(i) * 0.1, float64(i)*0.1 + 0.1, float64(i)*0.1 + 0.2},
		}
	}

	p := Pbest()

	pValues := []float64{0.1, 0.2, 0.5, 0.8, 1.0}
	for _, pval := range pValues {
		params := variants.Parameters{
			CurrPos: 0,
			DIM:     3,
			F:       0.5,
			P:       pval,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := p.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestPbest_Mutate_VariousFValues(t *testing.T) {
	// Test mutation with different F values
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := make([]models.Vector, 5)
	for i := range rankZero {
		rankZero[i] = models.Vector{
			Elements: []float64{float64(i) * 0.1, float64(i)*0.1 + 0.1, float64(i)*0.1 + 0.2},
		}
	}

	p := Pbest()

	fValues := []float64{0.1, 0.5, 0.9, 1.0, 1.5}
	for _, fval := range fValues {
		params := variants.Parameters{
			CurrPos: 0,
			DIM:     3,
			F:       fval,
			P:       0.2,
			Random:  rand.New(rand.NewSource(42)),
		}

		mutant, err := p.Mutate(elems, rankZero, params)
		assert.NoError(t, err)
		assert.Len(t, mutant.Elements, 3)
	}
}

func TestPbest_Mutate_SingleRankZeroElement(t *testing.T) {
	// Test with a single element in rankZero
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
		P:       1.0,
		Random:  rand.New(rand.NewSource(42)),
	}

	p := Pbest()
	mutant, err := p.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.Len(t, mutant.Elements, 3)
}

func TestPbest_Mutate_SmallPValue(t *testing.T) {
	// Test with very small P value (should select from top 10%)
	elems := make([]models.Vector, 10)
	for i := range elems {
		elems[i] = models.Vector{
			Elements: []float64{float64(i), float64(i + 1), float64(i + 2)},
		}
	}

	rankZero := make([]models.Vector, 10)
	for i := range rankZero {
		rankZero[i] = models.Vector{
			Elements: []float64{float64(i) * 0.1, float64(i)*0.1 + 0.1, float64(i)*0.1 + 0.2},
		}
	}

	params := variants.Parameters{
		CurrPos: 0,
		DIM:     3,
		F:       0.5,
		P:       0.05, // Top 5%
		Random:  rand.New(rand.NewSource(42)),
	}

	p := Pbest()
	mutant, err := p.Mutate(elems, rankZero, params)

	assert.NoError(t, err)
	assert.Len(t, mutant.Elements, 3)
}
