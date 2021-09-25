package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

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

func (r *followingTeacherRepository) CountFollowingTeachersByUserID(ctx context.Context, userID uint) (int, error) {
	count := struct{ Count int }{}
	query := `SELECT COUNT(*) AS count FROM following_teacher WHERE user_id = ?`
	if err := queries.Raw(query, userID).Bind(ctx, r.db, &count); err != nil {
		return 0, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("count failed"),
			errors.WithResource(errors.NewResource("following_teacher", "userID", userID)),
		)
	}
	return count.Count, nil
}

func (r *followingTeacherRepository) Create(ctx context.Context, followingTeacher *model2.FollowingTeacher) error {
	return followingTeacher.Insert(ctx, r.db, boil.Infer())
}

func (r *followingTeacherRepository) DeleteByUserIDAndTeacherIDs(ctx context.Context, userID uint, teacherIDs []uint) error {
	placeholder := strings.TrimRight(strings.Repeat("?,", len(teacherIDs)), ",")
	where := fmt.Sprintf("user_id = ? AND teacher_id IN (%s)", placeholder)
	_, err := model2.FollowingTeachers(qm.Where(where)).DeleteAll(ctx, r.db)
	return err
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

func (r *followingTeacherRepository) FindByUserID(
	ctx context.Context, userID uint,
) ([]*model2.FollowingTeacher, error) {
	fts, err := model2.FollowingTeachers(qm.Where("user_id = ?", userID)).All(ctx, r.db)
	if err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource("following_teacher", "userID", userID)),
		)
	}
	return fts, nil
}
