package de

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstants_Validate(t *testing.T) {
	validConstants := Constants{
		Executions:    5,
		Generations:   100,
		Dimensions:    10,
		ObjFuncAmount: 2,
	}

	t.Run("valid constants", func(t *testing.T) {
		err := validConstants.Validate()
		require.NoError(t, err)
	})

	t.Run("executions must be positive", func(t *testing.T) {
		c := validConstants
		c.Executions = 0
		assert.EqualError(t, c.Validate(), "executions must be positive")

		c.Executions = -1
		assert.EqualError(t, c.Validate(), "executions must be positive")
	})

	t.Run("generations must be positive", func(t *testing.T) {
		c := validConstants
		c.Generations = 0
		assert.EqualError(t, c.Validate(), "generations must be positive")

		c.Generations = -5
		assert.EqualError(t, c.Validate(), "generations must be positive")
	})

	t.Run("dimensions must be positive", func(t *testing.T) {
		c := validConstants
		c.Dimensions = 0
		assert.EqualError(t, c.Validate(), "dimensions must be positive")

		c.Dimensions = -2
		assert.EqualError(t, c.Validate(), "dimensions must be positive")
	})

	t.Run("obj func amount must be positive", func(t *testing.T) {
		c := validConstants
		c.ObjFuncAmount = 0
		assert.EqualError(t, c.Validate(), "objective function amount must be positive")

		c.ObjFuncAmount = -1
		assert.EqualError(t, c.Validate(), "objective function amount must be positive")
	})

	t.Run("minimum valid values", func(t *testing.T) {
		c := Constants{
			Executions:    1,
			Generations:   1,
			Dimensions:    1,
			ObjFuncAmount: 1,
		}
		require.NoError(t, c.Validate())
	})
}
