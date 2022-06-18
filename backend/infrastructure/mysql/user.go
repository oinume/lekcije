package mysql

// TODO: move this package to interface/mysql package

import (
	"context"
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.User {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateWithExec(ctx context.Context, exec repository.Executor, user *model2.User) error {
	return user.Insert(ctx, exec, boil.Infer())
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model2.User, error) {
	return r.FindByEmailWithExec(ctx, r.db, email)
}

func (r *userRepository) FindByEmailWithExec(ctx context.Context, exec repository.Executor, email string) (*model2.User, error) {
	return model2.Users(qm.Where("email = ?", email)).One(ctx, exec)
}

func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*model2.User, error) {
	return r.findByGoogleIDWithExec(ctx, r.db, googleID)
}

func (r *userRepository) FindByGoogleIDWithExec(ctx context.Context, exec repository.Executor, googleID string) (*model2.User, error) {
	return r.findByGoogleIDWithExec(ctx, exec, googleID)
}

func (r *userRepository) findByGoogleIDWithExec(ctx context.Context, exec repository.Executor, googleID string) (*model2.User, error) {
	query := `
		SELECT u.* FROM user AS u
		INNER JOIN user_google AS ug ON u.id = ug.user_id
		WHERE ug.google_id = ?
		LIMIT 1
	`
	u := &model2.User{}
	if err := queries.Raw(query, googleID).Bind(ctx, exec, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		// TODO: Wrap error with morikuni/failure
		return nil, errors.NewInternalError(
			errors.WithError(err),
			errors.WithResource(errors.NewResource("user_google", "google_id", googleID)),
		)
	}
	// TODO: expose User.doAfterSelectHooks in template

	return u, nil
}

func (r *userRepository) UpdateEmail(ctx context.Context, id uint, email string) error {
	const query = `UPDATE user SET email = ?, updated_at = NOW() WHERE id = ?`
	_, err := queries.Raw(query, email, id).ExecContext(ctx, r.db)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to update user email"),
			errors.WithResource(errors.NewResourceWithEntries(
				"user", []errors.ResourceEntry{
					{Key: "id", Value: id}, {Key: "email", Value: email},
				},
			)),
		)
	}
	return nil
}

func (r *userRepository) UpdateFollowedTeacherAt(ctx context.Context, id uint, time time.Time) error {
	const query = "UPDATE user SET followed_teacher_at = ? WHERE id = ?"
	_, err := queries.Raw(query, FormatDateTime(time), id).ExecContext(ctx, r.db)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to update user followed_teacher_at"),
		)
	}
	return nil
}

func (r *userRepository) UpdateOpenNotificationAt(ctx context.Context, id uint, time time.Time) error {
	const query = "UPDATE user SET open_notification_at = ? WHERE id = ?"
	_, err := queries.Raw(query, FormatDateTime(time), id).ExecContext(ctx, r.db)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to update user open_notification_at"),
		)
	}
	return nil
}
