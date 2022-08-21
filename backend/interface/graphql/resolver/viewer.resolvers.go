package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/interface/graphql/model"
)

// UpdateViewer is the resolver for the updateViewer field.
func (r *mutationResolver) UpdateViewer(ctx context.Context, input model.UpdateViewerInput) (*model.User, error) {
	user, err := context_data.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	if input.Email == nil {
		return toGraphQLUser(user), nil
	}

	email := *input.Email
	if !r.userUsecase.IsValidEmail(email) {
		return nil, fmt.Errorf("invalid email format") // TODO: invalid argument
	}
	duplicate, err := r.userUsecase.IsDuplicateEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, fmt.Errorf("email exists") // TODO: invalid argument
	}

	if err := r.userUsecase.UpdateEmail(ctx, uint(user.ID), email); err != nil {
		return nil, err
	}

	user.Email = email
	return toGraphQLUser(user), nil
}

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*model.User, error) {
	user, err := context_data.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	return toGraphQLUser(user), nil
}
