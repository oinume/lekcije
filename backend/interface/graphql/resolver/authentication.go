package resolver

import (
	"context"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
	"github.com/oinume/lekcije/backend/usecase"
)

func authenticateFromContext(ctx context.Context, userUsecase *usecase.User) (*model2.User, error) {
	apiToken, err := context_data.GetAPIToken(ctx)
	if err != nil {
		return nil, failure.New(errors.Unauthenticated)
	}
	user, err := userUsecase.FindLoggedInUser(ctx, apiToken)
	if err != nil {
		if failure.Is(err, errors.NotFound) {
			return nil, failure.Translate(err, errors.Unauthenticated, failure.Messagef("no user found"))
		}
		return nil, failure.Translate(err, errors.Internal)
	}
	return user, nil
}
