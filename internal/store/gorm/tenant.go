package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

// Tenant is a model for the tenant table.
type TenantModel struct {
	ID string `gorm:"primary_key,size:50"`
}

type tenantStore struct{ *gorm.DB }

func newTenantStore(db *gorm.DB) *tenantStore { return &tenantStore{db} }

func (st *tenantStore) Create(ctx context.Context, tnt *api.Tenant) error {
	tenant := TenantModel{ID: tnt.GetIds().TenantId}
	st.DB.Create(&tenant)
	return nil
}

func (st *tenantStore) Get(
	ctx context.Context, tntIDs *api.TenantIDs,
) (*api.Tenant, error) {
	var tnt TenantModel
	tx := st.First(&tnt, "id = ?", tntIDs.TenantId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.Tenant{Ids: &api.TenantIDs{TenantId: tnt.ID}}, nil
}

func (st *tenantStore) Delete(ctx context.Context, tntIDs *api.TenantIDs) error {
	var tenant TenantModel
	tx := st.First(&tenant, "id = ?", tntIDs.TenantId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
