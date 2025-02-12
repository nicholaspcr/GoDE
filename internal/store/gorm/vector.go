package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type vectorModel struct {
	gorm.Model
	ParetoID uint
	Pareto   paretoModel `gorm:"foreignKey:ParetoID"`

	Elements         []float64
	Objectives       []float64
	CrowdingDistance float64
}

type vectorStore struct{ *gorm.DB }

func newVectorStore(db *gorm.DB) *vectorStore { return &vectorStore{db} }

func (st *vectorStore) CreateVector(
	ctx context.Context, vec *api.Vector, paretoID uint,
) error {
	vector := vectorModel{
		ParetoID:         paretoID,
		Elements:         vec.Elements,
		Objectives:       vec.Objectives,
		CrowdingDistance: vec.CrowdingDistance,
	}
	tx := st.DB.WithContext(ctx).Create(&vector)
	return tx.Error
}

func (st *vectorStore) GetVector(
	ctx context.Context, vecIDs *api.VectorIDs,
) (*api.Vector, error) {
	return nil, nil
}

func (st *vectorStore) UpdateVector(
	ctx context.Context, vec *api.Vector, fields ...string,
) error {
	return nil
}

func (st *vectorStore) DeleteVector(
	ctx context.Context, vecIDs *api.VectorIDs,
) error {
	return nil
}
