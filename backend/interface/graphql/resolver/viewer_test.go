package resolver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	graphqlmodel "github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

func TestUpdateViewer(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	userUsecase := usecase.NewUser(repos.DB(), repos.User(), repos.UserGoogle())
	resolver := NewResolver(
		repos.FollowingTeacher(),
		repos.NotificationTimeSpan(),
		nil,
		repos.Teacher(),
		repos.User(),
		userUsecase,
	)

	type testCase struct {
		apiToken      string
		input         graphqlmodel.UpdateViewerInput
		wantResult    *graphqlmodel.User
		wantErrorCode failure.StringCode
	}
	tests := map[string]struct {
		setup func(ctx context.Context) *testCase
	}{
		"ok": {
			setup: func(ctx context.Context) *testCase {
				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)
				userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
					uat.UserID = user.ID
				})
				repos.CreateUserAPITokens(ctx, t, userAPIToken)
				newEmail := fmt.Sprintf("updated-%d@example.com", user.ID)
				return &testCase{
					apiToken: userAPIToken.Token,
					input: graphqlmodel.UpdateViewerInput{
						Email: &newEmail,
					},
					wantResult: &graphqlmodel.User{
						ID:           fmt.Sprint(user.ID),
						Email:        newEmail,
						ShowTutorial: !user.IsFollowedTeacher(),
					},
					wantErrorCode: "",
				}
			},
		},
		"invalid_email_format": {
			setup: func(ctx context.Context) *testCase {
				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)
				userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
					uat.UserID = user.ID
				})
				repos.CreateUserAPITokens(ctx, t, userAPIToken)
				newEmail := "invalid"
				return &testCase{
					apiToken: userAPIToken.Token,
					input: graphqlmodel.UpdateViewerInput{
						Email: &newEmail,
					},
					wantResult: &graphqlmodel.User{
						ID:           fmt.Sprint(user.ID),
						Email:        user.Email,
						ShowTutorial: !user.IsFollowedTeacher(),
					},
					wantErrorCode: errors.InvalidArgument,
				}
			},
		},
		"duplicate_email": {
			setup: func(ctx context.Context) *testCase {
				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)
				userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
					uat.UserID = user.ID
				})
				repos.CreateUserAPITokens(ctx, t, userAPIToken)
				newEmail := user.Email
				return &testCase{
					apiToken: userAPIToken.Token,
					input: graphqlmodel.UpdateViewerInput{
						Email: &newEmail,
					},
					wantResult: &graphqlmodel.User{
						ID:    fmt.Sprint(user.ID),
						Email: newEmail,
					},
					wantErrorCode: errors.InvalidArgument,
				}
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := test.setup(ctx)
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			got, err := resolver.Mutation().UpdateViewer(ctx, tc.input)
			if err != nil {
				if tc.wantErrorCode == "" {
					t.Fatalf("unexpected error: %v", err)
				} else {
					if !strings.Contains(err.Error(), tc.wantErrorCode.ErrorCode()) {
						t.Fatalf("err must contain code: wantErrorCode=%v, err=%v", tc.wantErrorCode, err)
					}
					return // OK
				}
			}
			fmt.Printf("user = %+v\n", got)

			if tc.wantErrorCode != "" {
				t.Fatalf("wantErrorCode is not empty but no error: wantErrorCode=%v", tc.wantErrorCode)
			}

			assertion.AssertEqual(t, tc.wantResult, got, "")
		})
	}
}
