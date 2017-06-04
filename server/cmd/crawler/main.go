package main

import (
	"flag"
	"log"
	"os"

	"github.com/oinume/lekcije/server/crawler"
)

func main() {
	m := &crawler.Main{}
	m.Concurrency = flag.Int("concurrency", 1, "Concurrency of crawler. (default=1)")
	m.ContinueOnError = flag.Bool("continue", true, "Continue to crawl if any error occurred. (default=true)")
	m.SpecifiedIDs = flag.String("ids", "", "Teacher IDs")
	m.Followed = flag.Bool("followed", false, "Crawl followed teachers")
	m.All = flag.Bool("all", false, "Crawl all teachers ordered by evaluation")
	m.New = flag.Bool("new", false, "Crawl all teachers ordered by new")
	m.LogLevel = flag.String("log-level", "info", "Log level")

	flag.Parse()
	if err := m.Run(); err != nil {
		log.Fatalf("err = %v", err) // TODO: Error handling
	}
	os.Exit(0)
}
