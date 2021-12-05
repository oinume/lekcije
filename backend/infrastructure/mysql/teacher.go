package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
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
