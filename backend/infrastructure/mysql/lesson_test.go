package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model2"
)

func Test_lessonRepository_FindAllByTeacherIDsDatetimeBetween(t *testing.T) {
	repo := mysql.NewLessonRepository(helper.DB(t).DB())
	repos := mysqltest.NewRepositories(helper.DB(t).DB())

	type testCase struct {
		teacherID uint
		fromDate  time.Time
		toDate    time.Time
		want      []*model2.Lesson
	}
	modeltest.NewFollowingTeacher()
	tests := map[string]struct {
		setup   func(ctx context.Context) *testCase
		wantErr bool
	}{
		"normal": {
			setup: func(ctx context.Context) *testCase {
				const teacherID = 1
				l1 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacherID
					l.Datetime = time.Date(2022, 11, 1, 10, 0, 0, 0, config.DefaultVars.LocalLocation)
				})
				l2 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacherID
					l.Datetime = time.Date(2022, 11, 2, 10, 30, 0, 0, config.DefaultVars.LocalLocation)
				})
				l3 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacherID
					l.Datetime = time.Date(2022, 11, 5, 10, 30, 0, 0, config.DefaultVars.LocalLocation)
				})
				lessons := []*model2.Lesson{l1, l2, l3}
				repos.CreateLessons(ctx, t, lessons...)

				return &testCase{
					teacherID: teacherID,
					fromDate:  time.Date(2022, 11, 1, 10, 0, 0, 0, config.DefaultVars.LocalLocation),
					toDate:    time.Date(2022, 11, 2, 10, 30, 0, 0, config.DefaultVars.LocalLocation),
					want:      []*model2.Lesson{l1, l2}, // l3 doesn't hit
				}
			},
		},
		"no records": {
			setup: func(ctx context.Context) *testCase {
				const teacherID = 2
				l1 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacherID
					l.Datetime = time.Date(2022, 11, 1, 10, 0, 0, 0, config.DefaultVars.LocalLocation)
				})
				lessons := []*model2.Lesson{l1}
				repos.CreateLessons(ctx, t, lessons...)

				return &testCase{
					teacherID: teacherID + 1,
					fromDate:  time.Date(2022, 12, 1, 10, 0, 0, 0, config.DefaultVars.LocalLocation),
					toDate:    time.Date(2022, 12, 2, 10, 30, 0, 0, config.DefaultVars.LocalLocation),
					want:      nil,
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			got, err := repo.FindAllByTeacherIDsDatetimeBetween(ctx, tc.teacherID, tc.fromDate, tc.toDate)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAllByTeacherIDsDatetimeBetween() error = %v, wantErr %v", err, tt.wantErr)
			}
			assertion.AssertEqual(t, len(tc.want), len(got), "length of want and got is not same", cmpopts.EquateApproxTime(10*time.Second))
			assertion.AssertEqual(t, tc.want, got, "", cmpopts.EquateApproxTime(10*time.Second))
		})
	}
}
