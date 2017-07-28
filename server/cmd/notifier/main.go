package main

import (
	"flag"
	"log"
	"os"

	"github.com/oinume/lekcije/server/notifier"
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
		log.Fatalf("err = %v", err) // TODO: Error handling
	}
	os.Exit(0)
}
