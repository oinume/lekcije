package log

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct{}

func FromContextZap(ctx context.Context) (*zap.Logger, bool) {
	v := ctx.Value(contextKey{})
	if v == nil {
		return nil, false
	}
	logger, ok := v.(*zap.Logger)
	if !ok {
		return nil, false
	}
	return logger, true
}

func MustFromContextZap(ctx context.Context) *zap.Logger {
	logger, ok := FromContextZap(ctx)
	if !ok {
		panic("failed to get logger from context")
	}
	return logger
}

func WithContextZap(parent context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(parent, contextKey{}, logger)
}
