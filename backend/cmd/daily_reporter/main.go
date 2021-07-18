package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/oinume/lekcije/backend/cli"
	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/model"
)

func main() {
	m := &dailyReporterMain{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type dailyReporterMain struct {
	outStream io.Writer
	errStream io.Writer
	db        *gorm.DB
}

func (m *dailyReporterMain) run(args []string) error {
	flagSet := flag.NewFlagSet("daily_reporter", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
		targetDate = flagSet.String("target-date", time.Now().UTC().Format("2006-01-02"), "Target date (YYYY-MM-DD)")
		//logLevel   = flagSet.String("log-level", "info", "Log level") // TODO: Move to config
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	if *targetDate == "" {
		return fmt.Errorf("-target-date is required")
	}
	date, err := time.Parse("2006-01-02", *targetDate)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", *targetDate)
	}

	config.MustProcessDefault()
	db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err != nil {
		cli.WriteError(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()
	m.db = db

	if err := m.createStatNewLessonNotifier(date); err != nil {
		return err
	}
	if err := m.createStatDailyUserNotificationEvent(date); err != nil {
		return err
	}

	return nil
}

func (m *dailyReporterMain) createStatNewLessonNotifier(date time.Time) error {
	service := model.NewEventLogEmailService(m.db)
	stats, err := service.FindStatDailyNotificationEventByDate(date)
	if err != nil {
		return err
	}
	statUUs, err := service.FindStatDailyNotificationEventUUCountByDate(date)
	if err != nil {
		return err
	}

	values := make(map[string]*model.StatDailyNotificationEvent, 100)
	for _, s := range stats {
		values[s.Event] = s
	}

	statDailyNotificationEventService := model.NewStatDailyNotificationEventService(m.db)
	for _, s := range statUUs {
		v := values[s.Event]
		v.UUCount = s.UUCount
		if err := statDailyNotificationEventService.CreateOrUpdate(v); err != nil {
			return err
		}
	}

	//statsNewLessonNotifierService := model.NewStatsNewLessonNotifierService(m.GormDB)
	//for _, s := range statUUs {
	//	v := values[s.Event]
	//	v.UUCount = s.UUCount
	//	if err := statsNewLessonNotifierService.CreateOrUpdate(v); err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (m *dailyReporterMain) createStatDailyUserNotificationEvent(date time.Time) error {
	service := model.NewStatDailyUserNotificationEventService(m.db)
	return service.CreateOrUpdate(date)
}
