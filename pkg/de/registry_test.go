package de

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	assert.NotNil(t, r)
	assert.Empty(t, r.List())
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()
	r.Register("gde3", AlgorithmMetadata{Description: "GDE3 algorithm"})

	assert.True(t, r.IsSupported("gde3"))
	meta, err := r.Get("gde3")
	require.NoError(t, err)
	assert.Equal(t, "gde3", meta.Name)
	assert.Equal(t, "GDE3 algorithm", meta.Description)
}

func TestRegistry_Register_overrides_name(t *testing.T) {
	r := NewRegistry()
	r.Register("myalgo", AlgorithmMetadata{Name: "wrong-name", Description: "desc"})

	meta, err := r.Get("myalgo")
	require.NoError(t, err)
	assert.Equal(t, "myalgo", meta.Name, "Register should override Name field")
}

func TestRegistry_IsSupported(t *testing.T) {
	r := NewRegistry()
	assert.False(t, r.IsSupported("nonexistent"))

	r.Register("algo1", AlgorithmMetadata{})
	assert.True(t, r.IsSupported("algo1"))
	assert.False(t, r.IsSupported("algo2"))
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	t.Run("returns error for missing algorithm", func(t *testing.T) {
		_, err := r.Get("missing")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "algorithm not found: missing")
	})

	t.Run("returns metadata for registered algorithm", func(t *testing.T) {
		r.Register("gde3", AlgorithmMetadata{Description: "Multi-objective DE"})
		meta, err := r.Get("gde3")
		require.NoError(t, err)
		assert.Equal(t, "Multi-objective DE", meta.Description)
	})
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()

	t.Run("empty registry", func(t *testing.T) {
		assert.Empty(t, r.List())
	})

	t.Run("returns sorted names", func(t *testing.T) {
		r.Register("zeta", AlgorithmMetadata{})
		r.Register("alpha", AlgorithmMetadata{})
		r.Register("mu", AlgorithmMetadata{})

		names := r.List()
		assert.Equal(t, []string{"alpha", "mu", "zeta"}, names)
	})
}

func TestRegistry_ListMetadata(t *testing.T) {
	r := NewRegistry()

	t.Run("empty registry", func(t *testing.T) {
		assert.Empty(t, r.ListMetadata())
	})

	t.Run("returns sorted metadata", func(t *testing.T) {
		r.Register("beta", AlgorithmMetadata{Description: "Beta desc"})
		r.Register("alpha", AlgorithmMetadata{Description: "Alpha desc"})

		metas := r.ListMetadata()
		require.Len(t, metas, 2)
		assert.Equal(t, "alpha", metas[0].Name)
		assert.Equal(t, "beta", metas[1].Name)
		assert.Equal(t, "Alpha desc", metas[0].Description)
		assert.Equal(t, "Beta desc", metas[1].Description)
	})
}
