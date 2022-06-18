package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type NotificationTimeSpan struct {
	notificationTimeSpanRepo repository.NotificationTimeSpan
}

func NewNotificationTimeSpan(notificationTimeSpanRepo repository.NotificationTimeSpan) *NotificationTimeSpan {
	return &NotificationTimeSpan{
		notificationTimeSpanRepo: notificationTimeSpanRepo,
	}
}

func (u *NotificationTimeSpan) FindByUserID(ctx context.Context, userID uint) ([]*model2.NotificationTimeSpan, error) {
	return u.notificationTimeSpanRepo.FindByUserID(ctx, userID)
}

func (u *NotificationTimeSpan) UpdateAll(ctx context.Context, userID uint, timeSpans []*model2.NotificationTimeSpan) error {
	return u.notificationTimeSpanRepo.UpdateAll(ctx, userID, timeSpans)
}
