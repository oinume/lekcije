package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type StatDailyNotificationEvent struct {
	Date    time.Time
	Event   string
	Count   uint32
	UUCount uint32
}

func (*StatDailyNotificationEvent) TableName() string {
	return "stat_daily_notification_event"
}

type StatDailyNotificationEventService struct {
	db *gorm.DB
}

func NewStatDailyNotificationEventService(db *gorm.DB) *StatDailyNotificationEventService {
	return &StatDailyNotificationEventService{db}
}

func (s *StatDailyNotificationEventService) CreateOrUpdate(v *StatDailyNotificationEvent) error {
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

func (s *StatDailyNotificationEventService) FindAllByDate(date time.Time) ([]*StatDailyNotificationEvent, error) {
	events := make([]*StatDailyNotificationEvent, 0, 1000)
	sql := fmt.Sprintf(`SELECT * FROM %s WHERE date = ?`, (&StatDailyNotificationEvent{}).TableName())
	if err := s.db.Raw(sql, date.Format("2006-01-02")).Scan(&events).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to FindAllByDate"),
		)
	}
	return events, nil

}
