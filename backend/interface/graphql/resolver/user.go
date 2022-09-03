package resolver

import (
	"context"

	"github.com/oinume/lekcije/backend/interface/graphql/model"
)

// ShowTutorial is the resolver for the showTutorial field.
func (r *userResolver) ShowTutorial(ctx context.Context, obj *model.User) (bool, error) {
	return obj.ShowTutorial, nil
}
