package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/notifier"
	"github.com/uber-go/zap"
)

var (
	dryRun      = flag.Bool("dry-run", false, "Don't update database with fetched lessons")
	concurrency = flag.Int("concurrency", 1, "concurrency of fetcher")
	logLevel    = flag.String("log-level", "info", "Log level")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("err = %v", err) // TODO: Error handling
	}
	os.Exit(0)
}

func run() error {
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

	var users []*model.User
	userSql := `SELECT * FROM user WHERE email_verified = 1`
	result := db.Raw(userSql).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return errors.InternalWrapf(result.Error, "")
	}

	notifier := notifier.NewNotifier(db, *concurrency, *dryRun)
	defer notifier.Close()
	for _, user := range users {
		if err := notifier.SendNotification(user); err != nil {
			return err
		}
	}

	return nil
}
