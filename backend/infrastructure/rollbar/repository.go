package rollbar

import (
	"context"
	"runtime"

	pkg_errors "github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"

	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/errors"
)

type errorRecorderRepository struct {
	client Client
}

func NewErrorRecorderRepository(client Client) repository.ErrorRecorder {
	client.SetStackTracer(StackTracer)
	return &errorRecorderRepository{
		client: client,
	}
}

func (r *errorRecorderRepository) Record(ctx context.Context, err error, userID string) {
	if userID != "" {
		ctx = rollbar.NewPersonContext(ctx, &rollbar.Person{
			Id: userID,
		})
	}
	r.client.ErrorWithStackSkipWithExtrasAndContext(ctx, "error", err, 0, nil)
}

func StackTracer(err error) ([]runtime.Frame, bool) {
	switch e := err.(type) {
	case *errors.AnnotatedError:
		if !e.OutputStackTrace() {
			return nil, false
		}
		return toFrames(e.StackTrace()), true
	case errors.StackTracer:
		return toFrames(e.StackTrace()), true
	default:
		return nil, false
	}
}

func toFrames(st pkg_errors.StackTrace) []runtime.Frame {
	frames := make([]runtime.Frame, 0, len(st))
	for _, f := range st {
		pc := uintptr(f)
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		frame := runtime.Frame{
			PC:       pc,
			Func:     fn,
			Function: fn.Name(),
			File:     file,
			Line:     line,
		}
		frames = append(frames, frame)
	}
	return frames
}
