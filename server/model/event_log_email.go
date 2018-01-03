package model

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/errors"
)

// TODO: Custom type
const (
	EmailTypeNewLessonNotifier = "new_lesson_notifier"
	EmailTypeFollowReminder    = "follow_reminder"
)

type EventLogEmail struct {
	Datetime   time.Time
	Event      string
	EmailType  string
	UserID     uint32
	UserAgent  string
	TeacherIDs string
	URL        string
}

func (*EventLogEmail) TableName() string {
	return "event_log_email"
}

type EventLogEmailService struct {
	db *gorm.DB
}

func NewEventLogEmailService(db *gorm.DB) *EventLogEmailService {
	return &EventLogEmailService{db}
}

func (s *EventLogEmailService) Create(e *EventLogEmail) error {
	return s.db.Create(e).Error
}

func (s *EventLogEmailService) FindStatsNewLessonNotifierByDate(date time.Time) ([]*StatsNewLessonNotifier, error) {
	sql := `
SELECT CAST(datetime AS DATE) AS date, event, COUNT(*) AS count
FROM event_log_email
WHERE
  datetime BETWEEN ? AND ?
  AND email_type = 'new_lesson_notifier'
GROUP BY event;
`
	d := date.Format("2006-01-02")
	values := make([]*StatsNewLessonNotifier, 0, 100)
	if err := s.db.Raw(strings.TrimSpace(sql), d+" 00:00:00", d+" 23:59:59").Scan(&values).Error; err != nil {
		return nil, errors.InternalWrapf(err, "error")
	}
	return values, nil
}

func (s *EventLogEmailService) FindStatsNewLessonNotifierUUCountByDate(date time.Time) ([]*StatsNewLessonNotifier, error) {
	sql := `
SELECT s.date, s.event, COUNT(*) AS uu_count
FROM (
  SELECT CAST(datetime AS DATE) AS date, event, user_id, COUNT(*) AS count
  FROM event_log_email
  WHERE
    datetime BETWEEN ? AND ?
    AND email_type = 'new_lesson_notifier'
  GROUP BY event, user_id
) AS s
GROUP BY s.event
`
	d := date.Format("2006-01-02")
	values := make([]*StatsNewLessonNotifier, 0, 100)
	if err := s.db.Raw(strings.TrimSpace(sql), d+" 00:00:00", d+" 23:59:59").Scan(&values).Error; err != nil {
		return nil, errors.InternalWrapf(err, "error")
	}
	return values, nil
}
