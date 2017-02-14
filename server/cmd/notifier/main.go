package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/notifier"
	"github.com/uber-go/zap"
	"github.com/pkg/profile"
)

var (
	dryRun      = flag.Bool("dry-run", false, "Don't update database with fetched lessons")
	concurrency = flag.Int("concurrency", 1, "concurrency of fetcher")
	logLevel    = flag.String("log-level", "info", "Log level")
	profileMode = flag.String("profile-mode", "", "block or cpu or mem")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("err = %v", err) // TODO: Error handling
	}
	os.Exit(0)
}

func run() error {
	switch *profileMode {
	case "block":
		defer profile.Start(profile.ProfilePath("."), profile.BlockProfile).Stop()
	case "cpu":
		defer profile.Start(profile.ProfilePath("."), profile.CPUProfile).Stop()
	case "mem":
		defer profile.Start(profile.ProfilePath("."), profile.MemProfile).Stop()
	}

	bootstrap.CheckCLIEnvVars()
	startedAt := time.Now().UTC()
	if *logLevel != "" {
		logger.App.SetLevel(logger.NewLevel(*logLevel))
	}
	logger.App.Info("notifier started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info("notifier finished", zap.Int("elapsed", int(elapsed)))
	}()

	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL, 1, !config.IsProductionEnv())
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO: Define method in UserService
	var users []*model.User
	userSQL := `SELECT * FROM user WHERE email_verified = 1`
	result := db.Raw(userSQL).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return errors.InternalWrapf(result.Error, "Failed to find Users")
	}

	mCountries, err := model.NewMCountryService(db).LoadAll()
	if err != nil {
		return errors.InternalWrapf(err, "Failed to load all MCountries")
	}
	fetcher := fetcher.NewTeacherLessonFetcher(nil, *concurrency, mCountries, logger.App)
	notifier := notifier.NewNotifier(db, fetcher, *dryRun)
	defer notifier.Close()
	for _, user := range users {
		if err := notifier.SendNotification(user); err != nil {
			return err
		}
	}

	return nil
}
