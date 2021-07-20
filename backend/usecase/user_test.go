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
		"normal": {
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
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := test.setup(ctx)

			user, userGoogle, err := userUsecase.CreateWithGoogle(ctx, tc.wantUser.Name, tc.wantUser.Email, tc.wantUserGoogle.GoogleID)
			if err != nil {
				t.Fatal(err)
			}
			// TODO: validate user, userGoogle
			assertion.AssertEqual(t, tc.wantUser, user, "", cmpopts.EquateApproxTime(10*time.Second))
			assertion.AssertEqual(t, tc.wantUserGoogle, userGoogle, "", cmpopts.EquateApproxTime(10*time.Second))
		})
	}
}
