package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type statNotifierRepository struct {
	db *sql.DB
}

func NewStatNotifierRepository(db *sql.DB) repository.StatNotifier {
	return &statNotifierRepository{db: db}
}

func (r *statNotifierRepository) CreateOrUpdate(ctx context.Context, statNotifier *model2.StatNotifier) error {
	return statNotifier.Upsert(ctx, r.db, boil.Infer(), boil.Infer())
}
