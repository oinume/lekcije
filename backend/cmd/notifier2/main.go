package main

import (
	"flag"
	"io"
	"os"

	"github.com/oinume/lekcije/backend/cli"
)

func main() {
	m := &notifierMain{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type notifierMain struct {
	outStream io.Writer
	errStream io.Writer
}

func (m *notifierMain) run(args []string) error {
	flagSet := flag.NewFlagSet("notifier", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
		concurrency          = flagSet.Int("concurrency", 1, "Concurrency of fetcher")
		dryRun               = flagSet.Bool("dry-run", false, "Don't update database with fetched lessons")
		fetcherCache         = flagSet.Bool("fetcher-cache", false, "Cache teacher and lesson data in Fetcher")
		notificationInterval = flagSet.Int("notification-interval", 0, "Notification interval")
		sendEmail            = flagSet.Bool("send-email", true, "Flag to send email")
		logLevel             = flagSet.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	return nil
}
