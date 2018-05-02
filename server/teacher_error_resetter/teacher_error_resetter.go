package teacher_error_resetter

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"go.uber.org/zap"
)

const fetchErrorCount = 5

type Main struct {
	Concurrency *int
	DryRun      *bool
	LogLevel    *string
	HTTPClient  *http.Client
	DB          *gorm.DB
}

func (m *Main) Run() error {
	config.MustProcessDefault()
	startedAt := time.Now().UTC()
	logger.App.Info("teacher_error_resetter started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info("teacher_error_resetter finished", zap.Int("elapsed", int(elapsed)))
	}()

	dbLogging := *m.LogLevel == "debug"
	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL(), 1, dbLogging)
	if err != nil {
		return err
	}
	defer db.Close()

	mCountryService := model.NewMCountryService(m.DB)
	mCountries, err := mCountryService.LoadAll()
	if err != nil {
		return err
	}

	teacherService := model.NewTeacherService(m.DB)
	teachers, err := teacherService.FindByFetchErrorCountGt(fetchErrorCount)
	if err != nil {
		return err
	}
	fetcher := fetcher.NewLessonFetcher(m.HTTPClient, *m.Concurrency, false, mCountries, logger.App)
	defer fetcher.Close()
	for _, t := range teachers {
		if _, _, err := fetcher.Fetch(t.ID); err != nil {
			logger.App.Error("fetcher.Fetch failed", zap.Uint32("id", t.ID), zap.Error(err))
			continue
		}
		if *m.DryRun {
			logger.App.Info("Skip teacher because of dry-run", zap.Uint32("id", t.ID))
			continue
		}
		if err := teacherService.ResetFetchErrorCount(t.ID); err != nil {
			logger.App.Error("teacherService.ResetFetchErrorCount failed", zap.Uint32("id", t.ID), zap.Error(err))
			continue
		}
	}

	return nil
}
