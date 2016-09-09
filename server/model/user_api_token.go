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
	UserId    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*UserApiToken) TableName() string {
	return "user_api_token"
}

type UserApiTokenServiceType struct {
	db *gorm.DB
}

var UserApiTokenService UserApiTokenServiceType

func (s *UserApiTokenServiceType) Create(userId uint32) (*UserApiToken, error) {
	apiToken := util.RandomString(64)
	userApiToken := UserApiToken{
		UserId: userId,
		Token:  apiToken,
	}
	if err := s.db.Create(&userApiToken).Error; err != nil {
		return nil, errors.InternalWrapf(err, "Failed to create UserApiToken")
	}
	return &userApiToken, nil
}

func (s *UserApiTokenServiceType) DeleteByUserIdAndToken(userId uint32, token string) error {
	result := s.db.Where("user_id = ? AND token = ?", userId, token).Delete(&UserApiToken{})
	if result.Error != nil {
		return errors.InternalWrapf(
			result.Error,
			"Failed to DeleteByUserIdAndToken: userId=%v, token=%v",
			userId, token,
		)
	}
	return nil
}
