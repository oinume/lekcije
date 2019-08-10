package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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
	Interval        *time.Duration
	LogLevel        *string
}

func (m *Main) Run() error {
	config.MustProcessDefault()
	if *m.Followed && *m.SpecifiedIDs != "" {
		return fmt.Errorf("can't specify -followed and -ids flags both")
	}

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
	defer db.Close()

	mCountryService := model.NewMCountryService(db)
	mCountries, err := mCountryService.LoadAll(ctx)
	if err != nil {
		return err
	}

	loader := m.createLoader(db)
	lessonFetcher := fetcher.NewLessonFetcher(nil, *m.Concurrency, false, mCountries, logger.App)
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
				// TODO: update lessons

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		time.Sleep(*m.Interval)
	}

	return nil
}

func (m *Main) createLoader(db *gorm.DB) TeacherIDLoader {
	var loader TeacherIDLoader
	if *m.SpecifiedIDs != "" {
		loader = &specificTeacherIDLoader{idString: *m.SpecifiedIDs}
	} else if *m.Followed {
		loader = &followedTeacherIDLoader{db: db}
	} else if *m.All {
		loader = NewScrapingTeacherIDLoader(ByRating, nil)
	} else if *m.New {
		loader = NewScrapingTeacherIDLoader(ByNew, nil)
	} else {
		loader = NewScrapingTeacherIDLoader(ByRating, nil)
	}
	return loader
}
