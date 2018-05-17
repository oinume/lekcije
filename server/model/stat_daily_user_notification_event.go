package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

type StatDailyUserNotificationEvent struct {
	Date   time.Time
	UserID uint32
	Event  string
	Count  uint32
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
SELECT
  IFNULL(ele.date, ?) AS date
  , u.id AS user_id
  , IFNULL(ele.event, 'open') AS event
  , IFNULL(ele.count, 0) AS count
FROM user AS u /* Select from user to insert zero count users */
LEFT JOIN (
  SELECT DATE(datetime) AS date, user_id, event, COUNT(*) AS count
  FROM event_log_email
  WHERE
    datetime BETWEEN ? AND ?
    AND event='open'
    GROUP BY date, user_id, event
) AS ele ON u.id = ele.user_id
ORDER BY user_id ASC
ON DUPLICATE KEY UPDATE count = IFNULL(ele.count, 0) 
`, tableName)
	values := []interface{}{
		date.Format("2006-01-02"),
		date.Format("2006-01-02") + " 00:00:00",
		date.Format("2006-01-02") + " 23:59:59",
	}
	if err := s.db.Exec(strings.TrimSpace(sql), values...).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource(tableName, "date", date)),
		)
	}
	return nil
}

func (s *StatDailyUserNotificationEventService) FindAllByDate(date time.Time) ([]*StatDailyUserNotificationEvent, error) {
	events := make([]*StatDailyUserNotificationEvent, 0, 1000)
	sql := fmt.Sprintf(`SELECT * FROM %s WHERE date = ?`, (&StatDailyUserNotificationEvent{}).TableName())
	if err := s.db.Raw(sql, date.Format("2006-01-02")).Scan(&events).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to FindAllByDate"),
		)
	}
	return events, nil
}
