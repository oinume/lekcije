package http

import (
	"fmt"
	"time"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
	api_v1 "github.com/oinume/lekcije/proto-gen/go/proto/api/v1"
)

func NotificationTimeSpansProto(timeSpans []*model2.NotificationTimeSpan) ([]*api_v1.NotificationTimeSpan, error) {
	values := make([]*api_v1.NotificationTimeSpan, 0, len(timeSpans))
	for _, v := range timeSpans {
		fromTime, err := time.Parse("15:04:05", v.FromTime)
		if err != nil {
			return nil, errors.NewInternalError(
				errors.WithError(err),
				errors.WithMessagef("Invalid time format: FromTime=%v", v.FromTime),
			)
		}
		toTime, err := time.Parse("15:04:05", v.ToTime)
		if err != nil {
			return nil, errors.NewInternalError(
				errors.WithError(err),
				errors.WithMessagef("Invalid time format: ToTime=%v", v.ToTime),
			)
		}
		values = append(values, &api_v1.NotificationTimeSpan{
			FromHour:   int32(fromTime.Hour()),
			FromMinute: int32(fromTime.Minute()),
			ToHour:     int32(toTime.Hour()),
			ToMinute:   int32(toTime.Minute()),
		})
	}
	return values, nil
}

func NotificationTimeSpansModel(timeSpans []*api_v1.NotificationTimeSpan, userID uint) []*model2.NotificationTimeSpan {
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

func TeachersProto(teachers []*model2.Teacher) []*api_v1.Teacher {
	values := make([]*api_v1.Teacher, len(teachers))
	for i, t := range teachers {
		values[i] = &api_v1.Teacher{
			Id:   uint32(t.ID),
			Name: t.Name,
		}
	}
	return values
}
