package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/crawler"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/uber-go/zap"
)

var (
	concurrency  = flag.Int("concurrency", 1, "concurrency of fetcher. (default: 1)")
	continueFlag = flag.Bool("continue", true, "Continue to fetch if any error occurred. (default: true)")
	ids          = flag.String("ids", "", "Teacher IDs")
	followed     = flag.Bool("followed", false, "Fetch followed teachers")
	allTeachers  = flag.Bool("all", false, "Fetch all teachers ordered by evaluation")
	newTeachers  = flag.Bool("new", false, "Fetch all teachers ordered by new")
	logLevel     = flag.String("log-level", "info", "Log level")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatalf("err = %v", err) // TODO: Error handling
	}
	os.Exit(0)
}

func run() error {
	bootstrap.CheckCLIEnvVars()
	if *followed && *ids != "" {
		return fmt.Errorf("Can't specify -followed and -ids flags both.")
	}

	startedAt := time.Now().UTC()
	if *logLevel != "" {
		logger.App.SetLevel(logger.NewLevel(*logLevel))
	}
	logger.App.Info("fetcher started")
	defer func() {
		elapsed := time.Now().UTC().Sub(startedAt) / time.Millisecond
		logger.App.Info("fetcher finished", zap.Int("elapsed", int(elapsed)))
	}()

	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL(), 1, !config.IsProductionEnv())
	if err != nil {
		return err
	}
	defer db.Close()

	mCountryService := model.NewMCountryService(db)
	mCountries, err := mCountryService.LoadAll()
	if err != nil {
		return err
	}

	// TODO: Loader must implement pagination
	var loader crawler.TeacherIDLoader
	if *ids != "" {
		loader = &crawler.SpecificTeacherIDLoader{IDString: *ids}
	} else if *followed {
		loader = &crawler.FollowedTeacherIDLoader{DB: db}
	} else if *allTeachers {
		loader = &crawler.ScrapingTeacherIDLoader{}
	}
	teacherIDs, err := loader.Load()
	if err != nil {
		return err
	}

	fetcher := fetcher.NewTeacherLessonFetcher(nil, *concurrency, false, mCountries, logger.App)
	teacherService := model.NewTeacherService(db)
	for _, id := range teacherIDs {
		teacher, _, err := fetcher.Fetch(id)
		if err != nil {
			if *continueFlag {
				logger.App.Error("Error during TeacherLessonFetcher.Fetch", zap.Error(err))
				continue
			} else {
				return err
			}
		}
		if err := teacherService.CreateOrUpdate(teacher); err != nil {
			if *continueFlag {
				logger.App.Error("Error during TeacherService.CreateOrUpdate", zap.Error(err))
			} else {
				return err
			}
		}
	}

	return nil
}
