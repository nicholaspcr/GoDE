package variants

import (
	"math/rand"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubVariant implements Interface for testing.
type stubVariant struct {
	name string
}

func (s *stubVariant) Name() string { return s.name }
func (s *stubVariant) Mutate(elems, rankZero []models.Vector, params Parameters) (models.Vector, error) {
	return models.Vector{}, nil
}

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	assert.NotNil(t, r)
	assert.Empty(t, r.List())
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()
	factory := func() Interface { return &stubVariant{name: "rand/1"} }
	meta := VariantMetadata{Description: "Random variant", Category: "rand"}

	r.Register("rand/1", factory, meta)

	assert.Contains(t, r.List(), "rand/1")
	got, ok := r.Get("rand/1")
	assert.True(t, ok)
	assert.Equal(t, "rand/1", got.Name)
	assert.Equal(t, "Random variant", got.Description)
	assert.Equal(t, "rand", got.Category)
}

func TestRegistry_Register_overrides_name(t *testing.T) {
	r := NewRegistry()
	factory := func() Interface { return &stubVariant{name: "test"} }
	meta := VariantMetadata{Name: "wrong", Description: "desc"}

	r.Register("correct", factory, meta)

	got, ok := r.Get("correct")
	assert.True(t, ok)
	assert.Equal(t, "correct", got.Name, "Register should override Name field")
}

func TestRegistry_Create(t *testing.T) {
	r := NewRegistry()

	t.Run("creates registered variant", func(t *testing.T) {
		r.Register("best/1", func() Interface {
			return &stubVariant{name: "best/1"}
		}, VariantMetadata{Category: "best"})

		v, err := r.Create("best/1")
		require.NoError(t, err)
		assert.Equal(t, "best/1", v.Name())
	})

	t.Run("returns error for unregistered variant", func(t *testing.T) {
		_, err := r.Create("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "variant does not exist")
	})
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()

	t.Run("empty registry", func(t *testing.T) {
		assert.Empty(t, r.List())
	})

	t.Run("returns sorted names", func(t *testing.T) {
		r.Register("rand/2", func() Interface { return &stubVariant{} }, VariantMetadata{})
		r.Register("best/1", func() Interface { return &stubVariant{} }, VariantMetadata{})
		r.Register("pbest", func() Interface { return &stubVariant{} }, VariantMetadata{})

		names := r.List()
		assert.Equal(t, []string{"best/1", "pbest", "rand/2"}, names)
	})
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	t.Run("returns false for missing variant", func(t *testing.T) {
		_, ok := r.Get("missing")
		assert.False(t, ok)
	})

	t.Run("returns metadata for registered variant", func(t *testing.T) {
		r.Register("test", func() Interface { return &stubVariant{} }, VariantMetadata{
			Description: "Test variant",
			Category:    "test",
		})

		meta, ok := r.Get("test")
		assert.True(t, ok)
		assert.Equal(t, "Test variant", meta.Description)
		assert.Equal(t, "test", meta.Category)
	})
}

func TestRegistry_ListMetadata(t *testing.T) {
	r := NewRegistry()

	t.Run("empty registry", func(t *testing.T) {
		assert.Empty(t, r.ListMetadata())
	})

	t.Run("returns sorted metadata", func(t *testing.T) {
		r.Register("zeta", func() Interface { return &stubVariant{} }, VariantMetadata{Description: "Z"})
		r.Register("alpha", func() Interface { return &stubVariant{} }, VariantMetadata{Description: "A"})

		metas := r.ListMetadata()
		require.Len(t, metas, 2)
		assert.Equal(t, "alpha", metas[0].Name)
		assert.Equal(t, "zeta", metas[1].Name)
	})
}

func TestGetStandardPValues(t *testing.T) {
	values := GetStandardPValues()
	assert.Equal(t, []float64{0.05, 0.10, 0.15, 0.20}, values)
	assert.Len(t, values, 4)
}

func TestValidateVectors(t *testing.T) {
	t.Run("valid vectors", func(t *testing.T) {
		vectors := []models.Vector{
			{Elements: []float64{1.0, 2.0, 3.0}},
			{Elements: []float64{4.0, 5.0, 6.0}},
		}
		err := ValidateVectors(vectors, []int{0, 1}, 3)
		assert.NoError(t, err)
	})

	t.Run("index out of range", func(t *testing.T) {
		vectors := []models.Vector{
			{Elements: []float64{1.0}},
		}
		err := ValidateVectors(vectors, []int{5}, 1)
		assert.ErrorIs(t, err, ErrInvalidVector)
	})

	t.Run("negative index", func(t *testing.T) {
		vectors := []models.Vector{
			{Elements: []float64{1.0}},
		}
		err := ValidateVectors(vectors, []int{-1}, 1)
		assert.ErrorIs(t, err, ErrInvalidVector)
	})

	t.Run("nil elements", func(t *testing.T) {
		vectors := []models.Vector{
			{Elements: nil},
		}
		err := ValidateVectors(vectors, []int{0}, 3)
		assert.ErrorIs(t, err, ErrInvalidVector)
	})

	t.Run("wrong dimension", func(t *testing.T) {
		vectors := []models.Vector{
			{Elements: []float64{1.0, 2.0}},
		}
		err := ValidateVectors(vectors, []int{0}, 3)
		assert.ErrorIs(t, err, ErrInvalidVector)
	})

	t.Run("empty indices", func(t *testing.T) {
		vectors := []models.Vector{
			{Elements: []float64{1.0}},
		}
		err := ValidateVectors(vectors, []int{}, 1)
		assert.NoError(t, err)
	})
}

func TestGenerateIndices_maxRetries(t *testing.T) {
	// When NP equals len(r), all slots are used - this forces the tight loop.
	// With startInd=0 and len(r)==NP, it should still work.
	r := make([]int, 5)
	random := rand.New(rand.NewSource(42))
	err := GenerateIndices(0, 5, r, random)
	assert.NoError(t, err)

	// Verify all indices are unique and within range
	seen := make(map[int]bool)
	for _, val := range r {
		assert.GreaterOrEqual(t, val, 0)
		assert.Less(t, val, 5)
		assert.False(t, seen[val], "duplicate index found")
		seen[val] = true
	}
}
