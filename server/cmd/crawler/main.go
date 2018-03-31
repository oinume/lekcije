package main

import (
	"flag"
	"os"
	"time"

	"github.com/oinume/lekcije/server/crawler"
	"github.com/oinume/lekcije/server/util"
)

func main() {
	m := &crawler.Main{}
	m.Concurrency = flag.Int("concurrency", 1, "Concurrency of crawler. (default=1)")
	m.ContinueOnError = flag.Bool("continue", true, "Continue to crawl if any error occurred. (default=true)")
	m.SpecifiedIDs = flag.String("ids", "", "Teacher IDs")
	m.Followed = flag.Bool("followed", false, "Crawl followed teachers")
	m.All = flag.Bool("all", false, "Crawl all teachers ordered by evaluation")
	m.New = flag.Bool("new", false, "Crawl all teachers ordered by new")
	m.Interval = flag.Duration("interval", 1*time.Second, "Fetch interval. (default=1s)")
	m.LogLevel = flag.String("log-level", "info", "Log level")

	flag.Parse()
	if err := m.Run(); err != nil {
		util.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
