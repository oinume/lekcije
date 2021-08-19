package rollbar

import (
	"context"

	"github.com/rollbar/rollbar-go"
)

//go:generate moq -out=client.moq.go . Client

// Client abstracts rollbar.Client
type Client interface {
	ErrorWithStackSkipWithExtrasAndContext(
		ctx context.Context,
		level string,
		err error,
		skip int,
		extras map[string]interface{},
	)

	SetStackTracer(stackTracer rollbar.StackTracerFunc)
}
