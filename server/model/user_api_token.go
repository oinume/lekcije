package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

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
