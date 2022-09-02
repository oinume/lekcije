package graphqltest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	interfacehttp "github.com/oinume/lekcije/backend/interface/http"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

var _ = `# @genqlient
mutation UpdateViewer($input: UpdateViewerInput!) {
  updateViewer(input: $input) {
    id
    email
  }
}
`

func TestUpdateViewer(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())

	type testCase struct {
		apiToken      string
		input         UpdateViewerInput
		wantResult    UpdateViewerUpdateViewerUser
		wantErrorCode failure.StringCode
		//wantResponse   *api_v1.GetMeResponse
		//wantStatusCode int
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
					input: UpdateViewerInput{
						Email: newEmail,
					},
					wantResult: UpdateViewerUpdateViewerUser{
						Id:    fmt.Sprint(user.ID),
						Email: newEmail,
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
					input: UpdateViewerInput{
						Email: newEmail,
					},
					wantResult: UpdateViewerUpdateViewerUser{
						Id:    fmt.Sprint(user.ID),
						Email: newEmail,
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
					input: UpdateViewerInput{
						Email: newEmail,
					},
					wantResult: UpdateViewerUpdateViewerUser{
						Id:    fmt.Sprint(user.ID),
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
			userUsecase := usecase.NewUser(repos.DB(), repos.User(), repos.UserGoogle())
			graphqlServer := newGraphQLServer(repos, userUsecase)
			server := httptest.NewServer(setAuthorizationContext(graphqlServer))
			t.Cleanup(func() { server.Close() })

			httpClient := server.Client()
			transport := httpClient.Transport
			httpClient.Transport = &authTransport{
				parent: transport,
				token:  tc.apiToken,
			}
			graphqlClient := graphql.NewClient(server.URL, httpClient)

			resp, err := UpdateViewer(ctx, graphqlClient, tc.input)
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

			assertion.AssertEqual(t, tc.wantResult, resp.GetUpdateViewer(), "")
		})
	}
}

func setAuthorizationContext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		auth, err := interfacehttp.ParseAuthorizationHeader(r.Header.Get("authorization"))
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(context_data.SetAPIToken(r.Context(), strings.TrimSpace(auth)))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
