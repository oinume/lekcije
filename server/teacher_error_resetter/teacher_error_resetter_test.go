package teacher_error_resetter

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
)

var helper *model.TestHelper

func TestMain(m *testing.M) {
	bootstrap.CheckCLIEnvVars()
	helper = model.NewTestHelper()
	// NOTE: Avoid "Failed to FindByPK: id=1: record not found"
	//helper.TruncateAllTables(helper.DB())
	os.Exit(m.Run())
}

func TestMain_Run(t *testing.T) {
	teacherService := model.NewTeacherService(helper.DB())
	helper.CreateRandomTeacher()
	teacher := &model.Teacher{
		ID:                uint32(util.RandomInt(99999)),
		Name:              "test",
		CountryID:         1,
		Gender:            "male",
		Birthday:          time.Now().UTC(),
		YearsOfExperience: 1,
		FavoriteCount:     10,
		Rating:            4.8,
		FetchErrorCount:   10,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
	if err := teacherService.CreateOrUpdate(teacher); err != nil {
		t.Fatalf("teacherService.CreateOrUpdate failed: err=%v", err)
	}

	mockTransport, err := fetcher.NewMockTransport("../fetcher/testdata/5982.html") // TODO: path
	if err != nil {
		t.Fatalf("fetcher.NewMockTransport failed: err=%v", err)
	}
	httpClient := &http.Client{
		Transport: mockTransport,
	}
	concurrency := 1
	logLevel := "debug"
	dryRun := false
	main := &Main{
		Concurrency: &concurrency,
		DryRun:      &dryRun,
		LogLevel:    &logLevel,
		HTTPClient:  httpClient,
		DB:          helper.DB(),
	}
	if err := main.Run(); err != nil {
		t.Fatalf("main.Run failed: err=%v", err)
	}

	gotTeacher, err := teacherService.FindByPK(teacher.ID)
	if err != nil {
		t.Fatalf("teacherService.FindByPK failed: err=%v", err)
	}
	if got, want := gotTeacher.FetchErrorCount, uint8(0); got != want {
		t.Errorf("FetchErrorCount doesn't match: got=%v, want=%v", got, want)
	}
}
