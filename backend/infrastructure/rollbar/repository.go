package rollbar

import (
	"context"

	"github.com/rollbar/rollbar-go"
	_ "github.com/rollbar/rollbar-go"

	"github.com/oinume/lekcije/backend/repository"
)

type errorRecorderRepository struct {
	client Client
}

func NewErrorRecorderRepository(client *rollbar.Client) repository.ErrorRecorder {
	//client := rollbar.New(token, environment, "beta" /* TODO: version */, "", "/")
	//client.SetStackTracer() TODO
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
