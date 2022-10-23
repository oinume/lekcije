package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type LessonStatusLog interface {
	Create(ctx context.Context, log *model2.LessonStatusLog) error
}
