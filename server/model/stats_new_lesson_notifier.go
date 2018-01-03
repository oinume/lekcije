package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type StatsNewLessonNotifier struct {
	Date    time.Time
	Event   string
	Count   uint32
	UUCount uint32
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

func (s *StatsNewLessonNotifierService) CreateOrUpdate(v *StatsNewLessonNotifier) error {
	date := v.Date.Format("2006-01-02")
	sql := fmt.Sprintf(`INSERT INTO %s VALUES (?, ?, ?, ?)`, v.TableName())
	sql += " ON DUPLICATE KEY UPDATE"
	sql += " count=?, uu_count=?"
	values := []interface{}{
		date, v.Event, v.Count, v.UUCount,
		v.Count, v.UUCount,
	}
	if err := s.db.Exec(sql, values...).Error; err != nil {
		return errors.InternalWrapf(err, "Failed to INSERT or UPDATE %s: date=%v", v.TableName(), date)
	}
	return nil
}
