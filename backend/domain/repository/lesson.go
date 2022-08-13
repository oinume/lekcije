package repository

import (
	"context"

	"github.com/oinume/lekcije/backend/model2"
)

type Lesson interface{}

type LessonFetcher interface {
	Close()
	Fetch(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error)
}
