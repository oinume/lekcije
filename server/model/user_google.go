package model

import (
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

func (s *UserGoogleService) FindOrCreate(googleID string, userID uint32) (*UserGoogle, error) {
	userGoogle := UserGoogle{
		GoogleID: googleID,
		UserID:   userID,
	}
	if err := s.db.FirstOrCreate(&userGoogle).Error; err != nil {
		return nil, errors.InternalWrapf(err, "Failed to find or create UserGoogle")
	}
	return &userGoogle, nil
}
