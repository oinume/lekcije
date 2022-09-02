package repository

import (
	"context"
	"time"

	"github.com/oinume/lekcije/backend/model2"
)

type User interface {
	CreateWithExec(ctx context.Context, exec Executor, user *model2.User) error
	FindByAPIToken(ctx context.Context, apiToken string) (*model2.User, error)
	FindByEmail(ctx context.Context, email string) (*model2.User, error)
	FindByEmailWithExec(ctx context.Context, exec Executor, email string) (*model2.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*model2.User, error)
	FindByGoogleIDWithExec(ctx context.Context, exec Executor, googleID string) (*model2.User, error)
	FindAllByEmailVerified(ctx context.Context, notificationInterval int) ([]*model2.User, error)
	UpdateEmail(ctx context.Context, id uint, email string) error
	UpdateFollowedTeacherAt(ctx context.Context, id uint, time time.Time) error
	UpdateOpenNotificationAt(ctx context.Context, id uint, time time.Time) error
}
