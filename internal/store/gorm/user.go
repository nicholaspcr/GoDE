package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api"
	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	Id       string
	Email    string
	Password string
}

type userStore struct {
	*gorm.DB
}

func newUserStore(db *gorm.DB) *userStore {
	return &userStore{db}
}

func (st *userStore) Create(ctx context.Context, usr *api.User) error {
	user := userModel{
		Id:       usr.GetId(),
		Email:    usr.GetEmail(),
		Password: usr.GetPassword(),
	}
	st.DB.Create(&user)
	return nil
}

func (st *userStore) Read(
	ctx context.Context, userID *api.UserID,
) (*api.User, error) {
	var user userModel
	tx := st.First(&user, "id = ?", userID.GetId())
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (st *userStore) Update(ctx context.Context, user *api.User) error {
	return nil
}

func (st *userStore) Delete(ctx context.Context, userID *api.UserID) error {
	return nil
}
