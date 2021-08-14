package usecase

import (
	"bytes"
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

func (er *ErrorRecorder) Record(ctx context.Context, err error, userID string) {
	fields := []zapcore.Field{
		zap.Error(err),
	}
	if e, ok := err.(errors.StackTracer); ok {
		b := &bytes.Buffer{}
		for _, f := range e.StackTrace() {
			fmt.Fprintf(b, "%+v\n", f)
		}
		fields = append(fields, zap.String("stacktrace", b.String()))
	}
	er.appLogger.Error("ErrorRecoder.Record", fields...)

	er.repo.Record(ctx, err, userID)
}
