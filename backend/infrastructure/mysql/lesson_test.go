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
	"github.com/oinume/lekcije/backend/randoms"
)

func Test_lessonRepository_FindAllByTeacherIDAndDatetimeAsMap(t *testing.T) {
	repo := mysql.NewLessonRepository(helper.DB(t).DB())
	repos := mysqltest.NewRepositories(helper.DB(t).DB())

	type testCase struct {
		teacherID  uint
		lessonArgs []*model2.Lesson
	}
	tests := map[string]struct {
		setup   func(ctx context.Context) *testCase
		wantErr bool
	}{
		"normal": {
			setup: func(ctx context.Context) *testCase {
				//boil.DebugMode = true
				//boil.DebugWriter = os.Stdout
				teacher := modeltest.NewTeacher()
				repos.CreateTeachers(ctx, t, teacher)

				l1 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacher.ID
					l.Datetime = time.Date(2022, 11, 1, 10, 0, 0, 0, time.UTC)
				})
				l2 := modeltest.NewLesson(func(l *model2.Lesson) {
					l.TeacherID = teacher.ID
					l.Datetime = time.Date(2022, 11, 1, 10, 30, 0, 0, time.UTC)
				})
				lessons := []*model2.Lesson{l1, l2}
				repos.CreateLessons(ctx, t, lessons...)

				return &testCase{
					teacherID:  teacher.ID,
					lessonArgs: lessons,
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			gotLessonsMap, err := repo.FindAllByTeacherIDAndDatetimeAsMap(ctx, tc.teacherID, tc.lessonArgs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAllByTeacherIDAndDatetimeAsMap() error = %v, wantErr %v", err, tt.wantErr)
			}

			assertion.AssertEqual(t, len(tc.lessonArgs), len(gotLessonsMap), "lesson length doesn't match")
			for _, l := range tc.lessonArgs {
				datetime := model2.LessonDatetime(l.Datetime).String()
				_, ok := gotLessonsMap[datetime]
				if !ok {
					t.Errorf("key %q must exist in result map", datetime)
				}
			}
		})
	}
}

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
				teacherID := uint(randoms.MustNewInt64(10000000))
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
				teacherID := uint(randoms.MustNewInt64(10000000))
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

func Test_lessonRepository_GetNewAvailableLessons(t *testing.T) {
	repo := mysql.NewLessonRepository(helper.DB(t).DB())

	type testCase struct {
		oldLessons []*model2.Lesson
		newLessons []*model2.Lesson
		want       []*model2.Lesson
	}
	tests := map[string]struct {
		setup func(ctx context.Context) *testCase
	}{
		"one available lessons": {
			setup: func(ctx context.Context) *testCase {
				teacherID := uint(randoms.MustNewInt64(100000))
				datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalLocation())
				lessons1 := newLessons(teacherID, datetime, "Reserved", 3)
				lessons2 := newLessons(teacherID, datetime, "Reserved", 3)
				lessons2[1].Status = "Available"
				return &testCase{
					oldLessons: lessons1,
					newLessons: lessons2,
					want: []*model2.Lesson{
						lessons2[1],
					},
				}
			},
		},
		"no available lessons": {
			setup: func(ctx context.Context) *testCase {
				teacherID := uint(randoms.MustNewInt64(100000))
				datetime := time.Date(2016, 10, 1, 14, 30, 0, 0, config.LocalLocation())
				lessons1 := newLessons(teacherID, datetime, "Reserved", 3)
				lessons2 := newLessons(teacherID, datetime, "Reserved", 3)
				// There are available lessons in old, means no difference between lessons1 and lessons2
				lessons1[0].Status = "Available"
				lessons2[0].Status = "Available"
				return &testCase{
					oldLessons: lessons1,
					newLessons: lessons2,
					want:       []*model2.Lesson{},
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			got := repo.GetNewAvailableLessons(ctx, tc.oldLessons, tc.newLessons)
			assertion.AssertEqual(t, len(tc.want), len(got), "length of lessons doesn't match")
			assertion.AssertEqual(t, tc.want, got, "")
		})
	}
}

func newLessons(teacherID uint, baseDatetime time.Time, status string, length int) []*model2.Lesson {
	lessons := make([]*model2.Lesson, length)
	now := time.Now().UTC()
	for i := range lessons {
		lessons[i] = &model2.Lesson{
			TeacherID: teacherID,
			Datetime:  baseDatetime.Add(time.Duration(i) * time.Hour),
			Status:    status,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return lessons
}

func Test_lessonRepository_FindOrCreate(t *testing.T) {
	db := helper.DB(t)
	repo := mysql.NewLessonRepository(db.DB())
	repos := mysqltest.NewRepositories(db.DB())

	type testCase struct {
		lesson *model2.Lesson
	}
	tests := map[string]struct {
		setup func(ctx context.Context) *testCase
	}{
		"create new one": {
			setup: func(ctx context.Context) *testCase {
				l := modeltest.NewLesson()
				return &testCase{
					lesson: l,
				}
			},
		},
		"exist": {
			setup: func(ctx context.Context) *testCase {
				l := modeltest.NewLesson()
				repos.CreateLessons(ctx, t, l)
				return &testCase{
					lesson: l,
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			got, err := repo.FindOrCreate(ctx, tc.lesson, true)
			if err != nil {
				t.Fatalf("FindOrCreate failed: unexpected error = %v", err)
			}
			assertion.AssertEqual(t, tc.lesson, got, "", cmpopts.EquateApproxTime(10*time.Second))
		})
	}
}

func Test_lessonRepository_UpdateStatus(t *testing.T) {
	db := helper.DB(t)
	repo := mysql.NewLessonRepository(db.DB())
	repos := mysqltest.NewRepositories(db.DB())

	type testCase struct {
		lesson    *model2.Lesson
		newStatus string
	}
	tests := map[string]struct {
		setup func(ctx context.Context) *testCase
	}{
		"normal": {
			setup: func(ctx context.Context) *testCase {
				l := modeltest.NewLesson(func(l *model2.Lesson) {
					l.Status = "available"
				})
				repos.CreateLessons(ctx, t, l)
				return &testCase{
					lesson:    l,
					newStatus: "reserved",
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tc := tt.setup(ctx)
			err := repo.UpdateStatus(ctx, tc.lesson.ID, tc.newStatus)
			if err != nil {
				t.Fatalf("UpdateStatus failed: unexpected error = %v", err)
			}
			got, err := repo.FindByID(ctx, tc.lesson.ID)
			if err != nil {
				t.Fatalf("FindByID failed: unexpected error = %v", err)
			}
			assertion.AssertEqual(t, tc.newStatus, got.Status, "status is not updated")
		})
	}

}
