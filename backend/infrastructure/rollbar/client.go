package rollbar

//go:generate moq -out=client.moq.go . Client

import (
	"context"

	"github.com/rollbar/rollbar-go"
)

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
