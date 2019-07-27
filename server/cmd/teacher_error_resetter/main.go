package main

import (
	"flag"
	"os"

	"github.com/oinume/lekcije/server/cli"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/teacher_error_resetter"
)

func main() {
	m := &teacher_error_resetter.Main{}
	m.Concurrency = flag.Int("concurrency", 1, "Concurrency of fetcher")
	m.DryRun = flag.Bool("dry-run", false, "Don't update database with fetched lessons")
	m.LogLevel = flag.String("log-level", "info", "Log level")
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
