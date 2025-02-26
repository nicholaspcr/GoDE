package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	Email    string `gorm:"index:user_email_index,not null,size:255"`
	Password string `gorm:"not null,size:255"`
}

type userStore struct{ *gorm.DB }

func newUserStore(db *gorm.DB) *userStore { return &userStore{db} }

func (st *userStore) CreateUser(
	ctx context.Context, usr *api.User,
) error {
	user := userModel{
		Email:    usr.GetIds().Email,
		Password: usr.Password,
	}
	tx := st.DB.WithContext(ctx).Create(&user)
	return tx.Error
}

func (st *userStore) GetUser(
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

func (st *userStore) UpdateUser(
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

func (st *userStore) DeleteUser(
	ctx context.Context, usrIDs *api.UserIDs,
) error {
	model := userModel{Email: usrIDs.Email}
	tx := st.DB.Delete(&model)
	return tx.Error
}
