package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/cli"
	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/infrastructure/dmm_eikaiwa"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/registry"
)

func main() {
	m := &teacherErrorResetterMain{
		outStream:  os.Stdout,
		errStream:  os.Stderr,
		db:         nil,
		httpClient: nil,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type teacherErrorResetterMain struct {
	outStream  io.Writer
	errStream  io.Writer
	db         *gorm.DB
	httpClient *http.Client
}

const fetchErrorCount = 5

func (m *teacherErrorResetterMain) run(args []string) error {
	flagSet := flag.NewFlagSet("teacher_error_resetter", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
		concurrency = flagSet.Int("concurrency", 1, "Concurrency of lessonFetcher")
		dryRun      = flagSet.Bool("dry-run", false, "Don't update database with fetched lessons")
		logLevel    = flag.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	config.MustProcessDefault()
	if m.db == nil {
		db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
		if err != nil {
			cli.WriteError(os.Stderr, err)
			os.Exit(1)
		}
		defer func() { _ = db.Close() }()
		m.db = db
	}

	startedAt := time.Now().UTC()
	appLogger := logger.NewAppLogger(os.Stderr, logger.NewLevel(*logLevel))
	appLogger.Info("teacher_error_resetter started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		appLogger.Info("teacher_error_resetter finished", zap.Int("elapsed", int(elapsed)))
	}()

	ctx := context.Background()
	teacherService := model.NewTeacherService(m.db)
	teachers, err := teacherService.FindByFetchErrorCountGt(fetchErrorCount)
	if err != nil {
		return err
	}

	mCountryList := registry.MustNewMCountryList(ctx, m.db.DB())
	lessonFetcher := dmm_eikaiwa.NewLessonFetcher(m.httpClient, *concurrency, false, mCountryList, appLogger)
	defer lessonFetcher.Close()
	for _, t := range teachers {
		if _, _, err := lessonFetcher.Fetch(ctx, uint(t.ID)); err != nil {
			appLogger.Error("lessonFetcher.Fetch failed", zap.Uint32("id", t.ID), zap.Error(err))
			continue
		}
		if *dryRun {
			appLogger.Info("Skip teacher because of dry-run", zap.Uint32("id", t.ID))
			continue
		}
		if err := teacherService.ResetFetchErrorCount(t.ID); err != nil {
			appLogger.Error("teacherService.ResetFetchErrorCount failed", zap.Uint32("id", t.ID), zap.Error(err))
			continue
		}
	}

	return nil
}
