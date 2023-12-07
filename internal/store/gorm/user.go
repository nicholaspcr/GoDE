package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api"
	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	Email    string `gorm:"index:user_email_index,unique,not null,size:255"`
	Password string `gorm:"not null,size:255"`
}

type userStore struct {
	*gorm.DB
}

func newUserStore(db *gorm.DB) *userStore {
	return &userStore{db}
}

func (st *userStore) Create(ctx context.Context, usr *api.User) error {
	user := userModel{
		Email:    usr.GetIds().Email,
		Password: usr.GetPassword(),
	}
	st.DB.Create(&user)
	return nil
}

func (st *userStore) Get(
	ctx context.Context, usrIDs *api.UserIDs,
) (*api.User, error) {
	var usr userModel
	tx := st.DB.First(&usr, "email = ?", usrIDs.Email)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.User{
		Ids:      &api.UserIDs{Email: usr.Email},
		Password: usr.Password,
	}, nil
}

func (st *userStore) Update(ctx context.Context, usr *api.User) error {
	// TODO: Update specific fields via:
	// db.Model(&user).
	//		Select("name").
	//		Updates(map[string]interface{}{
	//			"name": "hello",
	//			"age": 18,
	//			"active": false,
	//		})
	tx := st.DB.Updates(userModel{
		Email:    usr.GetIds().Email,
		Password: usr.GetPassword(),
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (st *userStore) Delete(ctx context.Context, usrIDs *api.UserIDs) error {
	var user userModel
	tx := st.DB.First(&user, "email = ?", usrIDs.Email)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
