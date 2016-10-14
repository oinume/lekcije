package model

import (
	"time"

	"github.com/jinzhu/gorm"
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
