package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const DefaultPlanID = uint8(4)

type Plan struct {
	ID                   uint8 `gorm:"primary_key"`
	Name                 string
	InternalName         string
	Price                int16
	NotificationInterval uint8
	ShowAd               bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (*Plan) TableName() string {
	return "plan"
}

type PlanService struct {
	db *gorm.DB
}

func NewPlanService(db *gorm.DB) *PlanService {
	return &PlanService{db: db}
}

func (s *PlanService) TableName() string {
	return (&Plan{}).TableName()
}

func (s *PlanService) FindByPk(id uint8) (*Plan, error) {
	plan := &Plan{}
	if err := s.db.First(plan, &Plan{ID: id}).Error; err != nil {
		return nil, err
	}
	return plan, nil
}
