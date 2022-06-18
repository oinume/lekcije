package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type teacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) repository.Teacher {
	return &teacherRepository{
		db: db,
	}
}

func (r *teacherRepository) Create(ctx context.Context, teacher *model2.Teacher) error {
	return teacher.Insert(ctx, r.db, boil.Infer())
}

func (r *teacherRepository) CreateOrUpdate(ctx context.Context, teacher *model2.Teacher) error {
	return teacher.Upsert(ctx, r.db, boil.Infer(), boil.Infer())
}

func (r *teacherRepository) FindByIDs(ctx context.Context, ids []uint) ([]*model2.Teacher, error) {
	return model2.Teachers(qm.WhereIn("id IN ?", fromUintSlice(ids)...)).All(ctx, r.db)
}

func fromUintSlice(values []uint) []interface{} {
	ret := make([]interface{}, len(values))
	for i, value := range values {
		ret[i] = interface{}(value)
	}
	return ret
}
