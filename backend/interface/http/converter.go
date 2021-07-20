package http

import (
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
