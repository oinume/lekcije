package mysql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type lessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) repository.Lesson {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) Create(ctx context.Context, lesson *model2.Lesson) error {
	return lesson.Insert(ctx, r.db, boil.Infer())
}

func (r *lessonRepository) FindAllByTeacherIDsDatetimeBetween(
	ctx context.Context, teacherID uint,
	fromDate, toDate time.Time,
) ([]*model2.Lesson, error) {
	_, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "lessonRepository.FindLessons")
	span.SetAttributes(attribute.KeyValue{
		Key:   "teacherID",
		Value: attribute.Int64Value(int64(teacherID)),
	})
	defer span.End()

	const midnightAdd = 2
	const format = "2006-01-02"
	toDateAdded := toDate.Add(24 * midnightAdd * time.Hour)
	lessons, err := model2.Lessons(
		qm.Where(
			"teacher_id = ? AND DATE(datetime) BETWEEN ? AND ?",
			teacherID, fromDate.Format(format), toDateAdded.Format(format),
		),
		qm.Limit(1000),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return lessons, nil
}

func (r *lessonRepository) GetNewAvailableLessons(ctx context.Context, oldLessons, newLessons []*model2.Lesson) []*model2.Lesson {
	// Pattern
	// 2016-01-01 00:00@Any -> Available
	oldLessonsMap := make(map[string]*model2.Lesson, len(oldLessons))
	newLessonsMap := make(map[string]*model2.Lesson, len(newLessons))
	availableLessons := make([]*model2.Lesson, 0, len(oldLessons)+len(newLessons))
	availableLessonsMap := make(map[string]*model2.Lesson, len(oldLessons)+len(newLessons))
	for _, l := range oldLessons {
		oldLessonsMap[model2.LessonDatetime(l.Datetime).String()] = l // TODO: Use LessonDatetime type as key
	}
	for _, l := range newLessons {
		newLessonsMap[model2.LessonDatetime(l.Datetime).String()] = l
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
		datetime := model2.LessonDatetime(l.Datetime).String()
		if _, ok := oldLessonsMap[datetime]; !ok && strings.ToLower(l.Status) == "available" {
			// not exists in oldLessons
			availableLessons = append(availableLessons, l)
			availableLessonsMap[datetime] = l
		}
	}

	// TODO: sort availableLessonsMap by datetime
	return availableLessons
}
