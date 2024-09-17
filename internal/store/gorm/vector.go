package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type vectorModel struct {
	gorm.Model
	ParetoD uint
	Pareto  paretoModel `gorm:"foreignKey:ParetoID"`
}

type vectorStore struct{ *gorm.DB }

func newVectorStore(db *gorm.DB) *vectorStore { return &vectorStore{db} }

func (st *vectorStore) CreateVector(
	ctx context.Context, usr *api.Vector,
) error {
	vector := vectorModel{}
	st.DB.WithContext(ctx).Create(&vector)
	return nil
}

func (st *vectorStore) GetVector(
	ctx context.Context, usrIDs *api.VectorIDs,
) (*api.Vector, error) {
	return nil, nil
}

func (st *vectorStore) UpdateVector(
	ctx context.Context, usr *api.Vector, fields ...string,
) error {
	return nil
}

func (st *vectorStore) DeleteVector(
	ctx context.Context, usrIDs *api.VectorIDs,
) error {
	return nil
}
