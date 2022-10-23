package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

func Test_Lesson_UpdateLessons(t *testing.T) {
	repos := mysqltest.NewRepositories(helper.DB(t).DB())
	uc := usecase.NewLesson(repos.Lesson(), repos.LessonStatusLog())

	type testCase struct {
		lessons []*model2.Lesson
		want    int
	}
	tests := map[string]struct {
		setup   func(ctx context.Context) *testCase
		wantErr bool
	}{
		"create ok": {
			setup: func(ctx context.Context) *testCase {
				teacher := modeltest.NewTeacher()
				repos.CreateTeachers(ctx, t, teacher)
				l1 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacher.ID
					l.Datetime = time.Date(2022, 10, 23, 10, 0, 0, 0, time.UTC)
				})
				l2 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacher.ID
					l.Datetime = time.Date(2022, 10, 23, 10, 30, 0, 0, time.UTC)
				})
				lessons := []*model2.Lesson{l1, l2}
				return &testCase{
					lessons: lessons,
					want:    len(lessons),
				}
			},
		},
		"update ok": {
			setup: func(ctx context.Context) *testCase {
				teacher := modeltest.NewTeacher()
				repos.CreateTeachers(ctx, t, teacher)
				l1 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacher.ID
					l.Datetime = time.Date(2022, 10, 23, 10, 0, 0, 0, time.UTC)
				})
				l2 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacher.ID
					l.Datetime = time.Date(2022, 10, 23, 10, 30, 0, 0, time.UTC)
				})
				repos.CreateLessons(ctx, t, l1, l2)

				l1.Status = "finished"
				l2.Status = "reserved" // TODO: enum
				return &testCase{
					lessons: []*model2.Lesson{l1, l2},
					want:    2,
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			got, err := uc.UpdateLessons(ctx, tc.lessons)
			if err != nil {
				t.Fatal(err)
			}
			assertion.AssertEqual(t, tc.want, got, "")
		})
	}
}
