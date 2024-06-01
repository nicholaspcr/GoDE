package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/internal/tenant"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	TenantID string `gorm:"index:user_tenant_index,unique,not null,size:50"`
	Email    string `gorm:"index:user_email_index,unique,not null,size:255"`
	Password string `gorm:"not null,size:255"`
	Name     string `gorm:"size:64"`
}

func (u *userModel) fillTenantID(tx *gorm.DB) error {
	ctx := tx.Statement.Context

	tnt := tenant.FromContext(ctx)
	u.TenantID = tnt.GetIds().TenantId

	return nil
}

func (u *userModel) BeforeCreate(tx *gorm.DB) error {
	return u.fillTenantID(tx)
}

func (u *userModel) BeforeUpdate(tx *gorm.DB) error {
	return u.fillTenantID(tx)
}

func (u *userModel) BeforeDelete(tx *gorm.DB) error {
	return u.fillTenantID(tx)
}

type userStore struct{ *gorm.DB }

func newUserStore(db *gorm.DB) userStore { return userStore{db} }

func (st *userStore) Create(ctx context.Context, usr *api.User) error {
	user := userModel{
		Email:    usr.GetIds().Email,
		Password: usr.Password,
		Name:     usr.Name,
	}
	st.DB.WithContext(ctx).Create(&user)
	return nil
}

func (st *userStore) Get(
	ctx context.Context, usrIDs *api.UserIDs,
) (*api.User, error) {
	var usr userModel

	tx := st.First(&usr, "email = ?", usrIDs.Email)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.User{
		Ids:      &api.UserIDs{Email: usr.Email},
		Password: usr.Password,
		Name:     usr.Name,
	}, nil
}

func (st *userStore) Update(
	ctx context.Context, usr *api.User, fields ...string,
) error {
	var model userModel

	tx := st.First(&model, "email = ?", usr.GetIds().Email)
	if tx.Error != nil {
		return tx.Error
	}

	columns := make(map[string]any)
	for _, field := range fields {
		switch field {
		default:
			return errors.ErrUnsupportedFieldMask
		case "email":
			columns[field] = usr.GetIds().Email
		case "password":
			columns[field] = usr.Password
		case "name":
			columns[field] = usr.Name
		}
	}

	tx = st.DB.Model(&model).Updates(columns)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (st *userStore) Delete(ctx context.Context, usrIDs *api.UserIDs) error {
	model := userModel{Email: usrIDs.Email}
	tx := st.DB.Delete(&model)
	return tx.Error
}
