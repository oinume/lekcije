package usecase_test

import (
	"bytes"
	"context"
	"testing"

	rollbar_go "github.com/rollbar/rollbar-go"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/infrastructure/rollbar"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/usecase"
)

func Test_ErrorRecorder_Record(t *testing.T) {
	tests := map[string]struct {
		clientMock *rollbar.ClientMock
	}{
		"ok": {
			clientMock: &rollbar.ClientMock{
				ErrorWithStackSkipWithExtrasAndContextFunc: func(ctx context.Context, level string, err error, skip int, extras map[string]interface{}) {
					// nop
				},
				SetStackTracerFunc: func(stackTracer rollbar_go.StackTracerFunc) {
					// nop
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var log bytes.Buffer
			errorRecorder := usecase.NewErrorRecorder(
				logger.NewAppLogger(&log, logger.NewLevel("info")),
				rollbar.NewErrorRecorderRepository(tt.clientMock),
			)
			err := errors.NewAnnotatedError(errors.CodeInternal)
			errorRecorder.Record(context.Background(), err, "1")
			assertion.AssertEqual(t, 1, len(tt.clientMock.ErrorWithStackSkipWithExtrasAndContextCalls()), "")
		})
	}
}
