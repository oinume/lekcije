package repository

//go:generate moq -out=lesson.moq.go . LessonFetcher

import (
	"context"
	"time"

	"github.com/oinume/lekcije/backend/model2"
)

type Lesson interface {
	FindAllByTeacherIDsDatetimeBetween(
		ctx context.Context, teacherID uint,
		fromDate, toDate time.Time,
	) ([]*model2.Lesson, error)
}

type LessonFetcher interface {
	Close()
	Fetch(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error)
}
