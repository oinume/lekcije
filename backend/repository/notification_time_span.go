package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type NotificationTimeSpan interface {
	FindByUserID(ctx context.Context, userID uint) ([]*model2.NotificationTimeSpan, error)
}
