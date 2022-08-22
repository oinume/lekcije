package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/morikuni/failure"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/errors"
)

type Middleware struct {
	appLogger *zap.Logger
}

func NewMiddleware(appLogger *zap.Logger) *Middleware {
	return &Middleware{
		appLogger: appLogger,
	}
}

func (m *Middleware) AroundOperations(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)
	m.appLogger.Info("graphql", zap.String("operation", oc.OperationName), zap.String("rawQuery", oc.RawQuery))
	return next(ctx)
}

func (m *Middleware) ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)

	var errorCode string
	code, ok := failure.CodeOf(e)
	if ok {
		errorCode = code.ErrorCode()
	} else {
		// Map unknown error to internal error
		errorCode = errors.Internal.ErrorCode()
	}

	err.Extensions = map[string]interface{}{
		"code": errorCode,
	}

	m.appLogger.Info("graphql", zap.String("code", errorCode), zap.Error(e))

	return err
}
