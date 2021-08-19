package repository

import "context"

type ErrorRecorder interface {
	Record(ctx context.Context, err error, userID string)
}

type NopErrorRecorder struct{}

func (ner *NopErrorRecorder) Record(ctx context.Context, err error, userID string) {
	// nop
}
