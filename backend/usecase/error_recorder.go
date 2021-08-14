package usecase

import (
	"context"

	"github.com/twitchtv/twirp"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/repository"
)

type ErrorRecorder struct {
	appLogger *zap.Logger
	repo      repository.ErrorRecorder
}

func NewErrorRecorder(appLogger *zap.Logger, repo repository.ErrorRecorder) *ErrorRecorder {
	return &ErrorRecorder{
		appLogger: appLogger, // TODO: remove and get logger from ctx
		repo:      repo,
	}
}

func (eh *ErrorRecorder) Record(ctx context.Context, err error, userID string) {
	switch err.(type) {
	case twirp.Error:
		eh.repo.Record(ctx, err, userID)
	case *errors.AnnotatedError:
		eh.repo.Record(ctx, err, userID)
	}
}
