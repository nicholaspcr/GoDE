package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type paretoModel struct {
	gorm.Model
	UserID uint
	User   userModel `gorm:"foreignKey:UserID"`

	Vectors []vectorModel
}

type paretoStore struct{ *gorm.DB }

func newParetoStore(db *gorm.DB) *paretoStore { return &paretoStore{db} }

func (st *paretoStore) CreatePareto(
	ctx context.Context, usr *api.Pareto,
) error {
	pareto := paretoModel{}
	st.DB.WithContext(ctx).Create(&pareto)
	return nil
}

func (st *paretoStore) GetPareto(
	ctx context.Context, usrIDs *api.ParetoIDs,
) (*api.Pareto, error) {
	return nil, nil
}

func (st *paretoStore) UpdatePareto(
	ctx context.Context, usr *api.Pareto, fields ...string,
) error {
	return nil
}

func (st *paretoStore) DeletePareto(
	ctx context.Context, usrIDs *api.ParetoIDs,
) error {
	return nil
}
