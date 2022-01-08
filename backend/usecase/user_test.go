package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

func Test_User_CreateWithGoogle(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	userUsecase := NewUser(repos.DB(), repos.User(), repos.UserGoogle())

	type testCase struct {
		wantUser       *model2.User
		wantUserGoogle *model2.UserGoogle
	}
	tests := map[string]struct {
		setup func(ctx context.Context) *testCase
	}{
		"ok_user_and_user_google_will_be_created": {
			setup: func(ctx context.Context) *testCase {
				now := time.Now()
				u := modeltest.NewUser(func(u *model2.User) {
					u.EmailVerified = 1
					u.CreatedAt = now
					u.UpdatedAt = now
				})
				ug := modeltest.NewUserGoogle(func(ug *model2.UserGoogle) {
					ug.UserID = u.ID
					ug.CreatedAt = now
					ug.UpdatedAt = now
				})
				return &testCase{
					wantUser:       u,
					wantUserGoogle: ug,
				}
			},
		},
		"ok_user_and_user_google_exist": {
			setup: func(ctx context.Context) *testCase {
				u := modeltest.NewUser()
				repos.CreateUsers(ctx, t, u)
				ug := modeltest.NewUserGoogle(func(ug *model2.UserGoogle) {
					ug.UserID = u.ID
				})
				repos.CreateUserGoogles(ctx, t, ug)
				return &testCase{
					wantUser:       u,
					wantUserGoogle: ug,
				}
			},
		},
		"ok_another_user_google_exist": {
			setup: func(ctx context.Context) *testCase {
				u := modeltest.NewUser()
				repos.CreateUsers(ctx, t, u)
				ug := modeltest.NewUserGoogle(func(ug *model2.UserGoogle) {
					ug.UserID = u.ID
				})
				repos.CreateUserGoogles(ctx, t, ug)

				now := time.Now().UTC()
				return &testCase{
					wantUser: u,
					wantUserGoogle: modeltest.NewUserGoogle(func(ug2 *model2.UserGoogle) {
						ug2.UserID = u.ID
						ug2.GoogleID = ug.GoogleID + "_another"
						ug2.CreatedAt = now
						ug2.UpdatedAt = now
					}),
				}
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := test.setup(ctx)

			user, userGoogle, err := userUsecase.CreateWithGoogle(ctx, tc.wantUser.Name, tc.wantUser.Email, tc.wantUserGoogle.GoogleID)
			if err != nil {
				t.Fatal(err)
			}
			assertion.AssertEqual(
				t, tc.wantUser, user, "",
				cmpopts.EquateApproxTime(10*time.Second),
				cmpopts.IgnoreFields(model2.User{}, "ID"),
			)
			if user.ID == 0 {
				t.Fatal("user.ID must be set")
			}

			assertion.AssertEqual(
				t, tc.wantUserGoogle, userGoogle, "",
				cmpopts.EquateApproxTime(10*time.Second),
				cmpopts.IgnoreFields(model2.UserGoogle{}, "UserID"),
			)
			if userGoogle.UserID == 0 {
				t.Fatal("userGoogle.ID must be set")
			}
		})
	}
}
