package mysql

import (
	"context"
	"database/sql"
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
