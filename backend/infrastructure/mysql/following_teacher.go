package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type followingTeacherRepository struct {
	db *sql.DB
}

func NewFollowingTeacherRepository(db *sql.DB) repository.FollowingTeacher {
	return &followingTeacherRepository{
		db: db,
	}
}

func (r *followingTeacherRepository) Create(ctx context.Context, followingTeacher *model2.FollowingTeacher) error {
	return followingTeacher.Insert(ctx, r.db, boil.Infer())
}

func (r *followingTeacherRepository) FindTeachersByUserID(ctx context.Context, userID uint) ([]*model2.Teacher, error) {
	query := `
		SELECT t.* FROM following_teacher AS ft
		INNER JOIN teacher AS t ON ft.teacher_id = t.id
		WHERE ft.user_id = ?
		ORDER BY ft.created_at DESC
	`
	teachers := make([]*model2.Teacher, 0, 100)
	if err := queries.Raw(query, userID).Bind(ctx, r.db, &teachers); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		// TODO: Wrap error with morikuni/failure
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource("following_teacher", "user_id", userID)),
		)
	}
	// TODO: expose FollowingTeacher.doAfterSelectHooks in template
	return teachers, nil
}
