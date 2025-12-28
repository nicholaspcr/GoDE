package gorm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type paretoModel struct {
	User userModel `gorm:"foreignKey:UserID"`
	gorm.Model
	MaxObjsJSON string        `gorm:"type:text"`
	Vectors     []vectorModel `gorm:"foreignKey:ParetoSetID"`
	UserID      uint
	Algorithm   string `gorm:"type:varchar(100);not null"`
	Problem     string `gorm:"type:varchar(100);not null"`
	Variant     string `gorm:"type:varchar(100);not null"`
}

func (paretoModel) TableName() string {
	return "pareto_sets"
}

// SetMaxObjs serializes float64 slice to JSON
func (p *paretoModel) SetMaxObjs(maxObjs []float64) error {
	data, err := json.Marshal(maxObjs)
	if err != nil {
		return err
	}
	p.MaxObjsJSON = string(data)
	return nil
}

// GetMaxObjs deserializes JSON to float64 slice
func (p *paretoModel) GetMaxObjs() ([]float64, error) {
	var maxObjs []float64
	if p.MaxObjsJSON == "" {
		return maxObjs, nil
	}
	err := json.Unmarshal([]byte(p.MaxObjsJSON), &maxObjs)
	return maxObjs, err
}

type paretoStore struct{ *gorm.DB }

func newParetoStore(db *gorm.DB) *paretoStore { return &paretoStore{db} }

func (st *paretoStore) CreatePareto(
	ctx context.Context, pareto *api.Pareto,
) error {
	// Create pareto model
	paretoModel := paretoModel{}

	// Set max objectives
	if err := paretoModel.SetMaxObjs(pareto.MaxObjs); err != nil {
		return err
	}

	// Set user ID if provided - return error if user not found
	if pareto.Ids != nil && pareto.Ids.UserId != "" {
		// Look up user by username/ID
		var user userModel
		tx := st.DB.WithContext(ctx).Where("username = ?", pareto.Ids.UserId).First(&user)
		if tx.Error != nil {
			if tx.Error == gorm.ErrRecordNotFound {
				return fmt.Errorf("user not found: %s", pareto.Ids.UserId)
			}
			return fmt.Errorf("failed to lookup user: %w", tx.Error)
		}
		paretoModel.UserID = user.ID
	}

	// Start transaction
	return st.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create pareto
		if err := tx.Create(&paretoModel).Error; err != nil {
			return err
		}

		// Prepare all vectors for batch insert
		vectorModels := make([]vectorModel, 0, len(pareto.Vectors))
		for _, vec := range pareto.Vectors {
			vectorModel := vectorModel{
				ParetoSetID:      paretoModel.ID,
				CrowdingDistance: vec.CrowdingDistance,
			}
			if err := vectorModel.SetElements(vec.Elements); err != nil {
				return err
			}
			if err := vectorModel.SetObjectives(vec.Objectives); err != nil {
				return err
			}
			vectorModels = append(vectorModels, vectorModel)
		}

		// Batch insert all vectors (100 per batch for optimal performance)
		if len(vectorModels) > 0 {
			if err := tx.CreateInBatches(vectorModels, 100).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (st *paretoStore) GetPareto(
	ctx context.Context, paretoIDs *api.ParetoIDs,
) (*api.Pareto, error) {
	var pareto paretoModel
	tx := st.DB.WithContext(ctx).Preload("Vectors").First(&pareto, paretoIDs.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Convert to API model
	maxObjs, err := pareto.GetMaxObjs()
	if err != nil {
		return nil, err
	}

	vectors := make([]*api.Vector, len(pareto.Vectors))
	for i, vec := range pareto.Vectors {
		elements, err := vec.GetElements()
		if err != nil {
			return nil, err
		}
		objectives, err := vec.GetObjectives()
		if err != nil {
			return nil, err
		}

		vectors[i] = &api.Vector{
			Ids:              &api.VectorIDs{Id: uint64(vec.ID)},
			Elements:         elements,
			Objectives:       objectives,
			CrowdingDistance: vec.CrowdingDistance,
		}
	}

	return &api.Pareto{
		Ids:     &api.ParetoIDs{Id: uint64(pareto.ID)},
		Vectors: vectors,
		MaxObjs: maxObjs,
	}, nil
}

func (st *paretoStore) UpdatePareto(
	ctx context.Context, pareto *api.Pareto, fields ...string,
) error {
	if pareto.Ids == nil || pareto.Ids.Id == 0 {
		return gorm.ErrRecordNotFound
	}

	var paretoModel paretoModel
	tx := st.DB.WithContext(ctx).First(&paretoModel, pareto.Ids.Id)
	if tx.Error != nil {
		return tx.Error
	}

	// Update fields
	updates := make(map[string]any)

	for _, field := range fields {
		if field == "max_objs" {
			if err := paretoModel.SetMaxObjs(pareto.MaxObjs); err != nil {
				return err
			}
			updates["max_objs_json"] = paretoModel.MaxObjsJSON
		}
	}

	if len(updates) == 0 {
		// Update all fields if none specified
		if err := paretoModel.SetMaxObjs(pareto.MaxObjs); err != nil {
			return err
		}
		tx = st.DB.WithContext(ctx).Save(&paretoModel)
	} else {
		tx = st.DB.WithContext(ctx).Model(&paretoModel).Updates(updates)
	}

	return tx.Error
}

func (st *paretoStore) DeletePareto(
	ctx context.Context, paretoIDs *api.ParetoIDs,
) error {
	return st.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete associated vectors first (cascade)
		if err := tx.Where("pareto_set_id = ?", paretoIDs.Id).Delete(&vectorModel{}).Error; err != nil {
			return err
		}

		// Delete pareto
		if err := tx.Delete(&paretoModel{}, paretoIDs.Id).Error; err != nil {
			return err
		}

		return nil
	})
}

// ListParetos returns paretos for a given user with pagination.
func (st *paretoStore) ListParetos(
	ctx context.Context, userIDs *api.UserIDs, limit, offset int,
) ([]*api.Pareto, int, error) {
	// Apply defaults and max limits
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Look up user by username
	var user userModel
	tx := st.DB.WithContext(ctx).Where("username = ?", userIDs.Username).First(&user)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}

	query := st.DB.WithContext(ctx).Where("user_id = ?", user.ID)

	// Get total count
	var totalCount int64
	if err := query.Model(&paretoModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	query = query.Order("created_at DESC").Limit(limit).Offset(offset)

	var paretos []paretoModel
	if err := query.Preload("Vectors").Find(&paretos).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*api.Pareto, len(paretos))
	for i, p := range paretos {
		maxObjs, err := p.GetMaxObjs()
		if err != nil {
			return nil, 0, err
		}

		vectors := make([]*api.Vector, len(p.Vectors))
		for j, vec := range p.Vectors {
			elements, err := vec.GetElements()
			if err != nil {
				return nil, 0, err
			}
			objectives, err := vec.GetObjectives()
			if err != nil {
				return nil, 0, err
			}

			vectors[j] = &api.Vector{
				Ids:              &api.VectorIDs{Id: uint64(vec.ID)},
				Elements:         elements,
				Objectives:       objectives,
				CrowdingDistance: vec.CrowdingDistance,
			}
		}

		result[i] = &api.Pareto{
			Ids:     &api.ParetoIDs{Id: uint64(p.ID)},
			Vectors: vectors,
			MaxObjs: maxObjs,
		}
	}

	return result, int(totalCount), nil
}

