package notifier

import (
	"fmt"
	"net/http"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/emailer"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/stopwatch"
	"github.com/pkg/profile"
	"go.uber.org/zap"
)

type Main struct {
	Concurrency          *int
	DryRun               *bool
	NotificationInterval *int
	SendEmail            *bool
	FetcherCache         *bool
	LogLevel             *string
	ProfileMode          *string
}

func (m *Main) Run() error {
	var sw stopwatch.Stopwatch
	switch *m.ProfileMode {
	case "block":
		defer profile.Start(profile.ProfilePath("."), profile.BlockProfile).Stop()
	case "cpu":
		defer profile.Start(profile.ProfilePath("."), profile.CPUProfile).Stop()
	case "mem":
		defer profile.Start(profile.ProfilePath("."), profile.MemProfile).Stop()
	case "trace":
		defer profile.Start(profile.ProfilePath("."), profile.TraceProfile).Stop()
	case "stopwatch":
		sw = stopwatch.NewSync()
		sw.Start()
	}

	bootstrap.CheckCLIEnvVars()
	startedAt := time.Now().UTC()
	//if *m.LogLevel != "" {
	//	//logger.App.SetLevel(logger.NewLevel(*m.LogLevel))
	//}
	logger.App.Info(fmt.Sprintf("notifier started (interval=%d)", *m.NotificationInterval))
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info(
			fmt.Sprintf("notifier finished (interval=%d)", *m.NotificationInterval),
			zap.Int("elapsed", int(elapsed)),
		)
	}()

	dbLogging := *m.LogLevel == "debug"
	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL(), 1, dbLogging)
	if err != nil {
		return err
	}
	defer db.Close()
	sw.Mark("model.OpenDB")

	if *m.NotificationInterval == 0 {
		return fmt.Errorf("-notification-interval is required")
	}
	users, err := model.NewUserService(db).FindAllEmailVerifiedIsTrue(*m.NotificationInterval)
	if err != nil {
		return err
	}
	mCountries, err := model.NewMCountryService(db).LoadAll()
	if err != nil {
		return err
	}
	fetcher := fetcher.NewLessonFetcher(nil, *m.Concurrency, *m.FetcherCache, mCountries, logger.App)

	var sender emailer.Sender
	if *m.SendEmail {
		sender = emailer.NewSendGridSender(http.DefaultClient)
	} else {
		sender = &emailer.NoSender{}
	}

	n := NewNotifier(db, fetcher, *m.DryRun, sender, sw)
	defer n.Close()
	for _, user := range users {
		if err := n.SendNotification(user); err != nil {
			return err
		}
	}

	return nil
}
