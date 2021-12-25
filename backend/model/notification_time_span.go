package model

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"go.opencensus.io/trace"

	"github.com/oinume/lekcije/backend/errors"
)

type NotificationTimeSpan struct {
	UserID    uint32
	Number    uint8
	FromTime  string
	ToTime    string
	CreatedAt time.Time
	from      time.Time
	to        time.Time
}

func (*NotificationTimeSpan) TableName() string {
	return "notification_time_span"
}

func (s *NotificationTimeSpan) ParseTime() error {
	f, err := time.Parse("15:04:05", s.FromTime)
	if err != nil {
		return err
	}
	s.from = f

	t, err := time.Parse("15:04:05", s.ToTime)
	if err != nil {
		return err
	}
	s.to = t
	return nil
}

func (s *NotificationTimeSpan) Within(t time.Time) bool {
	if err := s.ParseTime(); err != nil {
		return false
	}
	if s.from.Before(s.to) {
		if (t.After(s.from) || t.Equal(s.from)) && (t.Before(s.to) || t.Equal(s.to)) {
			return true
		}
	} else {
		// Add 24 hour to s.to if from > to. from=04:00, to=03:00 -> from=04:00, to=27:00
		toTime := s.to
		toTime = toTime.Add(time.Hour * 24)
		if (t.After(s.from) || t.Equal(s.from)) && (t.Before(toTime) || t.Equal(toTime)) {
			return true
		}
	}
	return false
}

type NotificationTimeSpanList []*NotificationTimeSpan

func (l NotificationTimeSpanList) Within(t time.Time) bool {
	target := t
	for _, timeSpan := range l {
		if timeSpan.Within(target) {
			return true
		}
	}
	return false
}

type NotificationTimeSpanService struct {
	db *gorm.DB
}

func NewNotificationTimeSpanService(db *gorm.DB) *NotificationTimeSpanService {
	return &NotificationTimeSpanService{db}
}

// FindByUserID is used from SendNotification on notifier.go
func (s *NotificationTimeSpanService) FindByUserID(
	ctx context.Context,
	userID uint32,
) ([]*NotificationTimeSpan, error) {
	_, span := trace.StartSpan(ctx, "NotificationTimeSpanService.FindByUserID")
	defer span.End()
	span.Annotatef([]trace.Attribute{
		trace.Int64Attribute("userID", int64(userID)),
	}, "userID:%d", userID)

	sql := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, (&NotificationTimeSpan{}).TableName())
	timeSpans := make([]*NotificationTimeSpan, 0, 10)
	if err := s.db.Raw(sql, userID).Scan(&timeSpans).Error; err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("FindByUserID select failed"),
			errors.WithResource(errors.NewResource((&NotificationTimeSpan{}).TableName(), "userID", userID)),
		)
	}
	return timeSpans, nil
}

// UpdateAll is used from TestNotifier_SendNotification
func (s *NotificationTimeSpanService) UpdateAll(userID uint32, timeSpans []*NotificationTimeSpan) error {
	for _, timeSpan := range timeSpans {
		if userID != timeSpan.UserID {
			return errors.NewInvalidArgumentError(
				errors.WithMessage("Given userID and userID of timeSpans must be same"),
			)
		}
	}

	tx := s.db.Begin()
	tableName := (&NotificationTimeSpan{}).TableName()
	sql := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ?`, tableName)
	if err := tx.Exec(sql, userID).Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("UpdateAll delete failed"),
			errors.WithResource(errors.NewResource(tableName, "userID", userID)),
		)
	}

	for _, timeSpan := range timeSpans {
		if err := tx.Create(timeSpan).Error; err != nil {
			return errors.NewInternalError(
				errors.WithError(err),
				errors.WithMessage("UpdateAll insert failed"),
				errors.WithResource(errors.NewResource(tableName, "userID", userID)),
			)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("UpdateAll commit failed"),
			errors.WithResource(errors.NewResource(tableName, "userID", userID)),
		)
	}

	return nil
}
