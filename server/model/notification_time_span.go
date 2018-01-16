package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
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

func NewNotificationTimeSpanService(db *gorm.DB) *NotificationTimeSpanService {
	return &NotificationTimeSpanService{db}
}

func (s *NotificationTimeSpanService) NewNotificationTimeSpansFromPB(
	userID uint32, args []*api_v1.NotificationTimeSpan,
) []*NotificationTimeSpan {
	values := make([]*NotificationTimeSpan, 0, len(args))
	now := time.Now().UTC()
	for i, v := range args {
		fromTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			int(v.FromHour), int(v.FromMinute), 0, 0, nil,
		)
		toTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			int(v.ToHour), int(v.ToMinute), 0, 0, nil,
		)
		values = append(values, &NotificationTimeSpan{
			UserID:   userID,
			Number:   uint8(i + 1),
			FromTime: fromTime,
			ToTime:   toTime,
		})
	}
	return values
}

func (s *NotificationTimeSpanService) UpdateAll(timeSpans []*NotificationTimeSpan) error {
	return nil
}
