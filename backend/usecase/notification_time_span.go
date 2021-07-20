package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type NotificationTimeSpan struct {
	notificationTimeSpanRepo repository.NotificationTimeSpan
}

func (u *NotificationTimeSpan) FindByUserID(ctx context.Context, userID uint) ([]*model2.NotificationTimeSpan, error) {
	return u.notificationTimeSpanRepo.FindByUserID(ctx, userID)
}
