package models

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Vector tests ---

func TestVector_Copy(t *testing.T) {
	v := &Vector{
		Elements:         []float64{1.0, 2.0, 3.0},
		Objectives:       []float64{0.1, 0.2},
		CrowdingDistance: 0.75,
	}
	c := v.Copy()

	assert.Equal(t, v.Elements, c.Elements)
	assert.Equal(t, v.Objectives, c.Objectives)
	assert.Equal(t, v.CrowdingDistance, c.CrowdingDistance)

	// Verify independence
	c.Elements[0] = 99.9
	assert.Equal(t, 1.0, v.Elements[0], "Copy should be independent")

	c.Objectives[0] = 99.9
	assert.Equal(t, 0.1, v.Objectives[0], "Copy objectives should be independent")
}

func TestVector_Copy_Empty(t *testing.T) {
	v := &Vector{}
	c := v.Copy()
	assert.Empty(t, c.Elements)
	assert.Empty(t, c.Objectives)
}

// --- Population tests ---

func TestPopulation_Copy(t *testing.T) {
	p := Population{
		{Elements: []float64{1.0, 2.0}, Objectives: []float64{0.5}, CrowdingDistance: 1.0},
		{Elements: []float64{3.0, 4.0}, Objectives: []float64{0.8}, CrowdingDistance: 2.0},
	}
	c := p.Copy()

	require.Len(t, c, 2)
	assert.Equal(t, p[0].Elements, c[0].Elements)
	assert.Equal(t, p[1].Elements, c[1].Elements)

	// Verify independence
	c[0].Elements[0] = 99.9
	assert.Equal(t, 1.0, p[0].Elements[0], "Copy should be independent from original")
}

func TestPopulation_Copy_Empty(t *testing.T) {
	var p Population
	c := p.Copy()
	assert.Empty(t, c)
}

// --- GeneratePopulation tests ---

func TestGeneratePopulation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	params := PopulationParams{
		FloorRange:     []float64{0.0, 0.0, 0.0},
		CeilRange:      []float64{1.0, 1.0, 1.0},
		DimensionSize:  3,
		PopulationSize: 10,
		ObjectivesSize: 2,
	}

	pop, err := GeneratePopulation(params, rng)
	require.NoError(t, err)
	require.Len(t, pop, 10)

	for _, v := range pop {
		assert.Len(t, v.Elements, 3)
		assert.Len(t, v.Objectives, 2)
		for i, el := range v.Elements {
			assert.GreaterOrEqual(t, el, params.FloorRange[i])
			assert.Less(t, el, params.CeilRange[i])
		}
	}
}

func TestGeneratePopulation_RangeMismatch(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	params := PopulationParams{
		FloorRange:    []float64{0.0, 0.0}, // length 2
		CeilRange:     []float64{1.0},      // length 1
		DimensionSize: 3,
		PopulationSize: 5,
	}

	_, err := GeneratePopulation(params, rng)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "floor range and ceil range must have the same size")
}

func TestGeneratePopulation_CustomRange(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	params := PopulationParams{
		FloorRange:     []float64{-5.0, 10.0},
		CeilRange:      []float64{-1.0, 20.0},
		DimensionSize:  2,
		PopulationSize: 20,
		ObjectivesSize: 3,
	}

	pop, err := GeneratePopulation(params, rng)
	require.NoError(t, err)
	require.Len(t, pop, 20)

	for _, v := range pop {
		assert.GreaterOrEqual(t, v.Elements[0], -5.0)
		assert.Less(t, v.Elements[0], -1.0)
		assert.GreaterOrEqual(t, v.Elements[1], 10.0)
		assert.Less(t, v.Elements[1], 20.0)
	}
}

func TestGeneratePopulation_ZeroSize(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	params := PopulationParams{
		FloorRange:     []float64{},
		CeilRange:      []float64{},
		DimensionSize:  0,
		PopulationSize: 0,
		ObjectivesSize: 0,
	}

	pop, err := GeneratePopulation(params, rng)
	require.NoError(t, err)
	assert.Empty(t, pop)
}
