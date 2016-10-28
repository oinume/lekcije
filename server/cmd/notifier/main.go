package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/notifier"
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
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	bootstrap.CheckCLIEnvVars()

	db, _, err := model.OpenDBAndSetToContext(context.Background(), bootstrap.CLIEnvVars.DBURL, !config.IsProductionEnv())
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
