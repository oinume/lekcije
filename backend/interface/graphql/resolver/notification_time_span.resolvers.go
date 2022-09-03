package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	graphqlmodel "github.com/oinume/lekcije/backend/interface/graphql/model"
)

// UpdateNotificationTimeSpans is the resolver for the updateNotificationTimeSpans field.
func (r *mutationResolver) UpdateNotificationTimeSpans(
	ctx context.Context,
	input graphqlmodel.UpdateNotificationTimeSpansInput,
) (*graphqlmodel.NotificationTimeSpanPayload, error) {
	user, err := authenticateFromContext(ctx, r.userUsecase)
	if err != nil {
		return nil, err
	}

	timeSpans := toModelNotificationTimeSpans(user.ID, input.TimeSpans)
	if err := r.notificationTimeSpanUsecase.UpdateAll(ctx, user.ID, timeSpans); err != nil {
		return nil, err
	}
	timeSpansGraphQL, err := toGraphQLNotificationTimeSpans(timeSpans)
	if err != nil {
		return nil, err
	}
	return &graphqlmodel.NotificationTimeSpanPayload{
		TimeSpans: timeSpansGraphQL,
	}, nil
}
