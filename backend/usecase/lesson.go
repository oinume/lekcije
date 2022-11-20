package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type Lesson struct {
	lessonRepo          repository.Lesson
	lessonStatusLogRepo repository.LessonStatusLog
}

func NewLesson(
	lessonRepo repository.Lesson,
	lessonStatusLogRepo repository.LessonStatusLog,
) *Lesson {
	return &Lesson{
		lessonRepo:          lessonRepo,
		lessonStatusLogRepo: lessonStatusLogRepo,
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

func (u *Lesson) UpdateLessons(ctx context.Context, lessons []*model2.Lesson) (int, error) {
	if len(lessons) == 0 {
		return 0, nil
	}

	existingLessons, err := u.lessonRepo.FindAllByTeacherIDAndDatetimeAsMap(ctx, lessons[0].TeacherID, lessons)
	if err != nil {
		return 0, err
	}

	rowsAffected := 0
	now := time.Now().UTC()
	for _, lesson := range lessons {
		lesson.Status = strings.ToLower(lesson.Status)
		if existing, ok := existingLessons[model2.LessonDatetime(lesson.Datetime).String()]; ok {
			if lesson.Status == existing.Status {
				continue
			}
			// UPDATE
			if err := u.lessonRepo.UpdateStatus(ctx, existing.ID, lesson.Status); err != nil {
				return 0, err
			}
			if err := u.createLessonStatusLog(ctx, existing.ID, lesson.Status, now); err != nil {
				return 0, err
			}
		} else {
			// INSERT
			dt := lesson.Datetime
			lesson.Datetime = time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), 0, time.UTC)
			l, err := u.lessonRepo.FindOrCreate(ctx, lesson, true)
			if err != nil {
				return 0, err
			}
			// TODO: transaction
			if err := u.createLessonStatusLog(ctx, l.ID, l.Status, now); err != nil {
				return 0, err
			}
		}
		rowsAffected++
	}

	return rowsAffected, nil
}

func (u *Lesson) createLessonStatusLog(ctx context.Context, lessonID uint64, status string, createdAt time.Time) error {
	log := &model2.LessonStatusLog{
		LessonID:  lessonID,
		Status:    status,
		CreatedAt: createdAt,
	}
	if err := u.lessonStatusLogRepo.Create(ctx, log); err != nil {
		return err
	}
	return nil
}
