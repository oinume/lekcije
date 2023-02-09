package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/rollbar/rollbar-go"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/cli"
	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/infrastructure/dmm_eikaiwa"
	"github.com/oinume/lekcije/backend/infrastructure/mysql"
	"github.com/oinume/lekcije/backend/infrastructure/open_telemetry"
	irollbar "github.com/oinume/lekcije/backend/infrastructure/rollbar"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/registry"
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
		concurrency          = flagSet.Int("concurrency", 1, "Concurrency of fetcher")
		dryRun               = flagSet.Bool("dry-run", false, "Don't update database with fetched lessons")
		fetcherCache         = flagSet.Bool("fetcher-cache", false, "Cache teacher and lesson data in Fetcher")
		notificationInterval = flagSet.Int("notification-interval", 0, "Notification interval")
		sendEmail            = flagSet.Bool("send-email", true, "Flag to send email")
		logLevel             = flagSet.String("log-level", "info", "Log level")
	)

	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	config.MustProcessDefault()
	startedAt := time.Now().UTC()
	appLogger := logger.NewAppLogger(os.Stderr, logger.NewLevel(*logLevel))
	appLogger.Info(fmt.Sprintf("notifier started (interval=%d)", *notificationInterval))

	const serviceName = "notifier"
	tracerProvider, err := open_telemetry.NewTracerProvider(serviceName, config.DefaultVars)
	if err != nil {
		log.Fatalf("NewTraceProvider failed: %v", err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tracerProvider)

	ctx, span := otel.Tracer(config.DefaultTracerName).Start(context.Background(), "main")
	defer span.End()
	db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	if *notificationInterval == 0 {
		return fmt.Errorf("-notification-interval is required")
	}
	users, err := mysql.NewUserRepository(db.DB()).FindAllByEmailVerifiedIsTrue(ctx, *notificationInterval)
	if err != nil {
		return err
	}
	mCountryList := registry.MustNewMCountryList(ctx, db.DB())
	lessonFetcher := dmm_eikaiwa.NewLessonFetcher(nil, *concurrency, *fetcherCache, mCountryList, appLogger)

	var sender emailer.Sender
	if *sendEmail {
		sender = emailer.NewSendGridSender(nil, appLogger)
	} else {
		sender = &emailer.NoSender{}
	}

	rollbarClient := rollbar.New(
		config.DefaultVars.RollbarAccessToken,
		config.DefaultVars.ServiceEnv,
		config.DefaultVars.VersionHash,
		"", "/",
	)
	errorRecorder := usecase.NewErrorRecorder(
		appLogger,
		irollbar.NewErrorRecorderRepository(rollbarClient),
	)
	lessonUsecase := registry.NewLessonUsecase(db.DB())
	notificationTimeSpanUsecase := registry.NewNotificationTimeSpanUsecase(db.DB())
	statNotifierUsecase := registry.NewStatNotifierUsecase(db.DB())
	teacherUsecase := registry.NewTeacherUsecase(db.DB())
	notifier := usecase.NewNotifier(
		appLogger, db, errorRecorder, lessonFetcher, *dryRun, lessonUsecase, notificationTimeSpanUsecase,
		statNotifierUsecase, teacherUsecase, sender, mysql.NewFollowingTeacherRepository(db.DB()),
	)
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

	elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
	appLogger.Info(
		fmt.Sprintf("notifier finished (interval=%d)", *notificationInterval),
		zap.Int("elapsed", int(elapsed)),
		zap.Int("usersCount", len(users)),
	)

	return nil
}
