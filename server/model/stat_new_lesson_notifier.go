package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type StatNewLessonNotifier struct {
	Date    time.Time
	Event   string
	Count   uint32
	UUCount uint32
}

func (*StatNewLessonNotifier) TableName() string {
	return "stat_new_lesson_notifier"
}

type StatNewLessonNotifierService struct {
	db *gorm.DB
}

func NewStatsNewLessonNotifierService(db *gorm.DB) *StatNewLessonNotifierService {
	return &StatNewLessonNotifierService{db}
}

func (s *StatNewLessonNotifierService) CreateOrUpdate(v *StatNewLessonNotifier) error {
	date := v.Date.Format("2006-01-02")
	sql := fmt.Sprintf(`INSERT INTO %s VALUES (?, ?, ?, ?)`, v.TableName())
	sql += " ON DUPLICATE KEY UPDATE"
	sql += " count=?, uu_count=?"
	values := []interface{}{
		date, v.Event, v.Count, v.UUCount,
		v.Count, v.UUCount,
	}
	if err := s.db.Exec(sql, values...).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource(v.TableName(), "date", date)),
		)
	}
	return nil
}
