package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rollbar/rollbar-go"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/cli"
	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/fetcher"
	irollbar "github.com/oinume/lekcije/backend/infrastructure/rollbar"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/notifier"
	"github.com/oinume/lekcije/backend/open_census"
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
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		appLogger.Info(
			fmt.Sprintf("notifier finished (interval=%d)", *notificationInterval),
			zap.Int("elapsed", int(elapsed)),
		)
	}()

	const serviceName = "notifier"
	exporter, flush, err := open_census.NewExporter(
		config.DefaultVars,
		serviceName,
		*notificationInterval == 10 || !config.IsProductionEnv(),
	)
	if err != nil {
		return fmt.Errorf("NewExporter failed: %v", err)
	}
	defer flush()
	trace.RegisterExporter(exporter)

	ctx, span := trace.StartSpan(context.Background(), "main")
	defer span.End()
	db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	if *notificationInterval == 0 {
		return fmt.Errorf("-notification-interval is required")
	}
	users, err := model.NewUserService(db).FindAllEmailVerifiedIsTrue(ctx, *notificationInterval)
	if err != nil {
		return err
	}
	mCountries, err := model.NewMCountryService(db).LoadAll(ctx)
	if err != nil {
		return err
	}
	lessonFetcher := fetcher.NewLessonFetcher(nil, *concurrency, *fetcherCache, mCountries, appLogger)

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
	n := notifier.NewNotifier(appLogger, db, errorRecorder, lessonFetcher, *dryRun, sender, nil)
	defer n.Close(ctx, &model.StatNotifier{
		Datetime:             startedAt,
		Interval:             uint8(*notificationInterval),
		Elapsed:              0,
		UserCount:            uint32(len(users)),
		FollowedTeacherCount: 0,
	})
	for _, user := range users {
		if err := n.SendNotification(ctx, user); err != nil {
			return err
		}
	}

	return nil
}
