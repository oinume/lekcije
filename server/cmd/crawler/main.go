package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/cli"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/crawler"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	m := &crawlerMain{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type crawlerMain struct {
	outStream io.Writer
	errStream io.Writer
	db        *gorm.DB
}

func (m *crawlerMain) run(args []string) error {
	flagSet := flag.NewFlagSet("crawler", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
		concurrency     = flagSet.Int("concurrency", 1, "Concurrency of crawler. (default=1)")
		continueOnError = flagSet.Bool("continue", true, "Continue to crawl if any error occurred. (default=true)")
		specifiedIDs    = flagSet.String("ids", "", "Teacher IDs")
		followedOnly    = flagSet.Bool("followedOnly", false, "Crawl followedOnly teachers")
		all             = flagSet.Bool("all", false, "Crawl all teachers ordered by evaluation")
		newOnly         = flagSet.Bool("new", false, "Crawl all teachers ordered by new")
		interval        = flagSet.Duration("interval", 1*time.Second, "Fetch interval. (default=1s)")
		//logLevel        = flag.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	if *followedOnly && *specifiedIDs != "" {
		return fmt.Errorf("can't specify -followedOnly and -ids flags both")
	}

	config.MustProcessDefault()

	ctx := context.Background()
	startedAt := time.Now().UTC()
	//if *m.LogLevel != "" {
	//	//logger.App.SetLevel(logger.NewLevel(*m.LogLevel))
	//}
	logger.App.Info("crawler started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info("crawler finished", zap.Int("elapsed", int(elapsed)))
	}()

	db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	mCountryService := model.NewMCountryService(db)
	mCountries, err := mCountryService.LoadAll(ctx)
	if err != nil {
		return err
	}

	loader := m.createLoader(db, *specifiedIDs, *followedOnly, *all, *newOnly)
	lessonFetcher := fetcher.NewLessonFetcher(nil, *concurrency, false, mCountries, logger.App)
	teacherService := model.NewTeacherService(db)
	for cursor := loader.GetInitialCursor(); cursor != ""; {
		var teacherIDs []uint32
		var err error
		teacherIDs, cursor, err = loader.Load(cursor)
		if err != nil {
			return err
		}

		// TODO: semaphore
		var g errgroup.Group
		for _, id := range teacherIDs {
			id := id
			g.Go(func() error {
				teacher, _, err := lessonFetcher.Fetch(ctx, id)
				if err != nil {
					if *continueOnError {
						logger.App.Error("Error during LessonFetcher.Fetch", zap.Error(err))
						return nil
					} else {
						return err
					}
				}
				if err := teacherService.CreateOrUpdate(teacher); err != nil {
					if *continueOnError {
						logger.App.Error("Error during TeacherService.CreateOrUpdate", zap.Error(err))
						return nil
					} else {
						return err
					}
				}
				// TODO: update lessons

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		time.Sleep(*interval)
	}

	return nil
}

func (m *crawlerMain) createLoader(
	db *gorm.DB,
	specifiedIDs string,
	followed bool,
	all bool,
	newOnly bool,
) crawler.TeacherIDLoader {
	var loader crawler.TeacherIDLoader
	if specifiedIDs != "" {
		loader = crawler.NewSpecificTeacherIDLoader(specifiedIDs)
	} else if followed {
		loader = crawler.NewFollowedTeacherIDLoader(db)
	} else if all {
		loader = crawler.NewScrapingTeacherIDLoader(crawler.ByRating, nil)
	} else if newOnly {
		loader = crawler.NewScrapingTeacherIDLoader(crawler.ByNew, nil)
	} else {
		loader = crawler.NewScrapingTeacherIDLoader(crawler.ByRating, nil)
	}
	return loader
}
