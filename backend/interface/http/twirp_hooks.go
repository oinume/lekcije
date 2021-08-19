package http

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/oinume/lekcije/backend/usecase"
)

type ErrorHandlerHooks struct{}

func NewErrorRecorderHooks(errorRecorder *usecase.ErrorRecorder) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			//log.Println("Error: " + string(twerr.Code()))
			errorRecorder.Record(ctx, twerr, "") // TODO: get user id from ctx
			return ctx
		},
	}
}
