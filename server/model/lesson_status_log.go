package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type LessonStatusLog struct {
	LessonID  uint64
	Status    string
	CreatedAt time.Time
}

func (*LessonStatusLog) TableName() string {
	return "lesson_status_log"
}

type LessonStatusLogService struct {
	db *gorm.DB
}

func NewLessonStatusLogService(db *gorm.DB) *LessonStatusLogService {
	return &LessonStatusLogService{db: db}
}

func (s *LessonStatusLogService) TableName() string {
	return (&LessonStatusLog{}).TableName()
}

func (s *LessonStatusLogService) Create(log *LessonStatusLog) error {
	if err := s.db.Create(log).Error; err != nil {
		return err
	}
	return nil
}
