package crawler

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/fetcher"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/uber-go/zap"
)

type Main struct {
	Concurrency     *int
	ContinueOnError *bool
	SpecifiedIDs    *string
	Followed        *bool
	All             *bool
	New             *bool
	LogLevel        *string
}

func (m *Main) Run() error {
	bootstrap.CheckCLIEnvVars()
	if *m.Followed && *m.SpecifiedIDs != "" {
		return fmt.Errorf("Can't specify -followed and -ids flags both.")
	}

	startedAt := time.Now().UTC()
	if *m.LogLevel != "" {
		logger.App.SetLevel(logger.NewLevel(*m.LogLevel))
	}
	logger.App.Info("crawler started")
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

	loader := m.createLoader(db)
	fetcher := fetcher.NewTeacherLessonFetcher(nil, *m.Concurrency, false, mCountries, logger.App)
	teacherService := model.NewTeacherService(db)
	for cursor := "a"; cursor != ""; {
		var teacherIDs []uint32
		var err error
		teacherIDs, cursor, err = loader.Load(cursor)
		if err != nil {
			return err
		}

		for _, id := range teacherIDs {
			teacher, _, err := fetcher.Fetch(id)
			if err != nil {
				if *m.ContinueOnError {
					logger.App.Error("Error during TeacherLessonFetcher.Fetch", zap.Error(err))
					continue
				} else {
					return err
				}
			}
			if err := teacherService.CreateOrUpdate(teacher); err != nil {
				if *m.ContinueOnError {
					logger.App.Error("Error during TeacherService.CreateOrUpdate", zap.Error(err))
				} else {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Main) createLoader(db *gorm.DB) teacherIDLoader {
	var loader teacherIDLoader
	if *m.SpecifiedIDs != "" {
		loader = &specificTeacherIDLoader{idString: *m.SpecifiedIDs}
	} else if *m.Followed {
		loader = &followedTeacherIDLoader{db: db}
	} else if *m.All {
		loader = &scrapingTeacherIDLoader{order: byRating}
	} else if *m.New {
		loader = &scrapingTeacherIDLoader{order: byNew}
	} else {
		loader = &scrapingTeacherIDLoader{order: byRating}
	}
	return loader
}
