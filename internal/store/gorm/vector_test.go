package gorm

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupVectorTestDB creates an in-memory DB with pareto + vector models.
func setupVectorTestDB(t *testing.T) *vectorStore {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	err = db.AutoMigrate(&paretoModel{}, &vectorModel{})
	require.NoError(t, err)
	return newVectorStore(db)
}

// createTestPareto inserts a minimal paretoModel and returns its ID.
func createTestPareto(t *testing.T, db *gorm.DB, _ string) uint {
	t.Helper()
	p := &paretoModel{} // UserID is uint; skip FK enforcement in SQLite
	require.NoError(t, db.Create(p).Error)
	return p.ID
}

func TestVectorModel_SetGetElements(t *testing.T) {
	v := &vectorModel{}
	elements := []float64{1.1, 2.2, 3.3}

	require.NoError(t, v.SetElements(elements))
	assert.NotEmpty(t, v.ElementsJSON)

	got, err := v.GetElements()
	require.NoError(t, err)
	assert.Equal(t, elements, got)
}

func TestVectorModel_SetGetElements_Empty(t *testing.T) {
	v := &vectorModel{}
	got, err := v.GetElements()
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestVectorModel_SetGetObjectives(t *testing.T) {
	v := &vectorModel{}
	objectives := []float64{0.5, 0.8}

	require.NoError(t, v.SetObjectives(objectives))
	assert.NotEmpty(t, v.ObjectivesJSON)

	got, err := v.GetObjectives()
	require.NoError(t, err)
	assert.Equal(t, objectives, got)
}

func TestVectorModel_SetGetObjectives_Empty(t *testing.T) {
	v := &vectorModel{}
	got, err := v.GetObjectives()
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestVectorStore_CreateAndGet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&paretoModel{}, &vectorModel{}))
	s := newVectorStore(db)
	ctx := context.Background()

	paretoID := createTestPareto(t, db, "user1")

	vec := &api.Vector{
		Elements:         []float64{1.0, 2.0, 3.0},
		Objectives:       []float64{0.5, 0.9},
		CrowdingDistance: 1.5,
	}

	require.NoError(t, s.CreateVector(ctx, vec, paretoID))

	// Get by ID 1 (first inserted)
	got, err := s.GetVector(ctx, &api.VectorIDs{Id: 1})
	require.NoError(t, err)
	assert.Equal(t, vec.Elements, got.Elements)
	assert.Equal(t, vec.Objectives, got.Objectives)
	assert.Equal(t, vec.CrowdingDistance, got.CrowdingDistance)
}

func TestVectorStore_GetVector_NotFound(t *testing.T) {
	s := setupVectorTestDB(t)
	ctx := context.Background()

	_, err := s.GetVector(ctx, &api.VectorIDs{Id: 999})
	require.Error(t, err)
}

func TestVectorStore_UpdateVector_AllFields(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&paretoModel{}, &vectorModel{}))
	s := newVectorStore(db)
	ctx := context.Background()

	paretoID := createTestPareto(t, db, "user1")

	vec := &api.Vector{
		Elements:         []float64{1.0, 2.0},
		Objectives:       []float64{0.5},
		CrowdingDistance: 1.0,
	}
	require.NoError(t, s.CreateVector(ctx, vec, paretoID))

	// Update with no specific fields = update all
	updated := &api.Vector{
		Ids:              &api.VectorIDs{Id: 1},
		Elements:         []float64{9.0, 8.0},
		Objectives:       []float64{0.2},
		CrowdingDistance: 3.0,
	}
	require.NoError(t, s.UpdateVector(ctx, updated))

	got, err := s.GetVector(ctx, &api.VectorIDs{Id: 1})
	require.NoError(t, err)
	assert.Equal(t, updated.Elements, got.Elements)
	assert.Equal(t, updated.Objectives, got.Objectives)
	assert.Equal(t, updated.CrowdingDistance, got.CrowdingDistance)
}

func TestVectorStore_UpdateVector_SpecificFields(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&paretoModel{}, &vectorModel{}))
	s := newVectorStore(db)
	ctx := context.Background()

	paretoID := createTestPareto(t, db, "user1")

	vec := &api.Vector{
		Elements:         []float64{1.0, 2.0},
		Objectives:       []float64{0.5},
		CrowdingDistance: 1.0,
	}
	require.NoError(t, s.CreateVector(ctx, vec, paretoID))

	// Update only crowding distance
	updated := &api.Vector{
		Ids:              &api.VectorIDs{Id: 1},
		Elements:         []float64{9.0, 8.0}, // should be ignored
		CrowdingDistance: 5.0,
	}
	require.NoError(t, s.UpdateVector(ctx, updated, "crowding_distance"))

	got, err := s.GetVector(ctx, &api.VectorIDs{Id: 1})
	require.NoError(t, err)
	assert.Equal(t, vec.Elements, got.Elements, "Elements should not be updated")
	assert.Equal(t, 5.0, got.CrowdingDistance)
}

func TestVectorStore_UpdateVector_NoID(t *testing.T) {
	s := setupVectorTestDB(t)
	ctx := context.Background()

	err := s.UpdateVector(ctx, &api.Vector{})
	require.Error(t, err)
}

func TestVectorStore_DeleteVector(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&paretoModel{}, &vectorModel{}))
	s := newVectorStore(db)
	ctx := context.Background()

	paretoID := createTestPareto(t, db, "user1")

	vec := &api.Vector{
		Elements:   []float64{1.0},
		Objectives: []float64{0.5},
	}
	require.NoError(t, s.CreateVector(ctx, vec, paretoID))

	require.NoError(t, s.DeleteVector(ctx, &api.VectorIDs{Id: 1}))

	_, err = s.GetVector(ctx, &api.VectorIDs{Id: 1})
	require.Error(t, err, "vector should not be found after deletion")
}
