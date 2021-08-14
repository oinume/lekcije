package repository

import "context"

type ErrorRecorder interface {
	Record(ctx context.Context, err error, userID string)
}
