package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

const DefaultPlanID = uint8(4)

type MPlan struct {
	ID                   uint8 `gorm:"primary_key"`
	Name                 string
	InternalName         string
	Price                int16
	NotificationInterval uint8
	ShowAd               bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (*MPlan) TableName() string {
	return "m_plan"
}

type MPlanService struct {
	db *gorm.DB
}

func NewPlanService(db *gorm.DB) *MPlanService {
	return &MPlanService{db: db}
}

func (s *MPlanService) TableName() string {
	return (&MPlan{}).TableName()
}

func (s *MPlanService) FindByPK(id uint8) (*MPlan, error) {
	plan := &MPlan{}
	if err := s.db.First(plan, &MPlan{ID: id}).Error; err != nil {
		return nil, errors.NotFoundWrapf(err, "MPlan not found for id = %v", id)
	}
	return plan, nil
}
