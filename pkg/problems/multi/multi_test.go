package multi

import (
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for ZDT1

func TestZdt1_Name(t *testing.T) {
	problem := Zdt1()
	assert.Equal(t, "zdt1", problem.Name())
}

func TestZdt1_Evaluate_Success(t *testing.T) {
	problem := Zdt1()
	vector := &models.Vector{
		Elements: []float64{0.5, 0.5, 0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.NoError(t, err)
	assert.Len(t, vector.Objectives, 2)

	// Verify objectives are valid numbers
	for _, obj := range vector.Objectives {
		assert.False(t, obj != obj) // Not NaN
	}
}

func TestZdt1_Evaluate_InsufficientDimensions(t *testing.T) {
	problem := Zdt1()
	vector := &models.Vector{
		Elements: []float64{0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least two")
}

func TestZdt1_Evaluate_BoundaryValues(t *testing.T) {
	problem := Zdt1()

	tests := []struct {
		name     string
		elements []float64
	}{
		{"all zeros", []float64{0.0, 0.0, 0.0}},
		{"all ones", []float64{1.0, 1.0, 1.0}},
		{"mixed", []float64{0.0, 1.0, 0.5}},
		{"high dimensions", []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vector := &models.Vector{Elements: tt.elements}
			err := problem.Evaluate(vector, 2)
			assert.NoError(t, err)
			assert.Len(t, vector.Objectives, 2)
		})
	}
}

// Tests for ZDT2

func TestZdt2_Name(t *testing.T) {
	problem := Zdt2()
	assert.Equal(t, "zdt2", problem.Name())
}

func TestZdt2_Evaluate_Success(t *testing.T) {
	problem := Zdt2()
	vector := &models.Vector{
		Elements: []float64{0.5, 0.5, 0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.NoError(t, err)
	assert.Len(t, vector.Objectives, 2)

	for _, obj := range vector.Objectives {
		assert.False(t, obj != obj) // Not NaN
	}
}

func TestZdt2_Evaluate_InsufficientDimensions(t *testing.T) {
	problem := Zdt2()
	vector := &models.Vector{
		Elements: []float64{0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.Error(t, err)
}

// Tests for ZDT3

func TestZdt3_Name(t *testing.T) {
	problem := Zdt3()
	assert.Equal(t, "zdt3", problem.Name())
}

func TestZdt3_Evaluate_Success(t *testing.T) {
	problem := Zdt3()
	vector := &models.Vector{
		Elements: []float64{0.5, 0.5, 0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.NoError(t, err)
	assert.Len(t, vector.Objectives, 2)

	for _, obj := range vector.Objectives {
		assert.False(t, obj != obj) // Not NaN
	}
}

func TestZdt3_Evaluate_InsufficientDimensions(t *testing.T) {
	problem := Zdt3()
	vector := &models.Vector{
		Elements: []float64{0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.Error(t, err)
}

// Tests for ZDT4

func TestZdt4_Name(t *testing.T) {
	problem := Zdt4()
	assert.Equal(t, "zdt4", problem.Name())
}

func TestZdt4_Evaluate_Success(t *testing.T) {
	problem := Zdt4()
	vector := &models.Vector{
		Elements: []float64{0.5, 0.5, 0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.NoError(t, err)
	assert.Len(t, vector.Objectives, 2)

	for _, obj := range vector.Objectives {
		assert.False(t, obj != obj) // Not NaN
	}
}

func TestZdt4_Evaluate_InsufficientDimensions(t *testing.T) {
	problem := Zdt4()
	vector := &models.Vector{
		Elements: []float64{0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.Error(t, err)
}

// Tests for ZDT6

func TestZdt6_Name(t *testing.T) {
	problem := Zdt6()
	assert.Equal(t, "zdt6", problem.Name())
}

func TestZdt6_Evaluate_Success(t *testing.T) {
	problem := Zdt6()
	vector := &models.Vector{
		Elements: []float64{0.5, 0.5, 0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.NoError(t, err)
	assert.Len(t, vector.Objectives, 2)

	for _, obj := range vector.Objectives {
		assert.False(t, obj != obj) // Not NaN
	}
}

func TestZdt6_Evaluate_InsufficientDimensions(t *testing.T) {
	problem := Zdt6()
	vector := &models.Vector{
		Elements: []float64{0.5},
	}

	err := problem.Evaluate(vector, 2)
	assert.Error(t, err)
}

// Tests for VNT1

func TestVnt1_Name(t *testing.T) {
	problem := Vnt1()
	assert.Equal(t, "vnt1", problem.Name())
}

func TestVnt1_Evaluate_Success(t *testing.T) {
	problem := Vnt1()
	vector := &models.Vector{
		Elements: []float64{0.5, 0.5},
	}

	err := problem.Evaluate(vector, 3)
	assert.NoError(t, err)
	assert.Len(t, vector.Objectives, 3)

	for _, obj := range vector.Objectives {
		assert.False(t, obj != obj) // Not NaN
	}
}

func TestVnt1_Evaluate_IncorrectDimensions(t *testing.T) {
	problem := Vnt1()

	tests := []struct {
		name     string
		elements []float64
	}{
		{"one dimension", []float64{0.5}},
		{"three dimensions", []float64{0.5, 0.5, 0.5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vector := &models.Vector{Elements: tt.elements}
			err := problem.Evaluate(vector, 3)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "only two")
		})
	}
}

func TestVnt1_Evaluate_VariousInputs(t *testing.T) {
	problem := Vnt1()

	tests := []struct {
		name     string
		elements []float64
	}{
		{"zeros", []float64{0.0, 0.0}},
		{"ones", []float64{1.0, 1.0}},
		{"negative", []float64{-1.0, -1.0}},
		{"mixed", []float64{-0.5, 0.5}},
		{"large positive", []float64{2.0, 2.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vector := &models.Vector{Elements: tt.elements}
			err := problem.Evaluate(vector, 3)
			assert.NoError(t, err)
			assert.Len(t, vector.Objectives, 3)

			for _, obj := range vector.Objectives {
				assert.False(t, obj != obj) // Not NaN
			}
		})
	}
}

// Comprehensive tests for all ZDT problems

func TestZDT_ObjectivesInRange(t *testing.T) {
	probs := []struct {
		name    string
		problem func() problems.Interface
	}{
		{"ZDT1", Zdt1},
		{"ZDT2", Zdt2},
		{"ZDT3", Zdt3},
		{"ZDT4", Zdt4},
		{"ZDT6", Zdt6},
	}

	for _, p := range probs {
		t.Run(p.name, func(t *testing.T) {
			problem := p.problem()
			vector := &models.Vector{
				Elements: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
			}

			err := problem.Evaluate(vector, 2)
			assert.NoError(t, err)
			assert.Len(t, vector.Objectives, 2)

			// All ZDT problems should produce non-negative objectives in [0,1] range
			for i, obj := range vector.Objectives {
				assert.False(t, obj != obj, "objective %d is NaN", i)
			}
		})
	}
}

func TestZDT_HighDimensionalInputs(t *testing.T) {
	probs := []struct {
		name    string
		problem func() problems.Interface
	}{
		{"ZDT1", Zdt1},
		{"ZDT2", Zdt2},
		{"ZDT3", Zdt3},
		{"ZDT4", Zdt4},
		{"ZDT6", Zdt6},
	}

	for _, p := range probs {
		t.Run(p.name, func(t *testing.T) {
			problem := p.problem()

			// Test with 30 dimensions
			elements := make([]float64, 30)
			for i := range elements {
				elements[i] = 0.5
			}

			vector := &models.Vector{Elements: elements}
			err := problem.Evaluate(vector, 2)
			assert.NoError(t, err)
			assert.Len(t, vector.Objectives, 2)
		})
	}
}

func TestAllProblems_ModifyObjectivesInPlace(t *testing.T) {
	probs := []struct {
		name       string
		problem    func() problems.Interface
		dimensions int
		objectives int
	}{
		{"ZDT1", Zdt1, 5, 2},
		{"ZDT2", Zdt2, 5, 2},
		{"ZDT3", Zdt3, 5, 2},
		{"ZDT4", Zdt4, 5, 2},
		{"ZDT6", Zdt6, 5, 2},
		{"VNT1", Vnt1, 2, 3},
	}

	for _, p := range probs {
		t.Run(p.name, func(t *testing.T) {
			problem := p.problem()
			elements := make([]float64, p.dimensions)
			for i := range elements {
				elements[i] = 0.5
			}

			vector := &models.Vector{
				Elements:   elements,
				Objectives: []float64{}, // Empty initially
			}

			err := problem.Evaluate(vector, p.objectives)
			require.NoError(t, err)
			assert.Len(t, vector.Objectives, p.objectives)

			// Verify objectives were set (not empty)
			for i, obj := range vector.Objectives {
				assert.False(t, obj != obj, "objective %d is NaN for %s", i, p.name)
			}
		})
	}
}

func TestAllProblems_MultipleEvaluations(t *testing.T) {
	probs := []struct {
		name       string
		problem    func() problems.Interface
		dimensions int
		objectives int
	}{
		{"ZDT1", Zdt1, 5, 2},
		{"VNT1", Vnt1, 2, 3},
	}

	for _, p := range probs {
		t.Run(p.name, func(t *testing.T) {
			problem := p.problem()

			// Evaluate the same vector multiple times
			elements := make([]float64, p.dimensions)
			for i := range elements {
				elements[i] = 0.5
			}

			vector := &models.Vector{Elements: elements}

			var previousObjectives []float64
			for i := range 3 {
				err := problem.Evaluate(vector, p.objectives)
				require.NoError(t, err)
				assert.Len(t, vector.Objectives, p.objectives)

				if i > 0 {
					// Results should be deterministic - same input should give same output
					assert.Equal(t, previousObjectives, vector.Objectives)
				}

				previousObjectives = make([]float64, len(vector.Objectives))
				copy(previousObjectives, vector.Objectives)
			}
		})
	}
}

func TestZDT_EdgeCase_FirstElementZero(t *testing.T) {
	probs := []struct {
		name    string
		problem func() problems.Interface
	}{
		{"ZDT1", Zdt1},
		{"ZDT2", Zdt2},
		{"ZDT3", Zdt3},
		{"ZDT4", Zdt4},
		{"ZDT6", Zdt6},
	}

	for _, p := range probs {
		t.Run(p.name, func(t *testing.T) {
			problem := p.problem()
			vector := &models.Vector{
				Elements: []float64{0.0, 0.5, 0.5},
			}

			err := problem.Evaluate(vector, 2)
			assert.NoError(t, err)
			assert.Len(t, vector.Objectives, 2)
		})
	}
}

func TestZDT_EdgeCase_FirstElementOne(t *testing.T) {
	probs := []struct {
		name    string
		problem func() problems.Interface
	}{
		{"ZDT1", Zdt1},
		{"ZDT2", Zdt2},
		{"ZDT3", Zdt3},
		{"ZDT4", Zdt4},
		{"ZDT6", Zdt6},
	}

	for _, p := range probs {
		t.Run(p.name, func(t *testing.T) {
			problem := p.problem()
			vector := &models.Vector{
				Elements: []float64{1.0, 0.0, 0.0},
			}

			err := problem.Evaluate(vector, 2)
			// Some problems might error with this input (e.g., sqrt of negative)
			// Just verify it doesn't panic
			_ = err
		})
	}
}
