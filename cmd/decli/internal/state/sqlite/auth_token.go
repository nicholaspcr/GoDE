package sqlite

import (
	"time"

	"gorm.io/gorm"
)

type authTokenModel struct {
	Token     string `gorm:"primaryKey"`
	CreatedAt time.Time
	Deleted   bool
}

type authTokenStore struct{ *gorm.DB }

// GetAuthToken gets the latest authentication token from the state store.
func (st *authTokenStore) GetAuthToken() (string, error) {
	var model authTokenModel

	res := st.Order("created_at DESC").First(&model)
	if res.Error != nil {
		return "", res.Error
	}
	return model.Token, nil
}

// InvalidateAuthToken invalidates the latest authentication token.
func (st *authTokenStore) InvalidateAuthToken() error {
	res := st.Model(&authTokenModel{}).
		Where("deleted = ?", false).
		Update("deleted", true)
	return res.Error
}

// SaveAuthToken saves the authentication token as the latest token.
func (st *authTokenStore) SaveAuthToken(token string) error {
	model := authTokenModel{Token: token, CreatedAt: time.Now()}
	res := st.Create(&model)
	return res.Error
}
