package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/morikuni/failure"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel"

	"github.com/oinume/lekcije/backend/domain/config"
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

func (r *userRepository) FindAllByEmailVerifiedIsTrue(ctx context.Context, notificationInterval int) ([]*model2.User, error) {
	_, span := otel.Tracer(config.DefaultTracerName).Start(ctx, "UserService.FindAllEmailVerifiedIsTrue")
	defer span.End()

	query := `
		SELECT u.* FROM (SELECT DISTINCT(user_id) FROM following_teacher) AS ft
		INNER JOIN user AS u ON ft.user_id = u.id
		INNER JOIN m_plan AS mp ON u.plan_id = mp.id
		WHERE
		  u.email_verified = 1
		  AND mp.notification_interval = ?
		ORDER BY u.open_notification_at DESC
`
	var users []*model2.User
	if err := queries.Raw(query, notificationInterval).Bind(ctx, r.db, &users); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return empty slice without error
		}
		return nil, failure.Wrap(err)
	}
	return users, nil
}

func (r *userRepository) FindByAPIToken(ctx context.Context, apiToken string) (*model2.User, error) {
	query := `
		SELECT u.* FROM user AS u
		INNER JOIN user_api_token AS uat ON u.id = uat.user_id
		WHERE uat.token = ?
		LIMIT 1
	`
	u := &model2.User{}
	if err := queries.Raw(query, apiToken).Bind(ctx, r.db, u); err != nil {
		if err == sql.ErrNoRows {
			return nil, failure.Translate(err, errors.NotFound)
		}
		return nil, failure.Wrap(err)
	}
	// TODO: expose User.doAfterSelectHooks in template
	return u, nil
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

func (r *userRepository) FindAllByEmailVerified(
	ctx context.Context, notificationInterval int,
) ([]*model2.User, error) {
	//TODO implement me
	panic("implement me")
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
