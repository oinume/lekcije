package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type StatDailyUserNotificationEvent struct {
	Date  time.Time
	Event string
	Count uint32

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

func (s *StatDailyUserNotificationEventService) CreateOrUpdate(date time.Time) error {
	tableName := (&StatDailyUserNotificationEvent{}).TableName()
	sql := fmt.Sprintf(`
INSERT INTO %s (date, user_id, event, count)
SELECT DATE(ele.datetime) AS date, ele.user_id, ele.event, COUNT(*) AS count
FROM event_log_email AS ele
WHERE
  ele.datetime BETWEEN ? AND ?
  AND ele.event='open'
GROUP BY ele.user_id
ON DUPLICATE KEY UPDATE count = ele.count 
`, tableName)
	values := []interface{}{
		date.Format("2006-01-02 00:00:00"),
		date.Format("2006-01-02 23:59:59"),
	}
	if err := s.db.Exec(strings.TrimSpace(sql), values...).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource(tableName, "date", date)),
		)
	}
	return nil
}
