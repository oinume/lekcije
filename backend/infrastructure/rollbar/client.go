package rollbar

import "context"

// Client abstracts rollbar.Client
type Client interface {
	ErrorWithStackSkipWithExtrasAndContext(
		ctx context.Context,
		level string,
		err error,
		skip int,
		extras map[string]interface{},
	)
}
