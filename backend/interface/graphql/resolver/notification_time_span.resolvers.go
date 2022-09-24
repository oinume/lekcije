package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/context_data"
	lerrors "github.com/oinume/lekcije/backend/errors"
	graphqlmodel "github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/model2"
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

	go func() {
		if err := r.gaMeasurementUsecase.SendEvent(
			ctx,
			context_data.MustGAMeasurementEvent(ctx),
			model2.GAMeasurementEventCategoryUser,
			"updateNotificationTimeSpan",
			fmt.Sprint(user.ID),
			0,
			uint32(user.ID),
		); err != nil {
			panic(err) // TODO: Better error handling
		}
	}()

	timeSpansGraphQL, err := toGraphQLNotificationTimeSpans(timeSpans)
	if err != nil {
		return nil, err
	}
	return &graphqlmodel.NotificationTimeSpanPayload{
		TimeSpans: timeSpansGraphQL,
	}, nil
}
