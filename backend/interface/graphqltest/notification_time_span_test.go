package graphqltest

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

var _ = `# @genqlient
mutation UpdateNotificationTimeSpans($input: UpdateNotificationTimeSpansInput!) {
  updateNotificationTimeSpans(input: $input) {
    timeSpans {
      fromHour
      fromMinute
      toHour
      toMinute
    }
  }
}
`

func TestUpdateNotificationTimeSpan(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())

	type testCase struct {
		apiToken      string
		input         UpdateNotificationTimeSpansInput
		wantResult    UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload
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
				return &testCase{
					apiToken: userAPIToken.Token,
					input: UpdateNotificationTimeSpansInput{
						TimeSpans: []NotificationTimeSpanInput{
							{
								FromHour:   13,
								FromMinute: 00,
								ToHour:     19,
								ToMinute:   30,
							},
							{
								FromHour:   21,
								FromMinute: 30,
								ToHour:     23,
								ToMinute:   30,
							},
						},
					},
					wantResult: UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayload{
						TimeSpans: []UpdateNotificationTimeSpansUpdateNotificationTimeSpansNotificationTimeSpanPayloadTimeSpansNotificationTimeSpan{
							{
								FromHour:   13,
								FromMinute: 00,
								ToHour:     19,
								ToMinute:   30,
							},
							{
								FromHour:   21,
								FromMinute: 30,
								ToHour:     23,
								ToMinute:   30,
							},
						},
					},
				}
			},
		},
		"invalid_argument_over_3": {
			setup: func(ctx context.Context) *testCase {
				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)
				userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
					uat.UserID = user.ID
				})
				repos.CreateUserAPITokens(ctx, t, userAPIToken)
				return &testCase{
					apiToken: userAPIToken.Token,
					input: UpdateNotificationTimeSpansInput{
						TimeSpans: []NotificationTimeSpanInput{
							{
								FromHour:   13,
								FromMinute: 00,
								ToHour:     19,
								ToMinute:   30,
							},
							{
								FromHour:   21,
								FromMinute: 30,
								ToHour:     23,
								ToMinute:   30,
							},
							{
								FromHour:   8,
								FromMinute: 30,
								ToHour:     9,
								ToMinute:   30,
							},
							{
								FromHour:   10,
								FromMinute: 30,
								ToHour:     11,
								ToMinute:   30,
							},
						},
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
			notificationTimeSpanUsecase := usecase.NewNotificationTimeSpan(repos.NotificationTimeSpan())
			userUsecase := usecase.NewUser(repos.DB(), repos.User(), repos.UserGoogle())
			graphqlServer := newGraphQLServer(repos, notificationTimeSpanUsecase, userUsecase)
			server := httptest.NewServer(setAuthorizationContext(graphqlServer))
			t.Cleanup(func() { server.Close() })

			httpClient := server.Client()
			transport := httpClient.Transport
			httpClient.Transport = &authTransport{
				parent: transport,
				token:  tc.apiToken,
			}
			graphqlClient := graphql.NewClient(server.URL, httpClient)

			resp, err := UpdateNotificationTimeSpans(ctx, graphqlClient, tc.input)
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

			if tc.wantErrorCode != "" {
				t.Fatalf("wantErrorCode is not empty but no error: wantErrorCode=%v", tc.wantErrorCode)
			}

			assertion.AssertEqual(t, tc.wantResult, resp.GetUpdateNotificationTimeSpans(), "")
		})
	}
}
