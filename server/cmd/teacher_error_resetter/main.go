package main

import (
	"flag"
	"os"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/teacher_error_resetter"
	"github.com/oinume/lekcije/server/util"
)

func main() {
	m := &teacher_error_resetter.Main{}
	m.Concurrency = flag.Int("concurrency", 1, "Concurrency of fetcher")
	m.DryRun = flag.Bool("dry-run", false, "Don't update database with fetched lessons")
	m.LogLevel = flag.String("log-level", "info", "Log level")
	flag.Parse()

	bootstrap.CheckCLIEnvVars()
	dbLogging := *m.LogLevel == "debug"
	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL(), 1, dbLogging)
	if err != nil {
		util.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	m.DB = db

	if err := m.Run(); err != nil {
		util.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
