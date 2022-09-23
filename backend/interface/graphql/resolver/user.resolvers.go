package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"time"

	lerrors "github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/interface/graphql/generated"
	"github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/model2"
)

// FollowingTeachers is the resolver for the followingTeachers field.
func (r *userResolver) FollowingTeachers(ctx context.Context, obj *model.User, first *int, after *string, last *int, before *string) (*model.FollowingTeacherConnection, error) {
	userID, err := strconv.ParseUint(obj.ID, 10, 32)
	if err != nil {
		return nil, err
	}
	fts, err := r.followingTeacherRepo.FindByUserID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}
	followingTeachers := make([]*model.FollowingTeacher, len(fts))
	teacherIDs := make([]uint, len(fts))
	for i, ft := range fts {
		teacherIDs[i] = ft.TeacherID
	}
	teachers, err := r.teacherRepo.FindByIDs(ctx, teacherIDs)
	if err != nil {
		return nil, err
	}
	teachersMap := make(map[uint]*model2.Teacher, len(teachers))
	for _, t := range teachers {
		teachersMap[t.ID] = t
	}

	for i, ft := range fts {
		t, ok := teachersMap[ft.TeacherID]
		if !ok {
			continue
		}
		followingTeachers[i] = &model.FollowingTeacher{
			ID: fmt.Sprintf("%d_%d", ft.UserID, ft.TeacherID),
			Teacher: &model.Teacher{
				ID:   fmt.Sprint(t.ID),
				Name: t.Name,
			},
			CreatedAt: ft.CreatedAt.String(),
		}
	}

	return &model.FollowingTeacherConnection{
		Nodes: followingTeachers,
	}, nil
}

// NotificationTimeSpans is the resolver for the notificationTimeSpans field.
func (r *userResolver) NotificationTimeSpans(ctx context.Context, obj *model.User) ([]*model.NotificationTimeSpan, error) {
	userID, err := strconv.ParseUint(obj.ID, 10, 32)
	if err != nil {
		return nil, err
	}
	timeSpans, err := r.notificationTimeSpanRepo.FindByUserID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}

	gqlTimeSpans := make([]*model.NotificationTimeSpan, len(timeSpans))
	for i, nts := range timeSpans {
		fromTime, err := time.Parse("15:04:05", nts.FromTime)
		if err != nil {
			return nil, lerrors.NewInternalError(
				lerrors.WithError(err),
				lerrors.WithMessagef("Invalid time format: FromTime=%v", nts.FromTime),
			)
		}
		toTime, err := time.Parse("15:04:05", nts.ToTime)
		if err != nil {
			return nil, lerrors.NewInternalError(
				lerrors.WithError(err),
				lerrors.WithMessagef("Invalid time format: ToTime=%v", nts.ToTime),
			)
		}
		gqlTimeSpans[i] = &model.NotificationTimeSpan{
			FromHour:   fromTime.Hour(),
			FromMinute: fromTime.Minute(),
			ToHour:     toTime.Hour(),
			ToMinute:   toTime.Minute(),
		}
	}

	return gqlTimeSpans, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
