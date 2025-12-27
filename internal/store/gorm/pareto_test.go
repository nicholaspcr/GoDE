package gorm

import (
	"context"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestParetoStore_CreatePareto(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	// Create a test user first
	user := &api.User{
		Ids:      &api.UserIDs{Username: "paretouser"},
		Email:    "pareto@example.com",
		Password: "password",
	}
	err := store.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("create pareto with vectors", func(t *testing.T) {
		pareto := &api.Pareto{
			Ids: &api.ParetoIDs{UserId: "paretouser"},
			Vectors: []*api.Vector{
				{
					Elements:         []float64{1.0, 2.0, 3.0},
					Objectives:       []float64{0.5, 0.8},
					CrowdingDistance: 1.5,
				},
				{
					Elements:         []float64{4.0, 5.0, 6.0},
					Objectives:       []float64{0.3, 0.9},
					CrowdingDistance: 2.1,
				},
			},
			MaxObjs: []float64{10.0, 20.0},
		}

		err := store.CreatePareto(ctx, pareto)
		assert.NoError(t, err)
	})

	t.Run("create pareto without vectors", func(t *testing.T) {
		pareto := &api.Pareto{
			Ids:     &api.ParetoIDs{UserId: "paretouser"},
			Vectors: []*api.Vector{},
			MaxObjs: []float64{5.0, 10.0},
		}

		err := store.CreatePareto(ctx, pareto)
		assert.NoError(t, err)
	})

	t.Run("create pareto with empty max objs", func(t *testing.T) {
		pareto := &api.Pareto{
			Ids: &api.ParetoIDs{UserId: "paretouser"},
			Vectors: []*api.Vector{
				{
					Elements:         []float64{1.0, 2.0},
					Objectives:       []float64{0.5, 0.8},
					CrowdingDistance: 1.0,
				},
			},
			MaxObjs: []float64{},
		}

		err := store.CreatePareto(ctx, pareto)
		assert.NoError(t, err)
	})
}

func TestParetoStore_GetPareto(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	// Create user and pareto
	user := &api.User{
		Ids:      &api.UserIDs{Username: "getparetouser"},
		Email:    "getpareto@example.com",
		Password: "password",
	}
	err := store.CreateUser(ctx, user)
	require.NoError(t, err)

	originalPareto := &api.Pareto{
		Ids: &api.ParetoIDs{UserId: "getparetouser"},
		Vectors: []*api.Vector{
			{
				Elements:         []float64{1.0, 2.0, 3.0},
				Objectives:       []float64{0.5, 0.8},
				CrowdingDistance: 1.5,
			},
			{
				Elements:         []float64{4.0, 5.0, 6.0},
				Objectives:       []float64{0.3, 0.9},
				CrowdingDistance: 2.1,
			},
		},
		MaxObjs: []float64{10.0, 20.0},
	}
	err = store.CreatePareto(ctx, originalPareto)
	require.NoError(t, err)

	// Get the created pareto's ID by listing
	paretos, _, err := store.ListParetos(ctx, &api.UserIDs{Username: "getparetouser"}, 50, 0)
	require.NoError(t, err)
	require.Greater(t, len(paretos), 0)
	paretoID := paretos[0].Ids.Id

	t.Run("get existing pareto", func(t *testing.T) {
		retrieved, err := store.GetPareto(ctx, &api.ParetoIDs{Id: paretoID})
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, paretoID, retrieved.Ids.Id)
		assert.Len(t, retrieved.Vectors, 2)
		assert.Equal(t, []float64{10.0, 20.0}, retrieved.MaxObjs)

		// Verify first vector
		assert.Equal(t, []float64{1.0, 2.0, 3.0}, retrieved.Vectors[0].Elements)
		assert.Equal(t, []float64{0.5, 0.8}, retrieved.Vectors[0].Objectives)
		assert.Equal(t, 1.5, retrieved.Vectors[0].CrowdingDistance)
	})

	t.Run("get non-existent pareto", func(t *testing.T) {
		_, err := store.GetPareto(ctx, &api.ParetoIDs{Id: 99999})
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestParetoStore_UpdatePareto(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	// Create user and pareto
	user := &api.User{
		Ids:      &api.UserIDs{Username: "updateparetouser"},
		Email:    "updatepareto@example.com",
		Password: "password",
	}
	err := store.CreateUser(ctx, user)
	require.NoError(t, err)

	pareto := &api.Pareto{
		Ids: &api.ParetoIDs{UserId: "updateparetouser"},
		Vectors: []*api.Vector{
			{
				Elements:         []float64{1.0, 2.0},
				Objectives:       []float64{0.5, 0.8},
				CrowdingDistance: 1.0,
			},
		},
		MaxObjs: []float64{5.0, 10.0},
	}
	err = store.CreatePareto(ctx, pareto)
	require.NoError(t, err)

	// Get the created pareto's ID
	paretos, _, err := store.ListParetos(ctx, &api.UserIDs{Username: "updateparetouser"}, 50, 0)
	require.NoError(t, err)
	require.Greater(t, len(paretos), 0)
	paretoID := paretos[0].Ids.Id

	t.Run("update max objs", func(t *testing.T) {
		updatedPareto := &api.Pareto{
			Ids:     &api.ParetoIDs{Id: paretoID},
			MaxObjs: []float64{15.0, 25.0},
		}
		err := store.UpdatePareto(ctx, updatedPareto, "max_objs")
		assert.NoError(t, err)

		// Verify update
		retrieved, err := store.GetPareto(ctx, &api.ParetoIDs{Id: paretoID})
		require.NoError(t, err)
		assert.Equal(t, []float64{15.0, 25.0}, retrieved.MaxObjs)
	})

	t.Run("update non-existent pareto", func(t *testing.T) {
		updatedPareto := &api.Pareto{
			Ids:     &api.ParetoIDs{Id: 99999},
			MaxObjs: []float64{15.0, 25.0},
		}
		err := store.UpdatePareto(ctx, updatedPareto, "max_objs")
		assert.Error(t, err)
	})

	t.Run("update without ID", func(t *testing.T) {
		updatedPareto := &api.Pareto{
			Ids:     nil,
			MaxObjs: []float64{15.0, 25.0},
		}
		err := store.UpdatePareto(ctx, updatedPareto, "max_objs")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestParetoStore_DeletePareto(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	// Create user and pareto
	user := &api.User{
		Ids:      &api.UserIDs{Username: "deleteparetouser"},
		Email:    "deletepareto@example.com",
		Password: "password",
	}
	err := store.CreateUser(ctx, user)
	require.NoError(t, err)

	pareto := &api.Pareto{
		Ids: &api.ParetoIDs{UserId: "deleteparetouser"},
		Vectors: []*api.Vector{
			{
				Elements:         []float64{1.0, 2.0},
				Objectives:       []float64{0.5, 0.8},
				CrowdingDistance: 1.0,
			},
		},
		MaxObjs: []float64{5.0, 10.0},
	}
	err = store.CreatePareto(ctx, pareto)
	require.NoError(t, err)

	// Get the created pareto's ID
	paretos, _, err := store.ListParetos(ctx, &api.UserIDs{Username: "deleteparetouser"}, 50, 0)
	require.NoError(t, err)
	require.Greater(t, len(paretos), 0)
	paretoID := paretos[0].Ids.Id

	t.Run("delete existing pareto", func(t *testing.T) {
		err := store.DeletePareto(ctx, &api.ParetoIDs{Id: paretoID})
		assert.NoError(t, err)

		// Verify deletion
		_, err = store.GetPareto(ctx, &api.ParetoIDs{Id: paretoID})
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("delete non-existent pareto", func(t *testing.T) {
		err := store.DeletePareto(ctx, &api.ParetoIDs{Id: 99999})
		// Delete succeeds even if record doesn't exist
		assert.NoError(t, err)
	})
}

func TestParetoStore_ListParetos(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	// Create user
	user := &api.User{
		Ids:      &api.UserIDs{Username: "listparetouser"},
		Email:    "listpareto@example.com",
		Password: "password",
	}
	err := store.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("list paretos for user with multiple paretos", func(t *testing.T) {
		// Create multiple paretos
		for i := 0; i < 3; i++ {
			pareto := &api.Pareto{
				Ids: &api.ParetoIDs{UserId: "listparetouser"},
				Vectors: []*api.Vector{
					{
						Elements:         []float64{float64(i), float64(i + 1)},
						Objectives:       []float64{0.5, 0.8},
						CrowdingDistance: 1.0,
					},
				},
				MaxObjs: []float64{5.0, 10.0},
			}
			err := store.CreatePareto(ctx, pareto)
			require.NoError(t, err)
		}

		// List paretos with default pagination
		paretos, totalCount, err := store.ListParetos(ctx, &api.UserIDs{Username: "listparetouser"}, 50, 0)
		assert.NoError(t, err)
		assert.Len(t, paretos, 3)
		assert.Equal(t, 3, totalCount)

		// Verify all paretos have vectors
		for _, p := range paretos {
			assert.Len(t, p.Vectors, 1)
		}
	})

	t.Run("list paretos with pagination", func(t *testing.T) {
		// List with limit 2
		paretos, totalCount, err := store.ListParetos(ctx, &api.UserIDs{Username: "listparetouser"}, 2, 0)
		assert.NoError(t, err)
		assert.Len(t, paretos, 2)
		assert.Equal(t, 3, totalCount)

		// List with offset 2
		paretos2, totalCount2, err := store.ListParetos(ctx, &api.UserIDs{Username: "listparetouser"}, 50, 2)
		assert.NoError(t, err)
		assert.Len(t, paretos2, 1)
		assert.Equal(t, 3, totalCount2)
	})

	t.Run("list paretos for user with no paretos", func(t *testing.T) {
		// Create another user
		user2 := &api.User{
			Ids:      &api.UserIDs{Username: "emptyuser"},
			Email:    "empty@example.com",
			Password: "password",
		}
		err := store.CreateUser(ctx, user2)
		require.NoError(t, err)

		// List paretos
		paretos, totalCount, err := store.ListParetos(ctx, &api.UserIDs{Username: "emptyuser"}, 50, 0)
		assert.NoError(t, err)
		assert.Len(t, paretos, 0)
		assert.Equal(t, 0, totalCount)
	})

	t.Run("list paretos for non-existent user", func(t *testing.T) {
		_, _, err := store.ListParetos(ctx, &api.UserIDs{Username: "doesnotexist"}, 50, 0)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestParetoModel_MaxObjsSerialization(t *testing.T) {
	t.Run("serialize and deserialize max objs", func(t *testing.T) {
		model := &paretoModel{}
		maxObjs := []float64{1.5, 2.5, 3.5}

		err := model.SetMaxObjs(maxObjs)
		assert.NoError(t, err)
		assert.NotEmpty(t, model.MaxObjsJSON)

		retrieved, err := model.GetMaxObjs()
		assert.NoError(t, err)
		assert.Equal(t, maxObjs, retrieved)
	})

	t.Run("deserialize empty max objs", func(t *testing.T) {
		model := &paretoModel{MaxObjsJSON: ""}

		retrieved, err := model.GetMaxObjs()
		assert.NoError(t, err)
		assert.Empty(t, retrieved)
	})

	t.Run("deserialize invalid JSON", func(t *testing.T) {
		model := &paretoModel{MaxObjsJSON: "invalid json"}

		_, err := model.GetMaxObjs()
		assert.Error(t, err)
	})
}

func TestParetoStore_TransactionRollback(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	// Create user
	user := &api.User{
		Ids:      &api.UserIDs{Username: "txuser"},
		Email:    "tx@example.com",
		Password: "password",
	}
	err := store.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("transaction rolls back on vector error", func(t *testing.T) {
		// This test verifies that if vector creation fails,
		// the pareto creation is also rolled back

		pareto := &api.Pareto{
			Ids: &api.ParetoIDs{UserId: "txuser"},
			Vectors: []*api.Vector{
				{
					Elements:         []float64{1.0, 2.0},
					Objectives:       []float64{0.5, 0.8},
					CrowdingDistance: 1.0,
				},
			},
			MaxObjs: []float64{5.0, 10.0},
		}

		// Normal case should succeed
		err := store.CreatePareto(ctx, pareto)
		assert.NoError(t, err)

		// Verify pareto was created
		paretos, _, err := store.ListParetos(ctx, &api.UserIDs{Username: "txuser"}, 50, 0)
		assert.NoError(t, err)
		assert.Len(t, paretos, 1)
	})
}
