package gde3

import (
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstants_Validate(t *testing.T) {
	validDE := de.Constants{
		Executions:    5,
		Generations:   100,
		Dimensions:    10,
		ObjFuncAmount: 2,
	}

	validGDE3 := Constants{
		DE: validDE,
		CR: 0.9,
		F:  0.5,
		P:  0.1,
	}

	t.Run("valid constants", func(t *testing.T) {
		require.NoError(t, validGDE3.Validate())
	})

	t.Run("propagates DE validation errors", func(t *testing.T) {
		c := validGDE3
		c.DE.Executions = 0
		assert.EqualError(t, c.Validate(), "executions must be positive")
	})

	t.Run("CR must be in [0, 1]", func(t *testing.T) {
		c := validGDE3
		c.CR = -0.1
		assert.ErrorContains(t, c.Validate(), "CR (crossover rate) must be in range [0, 1]")

		c.CR = 1.1
		assert.ErrorContains(t, c.Validate(), "CR (crossover rate) must be in range [0, 1]")
	})

	t.Run("CR boundary values", func(t *testing.T) {
		c := validGDE3
		c.CR = 0.0
		require.NoError(t, c.Validate(), "CR=0 should be valid")

		c.CR = 1.0
		require.NoError(t, c.Validate(), "CR=1 should be valid")
	})

	t.Run("F must be in (0, 2]", func(t *testing.T) {
		c := validGDE3
		c.F = 0.0
		assert.ErrorContains(t, c.Validate(), "F (scaling factor) must be in range (0, 2]")

		c.F = -0.5
		assert.ErrorContains(t, c.Validate(), "F (scaling factor) must be in range (0, 2]")

		c.F = 2.1
		assert.ErrorContains(t, c.Validate(), "F (scaling factor) must be in range (0, 2]")
	})

	t.Run("F boundary values", func(t *testing.T) {
		c := validGDE3
		c.F = 0.001
		require.NoError(t, c.Validate(), "F=0.001 should be valid")

		c.F = 2.0
		require.NoError(t, c.Validate(), "F=2 should be valid")
	})

	t.Run("P must be in (0, 1]", func(t *testing.T) {
		c := validGDE3
		c.P = 0.0
		assert.ErrorContains(t, c.Validate(), "P (selection parameter) must be in range (0, 1]")

		c.P = -0.1
		assert.ErrorContains(t, c.Validate(), "P (selection parameter) must be in range (0, 1]")

		c.P = 1.1
		assert.ErrorContains(t, c.Validate(), "P (selection parameter) must be in range (0, 1]")
	})

	t.Run("P boundary values", func(t *testing.T) {
		c := validGDE3
		c.P = 0.001
		require.NoError(t, c.Validate(), "P=0.001 should be valid")

		c.P = 1.0
		require.NoError(t, c.Validate(), "P=1 should be valid")
	})
}
