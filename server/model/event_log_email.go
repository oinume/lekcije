package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
  datetime DATETIME NOT NULL,
  event ENUM('click', 'delivered', 'open', 'deferred', 'drop', 'bounce', 'block') NOT NULL,
  email_type ENUM('new_lesson') NOT NULL,
  user_id int(10) unsigned NOT NULL,
  user_agent VARCHAR(255) DEFAULT NULL,
  teacher_ids TEXT DEFAULT NULL,
  url VARCHAR(255) DEFAULT NULL,
  KEY (`datetime`, `event`),
  KEY (`user_id`)
*/

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
