package repository

//go:generate moq -out=lesson.moq.go . LessonFetcher

import (
	"context"
	"time"

	"github.com/oinume/lekcije/backend/model2"
)

type Lesson interface {
	Create(ctx context.Context, lesson *model2.Lesson) error
	FindAllByTeacherIDsDatetimeBetween(
		ctx context.Context, teacherID uint,
		fromDate, toDate time.Time,
	) ([]*model2.Lesson, error)
	GetNewAvailableLessons(ctx context.Context, oldLessons, newLessons []*model2.Lesson) []*model2.Lesson
}

type LessonFetcher interface {
	Close()
	Fetch(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error)
}
