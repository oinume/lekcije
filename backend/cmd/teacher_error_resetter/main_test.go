package main

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/fetcher"
	"github.com/oinume/lekcije/backend/model"
)

var helper *model.TestHelper

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	helper = model.NewTestHelper()
	// NOTE: Avoid "Failed to FindByPK: id=1: record not found"
	helper.TruncateAllTables(nil)
	os.Exit(m.Run())
}

func Test_teacherErrorResetterMain_run(t *testing.T) {
	type fields struct {
		outStream io.Writer
		errStream io.Writer
	}

	tests := map[string]struct {
		args    []string
		fields  fields
		wantErr bool
	}{
		"normal": {
			args: []string{"teacher_error_resetter", "-concurrency=1"},
			fields: fields{
				outStream: os.Stdout,
				errStream: os.Stderr,
			},
			wantErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &teacherErrorResetterMain{
				outStream: test.fields.outStream,
				errStream: test.fields.errStream,
				db:        helper.DB(t),
			}

			teacherService := model.NewTeacherService(helper.DB(t))
			teacher := helper.CreateRandomTeacher(t)
			teacher.YearsOfExperience = 1
			teacher.FavoriteCount = 10
			teacher.Rating = 4.8
			teacher.FetchErrorCount = 10
			if err := teacherService.CreateOrUpdate(teacher); err != nil {
				t.Fatalf("teacherService.CreateOrUpdate failed: err=%v", err)
			}

			mockTransport, err := fetcher.NewMockTransport("../../fetcher/testdata/3986.html") // TODO: path
			if err != nil {
				t.Fatalf("fetcher.NewMockTransport failed: err=%v", err)
			}
			m.httpClient = &http.Client{
				Transport: mockTransport,
			}

			if err := m.run(test.args); (err != nil) != test.wantErr {
				t.Fatalf("teacherErrorResetterMain.run() error = %v, wantErr %v", err, test.wantErr)
			}

			gotTeacher, err := teacherService.FindByPK(teacher.ID)
			if err != nil {
				t.Fatalf("teacherService.FindByPK failed: err=%v", err)
			}
			if got, want := gotTeacher.FetchErrorCount, uint8(0); got != want {
				t.Errorf("FetchErrorCount doesn't match: got=%v, want=%v", got, want)
			}
		})
	}
}
