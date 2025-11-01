package gorm

import (
	"context"
	"encoding/json"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type paretoModel struct {
	User userModel `gorm:"foreignKey:UserID"`
	gorm.Model
	MaxObjsJSON string `gorm:"type:text"`
	Vectors     []vectorModel
	UserID      uint
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

	// Set user ID if provided
	if pareto.Ids != nil && pareto.Ids.UserId != "" {
		// Look up user by username/ID
		var user userModel
		tx := st.DB.WithContext(ctx).Where("username = ?", pareto.Ids.UserId).First(&user)
		if tx.Error == nil {
			paretoModel.UserID = user.ID
		}
	}

	// Start transaction
	return st.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create pareto
		if err := tx.Create(&paretoModel).Error; err != nil {
			return err
		}

		// Create vectors
		for _, vec := range pareto.Vectors {
			vectorModel := vectorModel{
				ParetoID:         paretoModel.ID,
				CrowdingDistance: vec.CrowdingDistance,
			}
			if err := vectorModel.SetElements(vec.Elements); err != nil {
				return err
			}
			if err := vectorModel.SetObjectives(vec.Objectives); err != nil {
				return err
			}
			if err := tx.Create(&vectorModel).Error; err != nil {
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
	updates := make(map[string]interface{})

	for _, field := range fields {
		switch field {
		case "max_objs":
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
		if err := tx.Where("pareto_id = ?", paretoIDs.Id).Delete(&vectorModel{}).Error; err != nil {
			return err
		}

		// Delete pareto
		if err := tx.Delete(&paretoModel{}, paretoIDs.Id).Error; err != nil {
			return err
		}

		return nil
	})
}

// ListParetos returns all paretos for a given user
func (st *paretoStore) ListParetos(
	ctx context.Context, userIDs *api.UserIDs,
) ([]*api.Pareto, error) {
	// Look up user by username
	var user userModel
	tx := st.DB.WithContext(ctx).Where("username = ?", userIDs.Username).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var paretos []paretoModel
	tx = st.DB.WithContext(ctx).Where("user_id = ?", user.ID).Preload("Vectors").Find(&paretos)
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := make([]*api.Pareto, len(paretos))
	for i, p := range paretos {
		maxObjs, err := p.GetMaxObjs()
		if err != nil {
			return nil, err
		}

		vectors := make([]*api.Vector, len(p.Vectors))
		for j, vec := range p.Vectors {
			elements, err := vec.GetElements()
			if err != nil {
				return nil, err
			}
			objectives, err := vec.GetObjectives()
			if err != nil {
				return nil, err
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

	return result, nil
}
