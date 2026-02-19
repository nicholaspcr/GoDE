package problems_test

import (
	"fmt"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockProblem is a minimal Interface implementation for testing.
type mockProblem struct{ name string }

func (m *mockProblem) Name() string                          { return m.name }
func (m *mockProblem) Evaluate(_ *models.Vector, _ int) error { return nil }

// mockFactory creates a mockProblem for the given name.
func mockFactory(name string) problems.ProblemFactory {
	return func(dim, objs int) (problems.Interface, error) {
		return &mockProblem{name: name}, nil
	}
}

// failingFactory always returns an error.
func failingFactory(dim, objs int) (problems.Interface, error) {
	return nil, fmt.Errorf("factory error")
}

func newRegistry() *problems.Registry {
	return problems.NewRegistry()
}

func TestRegistry_Register_And_Create(t *testing.T) {
	r := newRegistry()
	meta := problems.ProblemMetadata{Description: "test problem", MinDim: 2, MaxDim: 100, NumObjs: 2, Category: "multi"}

	r.Register("test-problem", mockFactory("test-problem"), meta)

	p, err := r.Create("test-problem", 5, 2)
	require.NoError(t, err)
	assert.Equal(t, "test-problem", p.Name())
}

func TestRegistry_Create_NotFound(t *testing.T) {
	r := newRegistry()

	_, err := r.Create("nonexistent", 5, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "problem does not exist")
}

func TestRegistry_Create_FactoryError(t *testing.T) {
	r := newRegistry()
	meta := problems.ProblemMetadata{Category: "test"}
	r.Register("failing", failingFactory, meta)

	_, err := r.Create("failing", 5, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "factory error")
}

func TestRegistry_List(t *testing.T) {
	r := newRegistry()
	meta := problems.ProblemMetadata{Category: "multi"}

	r.Register("bravo", mockFactory("bravo"), meta)
	r.Register("alpha", mockFactory("alpha"), meta)
	r.Register("charlie", mockFactory("charlie"), meta)

	names := r.List()
	assert.Equal(t, []string{"alpha", "bravo", "charlie"}, names, "List should return sorted names")
}

func TestRegistry_List_Empty(t *testing.T) {
	r := newRegistry()
	names := r.List()
	assert.Empty(t, names)
}

func TestRegistry_Get_Found(t *testing.T) {
	r := newRegistry()
	meta := problems.ProblemMetadata{
		Description: "ZDT1 benchmark",
		MinDim:      2,
		MaxDim:      30,
		NumObjs:     2,
		Category:    "multi",
	}
	r.Register("zdt1", mockFactory("zdt1"), meta)

	got, ok := r.Get("zdt1")
	require.True(t, ok)
	assert.Equal(t, "zdt1", got.Name)
	assert.Equal(t, "ZDT1 benchmark", got.Description)
	assert.Equal(t, 2, got.MinDim)
	assert.Equal(t, 30, got.MaxDim)
	assert.Equal(t, 2, got.NumObjs)
	assert.Equal(t, "multi", got.Category)
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := newRegistry()
	_, ok := r.Get("does-not-exist")
	assert.False(t, ok)
}

func TestRegistry_ListMetadata(t *testing.T) {
	r := newRegistry()
	r.Register("b", mockFactory("b"), problems.ProblemMetadata{Description: "B problem", Category: "many"})
	r.Register("a", mockFactory("a"), problems.ProblemMetadata{Description: "A problem", Category: "multi"})

	metas := r.ListMetadata()
	require.Len(t, metas, 2)
	assert.Equal(t, "a", metas[0].Name, "ListMetadata should be sorted by name")
	assert.Equal(t, "b", metas[1].Name)
	assert.Equal(t, "A problem", metas[0].Description)
}

func TestRegistry_ListMetadata_Empty(t *testing.T) {
	r := newRegistry()
	metas := r.ListMetadata()
	assert.Empty(t, metas)
}

func TestRegistry_Register_SetsName(t *testing.T) {
	r := newRegistry()
	// Register with mismatched name in metadata
	meta := problems.ProblemMetadata{Name: "wrong-name", Category: "multi"}
	r.Register("correct-name", mockFactory("correct-name"), meta)

	got, ok := r.Get("correct-name")
	require.True(t, ok)
	assert.Equal(t, "correct-name", got.Name, "Register should override metadata name with registered name")
}

func TestRegistry_Register_Overwrite(t *testing.T) {
	r := newRegistry()
	meta := problems.ProblemMetadata{Description: "first", Category: "multi"}
	r.Register("p", mockFactory("first"), meta)

	meta2 := problems.ProblemMetadata{Description: "second", Category: "many"}
	r.Register("p", mockFactory("second"), meta2)

	got, ok := r.Get("p")
	require.True(t, ok)
	assert.Equal(t, "second", got.Description, "Re-registering should overwrite")
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	r := newRegistry()
	meta := problems.ProblemMetadata{Category: "multi"}

	// Register several problems concurrently.
	done := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("prob-%d", i)
		go func(n string) {
			r.Register(n, mockFactory(n), meta)
			done <- struct{}{}
		}(name)
	}
	for i := 0; i < 10; i++ {
		<-done
	}

	names := r.List()
	assert.Len(t, names, 10)
}
