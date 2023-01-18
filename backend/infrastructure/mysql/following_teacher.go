package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/morikuni/failure"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
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
	_, err := model2.FindFollowingTeacher(ctx, r.db, followingTeacher.UserID, followingTeacher.TeacherID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return followingTeacher.Insert(ctx, r.db, boil.Infer())
	}
	// Do nothing
	return nil
}

func (r *followingTeacherRepository) DeleteByUserIDAndTeacherIDs(ctx context.Context, userID uint, teacherIDs []uint) error {
	placeholder := strings.TrimRight(strings.Repeat("?,", len(teacherIDs)), ",")
	where := fmt.Sprintf("user_id = ? AND teacher_id IN (%s)", placeholder)
	args := []interface{}{userID}
	for _, teacherID := range teacherIDs {
		args = append(args, teacherID)
	}
	_, err := model2.FollowingTeachers(qm.Where(where, args...)).DeleteAll(ctx, r.db)
	return err
}

func (r *followingTeacherRepository) FindTeacherIDsByUserID(
	ctx context.Context,
	userID uint,
	fetchErrorCount int,
	lastLessonAt time.Time,
) ([]uint, error) {
	_, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "followingTeacherRepository.FindTeacherIDsByUserID")
	span.SetAttributes(attribute.KeyValue{
		Key:   "userID",
		Value: attribute.Int64Value(int64(userID)),
	})
	defer span.End()

	values := make([]*model2.FollowingTeacher, 0, 1000)
	query := `
	SELECT ft.teacher_id FROM following_teacher AS ft
	INNER JOIN teacher AS t ON ft.teacher_id = t.id
	WHERE
      ft.user_id = ?
      AND t.fetch_error_count <= ?
      AND (t.last_lesson_at >= ? OR t.last_lesson_at = '0000-00-00 00:00:00')
	`
	if err := queries.Raw(query, userID, fetchErrorCount, FormatDateTime(lastLessonAt)).Bind(ctx, r.db, &values); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return empty slice without error
		}
		return nil, failure.Wrap(err, errors.NewUserIDContext(userID))
	}
	ids := make([]uint, len(values))
	for i, t := range values {
		ids[i] = t.TeacherID
	}
	return ids, nil
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
	fts, err := model2.FollowingTeachers(
		qm.Where("user_id = ?", userID),
		qm.OrderBy("created_at DESC"),
	).All(ctx, r.db)
	if err != nil {
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource("following_teacher", "userID", userID)),
		)
	}
	return fts, nil
}

func (r *followingTeacherRepository) FindByUserIDAndTeacherID(
	ctx context.Context, userID uint, teacherID uint,
) (*model2.FollowingTeacher, error) {
	return model2.FindFollowingTeacher(ctx, r.db, userID, teacherID)
}
