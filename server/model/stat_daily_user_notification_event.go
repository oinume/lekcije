package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type StatDailyUserNotificationEvent struct {
	Date    time.Time
	Event   string
	Count   uint32
	UUCount uint32
}

func (*StatDailyUserNotificationEvent) TableName() string {
	return "stat_daily_user_notification_event"
}

type StatDailyUserNotificationEventService struct {
	db *gorm.DB
}

func NewStatDailyUserNotificationEventService(db *gorm.DB) *StatDailyUserNotificationEventService {
	return &StatDailyUserNotificationEventService{db}
}

func (s *StatDailyUserNotificationEventService) CreateOrUpdate(v *StatDailyUserNotificationEvent) error {
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
