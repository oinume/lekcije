package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/oinume/lekcije/server/cli"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/emailer"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/notifier"
	"github.com/oinume/lekcije/server/open_census"
	"github.com/oinume/lekcije/server/stopwatch"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
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
		//logLevel             = flagSet.String("log-level", "info", "Log level")
	)

	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	sw := stopwatch.NewSync()
	sw.Start()

	config.MustProcessDefault()
	startedAt := time.Now().UTC()
	//if *logLevel != "" {
	//	logger.App.SetLevel(logger.NewLevel(*logLevel))
	//}
	logger.App.Info(fmt.Sprintf("notifier started (interval=%d)", *notificationInterval))
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info(
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
	lessonFetcher := fetcher.NewLessonFetcher(nil, *concurrency, *fetcherCache, mCountries, logger.App)

	var sender emailer.Sender
	if *sendEmail {
		sender = emailer.NewSendGridSender(nil)
	} else {
		sender = &emailer.NoSender{}
	}

	statNotifier := &model.StatNotifier{
		Datetime:             startedAt,
		Interval:             uint8(*notificationInterval),
		Elapsed:              0,
		UserCount:            uint32(len(users)),
		FollowedTeacherCount: 0,
	}
	n := notifier.NewNotifier(db, lessonFetcher, *dryRun, sender, sw, nil)
	defer n.Close(ctx, statNotifier)
	for _, user := range users {
		if err := n.SendNotification(ctx, user); err != nil {
			return err
		}
	}

	return nil
}
