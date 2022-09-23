package resolver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/morikuni/failure"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/domain/repository"
	graphqlmodel "github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

func TestCreateFollowingTeacher(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	followingTeacherUsecase := newFollowingTeacherUsecase(
		zap.NewNop(),
		repos,
		&repository.LessonFetcherMock{
			CloseFunc: func() {},
			FetchFunc: func(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error) {
				return &model2.Teacher{
					ID:   teacherID,
					Name: "mock teacher",
				}, nil, nil
			},
		},
	)
	notificationTimeSpanUsecase := usecase.NewNotificationTimeSpan(repos.NotificationTimeSpan())
	userUsecase := usecase.NewUser(repos.DB(), repos.User(), repos.UserGoogle())
	resolver := NewResolver(
		repos.FollowingTeacher(),
		followingTeacherUsecase,
		repos.NotificationTimeSpan(),
		notificationTimeSpanUsecase,
		repos.Teacher(),
		repos.User(),
		userUsecase,
	)

	type testCase struct {
		apiToken      string
		input         graphqlmodel.CreateFollowingTeacherInput
		wantResult    *graphqlmodel.CreateFollowingTeacherPayload
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
				const teacherID = "12345"
				return &testCase{
					apiToken: userAPIToken.Token,
					input: graphqlmodel.CreateFollowingTeacherInput{
						TeacherIDOrURL: "12345",
					},
					wantResult: &graphqlmodel.CreateFollowingTeacherPayload{
						ID: fmt.Sprintf("%v-%v", user.ID, teacherID),
					},
				}
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := test.setup(ctx)
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			gotResult, err := resolver.Mutation().CreateFollowingTeacher(ctx, tc.input)
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

			assertion.AssertEqual(t, tc.wantResult, gotResult, "")
		})
	}
}

func TestDeleteFollowingTeachers(t *testing.T) {
	helper := model.NewTestHelper()
	db := helper.DB(t)
	repos := mysqltest.NewRepositories(db.DB())
	followingTeacherUsecase := newFollowingTeacherUsecase(
		zap.NewNop(),
		repos,
		&repository.LessonFetcherMock{
			CloseFunc: func() {},
			FetchFunc: func(ctx context.Context, teacherID uint) (*model2.Teacher, []*model2.Lesson, error) {
				return &model2.Teacher{
					ID:   teacherID,
					Name: "mock teacher",
				}, nil, nil
			},
		},
	)
	notificationTimeSpanUsecase := usecase.NewNotificationTimeSpan(repos.NotificationTimeSpan())
	userUsecase := usecase.NewUser(repos.DB(), repos.User(), repos.UserGoogle())
	resolver := NewResolver(
		repos.FollowingTeacher(),
		followingTeacherUsecase,
		repos.NotificationTimeSpan(),
		notificationTimeSpanUsecase,
		repos.Teacher(),
		repos.User(),
		userUsecase,
	)

	type testCase struct {
		apiToken      string
		input         graphqlmodel.DeleteFollowingTeachersInput
		wantResult    *graphqlmodel.DeleteFollowingTeachersPayload
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
				fts := []*model2.FollowingTeacher{
					modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
						ft.UserID = user.ID
					}),
					modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
						ft.UserID = user.ID
					}),
				}
				repos.CreateFollowingTeachers(ctx, t, fts...)
				teacherIDs := make([]string, len(fts))
				for i, ft := range fts {
					teacherIDs[i] = fmt.Sprint(ft.TeacherID)
				}
				return &testCase{
					apiToken: userAPIToken.Token,
					input: graphqlmodel.DeleteFollowingTeachersInput{
						TeacherIds: teacherIDs,
					},
					wantResult: &graphqlmodel.DeleteFollowingTeachersPayload{
						TeacherIds: teacherIDs,
					},
				}
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := test.setup(ctx)
			ctx = context_data.SetAPIToken(ctx, tc.apiToken)
			gotResult, err := resolver.Mutation().DeleteFollowingTeachers(ctx, tc.input)
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

			assertion.AssertEqual(t, tc.wantResult, gotResult, "")
		})
	}
}

func newFollowingTeacherUsecase(appLogger *zap.Logger, repos *mysqltest.Repositories, lessonFetcher repository.LessonFetcher) *usecase.FollowingTeacher {
	return usecase.NewFollowingTeacher(
		appLogger,
		repos.DB(),
		repos.FollowingTeacher(),
		repos.User(),
		repos.Teacher(),
		lessonFetcher,
	)
}
