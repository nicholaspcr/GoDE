package gorm

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"gorm.io/gorm"
)

type UserModel struct {
	BaseModel

	ID       string      `gorm:"primary_key"`
	TenantID string      `gorm:"primary_key,index:user_email_index,default:default"`
	Tenant   TenantModel `gorm:"foreignKey:TenantID"`

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
	user := UserModel{
		ID:       usr.GetIds().UserId,
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
	tx := st.DB.First(&usr, "id = ?", usrIDs.UserId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &api.User{
		Ids:      &api.UserIDs{UserId: usr.ID},
		Email:    usr.Email,
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
	tx := st.DB.Updates(UserModel{
		ID:       usr.GetIds().UserId,
		Email:    usr.Email,
		Password: usr.Password,
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (st *userStore) Delete(ctx context.Context, usrIDs *api.UserIDs) error {
	var user UserModel
	tx := st.DB.First(&user, "id = ?", usrIDs.UserId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
