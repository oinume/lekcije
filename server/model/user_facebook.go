package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserFacebook struct {
	FacebookID string `gorm:"primary_key"`
	UserID     uint32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (*UserFacebook) TableName() string {
	return "user_facebook"
}

type UserFacebookService struct {
	db *gorm.DB
}

func NewUserFacebookService(db *gorm.DB) *UserFacebookService {
	return &UserFacebookService{db: db}
}
