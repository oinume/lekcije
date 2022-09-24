package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/morikuni/failure"

	"github.com/oinume/lekcije/backend/context_data"
	lerrors "github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/interface/graphql/model"
	"github.com/oinume/lekcije/backend/model2"
)

// CreateFollowingTeacher is the resolver for the createFollowingTeacher field.
func (r *mutationResolver) CreateFollowingTeacher(ctx context.Context, input model.CreateFollowingTeacherInput) (*model.CreateFollowingTeacherPayload, error) {
	user, err := authenticateFromContext(ctx, r.userUsecase)
	if err != nil {
		return nil, err
	}

	teacherIDOrURL := input.TeacherIDOrURL
	if teacherIDOrURL == "" {
		return nil, failure.New(lerrors.InvalidArgument, failure.Messagef("講師のURLまたはIDを入力して下さい"))
	}
	teacher, err := model2.NewTeacherFromIDOrURL(teacherIDOrURL)
	if err != nil {
		return nil, failure.New(lerrors.InvalidArgument, failure.Messagef("講師のURLまたはIDが正しくありません"))
	}

	followingTeacher, updateFollowedTeacherAt, err := r.followingTeacherUsecase.FollowTeacher(ctx, user, teacher)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := r.gaMeasurementUsecase.SendEvent(
			ctx, context_data.MustGAMeasurementEvent(ctx),
			model2.GAMeasurementEventCategoryFollowingTeacher,
			"follow", fmt.Sprint(teacher.ID), 1, uint32(user.ID),
		); err != nil {
			panic(err)
		}
		if updateFollowedTeacherAt {
			if err := r.gaMeasurementUsecase.SendEvent(
				ctx, context_data.MustGAMeasurementEvent(ctx),
				model2.GAMeasurementEventCategoryUser,
				"followFirstTime", fmt.Sprint(user.ID), 0, uint32(user.ID),
			); err != nil {
				panic(err)
			}
		}
	}()

	return &model.CreateFollowingTeacherPayload{
		ID:        followingTeacher.ID(),
		TeacherID: fmt.Sprint(teacher.ID),
	}, nil
}

// DeleteFollowingTeachers is the resolver for the deleteFollowingTeachers field.
func (r *mutationResolver) DeleteFollowingTeachers(ctx context.Context, input model.DeleteFollowingTeachersInput) (*model.DeleteFollowingTeachersPayload, error) {
	user, err := authenticateFromContext(ctx, r.userUsecase)
	if err != nil {
		return nil, err
	}
	ids := input.TeacherIds
	if len(ids) == 0 {
		return &model.DeleteFollowingTeachersPayload{}, nil
	}
	teacherIDs := make([]uint, len(ids))
	for i, id := range ids {
		tid, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return nil, failure.Translate(err, lerrors.Internal)
		}
		teacherIDs[i] = uint(tid)
	}

	if err := r.followingTeacherUsecase.DeleteFollowingTeachers(ctx, user.ID, teacherIDs); err != nil {
		return nil, err
	}

	go func() {
		if err := r.gaMeasurementUsecase.SendEvent(
			ctx, context_data.MustGAMeasurementEvent(ctx),
			model2.GAMeasurementEventCategoryFollowingTeacher,
			"unfollow",
			strings.Join(input.TeacherIds, ","),
			1,
			uint32(user.ID),
		); err != nil {
			panic(err)
		}
	}()

	return &model.DeleteFollowingTeachersPayload{
		TeacherIds: ids,
	}, nil
}

// FollowingTeachers is the resolver for the followingTeachers field.
func (r *queryResolver) FollowingTeachers(ctx context.Context) ([]*model.FollowingTeacher, error) {
	//teachers, err := r.followingTeacherRepo.FindTeachersByUserID()
	return []*model.FollowingTeacher{
		{
			ID: "1",
			Teacher: &model.Teacher{
				ID:   "12345",
				Name: "oinume",
			},
			CreatedAt: "",
		},
	}, nil
}
