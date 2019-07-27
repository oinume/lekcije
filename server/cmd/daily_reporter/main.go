package main

import (
	"flag"
	"os"
	"time"

	"github.com/oinume/lekcije/server/cli"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/daily_reporter"
	"github.com/oinume/lekcije/server/model"
)

func main() {
	m := &daily_reporter.Main{}
	m.TargetDate = flag.String("target-date", time.Now().UTC().Format("2006-01-02"), "Target date (YYYY-MM-DD)")
	m.LogLevel = flag.String("log-level", "info", "Log level") // TODO: Move to config
	flag.Parse()

	config.MustProcessDefault()
	db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err != nil {
		cli.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	m.DB = db

	if err := m.Run(); err != nil {
		cli.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
