package main

import (
	"context"
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
	dryRun   = flag.Bool("dry-run", false, "Don't update database with fetched lessons")
	logLevel = flag.String("log-level", "info", "Log level")
)

func main() {
	flag.Parse()
	err := run()
	if err != nil {
		log.Fatalf("err = %v", err) // TODO: Error handling
	}
	os.Exit(0)
}

func run() error {
	bootstrap.CheckCLIEnvVars()
	startedAt := time.Now().UTC()
	logger.AppLogger.Info("notifier started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.AppLogger.Info("notifier finished", zap.Int("elapsed", int(elapsed)))
	}()

	db, _, err := model.OpenDBAndSetToContext(
		context.Background(), bootstrap.CLIEnvVars.DBURL,
		1, !config.IsProductionEnv(),
	)
	if err != nil {
		return err
	}

	var users []*model.User
	userSql := `SELECT * FROM user WHERE email_verified = 1`
	result := db.Raw(userSql).Scan(&users)
	if result.Error != nil && !result.RecordNotFound() {
		return errors.InternalWrapf(result.Error, "")
	}

	notifier := notifier.NewNotifier(db, *dryRun)
	for _, user := range users {
		if err := notifier.SendNotification(user); err != nil {
			return err
		}
	}

	return nil
}
