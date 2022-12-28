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

type statNotifierRepository struct {
	db *sql.DB
}

func NewStatNotifierRepository(db *sql.DB) repository.StatNotifier {
	return &statNotifierRepository{db: db}
}

func (r *statNotifierRepository) CreateOrUpdate(ctx context.Context, statNotifier *model2.StatNotifier) error {
	_, err := model2.FindStatNotifier(ctx, r.db, statNotifier.Datetime, statNotifier.Interval)
	if err != nil {
		if !errors.IsNotFound(err) {
			return failure.Wrap(err)
		}
		return statNotifier.Insert(ctx, r.db, boil.Infer())
	}
	_, err = statNotifier.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return failure.Wrap(err)
	}
	return nil
}
