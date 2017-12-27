package teacher_error_resetter

import (
	"net/http"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"go.uber.org/zap"
)

const fetchErrorCount = 6

type Main struct {
	Concurrency *int
	DryRun      *bool
	LogLevel    *string
}

func (m *Main) Run() error {
	bootstrap.CheckCLIEnvVars()
	startedAt := time.Now().UTC()
	logger.App.Info("notifier started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info("notifier finished", zap.Int("elapsed", int(elapsed)))
	}()

	dbLogging := *m.LogLevel == "debug"
	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL(), 1, dbLogging)
	if err != nil {
		return err
	}
	defer db.Close()

	mCountryService := model.NewMCountryService(db)
	mCountries, err := mCountryService.LoadAll()
	if err != nil {
		return err
	}

	teacherService := model.NewTeacherService(db)
	teachers, err := teacherService.FindByFetchErrorCountGt(fetchErrorCount)
	if err != nil {
		return err
	}

	fetcher := fetcher.NewTeacherLessonFetcher(http.DefaultClient, *m.Concurrency, false, mCountries, logger.App)
	defer fetcher.Close()
	for _, t := range teachers {
		if _, _, err := fetcher.Fetch(t.ID); err != nil {
			// TODO: logging
			continue
		}
		if err := teacherService.ResetFetchErrorCount(t.ID); err != nil {
			// TODO: logging
			continue
		}
	}

	return nil
}
