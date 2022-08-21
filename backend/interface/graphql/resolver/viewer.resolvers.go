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
	panic(fmt.Errorf("not implemented: UpdateViewer - updateViewer"))
}

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*model.User, error) {
	user, err := context_data.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:           fmt.Sprint(user.ID),
		Email:        user.Email,
		ShowTutorial: !user.IsFollowedTeacher(),
	}, nil
}
