package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/errors"
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
		if err := wrapNotFound(result, userGoogle.TableName(), "userID", fmt.Sprint(userID)); err != nil {
			return nil, err
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithResource(errors.NewResource(userGoogle.TableName(), "userID", userID)),
		)
	}
	return userGoogle, nil
}
