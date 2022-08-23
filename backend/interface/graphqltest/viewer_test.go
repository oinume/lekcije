package graphqltest

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Khan/genqlient/graphql"

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
	//var log bytes.Buffer
	//appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))

	type testCase struct {
		apiToken   string
		input      UpdateViewerInput
		wantResult UpdateViewerUpdateViewerUser
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
				return &testCase{
					apiToken: userAPIToken.Token,
					input: UpdateViewerInput{
						Email: "xyz@example.com",
					},
					wantResult: UpdateViewerUpdateViewerUser{
						Id:    fmt.Sprint(user.ID),
						Email: user.Email,
					},
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
			//server := httptest.NewServer(interfacehttp.WithIDToken(interfacehttp.WithSessionID(graphqlServer.ServeHTTP,
			// config.ServiceEnvTest)))
			server := httptest.NewServer(graphqlServer)
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
				t.Fatalf("unexpected error: %v", err)
			}
			//if err != nil {
			//	if c.wantErrorCode == "" {
			//		t.Fatalf("unexpected error: %v", err)
			//	} else {
			//		if !strings.Contains(err.Error(), c.wantErrorCode.ErrorCode()) {
			//			t.Fatalf("err must contain code: wantErrorCode=%v, err=%v", c.wantErrorCode, err)
			//		}
			//		return // OK
			//	}
			//}
			//
			//if c.wantErrorCode != "" {
			//	t.Fatalf("wantErrorCode is not empty but no error: wantErrorCode=%v", c.wantErrorCode)
			//}
			assertion.AssertEqual(t, tc.wantResult, resp.GetUpdateViewer(), "")
		})
	}
}
