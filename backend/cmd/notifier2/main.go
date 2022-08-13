package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/cli"
	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/infrastructure/dmm_eikaiwa"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

func main() {
	m := &notifierMain{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type notifierMain struct {
	outStream io.Writer
	errStream io.Writer
}

func (m *notifierMain) run(args []string) error {
	flagSet := flag.NewFlagSet("notifier", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
		concurrency = flagSet.Int("concurrency", 1, "Concurrency of fetcher")
		//		dryRun               = flagSet.Bool("dry-run", false, "Don't update database with fetched lessons")
		fetcherCache         = flagSet.Bool("fetcher-cache", false, "Cache teacher and lesson data in Fetcher")
		notificationInterval = flagSet.Int("notification-interval", 0, "Notification interval")
		//		sendEmail            = flagSet.Bool("send-email", true, "Flag to send email")
		logLevel = flagSet.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}
	if *notificationInterval == 0 {
		return fmt.Errorf("-notification-interval is required")
	}

	ctx := context.Background()
	config.MustProcessDefault()

	startedAt := time.Now().UTC()
	appLogger := logger.NewAppLogger(os.Stderr, logger.NewLevel(*logLevel))
	appLogger.Info(fmt.Sprintf("notifier started (interval=%d)", *notificationInterval))
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		appLogger.Info(
			fmt.Sprintf("notifier finished (interval=%d)", *notificationInterval),
			zap.Int("elapsed", int(elapsed)),
		)
	}()

	gormDB, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err != nil {
		return err
	}
	defer func() { _ = gormDB.Close() }()

	dbRepo := mysql.NewDBRepository(gormDB.DB())
	//	lessonRepo := mysql.NewLessonRepository(gormDB.DB())
	userRepo := mysql.NewUserRepository(gormDB.DB())
	userGoogleRepo := mysql.NewUserGoogleRepository(gormDB.DB())
	userUsecase := usecase.NewUser(dbRepo, userRepo, userGoogleRepo)

	users, err := userUsecase.FindAllByEmailVerified(ctx, *notificationInterval)
	if err != nil {
		return err
	}
	mCountries, err := mysql.NewMCountryRepository(gormDB.DB()).FindAll(ctx)
	if err != nil {
		return err
	}
	mCountryList := model2.NewMCountryList(mCountries)

	lessonFetcher := dmm_eikaiwa.NewLessonFetcher(nil, *concurrency, *fetcherCache, mCountryList, appLogger)
	notificationUsecase := usecase.NewNotification()
	notifier := notificationUsecase.NewLessonNotifier(lessonFetcher)
	defer notifier.Close(ctx, &model2.StatNotifier{
		Datetime:             startedAt,
		Interval:             uint8(*notificationInterval),
		Elapsed:              0,
		UserCount:            uint(len(users)),
		FollowedTeacherCount: 0,
	})
	for _, user := range users {
		if err := notifier.SendNotification(ctx, user); err != nil {
			return err
		}
	}

	// lessonFetcher := fetcher.NewLessonFetcher(nil, *concurrency, *fetcherCache, mCountries, appLogger)

	return nil
}
