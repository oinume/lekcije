package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
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

func (u *Lesson) UpdateLessons(ctx context.Context, lessons []*model2.Lesson) (int64, error) {
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
		if l, ok := existingLessons[lesson.Datetime.Format(lessonTimeFormat)]; ok {
			if lesson.Status == l.Status {
				continue
			}
			// UPDATE
			values := &Lesson{Status: lesson.Status, UpdatedAt: now}
			if err := s.db.Model(lesson).Where("id = ?", l.ID).Updates(values).Error; err != nil {
				return 0, errors.NewInternalError(
					errors.WithError(err),
				)
			}
			rowsAffected++

			log := &LessonStatusLog{
				LessonID:  l.ID,
				Status:    lesson.Status,
				CreatedAt: now,
			}
			if err := NewLessonStatusLogService(s.db).Create(log); err != nil {
				return 0, err
			}
		} else {
			// INSERT
			dt := lesson.Datetime
			lesson.Datetime = time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), 0, time.UTC)
			where := Lesson{TeacherID: lesson.TeacherID, Datetime: lesson.Datetime}
			if err := s.db.Where(where).FirstOrCreate(lesson).Error; err != nil {
				return 0, errors.NewInternalError(
					errors.WithError(err),
					errors.WithMessage("FirstOrCreate failed"),
					errors.WithResource(errors.NewResourceWithEntries(
						s.TableName(),
						[]errors.ResourceEntry{
							{Key: "teacherID", Value: lesson.TeacherID},
							{Key: "datetime", Value: lesson.Datetime},
						},
					)),
				)
			}
			rowsAffected++

			log := &LessonStatusLog{
				LessonID:  lesson.ID,
				Status:    lesson.Status,
				CreatedAt: now,
			}
			if err := NewLessonStatusLogService(s.db).Create(log); err != nil {
				return 0, err
			}
		}
	}

	return int64(rowsAffected), nil
}
