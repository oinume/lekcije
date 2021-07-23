package http_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jinzhu/gorm"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/di"
	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	interface_http "github.com/oinume/lekcije/backend/interface/http"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/internal/twirptest"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	model2c "github.com/oinume/lekcije/backend/model2c"
	"github.com/oinume/lekcije/backend/usecase"
	"github.com/oinume/lekcije/backend/util"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

func Test_UserService_GetMe(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewUserServer(newUserService(db, appLogger))

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
				// TODO: Setup NotificationTimeSpan
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
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			tc := test.setup(ctx)

			client := twirptest.NewJSONClient()
			ctx = context_data.WithAPIToken(ctx, tc.apiToken)
			gotResponse := &api_v1.GetMeResponse{}
			statusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.UserPathPrefix+"GetMe",
				tc.request, gotResponse,
			)
			assertion.RequireEqual(t, tc.wantStatusCode, statusCode, "unexpected status code")
			if statusCode == http.StatusOK {
				assertion.AssertEqual(t, tc.wantResponse, gotResponse, "")
			} else {
				if err == nil {
					t.Fatal("err must not be nil")
				}
			}
		})
	}
}

func Test_UserService_UpdateMeEmail(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewUserServer(newUserService(db, appLogger))

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken       string
		request        *api_v1.UpdateMeEmailRequest
		user           *model2.User
		wantStatusCode int
		wantError      *twirptest.JSONError
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

				wantEmail := fmt.Sprintf("update-me-email-%s@example.com", util.RandomString(8))
				return &testCase{
					apiToken: userAPIToken.Token,
					user:     user,
					request: &api_v1.UpdateMeEmailRequest{
						Email: wantEmail,
					},
					wantStatusCode: http.StatusOK,
				}
			},
		},
		"invalid email": {
			setup: func(ctx context.Context) *testCase {
				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)
				userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
					uat.UserID = user.ID
				})
				repos.CreateUserAPITokens(ctx, t, userAPIToken)

				wantEmail := "invalid-email"
				return &testCase{
					apiToken: userAPIToken.Token,
					user:     user,
					request: &api_v1.UpdateMeEmailRequest{
						Email: wantEmail,
					},
					wantStatusCode: http.StatusBadRequest,
					wantError: &twirptest.JSONError{
						Code: string(twirp.InvalidArgument),
						Msg:  "email invalid email",
					},
				}
			},
		},
	}
	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			tc := test.setup(ctx)

			client := twirptest.NewJSONClient()
			ctx = context_data.WithAPIToken(ctx, tc.apiToken)
			ctx = context_data.WithGAMeasurementEvent(ctx, newGAMeasurementEvent())
			gotResponse := &api_v1.UpdateMeEmailResponse{}
			statusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.UserPathPrefix+"UpdateMeEmail",
				tc.request, gotResponse,
			)
			assertion.RequireEqual(t, tc.wantStatusCode, statusCode, "unexpected status code")
			if tc.wantStatusCode == http.StatusOK {
				gotUser, err := repos.User().FindByEmail(ctx, tc.request.Email)
				if err != nil {
					t.Fatalf("failed to find user by email: %v", err)
				}
				assertion.AssertEqual(t, tc.request.Email, gotUser.Email, "")
			} else {
				assertion.AssertEqual(
					t, tc.wantError, err, "",
					cmpopts.IgnoreFields(twirptest.JSONError{}, "Meta"),
				)
			}
		})
	}
}

func Test_UserService_UpdateMeNotificationTimeSpan(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewUserServer(newUserService(db, appLogger))

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken                  string
		request                   *api_v1.UpdateMeNotificationTimeSpanRequest
		user                      *model2.User
		wantStatusCode            int
		wantNotificationTimeSpans []*model2.NotificationTimeSpan
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
				timeSpans := make([]*model2.NotificationTimeSpan, 2)
				for i := range timeSpans {
					timeSpans[i] = modeltest.NewNotificationTimeSpan(func(nts *model2.NotificationTimeSpan) {
						nts.UserID = user.ID
						nts.Number = uint8(i) + 1
						nts.FromTime = fmt.Sprintf("%02d:00:00", nts.Number)
						nts.ToTime = fmt.Sprintf("%02d:30:00", nts.Number+1)
					})
				}
				repos.CreateNotificationTimeSpans(ctx, t, timeSpans...)

				timeSpansProto, err := interface_http.NotificationTimeSpansProto(timeSpans)
				if err != nil {
					t.Fatal(err)
				}
				return &testCase{
					apiToken: userAPIToken.Token,
					user:     user,
					request: &api_v1.UpdateMeNotificationTimeSpanRequest{
						NotificationTimeSpans: timeSpansProto,
					},
					wantStatusCode:            http.StatusOK,
					wantNotificationTimeSpans: timeSpans,
				}
			},
		},
	}
	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			tc := test.setup(ctx)

			client := twirptest.NewJSONClient()
			ctx = context_data.WithAPIToken(ctx, tc.apiToken)
			ctx = context_data.WithGAMeasurementEvent(ctx, newGAMeasurementEvent())
			gotResponse := &api_v1.UpdateMeNotificationTimeSpanResponse{}
			statusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.UserPathPrefix+"UpdateMeNotificationTimeSpan",
				tc.request, gotResponse,
			)
			assertion.RequireEqual(t, tc.wantStatusCode, statusCode, "unexpected status code")
			if statusCode == http.StatusOK {
				if err != nil {
					t.Fatal(err)
				}
				gotTimeSpans, err := repos.NotificationTimeSpan().FindByUserID(ctx, tc.user.ID)
				if err != nil {
					t.Fatalf("FindByUserID failed: %v", err)
				}
				assertion.AssertEqual(
					t, tc.wantNotificationTimeSpans, gotTimeSpans,
					"", cmpopts.EquateApproxTime(10*time.Second),
				)
			} else {
				if err == nil {
					t.Fatal("err must not be nil")
				}
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

func newGAMeasurementEvent() *model2c.GAMeasurementEvent {
	return &model2c.GAMeasurementEvent{
		UserAgentOverride: "ua",
		ClientID:          "test",
		DocumentHostName:  "localhost",
		DocumentPath:      "",
		DocumentTitle:     "",
		DocumentReferrer:  "",
		IPOverride:        "",
	}
}
