package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProblemByName(t *testing.T) {
	t.Run("valid problem names", func(t *testing.T) {
		validNames := []string{
			"zdt1", "zdt2", "zdt3", "zdt4", "zdt6",
			"vnt1",
			"dtlz1", "dtlz2", "dtlz3", "dtlz4", "dtlz5", "dtlz6", "dtlz7",
			"wfg1", "wfg2", "wfg3", "wfg4", "wfg5", "wfg6", "wfg7", "wfg8", "wfg9",
		}
		for _, name := range validNames {
			p, err := GetProblemByName(name)
			require.NoError(t, err, "problem %s should be found", name)
			assert.NotNil(t, p, "problem %s should not be nil", name)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		p, err := GetProblemByName("ZDT1")
		require.NoError(t, err)
		assert.NotNil(t, p)

		p, err = GetProblemByName("Dtlz1")
		require.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("invalid name", func(t *testing.T) {
		p, err := GetProblemByName("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, p)
		assert.Contains(t, err.Error(), "problem not found")
	})

	t.Run("empty name", func(t *testing.T) {
		p, err := GetProblemByName("")
		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestGetVariantByName(t *testing.T) {
	t.Run("valid variant names", func(t *testing.T) {
		validNames := []string{
			"rand/1", "rand/2",
			"best/1", "best/2",
			"pbest/1",
			"current-to-best/1",
		}
		for _, name := range validNames {
			v, err := GetVariantByName(name)
			require.NoError(t, err, "variant %s should be found", name)
			assert.NotNil(t, v, "variant %s should not be nil", name)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		v, err := GetVariantByName("RAND/1")
		require.NoError(t, err)
		assert.NotNil(t, v)
	})

	t.Run("invalid name", func(t *testing.T) {
		v, err := GetVariantByName("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, v)
		assert.Contains(t, err.Error(), "variant not found")
	})

	t.Run("empty name", func(t *testing.T) {
		v, err := GetVariantByName("")
		assert.Error(t, err)
		assert.Nil(t, v)
	})
}

func TestGetAllProblems(t *testing.T) {
	problems := GetAllProblems()
	assert.NotEmpty(t, problems)
	assert.Len(t, problems, len(problemSet))

	for _, p := range problems {
		assert.NotNil(t, p)
	}
}

func TestGetAllVariants(t *testing.T) {
	variants := GetAllVariants()
	assert.NotEmpty(t, variants)
	assert.Len(t, variants, len(variantSet))

	for _, v := range variants {
		assert.NotNil(t, v)
	}
}
