package de

import (
	"context"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestDominanceTest(t *testing.T) {
	tests := []struct {
		name     string
		x        []float64
		y        []float64
		expected int
	}{
		{
			name:     "x dominates y",
			x:        []float64{0, 0},
			y:        []float64{1, 1},
			expected: -1,
		},
		{
			name:     "y dominates x",
			x:        []float64{1, 1},
			y:        []float64{0, 0},
			expected: 1,
		},
		{
			name:     "no dominance",
			x:        []float64{1, 0},
			y:        []float64{0, 1},
			expected: 0,
		},
		{
			name:     "equal",
			x:        []float64{1, 1},
			y:        []float64{1, 1},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DominanceTest(tt.x, tt.y)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterDominated(t *testing.T) {
	elems := []models.Vector{
		{Objectives: []float64{1, 2}},
		{Objectives: []float64{2, 1}},
		{Objectives: []float64{3, 3}},
		{Objectives: []float64{0, 0}},
	}

	nonDominated, dominated := FilterDominated(elems)

	assert.Len(t, nonDominated, 1)
	assert.Len(t, dominated, 3)
}

func TestCalculateCrwdDist(t *testing.T) {
	elems := []models.Vector{
		{Objectives: []float64{1, 5}},
		{Objectives: []float64{2, 4}},
		{Objectives: []float64{3, 3}},
		{Objectives: []float64{4, 2}},
		{Objectives: []float64{5, 1}},
	}

	CalculateCrwdDist(elems)

	// Extremes should have max crowding distance
	assert.Equal(t, INF, elems[0].CrowdingDistance)
	assert.Equal(t, INF, elems[4].CrowdingDistance)

	// Check crowding distance of middle elements
	assert.InDelta(t, 1.0, elems[1].CrowdingDistance, 0.01)
	assert.InDelta(t, 1.0, elems[2].CrowdingDistance, 0.01)
	assert.InDelta(t, 1.0, elems[3].CrowdingDistance, 0.01)
}

func TestReduceByCrowdDistance(t *testing.T) {
	elems := []models.Vector{
		{Objectives: []float64{1, 5}},
		{Objectives: []float64{2, 4}},
		{Objectives: []float64{3, 3}},
		{Objectives: []float64{4, 2}},
		{Objectives: []float64{5, 1}},
		{Objectives: []float64{1, 1}},
		// Dominated elements
		{Objectives: []float64{0, 0}},
	}

	reduced, best := ReduceByCrowdDistance(context.Background(), elems, 5)

	assert.Len(t, reduced, 5)
	assert.Len(t, best, 1)
}
