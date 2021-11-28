package http

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/model2"
	api_v1 "github.com/oinume/lekcije/backend/proto_gen/proto/api/v1"
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

func UserProto(user *model.User) *api_v1.User {
	followedTeacherAt := time.Time{}
	if user.FollowedTeacherAt.Valid {
		followedTeacherAt = user.FollowedTeacherAt.Time
	}
	return &api_v1.User{
		Id:            int32(user.ID),
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		MPlan: &api_v1.MPlan{
			Id:   int32(user.PlanID),
			Name: "",
		},
		FollowedTeacherAt: timestamppb.New(followedTeacherAt),
	}
}

func User2Proto(user *model2.User) *api_v1.User {
	emailVerified := false
	if user.EmailVerified == 1 {
		emailVerified = true
	}
	followedTeacherAt := time.Time{}
	if user.FollowedTeacherAt.Valid {
		followedTeacherAt = user.FollowedTeacherAt.Time
	}
	return &api_v1.User{
		Id:            int32(user.ID),
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: emailVerified,
		MPlan: &api_v1.MPlan{
			Id:   int32(user.PlanID),
			Name: "",
		},
		FollowedTeacherAt: timestamppb.New(followedTeacherAt),
	}
}
