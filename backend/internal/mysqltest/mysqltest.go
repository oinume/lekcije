package mysqltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/repository"
)

type Repositories struct {
	sqlDB                *sql.DB
	db                   repository.DB
	notificationTimeSpan repository.NotificationTimeSpan
	user                 repository.User
	userAPIToken         repository.UserAPIToken
	userGoogle           repository.UserGoogle
}

func NewRepositories(sqlDB *sql.DB) *Repositories {
	return &Repositories{
		sqlDB:                sqlDB,
		db:                   mysql.NewDBRepository(sqlDB),
		notificationTimeSpan: mysql.NewNotificationTimeSpanRepository(sqlDB),
		user:                 mysql.NewUserRepository(sqlDB),
		userAPIToken:         mysql.NewUserAPITokenRepository(sqlDB),
		userGoogle:           mysql.NewUserGoogleRepository(sqlDB),
	}
}

func (r *Repositories) DB() repository.DB {
	return r.db
}

func (r *Repositories) NotificationTimeSpan() repository.NotificationTimeSpan {
	return r.notificationTimeSpan
}

func (r *Repositories) User() repository.User {
	return r.user
}

func (r *Repositories) UserAPIToken() repository.UserAPIToken {
	return r.userAPIToken
}

func (r *Repositories) UserGoogle() repository.UserGoogle {
	return r.userGoogle
}

func (r *Repositories) CreateNotificationTimeSpans(
	ctx context.Context, t *testing.T,
	timeSpans ...*model2.NotificationTimeSpan,
) {
	t.Helper()
	userIDCheck := make(map[uint]struct{}, len(timeSpans))
	for _, ts := range timeSpans {
		userIDCheck[ts.UserID] = struct{}{}
	}
	if len(userIDCheck) > 1 {
		t.Fatal("CreateNotificationTimeSpans failed because userID in timeSpans is not same")
	}
	for _, ts := range timeSpans {
		if err := r.notificationTimeSpan.UpdateAll(ctx, ts.UserID, timeSpans); err != nil {
			t.Fatalf("CreateNotificationTimeSpans failed: %v", err)
		}
	}
}

func (r *Repositories) CreateUsers(ctx context.Context, t *testing.T, users ...*model2.User) {
	t.Helper()
	for _, u := range users {
		if err := r.user.CreateWithExec(ctx, r.sqlDB, u); err != nil {
			t.Fatalf("CreateUsers failed: %v", err)
		}
	}
}

func (r *Repositories) CreateUserAPITokens(ctx context.Context, t *testing.T, userAPITokens ...*model2.UserAPIToken) {
	t.Helper()
	for _, uat := range userAPITokens {
		if err := r.userAPIToken.Create(ctx, uat); err != nil {
			t.Fatalf("CreateUserAPITokens failed: %v", err)
		}
	}
}

func (r *Repositories) CreateUserGoogles(ctx context.Context, t *testing.T, userGoogles ...*model2.UserGoogle) {
	t.Helper()
	for _, ug := range userGoogles {
		if err := r.userGoogle.CreateWithExec(ctx, r.sqlDB, ug); err != nil {
			t.Fatalf("CreateUserGoogles failed: %v", err)
		}
	}
}
