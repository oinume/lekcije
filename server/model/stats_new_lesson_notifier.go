package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type StatsNewLessonNotifier struct {
	Date  time.Time
	Event string
	Count uint32
}

func (*StatsNewLessonNotifier) TableName() string {
	return "stats_new_lesson_notifier"
}

type StatsNewLessonNotifierService struct {
	db *gorm.DB
}

func NewStatsNewLessonNotifierService(db *gorm.DB) *StatsNewLessonNotifierService {
	return &StatsNewLessonNotifierService{db}
}

func (s *StatsNewLessonNotifierService) Create(e *StatsNewLessonNotifierService) error {
	return s.db.Create(e).Error
}
