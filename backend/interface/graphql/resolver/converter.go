package resolver

import (
	"fmt"
	"time"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/errors"
	graphqlmodel "github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/model2"
)

func toGraphQLUser(user *model2.User) *graphqlmodel.User {
	return &graphqlmodel.User{
		ID:           fmt.Sprint(user.ID),
		Email:        user.Email,
		ShowTutorial: !user.IsFollowedTeacher(),
	}
}

func toModelNotificationTimeSpans(userID uint, timeSpans []*graphqlmodel.NotificationTimeSpanInput) []*model2.NotificationTimeSpan {
	values := make([]*model2.NotificationTimeSpan, len(timeSpans))
	for i, v := range timeSpans {
		fromTime := fmt.Sprintf("%v:%v", v.FromHour, v.FromMinute)
		toTime := fmt.Sprintf("%v:%v", v.ToHour, v.ToMinute)
		values[i] = &model2.NotificationTimeSpan{
			UserID:   userID,
			Number:   uint8(i + 1),
			FromTime: fromTime,
			ToTime:   toTime,
		}
	}
	return values
}

func toGraphQLNotificationTimeSpans(timeSpans []*model2.NotificationTimeSpan) ([]*graphqlmodel.NotificationTimeSpan, error) {
	values := make([]*graphqlmodel.NotificationTimeSpan, len(timeSpans))
	for i, v := range timeSpans {
		fromTime, err := time.Parse("15:04:05", v.FromTime)
		if err != nil {
			return nil, failure.Translate(err, errors.Internal, failure.Messagef("Invalid time format: FromTime=%v", v.FromTime))
		}
		toTime, err := time.Parse("15:04:05", v.ToTime)
		if err != nil {
			return nil, failure.Translate(err, errors.Internal, failure.Messagef("Invalid time format: ToTime=%v", v.ToTime))
		}
		values[i] = &graphqlmodel.NotificationTimeSpan{
			FromHour:   fromTime.Hour(),
			FromMinute: fromTime.Minute(),
			ToHour:     toTime.Hour(),
			ToMinute:   toTime.Minute(),
		}
	}
	return values, nil
}
