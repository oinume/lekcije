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
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
		return fmt.Errorf("can't specify -followed and -ids flags both")
	}

	startedAt := time.Now().UTC()
	//if *m.LogLevel != "" {
	//	//logger.App.SetLevel(logger.NewLevel(*m.LogLevel))
	//}
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
	fetcher := fetcher.NewLessonFetcher(nil, *m.Concurrency, false, mCountries, logger.App)
	teacherService := model.NewTeacherService(db)
	for cursor := loader.GetInitialCursor(); cursor != ""; {
		var teacherIDs []uint32
		var err error
		teacherIDs, cursor, err = loader.Load(cursor)
		if err != nil {
			return err
		}

		var g errgroup.Group
		for _, id := range teacherIDs {
			id := id
			g.Go(func() error {
				teacher, _, err := fetcher.Fetch(id)
				if err != nil {
					if *m.ContinueOnError {
						logger.App.Error("Error during LessonFetcher.Fetch", zap.Error(err))
						return nil
					} else {
						return err
					}
				}
				if err := teacherService.CreateOrUpdate(teacher); err != nil {
					if *m.ContinueOnError {
						logger.App.Error("Error during TeacherService.CreateOrUpdate", zap.Error(err))
						return nil
					} else {
						return err
					}
				}

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		time.Sleep(1 * time.Second) // TODO: -interval flag
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
		loader = newScrapingTeacherIDLoader(byRating, nil)
	} else if *m.New {
		loader = newScrapingTeacherIDLoader(byNew, nil)
	} else {
		loader = newScrapingTeacherIDLoader(byRating, nil)
	}
	return loader
}
