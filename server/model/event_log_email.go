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
	EmailTypeRegistration      = "registration"
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

func (s *EventLogEmailService) FindStatsNewLessonNotifierByDate(date time.Time) ([]*StatNewLessonNotifier, error) {
	sql := `
SELECT CAST(datetime AS DATE) AS date, event, COUNT(*) AS count
FROM event_log_email
WHERE
  datetime BETWEEN ? AND ?
  AND email_type = 'new_lesson_notifier'
GROUP BY event;
`
	d := date.Format("2006-01-02")
	values := make([]*StatNewLessonNotifier, 0, 100)
	if err := s.db.Raw(strings.TrimSpace(sql), d+" 00:00:00", d+" 23:59:59").Scan(&values).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("select failed"),
			errors.WithResource(errors.NewResource("event_log_email", "date", date.Format("2006-01-02"))),
		)
	}
	return values, nil
}

func (s *EventLogEmailService) FindStatsNewLessonNotifierUUCountByDate(date time.Time) ([]*StatNewLessonNotifier, error) {
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
	values := make([]*StatNewLessonNotifier, 0, 100)
	if err := s.db.Raw(strings.TrimSpace(sql), d+" 00:00:00", d+" 23:59:59").Scan(&values).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("select failed"),
			errors.WithResource(errors.NewResource("event_log_email", "date", date.Format("2006-01-02"))),
		)
	}
	return values, nil
}

func (s *EventLogEmailService) FindStatDailyNotificationEventByDate(date time.Time) ([]*StatDailyNotificationEvent, error) {
	sql := `
SELECT CAST(datetime AS DATE) AS date, event, COUNT(*) AS count
FROM event_log_email
WHERE
  datetime BETWEEN ? AND ?
  AND email_type = 'new_lesson_notifier'
GROUP BY date, event;
`
	d := date.Format("2006-01-02")
	values := make([]*StatDailyNotificationEvent, 0, 100)
	if err := s.db.Raw(strings.TrimSpace(sql), d+" 00:00:00", d+" 23:59:59").Scan(&values).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("select failed"),
			errors.WithResource(errors.NewResource("event_log_email", "date", date.Format("2006-01-02"))),
		)
	}
	return values, nil
}

func (s *EventLogEmailService) FindStatDailyNotificationEventUUCountByDate(date time.Time) ([]*StatDailyNotificationEvent, error) {
	sql := `
SELECT s.date, s.event, COUNT(*) AS uu_count
FROM (
  SELECT CAST(datetime AS DATE) AS date, event, user_id, COUNT(*) AS count
  FROM event_log_email
  WHERE
    datetime BETWEEN ? AND ?
    AND email_type = 'new_lesson_notifier'
  GROUP BY date, event, user_id
) AS s
GROUP BY s.date, s.event
`
	d := date.Format("2006-01-02")
	values := make([]*StatDailyNotificationEvent, 0, 1000)
	if err := s.db.Raw(strings.TrimSpace(sql), d+" 00:00:00", d+" 23:59:59").Scan(&values).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("select failed"),
			errors.WithResource(errors.NewResource("event_log_email", "date", date.Format("2006-01-02"))),
		)
	}
	return values, nil
}

func (s *EventLogEmailService) FindStatDailyUserNotificationEvent(date time.Time) ([]*StatDailyUserNotificationEvent, error) {
	return nil, nil
}
