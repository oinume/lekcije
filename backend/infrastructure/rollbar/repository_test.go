package rollbar_test

import (
	"context"
	"testing"

	rollbar_go "github.com/rollbar/rollbar-go"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/infrastructure/rollbar"
	"github.com/oinume/lekcije/backend/internal/assertion"
)

func Test_errorRecorderRepository_Record(t *testing.T) {
	tests := map[string]struct {
		client *rollbar.ClientMock
	}{
		"normal": {
			client: &rollbar.ClientMock{
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
			r := rollbar.NewErrorRecorderRepository(tt.client)
			ctx := context.Background()
			err := errors.NewAnnotatedError(errors.CodeInternal)
			r.Record(ctx, err, "1")
			assertion.AssertEqual(t, 1, len(tt.client.ErrorWithStackSkipWithExtrasAndContextCalls()), "")
		})
	}
}
