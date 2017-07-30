package model

import (
	"time"

	"github.com/jinzhu/gorm"
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
