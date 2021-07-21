package http_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/di"
	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	interface_http "github.com/oinume/lekcije/backend/interface/http"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

func Test_UserService_GetMe(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	service := newUserService(db, appLogger)
	handler := api_v1.NewUserServer(service)

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken       string
		request        *api_v1.GetMeRequest
		wantResponse   *api_v1.GetMeResponse
		wantStatusCode int
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
					request:  &api_v1.GetMeRequest{},
					wantResponse: &api_v1.GetMeResponse{
						UserId: int32(user.ID),
						Email:  user.Email,
					},
					wantStatusCode: http.StatusOK,
				}
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := test.setup(ctx)
			var body bytes.Buffer
			marshaler := &jsonpb.Marshaler{}
			if err := marshaler.Marshal(&body, tc.request); err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			req, err := http.NewRequest("POST", api_v1.UserPathPrefix+"GetMe", &body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			ctx = context_data.WithAPIToken(ctx, tc.apiToken)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tc.wantStatusCode {
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Fatalf("want %d but got %d\n%s", tc.wantStatusCode, resp.StatusCode, string(b))
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				gotResponse := &api_v1.GetMeResponse{}
				unmarshaler := &jsonpb.Unmarshaler{}
				if err := unmarshaler.Unmarshal(resp.Body, gotResponse); err != nil {
					t.Fatal(err)
				}
				assertion.AssertEqual(t, tc.wantResponse, gotResponse, "")
			}
		})
	}
}

func newUserService(db *gorm.DB, appLogger *zap.Logger) api_v1.User {
	gaMeasurement := usecase.NewGAMeasurement(ga_measurement.NewGAMeasurementRepository(ga_measurement.NewFakeClient()))
	return interface_http.NewUserService(
		db, appLogger,
		gaMeasurement,
		di.NewNotificationTimeSpanUsecase(db.DB()),
		di.NewUserUsecase(db.DB()),
	)
}
