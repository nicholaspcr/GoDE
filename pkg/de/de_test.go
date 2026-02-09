package de

import (
	"context"
	"errors"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockAlgorithm is a test double for Algorithm.
type mockAlgorithm struct {
	executeFunc func(ctx context.Context, pareto chan<- []models.Vector, maxObj chan<- []float64) error
}

func (m *mockAlgorithm) Execute(ctx context.Context, pareto chan<- []models.Vector, maxObj chan<- []float64) error {
	return m.executeFunc(ctx, pareto, maxObj)
}

func TestNew(t *testing.T) {
	algo := &mockAlgorithm{executeFunc: func(ctx context.Context, p chan<- []models.Vector, m chan<- []float64) error {
		return nil
	}}

	t.Run("succeeds with algorithm set", func(t *testing.T) {
		d, err := New(Config{}, WithAlgorithm(algo))
		require.NoError(t, err)
		assert.NotNil(t, d)
	})

	t.Run("fails without algorithm", func(t *testing.T) {
		d, err := New(Config{})
		assert.EqualError(t, err, "no algorithm set")
		assert.Nil(t, d)
	})

	t.Run("applies all options", func(t *testing.T) {
		cfg := Config{ParetoChannelLimiter: 5, MaxChannelLimiter: 5, ResultLimiter: 10}
		d, err := New(cfg,
			WithAlgorithm(algo),
			WithExecutions(3),
			WithGenerations(50),
			WithDimensions(10),
			WithObjFuncAmount(2),
		)
		require.NoError(t, err)
		assert.Equal(t, 3, d.constants.Executions)
		assert.Equal(t, 50, d.constants.Generations)
		assert.Equal(t, 10, d.constants.Dimensions)
		assert.Equal(t, 2, d.constants.ObjFuncAmount)
		assert.Equal(t, 5, d.config.ParetoChannelLimiter)
	})

	t.Run("progress callback is set", func(t *testing.T) {
		called := false
		cb := func(gen, total, size int, pareto []models.Vector) { called = true }
		d, err := New(Config{}, WithAlgorithm(algo), WithProgressCallback(cb))
		require.NoError(t, err)
		assert.NotNil(t, d.progressCallback)
		d.progressCallback(1, 10, 5, nil)
		assert.True(t, called)
	})
}

func TestExecute(t *testing.T) {
	t.Run("returns results from successful execution", func(t *testing.T) {
		algo := &mockAlgorithm{
			executeFunc: func(ctx context.Context, pareto chan<- []models.Vector, maxObj chan<- []float64) error {
				pareto <- []models.Vector{
					{Elements: []float64{1.0}, Objectives: []float64{0.5, 0.5}},
				}
				maxObj <- []float64{1.0, 1.0}
				return nil
			},
		}

		d, err := New(
			Config{ParetoChannelLimiter: 10, MaxChannelLimiter: 10, ResultLimiter: 100},
			WithAlgorithm(algo),
			WithExecutions(1),
		)
		require.NoError(t, err)

		pareto, maxObjs, err := d.Execute(context.Background())
		require.NoError(t, err)
		assert.NotEmpty(t, pareto)
		assert.Len(t, maxObjs, 1)
	})

	t.Run("returns error when all executions fail", func(t *testing.T) {
		algo := &mockAlgorithm{
			executeFunc: func(ctx context.Context, pareto chan<- []models.Vector, maxObj chan<- []float64) error {
				return errors.New("execution failed")
			},
		}

		d, err := New(
			Config{ParetoChannelLimiter: 10, MaxChannelLimiter: 10, ResultLimiter: 100},
			WithAlgorithm(algo),
			WithExecutions(2),
		)
		require.NoError(t, err)

		_, _, err = d.Execute(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "execution failed")
	})

	t.Run("returns error when context already cancelled", func(t *testing.T) {
		algo := &mockAlgorithm{
			executeFunc: func(ctx context.Context, pareto chan<- []models.Vector, maxObj chan<- []float64) error {
				return nil
			},
		}

		d, err := New(
			Config{ParetoChannelLimiter: 10, MaxChannelLimiter: 10, ResultLimiter: 100},
			WithAlgorithm(algo),
			WithExecutions(1),
		)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err = d.Execute(ctx)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("handles partial execution failures", func(t *testing.T) {
		callCount := 0
		algo := &mockAlgorithm{
			executeFunc: func(ctx context.Context, pareto chan<- []models.Vector, maxObj chan<- []float64) error {
				n := FromContextExecutionNumber(ctx)
				if n == 0 {
					callCount++
					return errors.New("first execution failed")
				}
				pareto <- []models.Vector{
					{Elements: []float64{1.0}, Objectives: []float64{0.3, 0.7}},
				}
				maxObj <- []float64{1.0, 1.0}
				callCount++
				return nil
			},
		}

		d, err := New(
			Config{ParetoChannelLimiter: 10, MaxChannelLimiter: 10, ResultLimiter: 100},
			WithAlgorithm(algo),
			WithExecutions(2),
		)
		require.NoError(t, err)

		pareto, maxObjs, err := d.Execute(context.Background())
		require.NoError(t, err)
		assert.NotEmpty(t, pareto)
		assert.Len(t, maxObjs, 1)
	})

	t.Run("multiple executions merge pareto results", func(t *testing.T) {
		algo := &mockAlgorithm{
			executeFunc: func(ctx context.Context, pareto chan<- []models.Vector, maxObj chan<- []float64) error {
				n := FromContextExecutionNumber(ctx)
				pareto <- []models.Vector{
					{Elements: []float64{float64(n)}, Objectives: []float64{float64(n) * 0.1, 1.0 - float64(n)*0.1}},
				}
				maxObj <- []float64{1.0, 1.0}
				return nil
			},
		}

		d, err := New(
			Config{ParetoChannelLimiter: 10, MaxChannelLimiter: 10, ResultLimiter: 100},
			WithAlgorithm(algo),
			WithExecutions(3),
		)
		require.NoError(t, err)

		pareto, maxObjs, err := d.Execute(context.Background())
		require.NoError(t, err)
		assert.NotEmpty(t, pareto)
		assert.Len(t, maxObjs, 3)
	})
}

func TestFilterCollectedPareto(t *testing.T) {
	d := &de{
		config: Config{ResultLimiter: 5},
	}

	t.Run("merges and filters multiple pareto sets", func(t *testing.T) {
		allPareto := [][]models.Vector{
			{
				{Objectives: []float64{1.0, 5.0}},
				{Objectives: []float64{3.0, 3.0}},
			},
			{
				{Objectives: []float64{2.0, 4.0}},
				{Objectives: []float64{5.0, 1.0}},
			},
		}

		result := d.filterCollectedPareto(context.Background(), allPareto)
		assert.NotEmpty(t, result)
		assert.LessOrEqual(t, len(result), 5)
	})

	t.Run("handles empty pareto sets", func(t *testing.T) {
		result := d.filterCollectedPareto(context.Background(), nil)
		assert.Empty(t, result)
	})

	t.Run("handles single pareto set", func(t *testing.T) {
		allPareto := [][]models.Vector{
			{
				{Objectives: []float64{1.0, 2.0}},
				{Objectives: []float64{2.0, 1.0}},
			},
		}

		result := d.filterCollectedPareto(context.Background(), allPareto)
		assert.NotEmpty(t, result)
	})
}
