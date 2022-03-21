package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/oinume/lekcije/backend/context_data"
	"github.com/oinume/lekcije/backend/graph/model"
)

func (r *queryResolver) Viewer(ctx context.Context) (*model.User, error) {
	user, err := context_data.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    fmt.Sprint(user.ID),
		Email: user.Email,
	}, nil
}
