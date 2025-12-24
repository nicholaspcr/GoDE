package gorm

import (
	"context"
	"encoding/json"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type vectorModel struct {
	gorm.Model
	ElementsJSON     string      `gorm:"type:text"`
	ObjectivesJSON   string      `gorm:"type:text"`
	Pareto           paretoModel `gorm:"foreignKey:ParetoID"`
	ParetoID         uint
	CrowdingDistance float64
}

// SetElements serializes float64 slice to JSON
func (v *vectorModel) SetElements(elements []float64) error {
	data, err := json.Marshal(elements)
	if err != nil {
		return err
	}
	v.ElementsJSON = string(data)
	return nil
}

// GetElements deserializes JSON to float64 slice
func (v *vectorModel) GetElements() ([]float64, error) {
	var elements []float64
	if v.ElementsJSON == "" {
		return elements, nil
	}
	err := json.Unmarshal([]byte(v.ElementsJSON), &elements)
	return elements, err
}

// SetObjectives serializes float64 slice to JSON
func (v *vectorModel) SetObjectives(objectives []float64) error {
	data, err := json.Marshal(objectives)
	if err != nil {
		return err
	}
	v.ObjectivesJSON = string(data)
	return nil
}

// GetObjectives deserializes JSON to float64 slice
func (v *vectorModel) GetObjectives() ([]float64, error) {
	var objectives []float64
	if v.ObjectivesJSON == "" {
		return objectives, nil
	}
	err := json.Unmarshal([]byte(v.ObjectivesJSON), &objectives)
	return objectives, err
}

type vectorStore struct{ *gorm.DB }

func newVectorStore(db *gorm.DB) *vectorStore { return &vectorStore{db} }

func (st *vectorStore) CreateVector(
	ctx context.Context, vec *api.Vector, paretoID uint,
) error {
	vector := vectorModel{
		ParetoID:         paretoID,
		CrowdingDistance: vec.CrowdingDistance,
	}

	if err := vector.SetElements(vec.Elements); err != nil {
		return err
	}
	if err := vector.SetObjectives(vec.Objectives); err != nil {
		return err
	}

	tx := st.DB.WithContext(ctx).Create(&vector)
	return tx.Error
}

func (st *vectorStore) GetVector(
	ctx context.Context, vecIDs *api.VectorIDs,
) (*api.Vector, error) {
	var vector vectorModel
	tx := st.DB.WithContext(ctx).First(&vector, vecIDs.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	elements, err := vector.GetElements()
	if err != nil {
		return nil, err
	}

	objectives, err := vector.GetObjectives()
	if err != nil {
		return nil, err
	}

	return &api.Vector{
		Ids:              &api.VectorIDs{Id: uint64(vector.ID)},
		Elements:         elements,
		Objectives:       objectives,
		CrowdingDistance: vector.CrowdingDistance,
	}, nil
}

func (st *vectorStore) UpdateVector(
	ctx context.Context, vec *api.Vector, fields ...string,
) error {
	if vec.Ids == nil || vec.Ids.Id == 0 {
		return gorm.ErrRecordNotFound
	}

	var vector vectorModel
	tx := st.DB.WithContext(ctx).First(&vector, vec.Ids.Id)
	if tx.Error != nil {
		return tx.Error
	}

	// Update fields
	updates := make(map[string]any)

	for _, field := range fields {
		switch field {
		case "elements":
			if err := vector.SetElements(vec.Elements); err != nil {
				return err
			}
			updates["elements_json"] = vector.ElementsJSON
		case "objectives":
			if err := vector.SetObjectives(vec.Objectives); err != nil {
				return err
			}
			updates["objectives_json"] = vector.ObjectivesJSON
		case "crowding_distance":
			updates["crowding_distance"] = vec.CrowdingDistance
		}
	}

	if len(updates) == 0 {
		// Update all fields if none specified
		if err := vector.SetElements(vec.Elements); err != nil {
			return err
		}
		if err := vector.SetObjectives(vec.Objectives); err != nil {
			return err
		}
		vector.CrowdingDistance = vec.CrowdingDistance
		tx = st.DB.WithContext(ctx).Save(&vector)
	} else {
		tx = st.DB.WithContext(ctx).Model(&vector).Updates(updates)
	}

	return tx.Error
}

func (st *vectorStore) DeleteVector(
	ctx context.Context, vecIDs *api.VectorIDs,
) error {
	tx := st.DB.WithContext(ctx).Delete(&vectorModel{}, vecIDs.Id)
	return tx.Error
}
