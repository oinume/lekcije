package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/morikuni/failure"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
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

func (r *teacherRepository) IncrementFetchErrorCount(ctx context.Context, id uint, value int) error {
	sql := `UPDATE teacher SET fetch_error_count = fetch_error_count + ?, updated_at = NOW() WHERE id = ?`
	if _, err := queries.Raw(sql, value, id).ExecContext(ctx, r.db); err != nil {
		return failure.Wrap(err, failure.Context{"id": fmt.Sprint(id)})
	}
	return nil
}

func fromUintSlice(values []uint) []interface{} {
	ret := make([]interface{}, len(values))
	for i, value := range values {
		ret[i] = interface{}(value)
	}
	return ret
}
