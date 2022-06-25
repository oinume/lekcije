package usecase

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type Notification struct {
}

func NewNotification() *Notification {
	return &Notification{}
}

func (n *Notification) NewLessonNotifier() *LessonNotifier {
	return &LessonNotifier{}
}

type LessonNotifier struct{}

func (ln *LessonNotifier) SendNotification(ctx context.Context, user *model2.User) error {
	return nil
}

func (ln *LessonNotifier) Close(ctx context.Context, stat *model2.StatNotifier) {
}
