package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	Username string `gorm:"index:username_index,not null,size:64"`
	Email    string `gorm:"index:user_email_index,not null,size:256"`
	Password string `gorm:"not null,size:256"`
}

type userStore struct{ *gorm.DB }

func newUserStore(db *gorm.DB) *userStore { return &userStore{db} }

func (st *userStore) CreateUser(
	ctx context.Context, usr *api.User,
) error {
	user := userModel{
		Username: usr.GetIds().Username,
		Email:    usr.Email,
		Password: usr.Password,
	}
	tx := st.DB.WithContext(ctx).Create(&user)
	return tx.Error
}

func (st *userStore) GetUser(
	ctx context.Context, usrIDs *api.UserIDs,
) (*api.User, error) {
	var usr userModel

	tx := st.First(&usr, "username = ?", usrIDs.Username)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.User{
		Ids:      &api.UserIDs{Username: usr.Username},
		Password: usr.Password,
	}, nil
}

func (st *userStore) UpdateUser(
	ctx context.Context, usr *api.User, fields ...string,
) error {
	var model userModel

	tx := st.First(&model, "username = ?", usr.GetIds().Username)
	if tx.Error != nil {
		return tx.Error
	}

	columns := make(map[string]any)
	for _, field := range fields {
		switch field {
		default:
			return errors.ErrUnsupportedFieldMask
		case "email":
			columns[field] = usr.Email
		case "password":
			columns[field] = usr.Password
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
	model := userModel{Username: usrIDs.Username}
	tx := st.DB.Delete(&model)
	return tx.Error
}
