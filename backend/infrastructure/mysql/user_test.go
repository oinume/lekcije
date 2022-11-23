package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model2"
)

func Test_userRepository_FindAllByEmailVerifiedIsTrue(t *testing.T) {
	repo := mysql.NewUserRepository(helper.DB(t).DB())
	repos := mysqltest.NewRepositories(helper.DB(t).DB())

	type testCase struct {
		wantUsers []*model2.User
	}
	tests := map[string]struct {
		setup   func(ctx context.Context) *testCase
		wantErr bool
	}{
		"normal": {
			setup: func(ctx context.Context) *testCase {
				helper.TruncateAllTables(t)
				//boil.DebugMode = true
				//boil.DebugWriter = os.Stdout
				user1 := modeltest.NewUser(func(u *model2.User) {
					u.EmailVerified = 1
				})
				user2 := modeltest.NewUser(func(u *model2.User) {
					u.EmailVerified = 1
				})
				user3 := modeltest.NewUser(func(u *model2.User) {
					u.EmailVerified = 1
				})
				repos.CreateUsers(ctx, t, user1, user2, user3)

				ft1 := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
					ft.UserID = user1.ID
				})
				ft2 := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
					ft.UserID = user2.ID
				})
				repos.CreateFollowingTeachers(ctx, t, ft1)
				repos.CreateFollowingTeachers(ctx, t, ft2)

				return &testCase{
					wantUsers: []*model2.User{user1, user2},
				}
			},
		},
		"no email verified users": {
			setup: func(ctx context.Context) *testCase {
				helper.TruncateAllTables(t)
				user1 := modeltest.NewUser(func(u *model2.User) {
					u.EmailVerified = 0
				})
				repos.CreateUsers(ctx, t, user1)
				ft1 := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
					ft.UserID = user1.ID
				})
				repos.CreateFollowingTeachers(ctx, t, ft1)

				return &testCase{
					wantUsers: []*model2.User{},
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			gotUsers, err := repo.FindAllByEmailVerifiedIsTrue(ctx, 10)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAllByEmailVerifiedIsTrue: error = %v, wantErr %v", err, tt.wantErr)
			}

			assertion.AssertEqual(t, len(tc.wantUsers), len(gotUsers), "user length doesn't match")
			for i := range tc.wantUsers {
				assertion.AssertEqual(t, tc.wantUsers[i], gotUsers[i], "", cmpopts.EquateApproxTime(10*time.Second))
			}
		})
	}
}
