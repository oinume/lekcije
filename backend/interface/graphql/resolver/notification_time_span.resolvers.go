package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/morikuni/failure"

	lerrors "github.com/oinume/lekcije/backend/errors"
	graphqlmodel "github.com/oinume/lekcije/backend/interface/graphql/model"
)

// UpdateNotificationTimeSpans is the resolver for the updateNotificationTimeSpans field.
func (r *mutationResolver) UpdateNotificationTimeSpans(ctx context.Context, input graphqlmodel.UpdateNotificationTimeSpansInput) (*graphqlmodel.NotificationTimeSpanPayload, error) {
	user, err := authenticateFromContext(ctx, r.userUsecase)
	if err != nil {
		return nil, err
	}

	if len(input.TimeSpans) > 3 {
		return nil, failure.New(lerrors.InvalidArgument, failure.Messagef("レッスン希望時間帯は3つまで登録可能です"))
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
