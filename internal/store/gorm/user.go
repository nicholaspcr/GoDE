package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/internal/tenant"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type UserModel struct {
	BaseModel

	ID     string      `gorm:"primary_key"`
	Tenant TenantModel `gorm:"foreignKey:ID"`

	Email    string `gorm:"index:user_email_index,unique,not null,size:255"`
	Password string `gorm:"not null,size:255"`
}

type userStore struct{ *gorm.DB }

func newUserStore(db *gorm.DB) *userStore { return &userStore{db} }

func (st *userStore) Create(ctx context.Context, usr *api.User) error {
	tnt, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	user := UserModel{
		ID:       usr.GetIds().UserId,
		Tenant:   TenantModel{ID: tnt.GetIds().TenantId},
		Email:    usr.Email,
		Password: usr.Password,
	}
	st.DB.Create(&user)
	return nil
}

func (st *userStore) Get(
	ctx context.Context, usrIDs *api.UserIDs,
) (*api.User, error) {
	var usr UserModel
	tx := st.First(&usr, "id = ?", usrIDs.UserId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.User{
		Ids:      &api.UserIDs{UserId: usr.ID},
		Email:    usr.Email,
		Password: usr.Password,
	}, nil
}

func (st *userStore) Update(ctx context.Context, usr *api.User, fields ...string) error {
	tnt, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	model := UserModel{
		ID:     usr.GetIds().UserId,
		Tenant: TenantModel{ID: tnt.GetIds().TenantId},
	}

	for _, field := range fields {
		switch field {
		default:
			return errors.ErrUnsupportedFieldMask
		case "email":
			model.Email = usr.Email
		case "password":
			model.Password = usr.Password
		}
	}

	tx := st.DB.Model(&model).Updates(model)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (st *userStore) Delete(ctx context.Context, usrIDs *api.UserIDs) error {
	tnt, err := tenant.FromContext(ctx)
	if err != nil {
		return err
	}
	model := UserModel{
		ID:     usrIDs.UserId,
		Tenant: TenantModel{ID: tnt.GetIds().TenantId},
	}
	tx := st.DB.Delete(&model)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
