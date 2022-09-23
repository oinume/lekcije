package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/morikuni/failure"

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

	followingTeacher, _, err := r.followingTeacherUsecase.FollowTeacher(ctx, user, teacher)
	if err != nil {
		return nil, err
	}

	//go func() {
	//	if err := r.gaMeasurementUsecase.SendEvent(
	//		ctx, context_data.MustGAMeasurementEvent(ctx),
	//		model2.GAMeasurementEventCategoryFollowingTeacher,
	//		"follow", fmt.Sprint(teacher.ID), 1, user.ID,
	//	); err != nil {
	//		panic(err)
	//	}
	//	if updateFollowedTeacherAt {
	//		if err := s.gaMeasurementUsecase.SendEvent(
	//			ctx, context_data.MustGAMeasurementEvent(ctx),
	//			model2.GAMeasurementEventCategoryUser,
	//			"followFirstTime", fmt.Sprint(user.ID), 0, user.ID,
	//		); err != nil {
	//			panic(err)
	//		}
	//	}
	//}()
	return &model.CreateFollowingTeacherPayload{
		ID: followingTeacher.ID(),
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
