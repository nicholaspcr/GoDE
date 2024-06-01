package gorm

import (
	"context"
	"time"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

// tenantModel wraps the tenant_id column into a separate type.
type tenantModel struct {
	ID        string `gorm:"primary_key,size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type tenantStore struct{ *gorm.DB }

func newTenantStore(db *gorm.DB) *tenantStore { return &tenantStore{db} }

func (st *tenantStore) CreateTenant(
	ctx context.Context, tnt *api.Tenant,
) error {
	t := tenantModel{ID: tnt.GetIds().TenantId}
	st.DB.Create(&t)
	return nil
}

func (st *tenantStore) GetTenant(
	ctx context.Context, tntIDs *api.TenantIDs,
) (*api.Tenant, error) {
	var tnt tenantModel
	tx := st.First(&tnt, "id = ?", tntIDs.TenantId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.Tenant{Ids: &api.TenantIDs{TenantId: tnt.ID}}, nil
}

func (st *tenantStore) DeleteTenant(
	ctx context.Context, tntIDs *api.TenantIDs,
) error {
	var tnt tenantModel
	tx := st.First(&tnt, "id = ?", tntIDs.TenantId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
