package dmm_eikaiwa

import (
	"context"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/model2"
)

type lessonFetcher struct {
	httpClient *http.Client
	semaphore  chan struct{}
	caching    bool
	//	cache      map[uint32]*teacherLessons
	cacheLock  *sync.RWMutex
	logger     *zap.Logger
	mCountries []*model2.MCountry
}

func NewLessonFetcher(
	httpClient *http.Client,
	concurrency int,
	caching bool,
	mCountries []*model2.MCountry,
	log *zap.Logger,
) repository.LessonFetcher {
	return &lessonFetcher{
		httpClient: httpClient,
		caching:    caching,
		mCountries: mCountries, // TODO: Define model2.MCountries
		logger:     log,
	}
}

func (f *lessonFetcher) Fetch(ctx context.Context, teacherID uint32) (*model2.Teacher, []*model2.Lesson, error) {
	// TODO: Use otel for http client tracing
	panic("implement")
}
