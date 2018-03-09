package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type UserGoogle struct {
	GoogleID  string `gorm:"primary_key"`
	UserID    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*UserGoogle) TableName() string {
	return "user_google"
}

type UserGoogleService struct {
	db *gorm.DB
}

func NewUserGoogleService(db *gorm.DB) *UserGoogleService {
	return &UserGoogleService{db: db}
}

func (s *UserGoogleService) FindByUserID(userID uint32) (*UserGoogle, error) {
	userGoogle := &UserGoogle{}
	if result := s.db.First(userGoogle, &UserGoogle{UserID: userID}); result.Error != nil {
		if result.RecordNotFound() {
			return nil, errors.NewStandardError(
				errors.CodeNotFound,
				errors.WithError(result.Error),
				errors.WithResource(userGoogle.TableName(), "userID", fmt.Sprint(userID)),
			)
		} else {
			return nil, errors.NewStandardError(
				errors.CodeInternal,
				errors.WithError(result.Error),
				errors.WithResource(userGoogle.TableName(), "userID", fmt.Sprint(userID)),
			)
		}
	}
	return userGoogle, nil
}
