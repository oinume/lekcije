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
	api_v1 "github.com/oinume/lekcije/backend/proto_gen/proto/api/v1"
	"github.com/oinume/lekcije/backend/randoms"
	"github.com/oinume/lekcije/backend/usecase"
)

func Test_MeService_GetMe(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewMeServer(newMeService(db, appLogger))

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
						UserId:       int32(user.ID),
						Email:        user.Email,
						User:         interface_http.User2Proto(user),
						ShowTutorial: !user.IsFollowedTeacher(),
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
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			gotResponse := &api_v1.GetMeResponse{}
			statusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.MePathPrefix+"GetMe",
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

func Test_MeService_CreateFollowingTeacher(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewMeServer(newMeService(db, appLogger))

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken             string
		request              *api_v1.CreateFollowingTeacherRequest
		user                 *model2.User
		wantFollowingTeacher *model2.FollowingTeacher
		wantStatusCode       int
	}
	tests := map[string]struct {
		setup func(ctx context.Context) *testCase
	}{
		"ok_with_id": {
			setup: func(ctx context.Context) *testCase {
				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)
				userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
					uat.UserID = user.ID
				})
				repos.CreateUserAPITokens(ctx, t, userAPIToken)

				const teacherID = 46602
				ft := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
					ft.UserID = user.ID
					ft.TeacherID = teacherID
				})
				return &testCase{
					apiToken: userAPIToken.Token,
					request: &api_v1.CreateFollowingTeacherRequest{
						TeacherIdOrUrl: fmt.Sprint(teacherID),
					},
					user:                 user,
					wantFollowingTeacher: ft,
					wantStatusCode:       http.StatusOK,
				}
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			// t.Parallel()

			ctx := context.Background()
			tc := test.setup(ctx)

			client := twirptest.NewJSONClient()
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			ctx = context_data.WithGAMeasurementEvent(ctx, newGAMeasurementEvent())
			gotResponse := &api_v1.CreateFollowingTeacherResponse{}
			gotStatusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.MePathPrefix+"CreateFollowingTeacher",
				tc.request, gotResponse,
			)
			t.Logf("SendRequest err=%+v\n", err)
			if gotStatusCode == http.StatusOK {
				if err != nil {
					t.Fatalf("SendRequest failed: %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("error must not be nil")
				}
			}
			assertion.RequireEqual(t, tc.wantStatusCode, gotStatusCode, "unexpected status code")
			if gotStatusCode != http.StatusOK {
				return
			}

			gotFollowingTeachers, err := repos.FollowingTeacher().FindByUserID(ctx, tc.wantFollowingTeacher.UserID)
			if err != nil {
				t.Fatal(err)
			}
			assertion.RequireEqual(t, 1, len(gotFollowingTeachers), "unexpected gotFollowingTeachers length")
			assertion.AssertEqual(
				t, tc.wantFollowingTeacher, gotFollowingTeachers[0], "",
				cmpopts.IgnoreFields(model2.FollowingTeacher{}, "CreatedAt", "UpdatedAt"),
			)

			gotUser, err := repos.User().FindByEmail(ctx, tc.user.Email)
			if err != nil {
				t.Fatal(err)
			}
			if gotUser.FollowedTeacherAt.Time.IsZero() {
				t.Fatalf("user.FollowedTeacherAt must be updated but zero: %v", gotUser.FollowedTeacherAt)
			}
		})
	}
}

func Test_MeService_DeleteFollowingTeacher(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewMeServer(newMeService(db, appLogger))

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken       string
		request        *api_v1.DeleteFollowingTeachersRequest
		user           *model2.User
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

				fts := []*model2.FollowingTeacher{
					modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
						ft.UserID = user.ID
					}),
					modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
						ft.UserID = user.ID
					}),
				}
				teacherIDs := make([]uint32, len(fts))
				for i, ft := range fts {
					teacherIDs[i] = uint32(ft.TeacherID)
				}
				return &testCase{
					apiToken: userAPIToken.Token,
					request: &api_v1.DeleteFollowingTeachersRequest{
						TeacherIds: teacherIDs,
					},
					user:           user,
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
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			ctx = context_data.WithGAMeasurementEvent(ctx, newGAMeasurementEvent())
			gotResponse := &api_v1.DeleteFollowingTeachersResponse{}
			gotStatusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.MePathPrefix+"DeleteFollowingTeachers",
				tc.request, gotResponse,
			)
			assertion.RequireEqual(
				t, tc.wantStatusCode, gotStatusCode,
				fmt.Sprintf("unexpected status code:%v", err),
			)

			if gotStatusCode == http.StatusOK {
				got, err := repos.FollowingTeacher().FindByUserID(ctx, tc.user.ID)
				if err != nil {
					t.Fatalf("unexpected error:%v", err)
				}
				assertion.AssertEqual(t, 0, len(got), "unexpected FollowingTeacher length")
			} else {
				if err == nil {
					t.Fatal("error must not be nil")
				}
			}
		})
	}
}

