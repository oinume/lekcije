package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/context_data"
	lerrors "github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/model2"
)

// UpdateViewer is the resolver for the updateViewer field.
func (r *mutationResolver) UpdateViewer(ctx context.Context, input model.UpdateViewerInput) (*model.User, error) {
	user, err := authenticateFromContext(ctx, r.userUsecase)
	if err != nil {
		return nil, err
	}
	if input.Email == nil {
		return toGraphQLUser(user), nil
	}

	email := *input.Email
	if !r.userUsecase.IsValidEmail(email) {
		// TODO: i18n
		return nil, failure.New(lerrors.InvalidArgument, failure.Messagef("正しいメールアドレスを入力してください。"))
	}
	duplicate, err := r.userUsecase.IsDuplicateEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if duplicate {
		// TODO: i18n
		return nil, failure.New(lerrors.InvalidArgument, failure.Messagef("メールアドレスはすでに登録されています。"))
	}

	if err := r.userUsecase.UpdateEmail(ctx, user.ID, email); err != nil {
		return nil, err
	}
	user.Email = email

	go func() {
		if err := r.gaMeasurementUsecase.SendEvent(
			ctx,
			context_data.MustGAMeasurementEvent(ctx),
			model2.GAMeasurementEventCategoryUser,
			"update",
			fmt.Sprint(user.ID),
			0,
			uint32(user.ID),
		); err != nil {
			panic(err) // TODO: Better error handling
		}
	}()

	return toGraphQLUser(user), nil
}

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*model.User, error) {
	user, err := authenticateFromContext(ctx, r.userUsecase)
	if err != nil {
		return nil, err
	}
	return toGraphQLUser(user), nil
}
