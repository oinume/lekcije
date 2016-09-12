package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/util"
)

const UserApiTokenExpiration = time.Hour * 24 * 30

type UserApiToken struct {
	Token     string `gorm:"primary_key;AUTO_INCREMENT"`
	UserID    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*UserApiToken) TableName() string {
	return "user_api_token"
}

type UserApiTokenService struct {
	db *gorm.DB
}

func NewUserApiTokenService(db *gorm.DB) *UserApiTokenService {
	return &UserApiTokenService{db: db}
}

func (s *UserApiTokenService) Create(userID uint32) (*UserApiToken, error) {
	apiToken := util.RandomString(64)
	userApiToken := UserApiToken{
		UserID: userID,
		Token:  apiToken,
	}
	if err := s.db.Create(&userApiToken).Error; err != nil {
		return nil, errors.InternalWrapf(err, "Failed to create UserApiToken")
	}
	return &userApiToken, nil
}

func (s *UserApiTokenService) DeleteByUserIDAndToken(userID uint32, token string) error {
	result := s.db.Where("user_id = ? AND token = ?", userID, token).Delete(&UserApiToken{})
	if result.Error != nil {
		return errors.InternalWrapf(
			result.Error,
			"Failed to DeleteByUserIDAndToken: userID=%v, token=%v",
			userID, token,
		)
	}
	return nil
}
