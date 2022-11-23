package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/internal/slice_util"
	"github.com/oinume/lekcije/backend/model2"
)

func Test_followingTeacherRepository_FindTeacherIDsByUserID(t *testing.T) {
	repo := mysql.NewFollowingTeacherRepository(helper.DB(t).DB())
	repos := mysqltest.NewRepositories(helper.DB(t).DB())

	type testCase struct {
		userID         uint
		lastLessonAt   time.Time
		wantTeacherIDs []uint
	}

	tests := map[string]struct {
		setup   func(ctx context.Context) *testCase
		wantErr bool
	}{
		"normal": {
			setup: func(ctx context.Context) *testCase {
				helper.TruncateAllTables(t)
				//boil.DebugMode = true
				//boil.DebugWriter = os.Stdout

				user := modeltest.NewUser()
				repos.CreateUsers(ctx, t, user)

				now := time.Now().UTC()
				teacher1 := modeltest.NewTeacher(func(t *model2.Teacher) {
					t.LastLessonAt = now.Add(1 * time.Hour)
				})
				ft1 := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
					ft.UserID = user.ID
					ft.TeacherID = teacher1.ID
				})
				teacher2 := modeltest.NewTeacher(func(t *model2.Teacher) {
					t.LastLessonAt = now.Add(24 * time.Hour)
				})
				ft2 := modeltest.NewFollowingTeacher(func(ft *model2.FollowingTeacher) {
					ft.UserID = user.ID
					ft.TeacherID = teacher2.ID
				})
				repos.CreateTeachers(ctx, t, teacher1, teacher2)
				repos.CreateFollowingTeachers(ctx, t, ft1, ft2)

				teacherIDs := []uint{teacher1.ID, teacher2.ID}
				slice_util.Sort(teacherIDs)
				return &testCase{
					userID:         user.ID,
					lastLessonAt:   now,
					wantTeacherIDs: teacherIDs,
				}
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			gotTeacherIDs, err := repo.FindTeacherIDsByUserID(ctx, tc.userID, 5, tc.lastLessonAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindTeacherIDsByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			slice_util.Sort(gotTeacherIDs)
			assertion.AssertEqual(t, tc.wantTeacherIDs, gotTeacherIDs, "")
		})
	}
}
