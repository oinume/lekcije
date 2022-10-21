package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/goenum"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/util"
)

const (
	lessonTimeFormat = "2006-01-02 15:04"
)

type Lesson struct {
	ID        uint64 `gorm:"primary_key"`
	TeacherID uint32
	Datetime  time.Time
	Status    string // TODO: enum
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Lesson) TableName() string {
	return "lesson"
}

func (l *Lesson) String() string {
	return fmt.Sprintf(
		"TeacherID=%v, Datetime=%v, Status=%v",
		l.TeacherID, l.Datetime.Format(lessonTimeFormat), l.Status,
	)
}

type LessonStatus struct {
	Finished  int `goenum:"終了"`
	Reserved  int `goenum:"予約済"`
	Available int `goenum:"予約可"`
	Cancelled int `goenum:"休講"`
}

var LessonStatuses = goenum.EnumerateStruct(&LessonStatus{
	Finished:  1,
	Reserved:  2,
	Available: 3,
	Cancelled: 4,
})

type LessonService struct {
	db *gorm.DB
}

func NewLessonService(db *gorm.DB) *LessonService {
	return &LessonService{db: db}
}

func (s *LessonService) TableName() string {
	return (&Lesson{}).TableName()
}

func (s *LessonService) UpdateLessons(lessons []*Lesson) (int64, error) {
	if len(lessons) == 0 {
		return 0, nil
	}

	existingLessons, err := s.FindLessonsByTeacherIDAndDatetimeAsMap(lessons[0].TeacherID, lessons)
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

func (s *LessonService) FindLessonsByTeacherIDAndDatetimeAsMap(
	teacherID uint32,
	lessonsArgs []*Lesson,
) (map[string]*Lesson, error) {
	if len(lessonsArgs) == 0 {
		return nil, nil
	}

	datetimes := make([]string, len(lessonsArgs))
	for i, l := range lessonsArgs {
		datetimes[i] = l.Datetime.Format(dbDatetimeFormat)
	}

	placeholder := Placeholders(util.StringToInterfaceSlice(datetimes...))
	lessons := make([]*Lesson, 0, len(lessonsArgs))
	sql := strings.TrimSpace(fmt.Sprintf(`
SELECT * FROM %s
WHERE
  teacher_id = ?
  AND datetime IN (%s)
`, s.TableName(), placeholder))

	values := []interface{}{teacherID}
	values = append(values, util.StringToInterfaceSlice(datetimes...)...)
	result := s.db.Raw(sql, values...).Scan(&lessons)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to find lessonsArgs"),
			errors.WithResource(errors.NewResource(s.TableName(), "teacherID", teacherID)),
		)
	}

	ret := make(map[string]*Lesson, len(lessons))
	for _, l := range lessons {
		ret[l.Datetime.Format(lessonTimeFormat)] = l
	}
	return ret, nil
}

func (s *LessonService) FindLessons(
	ctx context.Context,
	teacherID uint32,
	fromDate,
	toDate time.Time,
) ([]*Lesson, error) {
	_, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "LessonService.FindLessons")
	span.SetAttributes(attribute.KeyValue{
		Key:   "teacherID",
		Value: attribute.Int64Value(int64(teacherID)),
	})
	defer span.End()

	lessons := make([]*Lesson, 0, 1000)
	sql := strings.TrimSpace(fmt.Sprintf(`
SELECT * FROM %s
WHERE
  teacher_id = ?
  AND DATE(datetime) BETWEEN ? AND ?
ORDER BY datetime ASC
LIMIT 1000
	`, s.TableName()))

	toDateAdded := toDate.Add(24 * 2 * time.Hour)
	result := s.db.Raw(sql, teacherID, fromDate.Format("2006-01-02"), toDateAdded.Format("2006-01-02")).Scan(&lessons)
	if result.Error != nil {
		if result.RecordNotFound() {
			return lessons, nil
		}
		return nil, errors.NewInternalError(
			errors.WithError(result.Error),
			errors.WithMessage("Failed to find lessons"),
			errors.WithResource(errors.NewResource(s.TableName(), "teacherID", teacherID)),
		)
	}

	return lessons, nil
}

func (s *LessonService) GetNewAvailableLessons(ctx context.Context, oldLessons, newLessons []*Lesson) []*Lesson {
	_, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "LessonService.GetNewAvailableLessons")
	defer span.End()

	// Pattern
	// 2016-01-01 00:00@Any -> Available
	oldLessonsMap := make(map[string]*Lesson, len(oldLessons))
	newLessonsMap := make(map[string]*Lesson, len(newLessons))
	availableLessons := make([]*Lesson, 0, len(oldLessons)+len(newLessons))
	availableLessonsMap := make(map[string]*Lesson, len(oldLessons)+len(newLessons))
	for _, l := range oldLessons {
		oldLessonsMap[l.Datetime.Format(lessonTimeFormat)] = l
	}
	for _, l := range newLessons {
		newLessonsMap[l.Datetime.Format(lessonTimeFormat)] = l
	}
	for datetime, oldLesson := range oldLessonsMap {
		newLesson, newLessonExists := newLessonsMap[datetime]
		oldStatus := strings.ToLower(oldLesson.Status)
		newStatus := strings.ToLower(newLesson.Status)
		if newLessonExists && oldStatus != "available" && newStatus == "available" {
			// exists in oldLessons and newLessons and "any status" -> "available"
			availableLessons = append(availableLessons, newLesson)
			availableLessonsMap[datetime] = newLesson
		}
	}
	for _, l := range newLessons {
		datetime := l.Datetime.Format(lessonTimeFormat)
		if _, ok := oldLessonsMap[datetime]; !ok && strings.ToLower(l.Status) == "available" {
			// not exists in oldLessons
			availableLessons = append(availableLessons, l)
			availableLessonsMap[datetime] = l
		}
	}

	// TODO: sort availableLessonsMap by datetime
	return availableLessons
}
