package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/internal/modeltest"
	"github.com/oinume/lekcije/backend/internal/mysqltest"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
)

func TestMain(m *testing.M) {

}

func Test_lessonRepository_FindAllByTeacherIDsDatetimeBetween(t *testing.T) {
	helper := model.NewTestHelper()
	repo := mysql.NewLessonRepository(helper.DB(t).DB())
	repos := mysqltest.NewRepositories(helper.DB(t).DB())

	type args struct {
		teacherID uint
		fromDate  time.Time
		toDate    time.Time
	}
	modeltest.NewFollowingTeacher()
	tests := map[string]struct {
		setup   func()
		args    args
		want    []*model2.Lesson
		wantErr bool
	}{
		"normal": {
			setup: func() {
				// TODO: setup fixed date fixture
				lessons := []*model2.Lesson{
					modeltest.NewLesson(),
					modeltest.NewLesson(),
				}
				repos.CreateLessons(context.Background(), t, lessons...)
			},
			args: args{},
		},
		"no records": {},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			tt.setup()
			got, err := repo.FindAllByTeacherIDsDatetimeBetween(ctx, tt.args.teacherID, tt.args.fromDate, tt.args.toDate)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAllByTeacherIDsDatetimeBetween() error = %v, wantErr %v", err, tt.wantErr)
			}
			assertion.AssertEqual(t, tt.want, got, "")
		})
	}
}