// CreateParetoSet creates a pareto set with vectors and max objectives.
func (st *paretoStore) CreateParetoSet(ctx context.Context, paretoSet *store.ParetoSet) error {
	// Create pareto model
	paretoModel := paretoModel{
		Algorithm: paretoSet.Algorithm,
		Problem:   paretoSet.Problem,
		Variant:   paretoSet.Variant,
	}

	// Convert max objectives
	flatMaxObjs := make([]float64, 0)
	for _, maxObj := range paretoSet.MaxObjectives {
		if maxObj != nil && maxObj.Values != nil {
			flatMaxObjs = append(flatMaxObjs, maxObj.Values...)
		}
	}

	if err := paretoModel.SetMaxObjs(flatMaxObjs); err != nil {
		return err
	}

	// Look up user - return error if not found to prevent orphaned pareto records
	var user userModel
	tx := st.DB.WithContext(ctx).Where("username = ?", paretoSet.UserID).First(&user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found: %s", paretoSet.UserID)
		}
		return fmt.Errorf("failed to lookup user: %w", tx.Error)
	}
	paretoModel.UserID = user.ID

	// Start transaction
	return st.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create pareto
		if err := tx.Create(&paretoModel).Error; err != nil {
			return err
		}

		// Set the ID back to the paretoSet
		paretoSet.ID = uint64(paretoModel.ID)

		// Prepare all vectors for batch insert
		vectorModels := make([]vectorModel, 0, len(paretoSet.Vectors))
		for _, vec := range paretoSet.Vectors {
			vectorModel := vectorModel{
				ParetoSetID:      paretoModel.ID,
				CrowdingDistance: vec.CrowdingDistance,
			}
			if err := vectorModel.SetElements(vec.Elements); err != nil {
				return err
			}
			if err := vectorModel.SetObjectives(vec.Objectives); err != nil {
				return err
			}
			vectorModels = append(vectorModels, vectorModel)
		}

		// Batch insert all vectors (100 per batch for optimal performance)
		if len(vectorModels) > 0 {
			if err := tx.CreateInBatches(vectorModels, 100).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetParetoSetByID retrieves a pareto set by its ID.
func (st *paretoStore) GetParetoSetByID(ctx context.Context, id uint64) (*store.ParetoSet, error) {
	var paretoModel paretoModel
	tx := st.DB.WithContext(ctx).Preload("Vectors").Preload("User").First(&paretoModel, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, store.ErrParetoSetNotFound
		}
		return nil, tx.Error
	}

	// Convert to store.ParetoSet
	maxObjs, err := paretoModel.GetMaxObjs()
	if err != nil {
		return nil, err
	}

	vectors := make([]*api.Vector, len(paretoModel.Vectors))
	for i, vec := range paretoModel.Vectors {
		elements, err := vec.GetElements()
		if err != nil {
			return nil, err
		}
		objectives, err := vec.GetObjectives()
		if err != nil {
			return nil, err
		}

		vectors[i] = &api.Vector{
			Elements:         elements,
			Objectives:       objectives,
			CrowdingDistance: vec.CrowdingDistance,
		}
	}

	// Convert flat max objs to store.MaxObjectives
	maxObjectives := []*store.MaxObjectives{
		{Values: maxObjs},
	}

	return &store.ParetoSet{
		ID:            uint64(paretoModel.ID),
		UserID:        paretoModel.User.Username,
		Algorithm:     paretoModel.Algorithm,
		Problem:       paretoModel.Problem,
		Variant:       paretoModel.Variant,
		Vectors:       vectors,
		MaxObjectives: maxObjectives,
		CreatedAt:     paretoModel.CreatedAt,
	}, nil
}
