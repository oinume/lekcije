package mysql

import (
	"context"
	"database/sql"

	"github.com/morikuni/failure"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
)

type lessonStatusLogRepo struct {
	db *sql.DB
}

func NewLessonStatusLogRepository(db *sql.DB) repository.LessonStatusLog {
	return &lessonStatusLogRepo{db: db}
}

func (r *lessonStatusLogRepo) Create(ctx context.Context, log *model2.LessonStatusLog) error {
	if err := log.Insert(ctx, r.db, boil.Infer()); err != nil {
		return failure.Translate(err, errors.Internal)
	}
	return nil
}
