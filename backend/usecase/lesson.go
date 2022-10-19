package usecase

import (
	"context"
	"time"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type Lesson struct {
	lessonRepo repository.Lesson
}

func NewLesson(lessonRepo repository.Lesson) *Lesson {
	return &Lesson{
		lessonRepo: lessonRepo,
	}
}

func (u *Lesson) FindLessons(
	ctx context.Context,
	teacherID uint, fromDate, toDate time.Time,
) ([]*model2.Lesson, error) {
	return u.lessonRepo.FindAllByTeacherIDsDatetimeBetween(ctx, teacherID, fromDate, toDate)
}

func (u *Lesson) GetNewAvailableLessons(ctx context.Context, oldLessons, newLessons []*model2.Lesson) []*model2.Lesson {
	return u.lessonRepo.GetNewAvailableLessons(ctx, oldLessons, newLessons)
}