func Test_MeService_ListFollowingTeachers(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewMeServer(newMeService(db, appLogger))

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken       string
		wantResponse   *api_v1.ListFollowingTeachersResponse
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

				teachers := make([]*model2.Teacher, 2)
				for i := range teachers {
					teachers[i] = modeltest.NewTeacher()
				}
				repos.CreateTeachers(ctx, t, teachers...)

				fts := make([]*model2.FollowingTeacher, len(teachers))
				for i, teacher := range teachers {
					fts[i] = modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
						ft.UserID = user.ID
						ft.TeacherID = teacher.ID
					})
				}
				repos.CreateFollowingTeachers(ctx, t, fts...)

				return &testCase{
					apiToken: userAPIToken.Token,
					wantResponse: &api_v1.ListFollowingTeachersResponse{
						Teachers: interface_http.TeachersProto(teachers),
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
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			gotResponse := &api_v1.ListFollowingTeachersResponse{}
			gotStatusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.MePathPrefix+"ListFollowingTeachers",
				&api_v1.ListFollowingTeachersRequest{}, gotResponse,
			)
			assertion.RequireEqual(t, tc.wantStatusCode, gotStatusCode, "unexpected status code")

			if gotStatusCode == http.StatusOK {
				sortOpt := cmpopts.SortSlices(func(i, j *api_v1.Teacher) bool {
					return i.Id < j.Id
				})
				assertion.AssertEqual(
					t, tc.wantResponse.Teachers, gotResponse.Teachers, "",
					sortOpt, cmpopts.IgnoreUnexported(api_v1.Teacher{}),
				)
			} else {
				if err == nil {
					t.Fatal("error must not be nil")
				}
			}
		})
	}
}

func Test_MeService_UpdateMeEmail(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewMeServer(newMeService(db, appLogger))

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken       string
		request        *api_v1.UpdateEmailRequest
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

				wantEmail := fmt.Sprintf("update-me-email-%s@example.com", randoms.MustNewString(8))
				return &testCase{
					apiToken: userAPIToken.Token,
					user:     user,
					request: &api_v1.UpdateEmailRequest{
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

				wantEmail := "invalid-email@"
				return &testCase{
					apiToken: userAPIToken.Token,
					user:     user,
					request: &api_v1.UpdateEmailRequest{
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
		"duplicate email": {
			setup: func(ctx context.Context) *testCase {
				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)
				userAPIToken := modeltest.NewUserAPIToken(func(uat *model2.UserAPIToken) {
					uat.UserID = user.ID
				})
				repos.CreateUserAPITokens(ctx, t, userAPIToken)

				return &testCase{
					apiToken: userAPIToken.Token,
					user:     user,
					request: &api_v1.UpdateEmailRequest{
						Email: user.Email,
					},
					wantStatusCode: http.StatusBadRequest,
					wantError: &twirptest.JSONError{
						Code: string(twirp.InvalidArgument),
						Msg:  "email email exists",
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
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			ctx = context_data.WithGAMeasurementEvent(ctx, newGAMeasurementEvent())
			gotResponse := &api_v1.UpdateEmailResponse{}
			statusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.MePathPrefix+"UpdateEmail",
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

func Test_MeService_UpdateMeNotificationTimeSpan(t *testing.T) {
	t.Parallel()

	helper := model.NewTestHelper()
	db := helper.DB(t)
	var log bytes.Buffer
	appLogger := logger.NewAppLogger(&log, logger.NewLevel("info"))
	handler := api_v1.NewMeServer(newMeService(db, appLogger))

	repos := mysqltest.NewRepositories(db.DB())
	type testCase struct {
		apiToken                  string
		request                   *api_v1.UpdateNotificationTimeSpanRequest
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
					request: &api_v1.UpdateNotificationTimeSpanRequest{
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
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			ctx = context_data.WithGAMeasurementEvent(ctx, newGAMeasurementEvent())
			gotResponse := &api_v1.UpdateNotificationTimeSpanResponse{}
			statusCode, err := client.SendRequest(
				ctx, t, handler, api_v1.MePathPrefix+"UpdateNotificationTimeSpan",
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

func newMeService(db *gorm.DB, appLogger *zap.Logger) api_v1.Me {
	gaMeasurement := usecase.NewGAMeasurement(ga_measurement.NewGAMeasurementRepository(ga_measurement.NewFakeClient()))
	return interface_http.NewMeService(
		db, appLogger,
		di.NewFollowingTeacherUsecase(appLogger, db.DB()),
		gaMeasurement,
		di.NewNotificationTimeSpanUsecase(db.DB()),
		di.NewUserUsecase(db.DB()),
	)
}

func newGAMeasurementEvent() *model2.GAMeasurementEvent {
	return &model2.GAMeasurementEvent{
		UserAgentOverride: "ua",
		ClientID:          "test",
		DocumentHostName:  "localhost",
		DocumentPath:      "",
		DocumentTitle:     "",
		DocumentReferrer:  "",
		IPOverride:        "",
	}
}
