package main

import (
	"flag"
	"os"

	"github.com/oinume/lekcije/server/notifier"
	"github.com/oinume/lekcije/server/util"
)

func main() {
	m := &notifier.Main{}
	m.Concurrency = flag.Int("concurrency", 1, "Concurrency of fetcher")
	m.DryRun = flag.Bool("dry-run", false, "Don't update database with fetched lessons")
	m.FetcherCache = flag.Bool("fetcher-cache", false, "Cache teacher and lesson data in Fetcher")
	m.NotificationInterval = flag.Int("notification-interval", 0, "Notification interval")
	m.SendEmail = flag.Bool("send-email", true, "Flag to send email")
	m.LogLevel = flag.String("log-level", "info", "Log level")
	m.ProfileMode = flag.String("profile-mode", "", "block|cpu|mem|trace")

	flag.Parse()
	if err := m.Run(); err != nil {
		util.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
