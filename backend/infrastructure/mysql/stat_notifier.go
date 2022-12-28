package mysql

import (
	"context"
	"database/sql"
	"fmt"

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
	err := transaction(ctx, r.db, func(exec repository.Executor) error {
		_, err := model2.FindStatNotifier(ctx, exec, statNotifier.Datetime, statNotifier.Interval)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			fmt.Printf("INSERT: err=%v\n", err)
			return statNotifier.Insert(ctx, exec, boil.Infer())
		}
		_, err = statNotifier.Update(ctx, exec, boil.Infer())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return failure.Wrap(err, failure.Message("statNotifierRepository.CreateOrUpdate failed"))
	}
	return nil
}
