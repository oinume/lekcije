package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type NotificationTimeSpan struct {
	UserID    uint32
	Number    uint8
	FromTime  time.Time
	ToTime    time.Time
	CreatedAt time.Time
}

func (*NotificationTimeSpan) TableName() string {
	return "notification_time_span"
}

type NotificationTimeSpanService struct {
	db *gorm.DB
}

func NewNotificationTimeSpanService(db *gorm.DB) *EventLogEmailService {
	return &EventLogEmailService{db}
}
